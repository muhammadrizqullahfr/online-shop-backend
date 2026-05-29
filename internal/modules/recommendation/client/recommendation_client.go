package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/RakaMurdiarta/online-shop-system/internal/config"
	"github.com/RakaMurdiarta/online-shop-system/internal/models"
)

// FastAPIProduct represents the product format expected by FastAPI
type FastAPIProduct struct {
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
}

type RecommendRequest struct {
	Query  string `json:"query"`
	UserID *int   `json:"user_id,omitempty"`
	Limit  int    `json:"limit"`
}

type RecommendResult struct {
	ID string `json:"id"`
}

type RecommendResponse struct {
	Results []RecommendResult `json:"results"`
}

type TrackRequest struct {
	UserID    int    `json:"user_id"`
	ProductID string `json:"product_id"`
	Action    string `json:"action"` // e.g., "click", "view", "purchase"
}

type ChatbotRequest struct {
	Query  string `json:"query"`
	UserID *int   `json:"user_id,omitempty"`
}

type ChatbotResponse struct {
	Response       string           `json:"response"`
	SourceProducts []FastAPIProduct `json:"source_products"`
}

type RecommendationClient interface {
	SyncProductToFastAPI(p models.Product) error
	BulkSyncProductsToFastAPI(products []models.Product) error
	GetRecommendations(query string, userID *int, limit int) ([]string, error)
	SendTrackEvent(tr TrackRequest) error
	QueryChatbot(cr ChatbotRequest) (ChatbotResponse, error)
}

type recommendationClientImpl struct {
	cfg           *config.Config
	httpClient    *http.Client
	chatbotClient *http.Client
}

func NewRecommendationClient(cfg *config.Config) RecommendationClient {
	return &recommendationClientImpl{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		chatbotClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (c *recommendationClientImpl) toFastAPIProduct(p models.Product) FastAPIProduct {
	var price float64
	if p.Price != nil {
		price = *p.Price
	}

	tags := []string(p.Tags)
	if tags == nil {
		tags = []string{}
	}

	return FastAPIProduct{
		ID:            p.ID.String(),
		Name:          p.Name,
		Description:   p.Description,
		CategoryID:    int(p.CategoryID),
		Price:         price,
		OriginalPrice: p.OriginalPrice,
		Image:         p.Image,
		Rating:        p.Rating,
		Reviews:       p.Reviews,
		Stock:         p.Stock,
		InStock:       p.InStock,
		Featured:      p.Featured,
		Tags:          tags,
	}
}

func (c *recommendationClientImpl) SyncProductToFastAPI(p models.Product) error {
	if c.cfg.RecommendationURL == "" {
		return errors.New("recommendation URL is not configured")
	}

	url := fmt.Sprintf("%s/api/v1/products/ingest", c.cfg.RecommendationURL)
	fastProduct := c.toFastAPIProduct(p)

	payloadBytes, err := json.Marshal(fastProduct)
	if err != nil {
		return fmt.Errorf("failed to marshal product to JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.cfg.RecommendationAPIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("connection to recommendation service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("recommendation service responded with error status: %d", resp.StatusCode)
	}

	return nil
}

func (c *recommendationClientImpl) BulkSyncProductsToFastAPI(products []models.Product) error {
	if c.cfg.RecommendationURL == "" {
		return errors.New("recommendation URL is not configured")
	}

	url := fmt.Sprintf("%s/api/v1/products/ingest/bulk", c.cfg.RecommendationURL)

	fastProducts := make([]FastAPIProduct, 0, len(products))
	for _, p := range products {
		fastProducts = append(fastProducts, c.toFastAPIProduct(p))
	}

	payload := map[string][]FastAPIProduct{
		"products": fastProducts,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal bulk products to JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create bulk HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.cfg.RecommendationAPIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to contact recommendation service for bulk ingest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("recommendation service bulk ingest error status: %d", resp.StatusCode)
	}

	return nil
}

func (c *recommendationClientImpl) GetRecommendations(query string, userID *int, limit int) ([]string, error) {
	if c.cfg.RecommendationURL == "" {
		return nil, errors.New("recommendation URL is not configured")
	}

	url := fmt.Sprintf("%s/api/v1/recommend/fts", c.cfg.RecommendationURL)

	reqPayload := RecommendRequest{
		Query:  query,
		UserID: userID,
		Limit:  limit,
	}

	payloadBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query to JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to build recommendation HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.cfg.RecommendationAPIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to contact recommendation service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("recommendation service returned error status: %d", resp.StatusCode)
	}

	var recResp RecommendResponse
	if err := json.NewDecoder(resp.Body).Decode(&recResp); err != nil {
		return nil, errors.New("failed to decode JSON response from recommendation service")
	}

	ids := make([]string, 0, len(recResp.Results))
	for _, res := range recResp.Results {
		ids = append(ids, res.ID)
	}

	return ids, nil
}
func (c *recommendationClientImpl) SendTrackEvent(tr TrackRequest) error {
	if c.cfg.RecommendationURL == "" {
		return errors.New("recommendation URL is not configured")
	}

	url := fmt.Sprintf("%s/api/v1/recommend/track", c.cfg.RecommendationURL)

	payload := map[string]any{
		"user_id":          tr.UserID,
		"product_id":       tr.ProductID,
		"interaction_type": tr.Action,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal tracking payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create tracking HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.cfg.RecommendationAPIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to contact recommendation service for tracking: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("recommendation service track responded with error status: %d", resp.StatusCode)
	}

	return nil
}

func (c *recommendationClientImpl) QueryChatbot(cr ChatbotRequest) (ChatbotResponse, error) {
	var chatResp ChatbotResponse
	if c.cfg.RecommendationURL == "" {
		return chatResp, errors.New("recommendation URL is not configured")
	}

	url := fmt.Sprintf("%s/api/v1/chatbot/query", c.cfg.RecommendationURL)

	payloadBytes, err := json.Marshal(cr)
	if err != nil {
		return chatResp, fmt.Errorf("failed to marshal chatbot query: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return chatResp, fmt.Errorf("failed to build chatbot HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.cfg.RecommendationAPIKey)

	resp, err := c.chatbotClient.Do(req)
	if err != nil {
		return chatResp, fmt.Errorf("failed to contact chatbot service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return chatResp, fmt.Errorf("chatbot service returned error status: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return chatResp, errors.New("failed to decode chatbot JSON response")
	}

	return chatResp, nil
}
