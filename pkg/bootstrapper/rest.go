package bootstrapper

import (
	"net/http"

	"github.com/RakaMurdiarta/online-shop-system/internal/config"
	"github.com/RakaMurdiarta/online-shop-system/internal/middlewares"
	arp "github.com/RakaMurdiarta/online-shop-system/internal/modules/articles/provider"
	ap "github.com/RakaMurdiarta/online-shop-system/internal/modules/auth/provider"
	authServiceImpl "github.com/RakaMurdiarta/online-shop-system/internal/modules/auth/services/impl"
	cp "github.com/RakaMurdiarta/online-shop-system/internal/modules/cart/provider"
	cartRepoImpl "github.com/RakaMurdiarta/online-shop-system/internal/modules/cart/repository/impl"
	cartServiceImpl "github.com/RakaMurdiarta/online-shop-system/internal/modules/cart/services/impl"
	fp "github.com/RakaMurdiarta/online-shop-system/internal/modules/feedback/provider"
	mp "github.com/RakaMurdiarta/online-shop-system/internal/modules/mailer/provider"
	op "github.com/RakaMurdiarta/online-shop-system/internal/modules/orders/provider"
	orderRepoImpl "github.com/RakaMurdiarta/online-shop-system/internal/modules/orders/repository/impl"
	orderServiceImpl "github.com/RakaMurdiarta/online-shop-system/internal/modules/orders/services/impl"
	pp "github.com/RakaMurdiarta/online-shop-system/internal/modules/products/provider"
	productRepoImpl "github.com/RakaMurdiarta/online-shop-system/internal/modules/products/repository/impl"
	productServiceImpl "github.com/RakaMurdiarta/online-shop-system/internal/modules/products/services/impl"
	usp "github.com/RakaMurdiarta/online-shop-system/internal/modules/upload/provider"
	up "github.com/RakaMurdiarta/online-shop-system/internal/modules/users/provider"
	userRepoImpl "github.com/RakaMurdiarta/online-shop-system/internal/modules/users/repository/impl"
	userServiceImpl "github.com/RakaMurdiarta/online-shop-system/internal/modules/users/services/impl"
	"github.com/RakaMurdiarta/online-shop-system/pkg/database"
	"github.com/RakaMurdiarta/online-shop-system/pkg/mailer"
	"github.com/RakaMurdiarta/online-shop-system/pkg/shared"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

type Server struct {
	DB            *gorm.DB
	e             *echo.Echo
	conf          *config.Config
	storageClient *shared.SupabaseStorageClient
	mailTransport mailer.Transport
}

func NewServer(e *echo.Echo, c *config.Config, db *gorm.DB, storageClient *shared.SupabaseStorageClient, mailTransport mailer.Transport) *Server {
	return &Server{e: e, conf: c, DB: db, storageClient: storageClient, mailTransport: mailTransport}
}

func (s *Server) InitAPI() {
	txManager := database.NewTransactionManager(s.DB)
	private, public := s.initInternalRoute()
	xenditClient := shared.NewXenditClient(s.conf.XenditSecretKey)

	userRepo := userRepoImpl.NewUserRepository(txManager)
	categoryRepo := productRepoImpl.NewCategoryRepository(txManager)
	productRepo := productRepoImpl.NewProductRepository(txManager)
	cartRepo := cartRepoImpl.NewCartRepository(txManager)
	orderRepo := orderRepoImpl.NewNewOrderRepository(txManager)

	categoryService := productServiceImpl.NewCategoryService(categoryRepo, txManager, s.conf)
	productService := productServiceImpl.NewProductService(productRepo, categoryRepo)
	cartService := cartServiceImpl.NewNewCartService(cartRepo, productRepo)
	orderService := orderServiceImpl.NewOrderService(orderRepo, xenditClient)
	orderCallbackService := orderServiceImpl.NewOrderCallbackService(orderRepo, xenditClient)
	authService := authServiceImpl.NewAuthService(userRepo, s.conf)
	userService := userServiceImpl.NewUserService(userRepo)

	ap.AuthProvide(private, public, s.conf, userRepo, authService)
	arp.ArticleProvider(txManager, private, public)
	pp.ProductProvider(s.DB, private, public, txManager, userRepo, categoryRepo, categoryService, productService, s.conf, s.storageClient)
	cp.CartProvider(txManager, private, productRepo, cartRepo, cartService)
	op.OrderProvider(txManager, private, public, orderRepo, orderService, orderCallbackService)
	up.UserProvider(private, txManager, s.conf, userService)
	usp.UploadProvider(private, s.storageClient)

	mailerService := mp.MailerProvider(s.mailTransport)
	fp.FeedbackProvider(txManager, public, mailerService)
}

func (s *Server) initInternalRoute() (keyWithJWT *echo.Group, v1 *echo.Group) {

	s.e.GET("/health", healthFunc)

	//TODO: SECURITY HEADER like Helmet
	//TODO: HOST Validate

	api := s.e.Group("/api")
	s.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "https://mart-volt-shop.vercel.app"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "ngrok-skip-browser-warning"},
		AllowCredentials: true,
	}))
	s.e.Use(middleware.Recover())             //handling panic error
	s.e.Use(middleware.RemoveTrailingSlash()) // remove / in the end or endpoint url
	s.e.Use(middleware.RequestLogger())       // middleware logger -> HTTP Request Logger
	v1 = api.Group("/v1")

	keyWithJWT = v1.Group("")
	keyWithJWT.Use(middlewares.AuthMiddleware(s.conf.JwtSecretKey))

	return keyWithJWT, v1

}

func healthFunc(c *echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
