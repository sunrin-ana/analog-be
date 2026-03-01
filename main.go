package main

import (
	"analog-be/controller"
	"analog-be/entity"
	"analog-be/interceptor"
	"analog-be/pkg"
	"analog-be/repository"
	"analog-be/routes"
	"analog-be/service"
	"context"
	"crypto/tls"
	"database/sql"
	"os"
	"time"

	"github.com/NARUBROWN/spine/pkg/boot"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "analog-be/docs"

	"github.com/NARUBROWN/spine"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.uber.org/zap"
)

// @title AnAlog API
// @version 1.0.0
// @description Spine 기반 REST API

// @host localhost:8080
// @BasePath /api
func main() {
	if err := pkg.InitLogger(); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer pkg.Logger.Sync()

	logger := pkg.GetLogger()
	logger.Info("Starting Analog Server")

	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	ValidateEnvVars(logger)

	app := spine.New()

	db := newBunDB()

	db.RegisterModel(
		// relation
		(*entity.LogToUser)(nil),
		(*entity.LogToTopic)(nil),

		(*entity.Log)(nil),
		(*entity.Topic)(nil),
		(*entity.User)(nil),
		(*entity.Comment)(nil),
		(*entity.OAuthState)(nil),
		(*entity.Session)(nil),
	)

	app.Constructor(
		// 디비
		func() *bun.DB { return db },

		// 기타
		func() *zap.Logger { return logger },

		// 레포지토리
		repository.NewUserRepository,
		repository.NewLogRepository,
		repository.NewCommentRepository,
		repository.NewTopicRepository,
		repository.NewTokenRepository,

		// 서비스
		service.NewLogService,
		service.NewUserService,
		service.NewAnAccountOAuthService,
		service.NewCommentService,
		service.NewTopicService,
		service.NewAnAmericanoService,
		service.NewFeedService,
		service.NewTokenService,

		// 컨트롤러
		controller.NewHealthController,
		controller.NewLogController,
		controller.NewUserController,
		controller.NewAuthController,
		controller.NewTopicController,
		controller.NewFeedController,

		// 인터셉터
		interceptor.NewTxInterceptor,
		interceptor.NewAuthInterceptor,
	)

	// 전역 인터셉터
	app.Interceptor(
		interceptor.NewCORSInterceptor(),
		interceptor.NewRateLimitInterceptor(),
		interceptor.NewLoggingInterceptor(),
		interceptor.NewErrorInterceptor(),
	)

	routes.RegisterHealthRoutes(app)
	routes.RegisterLogRoutes(app)
	routes.RegisterUserRoutes(app)
	routes.RegisterAuthRoutes(app)
	routes.RegisterTopicRoutes(app)
	routes.RegisterFeedRoutes(app)

	app.Transport(func(t any) {
		e := t.(*echo.Echo)
		e.GET("/docs/*", echo.WrapHandler(httpSwagger.WrapHandler))
	})

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Server starting", zap.String("port", port))
	app.Run(boot.Options{
		Address:                ":" + port,
		EnableGracefulShutdown: true,
		HTTP: &boot.HTTPOptions{
			GlobalPrefix: "/api",
		},
	})
}

func newBunDB() *bun.DB {
	host := GetEnv("DB_HOST", "localhost")
	port := GetEnv("DB_PORT", "5437")
	user := GetEnv("DB_USER", "test")
	password := GetEnv("DB_PASSWORD", "test")
	database := GetEnv("DB_NAME", "test")

	tlsConfig := &tls.Config{
		InsecureSkipVerify: os.Getenv("DB_TLS_SKIP_VERIFY") == "true",
	}

	pgconn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(host+":"+port),
		pgdriver.WithTLSConfig(tlsConfig),
		pgdriver.WithUser(user),
		pgdriver.WithPassword(password),
		pgdriver.WithDatabase(database),
		pgdriver.WithApplicationName("analog"),
		pgdriver.WithTimeout(10*time.Second),
		pgdriver.WithDialTimeout(5*time.Second),
		pgdriver.WithReadTimeout(10*time.Second),
		pgdriver.WithWriteTimeout(10*time.Second),
		pgdriver.WithInsecure(true),
	)

	sqldb := sql.OpenDB(pgconn)

	sqldb.SetMaxOpenConns(25)
	sqldb.SetMaxIdleConns(5)
	sqldb.SetConnMaxLifetime(5 * time.Minute)
	sqldb.SetConnMaxIdleTime(10 * time.Minute)

	db := bun.NewDB(sqldb, pgdialect.New())

	if os.Getenv("DEBUG") == "true" {
		db.AddQueryHook(&debugHook{})
	}

	return db
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

type debugHook struct{}

func (h *debugHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	return ctx
}

func (h *debugHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	logger := pkg.GetLogger()
	logger.Debug("Database query",
		zap.String("query", event.Query),
		zap.Duration("duration", time.Since(event.StartTime)),
	)
}

func ValidateEnvVars(logger *zap.Logger) {
	required := []string{
		"AN_ACCOUNT_BASE_URL",
		"AN_ACCOUNT_CLIENT_ID",
		"AN_ACCOUNT_CLIENT_SECRET",
	}

	for _, key := range required {
		if os.Getenv(key) == "" {
			logger.Fatal("Missing required environment variable",
				zap.String("variable", key))
		}
	}

	if os.Getenv("DB_TLS_SKIP_VERIFY") == "true" {
		logger.Warn("Database TLS certificate validation is disabled - only use in development")
	}
}
