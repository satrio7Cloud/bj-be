package di

import (
	"fmt"

	"bj/internal/config"
	"bj/internal/constants"
	"bj/internal/database"
	"bj/internal/pkg/minio"
	"bj/internal/pkg/qris"
	"bj/internal/pkg/redis"
	"bj/internal/shared/storage"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	redisPkg "github.com/redis/go-redis/v9"

	adminHandler "bj/internal/domain/user/handler"
	adminRepository "bj/internal/domain/user/repository"
	adminUsecase "bj/internal/domain/user/usecase"

	brandHandler "bj/internal/domain/brand/handler"
	brandRepository "bj/internal/domain/brand/repository"
	brandUsecase "bj/internal/domain/brand/usecase"

	customerHandler "bj/internal/domain/customer/handler"
	customerRepository "bj/internal/domain/customer/repository"
	customerUsecase "bj/internal/domain/customer/usecase"

	documentHandler "bj/internal/domain/document/handler"
	documentRepository "bj/internal/domain/document/repository"
	documentUsecase "bj/internal/domain/document/usecase"

	modelHandler "bj/internal/domain/model/handler"
	modelRepository "bj/internal/domain/model/repository"
	modelUsecase "bj/internal/domain/model/usecase"

	orderHandler "bj/internal/domain/order/handler"
	orderRepository "bj/internal/domain/order/repository"
	orderUsecase "bj/internal/domain/order/usecase"

	orderlogRepository "bj/internal/domain/orderlog/repository"

	serviceHandler "bj/internal/domain/service/handler"
	serviceRepository "bj/internal/domain/service/repository"
	serviceUsecase "bj/internal/domain/service/usecase"

	vehicleHandler "bj/internal/domain/vehicle/handler"
	vehicleRepository "bj/internal/domain/vehicle/repository"
	vehicleUsecase "bj/internal/domain/vehicle/usecase"

	dashboardHandler "bj/internal/domain/dashboard/handler"
	dashboardRepository "bj/internal/domain/dashboard/repository"
	dashboardUsecase "bj/internal/domain/dashboard/usecase"

	invoiceHandler "bj/internal/domain/invoice/handler"
	invoiceRepository "bj/internal/domain/invoice/repository"
	invoiceUsecase "bj/internal/domain/invoice/usecase"

	notificationHandler "bj/internal/domain/notification/handler"
	notificationRepository "bj/internal/domain/notification/repository"
	notificationUsecase "bj/internal/domain/notification/usecase"

	biteshipHandler "bj/internal/domain/biteship/handler"
	biteshipService "bj/internal/domain/biteship/service"
	biteshipUsecase "bj/internal/domain/biteship/usecase"

	paymentHandler "bj/internal/domain/payment/handler"
	paymentUsecase "bj/internal/domain/payment/usecase"

	companyHandler "bj/internal/domain/company/handler"
	companyRepository "bj/internal/domain/company/repository"
	companyUsecase "bj/internal/domain/company/usecase"

	"gorm.io/gorm"
)

// Container adalah pusat Dependency Injection
type Container struct {
	Config *config.Config
	DB     *gorm.DB
	Redis  *redisPkg.Client

	// Admin Domain
	AdminRepo adminRepository.AdminRepository
	AdminUC   adminUsecase.AdminUsecase
	AdminHdl  *adminHandler.AdminHandler

	// Brand Domain
	BrandRepo brandRepository.BrandRepository
	BrandUC   brandUsecase.BrandUsecase
	BrandHdl  *brandHandler.BrandHandler

	// Customer Domain
	CustomerRepo customerRepository.CustomerRepository
	CustomerUC   customerUsecase.CustomerUsecase
	CustomerHdl  *customerHandler.CustomerHandler

	// Document Domain
	DocumentRepo   documentRepository.DocumentRepository
	StorageService storage.StorageService
	DocumentUC     documentUsecase.DocumentUsecase
	DocumentHdl    *documentHandler.DocumentHandler

	// Model Domain
	ModelRepo modelRepository.ModelRepository
	ModelUC   modelUsecase.ModelUsecase
	ModelHdl  *modelHandler.ModelHandler

	// Order Domain
	OrderRepo orderRepository.OrderRepository
	OrderUC   orderUsecase.OrderUsecase
	OrderHdl  *orderHandler.OrderHandler

	// Service Domain
	ServiceRepo serviceRepository.ServiceRepository
	ServiceUC   serviceUsecase.ServiceUsecase
	ServiceHdl  *serviceHandler.ServiceHandler

	// Vehicle Domain
	VehicleRepo vehicleRepository.VehicleRepository
	VehicleUC   vehicleUsecase.VehicleUsecase
	VehicleHdl  *vehicleHandler.VehicleHandler

	// Dashboard Domain
	DashboardRepo dashboardRepository.DashboardRepository
	DashboardUC   dashboardUsecase.DashboardUsecase
	DashboardHdl  *dashboardHandler.DashboardHandler

	// Invoice Domain
	InvoiceRepo invoiceRepository.InvoiceRepository
	InvoiceUC   invoiceUsecase.InvoiceUsecase
	InvoiceHdl  *invoiceHandler.InvoiceHandler

	// Notification Domain
	NotificationRepo notificationRepository.NotificationRepository
	NotificationUC   notificationUsecase.NotificationUsecase
	NotificationHdl  *notificationHandler.NotificationHandler

	// Biteship Domain
	BiteshipSvc biteshipService.BiteshipService
	BiteshipUC  biteshipUsecase.BiteshipUsecase
	BiteshipHdl *biteshipHandler.BiteshipHandler

	// Payment Domain
	PaymentUC  paymentUsecase.PaymentUsecase
	PaymentHdl *paymentHandler.PaymentHandler

	// Company Profile Domain
	CompanyRepo companyRepository.CompanyProfileRepository
	CompanyUC   companyUsecase.CompanyProfileUsecase
	CompanyHdl  *companyHandler.CompanyProfileHandler
}

func NewContainer(cfg *config.Config) (*Container, error) {
	if cfg == nil {
		cfg = config.LoadConfig()
	}

	db := database.ConnectSQLServer(cfg)
	if db == nil {
		return nil, fmt.Errorf(constants.MsgDbDisconnected)
	}

	// Setup MinIO & Storage Service
	minioClient, err := minio.NewMinioClient(cfg)
	if err != nil {
		return nil, fmt.Errorf(constants.MsgMinioInitFailed, err)
	}
	storageService := storage.NewMinioStorage(minioClient, cfg)

	// Setup Redis
	redisClient, err := redis.NewRedisClient(cfg)
	if err != nil {
		return nil, fmt.Errorf(constants.MsgRedisInitFailed, err)
	}

	// Setup Qris Generator
	qrisGenerator := qris.NewQrisGenerator(cfg)

	// Setup Midtrans Snap Client
	var snapClient snap.Client
	midtransEnv := midtrans.Sandbox
	if cfg.Midtrans.Env == "production" {
		midtransEnv = midtrans.Production
	}
	snapClient.New(cfg.Midtrans.ServerKey, midtransEnv)

	// Setup Admin Domain
	adminRepo := adminRepository.NewAdminRepository(db)
	adminUC := adminUsecase.NewAdminUsecase(adminRepo)
	adminHdl := adminHandler.NewAdminHandler(adminUC)

	brandRepo := brandRepository.NewBrandRepository(db)
	brandUC := brandUsecase.NewBrandUsecase(brandRepo)
	brandHdl := brandHandler.NewBrandHandler(brandUC)

	customerRepo := customerRepository.NewCustomerRepository(db)
	customerUC := customerUsecase.NewCustomerUsecase(customerRepo)
	customerHdl := customerHandler.NewCustomerHandler(customerUC)

	modelRepo := modelRepository.NewModelRepository(db)
	modelUC := modelUsecase.NewModelUsecase(modelRepo)
	modelHdl := modelHandler.NewModelHandler(modelUC)

	vehicleRepo := vehicleRepository.NewVehicleRepository(db)
	vehicleUC := vehicleUsecase.NewVehicleUsecase(vehicleRepo)
	vehicleHdl := vehicleHandler.NewVehicleHandler(vehicleUC)

	biteshipSvc := biteshipService.NewBiteshipService(cfg)
	biteshipUC := biteshipUsecase.NewBiteshipUsecase(
		biteshipSvc,
		cfg.Biteship.DefaultCouriers,
		cfg.Biteship.OriginPostalCode,
		cfg.Biteship.DefaultItemName,
		cfg.Biteship.DefaultItemWeight,
		cfg.Biteship.DefaultItemValue,
	)
	biteshipHdl := biteshipHandler.NewBiteshipHandler(biteshipUC)

	orderRepo := orderRepository.NewOrderRepository(db)
	orderLogRepo := orderlogRepository.NewOrderLogRepository(db)
	serviceRepo := serviceRepository.NewServiceRepository(db)
	orderUC := orderUsecase.NewOrderUsecase(orderRepo, orderLogRepo, serviceRepo, customerRepo, vehicleRepo, biteshipUC, storageService, snapClient)
	orderHdl := orderHandler.NewOrderHandler(orderUC)

	serviceUC := serviceUsecase.NewServiceUsecase(serviceRepo, redisClient)
	serviceHdl := serviceHandler.NewServiceHandler(serviceUC)

	documentRepo := documentRepository.NewDocumentRepository(db)
	documentUC := documentUsecase.NewDocumentUsecase(documentRepo, orderRepo, serviceRepo, storageService)
	documentHdl := documentHandler.NewDocumentHandler(documentUC)

	dashboardRepo := dashboardRepository.NewDashboardRepository(db)
	dashboardUC := dashboardUsecase.NewDashboardUsecase(dashboardRepo)
	dashboardHdl := dashboardHandler.NewDashboardHandler(dashboardUC)

	invoiceRepo := invoiceRepository.NewInvoiceRepository(db)
	invoiceUC := invoiceUsecase.NewInvoiceUsecase(invoiceRepo, orderRepo, customerRepo, vehicleRepo, serviceRepo, qrisGenerator, cfg)
	invoiceHdl := invoiceHandler.NewInvoiceHandler(invoiceUC)

	notificationRepo := notificationRepository.NewNotificationRepository(db)
	notificationUC := notificationUsecase.NewNotificationUsecase(notificationRepo, cfg.NotificationCron)
	notificationHdl := notificationHandler.NewNotificationHandler(notificationUC)

	paymentUC := paymentUsecase.NewPaymentUsecase(orderUC)
	paymentHdl := paymentHandler.NewPaymentHandler(paymentUC)

	companyRepo := companyRepository.NewCompanyProfileRepository(db, cfg.Office)
	companyUC := companyUsecase.NewCompanyProfileUsecase(companyRepo)
	companyHdl := companyHandler.NewCompanyProfileHandler(companyUC)

	// Start Cronjob Scheduler in background
	go notificationUC.StartNotificationScheduler()

	// Return Container
	return &Container{
		Config: cfg,
		DB:     db,
		Redis:  redisClient,

		AdminRepo:        adminRepo,
		AdminUC:          adminUC,
		AdminHdl:         adminHdl,
		BrandRepo:        brandRepo,
		BrandUC:          brandUC,
		BrandHdl:         brandHdl,
		CustomerRepo:     customerRepo,
		CustomerUC:       customerUC,
		CustomerHdl:      customerHdl,
		DocumentRepo:     documentRepo,
		StorageService:   storageService,
		DocumentUC:       documentUC,
		DocumentHdl:      documentHdl,
		ModelRepo:        modelRepo,
		ModelUC:          modelUC,
		ModelHdl:         modelHdl,
		OrderRepo:        orderRepo,
		OrderUC:          orderUC,
		OrderHdl:         orderHdl,
		ServiceRepo:      serviceRepo,
		ServiceUC:        serviceUC,
		ServiceHdl:       serviceHdl,
		VehicleRepo:      vehicleRepo,
		VehicleUC:        vehicleUC,
		VehicleHdl:       vehicleHdl,
		DashboardRepo:    dashboardRepo,
		DashboardUC:      dashboardUC,
		DashboardHdl:     dashboardHdl,
		InvoiceRepo:      invoiceRepo,
		InvoiceUC:        invoiceUC,
		InvoiceHdl:       invoiceHdl,
		NotificationRepo: notificationRepo,
		NotificationUC:   notificationUC,
		NotificationHdl:  notificationHdl,
		BiteshipSvc:      biteshipSvc,
		BiteshipUC:       biteshipUC,
		BiteshipHdl:      biteshipHdl,
		PaymentUC:        paymentUC,
		PaymentHdl:       paymentHdl,
		CompanyRepo:      companyRepo,
		CompanyUC:        companyUC,
		CompanyHdl:       companyHdl,
	}, nil
}
