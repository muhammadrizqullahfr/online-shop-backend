package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/RakaMurdiarta/online-shop-system/internal/models"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/products/delivery"
	"github.com/RakaMurdiarta/online-shop-system/internal/modules/recommendation/client"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/customs"
	"github.com/RakaMurdiarta/online-shop-system/pkg/common/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RecommendationHandler struct {
	client client.RecommendationClient
	db     *gorm.DB
}

func NewRecommendationHandler(client client.RecommendationClient, db *gorm.DB) *RecommendationHandler {
	return &RecommendationHandler{
		client: client,
		db:     db,
	}
}

// RecommendProductHandler handles recommendations with a local DB search fallback
func (h *RecommendationHandler) RecommendProductHandler(c *echo.Context) error {
	query := c.QueryParam("q")
	userIDStr := c.QueryParam("user_id")

	var userID *int
	if userIDStr != "" {
		if id, err := strconv.Atoi(userIDStr); err == nil {
			userID = &id
		}
	}

	log.Printf("[Golang] Menerima kueri pencarian rekomendasi: %q, UserID: %v", query, userIDStr)

	// 1. Panggil Recommendation Client (FastAPI)
	recommendedIDs, err := h.client.GetRecommendations(query, userID, 10)
	if err != nil {
		log.Printf("[Golang] Recommendation Client error: %v. Mengaktifkan Fallback pencarian database lokal...", err)

		// Fallback local GORM search
		var products []models.Product
		dbErr := h.db.
			Preload("Category").
			Preload("Variants").
			Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").
			Order("rating DESC").
			Limit(10).
			Find(&products).Error

		if dbErr != nil {
			return c.JSON(http.StatusInternalServerError, response.NewResponseError(
				"Gagal memproses pencarian fallback",
				*customs.NewErrorValue("database", dbErr.Error()),
			))
		}

		responseItems := make([]delivery.NewProductResponse, 0, len(products))
		for i := range products {
			responseItems = append(responseItems, *delivery.ToNewProductResponse(&products[i]))
		}

		return c.JSON(http.StatusOK, map[string]any{
			"message":     "Layanan ML mengalami kendala, menyajikan fallback database lokal",
			"products":    responseItems,
			"is_fallback": true,
		})
	}

	// 2. Sukses menarik ID terurut dari FastAPI. Tarik data lengkap produk dari PostgreSQL
	if len(recommendedIDs) == 0 {
		return c.JSON(http.StatusOK, map[string]any{
			"message":     "Sukses memproses rekomendasi kosong dari FastAPI",
			"products":    []delivery.NewProductResponse{},
			"is_fallback": false,
		})
	}

	var products []models.Product
	dbErr := h.db.
		Preload("Category").
		Preload("Variants").
		Where("id IN ?", recommendedIDs).
		Find(&products).Error

	if dbErr != nil {
		return c.JSON(http.StatusInternalServerError, response.NewResponseError(
			"Gagal memuat detail produk rekomendasi",
			*customs.NewErrorValue("database", dbErr.Error()),
		))
	}

	// Buat map pencocokan ID agar urutan produk sesuai dengan rekomendasi FastAPI
	productMap := make(map[string]*models.Product)
	for i := range products {
		productMap[products[i].ID.String()] = &products[i]
	}

	orderedResponse := make([]delivery.NewProductResponse, 0, len(recommendedIDs))
	for _, id := range recommendedIDs {
		if p, ok := productMap[id]; ok {
			orderedResponse = append(orderedResponse, *delivery.ToNewProductResponse(p))
		}
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message":     "Sukses memproses rekomendasi reranked dari FastAPI",
		"products":    orderedResponse,
		"is_fallback": false,
	})
}

// TrackInteractionHandler handles recording buyer interaction tracking event
func (h *RecommendationHandler) TrackInteractionHandler(c *echo.Context) error {
	var tr client.TrackRequest
	if err := c.Bind(&tr); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Payload tracking tidak valid",
			*customs.NewErrorValue("validation", err.Error()),
		))
	}

	log.Printf("[Golang] Menerima event pelacakan: User %d - Produk %s - Aksi %s", tr.UserID, tr.ProductID, tr.Action)

	err := h.client.SendTrackEvent(tr)
	if err != nil {
		log.Printf("[Golang] Gagal meneruskan event pelacakan: %v", err)
		return c.JSON(http.StatusInternalServerError, response.NewResponseError(
			"Gagal mencatat event interaksi ke microservice",
			*customs.NewErrorValue("microservice", err.Error()),
		))
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "Event interaksi berhasil dicatat",
		"status":  "success",
	})
}

// ChatbotQueryHandler handles shopping chatbot query with AI-selected source products
func (h *RecommendationHandler) ChatbotQueryHandler(c *echo.Context) error {
	var cr client.ChatbotRequest
	if err := c.Bind(&cr); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewResponseError(
			"Payload kueri chatbot tidak valid",
			*customs.NewErrorValue("validation", err.Error()),
		))
	}

	log.Printf("[Golang] Menerima kueri chatbot pembeli: %q", cr.Query)

	chatResp, err := h.client.QueryChatbot(cr)
	if err != nil {
		log.Printf("[Golang] Chatbot AI mengalami kendala: %v. Mengaktifkan fallback tanggapan...", err)

		return c.JSON(http.StatusOK, map[string]any{
			"response":        "Maaf ya Kak, asisten belanja AI kami sedang beristirahat sebentar. Kakak tetap dapat mencari katalog produk terbaik kami secara langsung melalui fitur pencarian!",
			"source_products": []delivery.NewProductResponse{},
			"is_fallback":     true,
		})
	}

	// Ubah FastAPIProduct kembali ke model product response untuk keseragaman API frontend
	sourceIDs := make([]string, 0, len(chatResp.SourceProducts))
	for _, fp := range chatResp.SourceProducts {
		sourceIDs = append(sourceIDs, fp.ID)
	}

	var products []models.Product
	if len(sourceIDs) > 0 {
		_ = h.db.
			Preload("Category").
			Preload("Variants").
			Where("id IN ?", sourceIDs).
			Find(&products).Error
	}

	productMap := make(map[string]*models.Product)
	for i := range products {
		productMap[products[i].ID.String()] = &products[i]
	}

	orderedResponse := make([]delivery.NewProductResponse, 0, len(sourceIDs))
	for _, id := range sourceIDs {
		if p, ok := productMap[id]; ok {
			orderedResponse = append(orderedResponse, *delivery.ToNewProductResponse(p))
		}
	}

	return c.JSON(http.StatusOK, map[string]any{
		"response":        chatResp.Response,
		"source_products": orderedResponse,
		"is_fallback":     false,
	})
}

// CreateProductHandler handles simulation of product creation
func (h *RecommendationHandler) CreateProductHandler(c *echo.Context) error {
	var req struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Price       float64  `json:"price"`
		Rating      float64  `json:"rating"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	priceVal := req.Price
	p := models.Product{
		Name:          req.Name,
		Description:   req.Description,
		Price:         &priceVal,
		Rating:        req.Rating,
		CategoryID:    1, // Default to first category
		Image:         "https://images.unsplash.com/photo-1559056199-641a0ac8b55e?w=500",
		InStock:       true,
		Stock:         50,
	}

	if err := h.db.Create(&p).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}

	// Sync
	_ = h.client.SyncProductToFastAPI(p)

	return c.JSON(http.StatusCreated, map[string]any{
		"message":      "Produk berhasil disimpan dan disinkronisasikan ke indeks AI FastAPI.",
		"product":      p,
		"sync_success": true,
	})
}

// BulkCreateProductHandler handles bulk product creation and syncing with GORM conflict resolution
func (h *RecommendationHandler) BulkCreateProductHandler(c *echo.Context) error {
	var req struct {
		Products []struct {
			ID            string   `json:"id"`
			Name          string   `json:"name"`
			Description   string   `json:"description"`
			CategoryID    int      `json:"category_id"`
			Price         float64  `json:"price"`
			OriginalPrice *float64 `json:"original_price"`
			Image         string   `json:"image"`
			Rating        float64  `json:"rating"`
			Reviews       int      `json:"reviews"`
			Stock         int      `json:"stock"`
			InStock       bool     `json:"in_stock"`
			Featured      bool     `json:"featured"`
			Tags          []string `json:"tags"`
		} `json:"products"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	var products []models.Product
	for _, rp := range req.Products {
		priceVal := rp.Price
		prodID, _ := uuid.Parse(rp.ID)
		if prodID == uuid.Nil {
			prodID = uuid.New()
		}
		p := models.Product{
			ID:            prodID,
			Name:          rp.Name,
			Description:   rp.Description,
			Price:         &priceVal,
			OriginalPrice: rp.OriginalPrice,
			Image:         rp.Image,
			Rating:        rp.Rating,
			Reviews:       rp.Reviews,
			Stock:         rp.Stock,
			InStock:       rp.InStock,
			Featured:      rp.Featured,
			Tags:          models.StringSlice(rp.Tags),
		}
		if rp.CategoryID > 0 {
			p.CategoryID = uint(rp.CategoryID)
		} else {
			p.CategoryID = 1
		}
		products = append(products, p)
	}

	// Use OnConflict clause to avoid duplicate key conflicts (updates existing products instead of failing)
	if err := h.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}

	// Bulk Sync to FastAPI
	_ = h.client.BulkSyncProductsToFastAPI(products)

	return c.JSON(http.StatusCreated, map[string]any{
		"message":      "Daftar produk berhasil disimpan dan disinkronisasikan.",
		"products":     products,
		"sync_success": true,
	})
}
