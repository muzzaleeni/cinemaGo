package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/forChin/my-project/internal/config"
	"github.com/forChin/my-project/internal/handler"
	"github.com/forChin/my-project/internal/service"
	"github.com/forChin/my-project/internal/store"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger, err := newLogger(cfg)
	if err != nil {
		return fmt.Errorf("create logger: %w", err)
	}
	defer logger.Sync()

	sugarLogger := logger.Sugar()

	setCustomConverters()

	dbURL := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d database=%s sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost,
		cfg.DBPort, cfg.DBName, cfg.DBSSLMode,
	)
	connConfig, err := pgx.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("parse connection config: %w", err)
	}

	// setting describe mode for pgBouncer (https://github.com/jackc/pgx/issues/602)
	connConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	connStr := stdlib.RegisterConnConfig(connConfig)
	db, err := sqlx.Connect("pgx", connStr)
	if err != nil {
		return fmt.Errorf("connect to db: %w", err)
	}

	app := fiber.New(fiber.Config{
		BodyLimit:    cfg.HTTPReqBodySizeLimit, // max request body size
		ErrorHandler: newCustomHTTPErrorHandler(sugarLogger),
	})
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,

		// handler for logging occured panics
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			buf := make([]byte, 1024)
			buf = buf[:runtime.Stack(buf, false)]
			sugarLogger.Errorf("panic: %v\n%s\n", e, buf)
		},
	}))
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// endpoints
	userStore := store.NewUserStore(db)
	userService := service.NewUserService(userStore)
	userHandler := handler.NewUserHandler(userService, sugarLogger)
	v1.Get("/users", userHandler.GetAll)
	v1.Post("/users", userHandler.Create)
	v1.Patch("/users/:id", userHandler.Update)
	v1.Delete("/users/:id", userHandler.Delete)

	addr := fmt.Sprintf("%s:%d", cfg.HTTPHost, cfg.HTTPPort)
	err = app.Listen(addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	return nil
}

func loadConfig() (*config.Config, error) {
	// load
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// print config to stdout
	configBytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	fmt.Println("Configuration:", string(configBytes))

	return cfg, nil
}

// newLogger creates new logger with rotation.
func newLogger(cfg *config.Config) (*zap.Logger, error) {
	zcfg := zap.NewProductionConfig()
	zcfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// logrotate
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.LoggerOutput,
		MaxSize:    cfg.LoggerRotateMaxSize,
		MaxBackups: cfg.LoggerRotateMaxBackups,
		MaxAge:     cfg.LoggerRotateMaxAge,
		Compress:   cfg.LoggerRotateWithCompress,
	})

	lvl, err := zapcore.ParseLevel(cfg.LoggerLevel)
	if err != nil {
		return nil, fmt.Errorf("parse level: %w", err)
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zcfg.EncoderConfig),
		zapcore.AddSync(w),
		lvl,
	)

	return zap.New(core, zap.AddCaller()), nil
}

// newCustomHTTPErrorHandler returns handler that will be
// executed when an error is returned from our http handlers.
func newCustomHTTPErrorHandler(logger *zap.SugaredLogger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		var ferr *fiber.Error
		if errors.As(err, &ferr) {
			if ferr.Code == http.StatusInternalServerError {
				// Don't need to return message when internal server error occured.
				return c.Status(http.StatusInternalServerError).Send(nil)
			}

			if ferr.Message == utils.StatusMessage(ferr.Code) {
				// Don't need to return default (filled by fiber) message.
				return c.Status(ferr.Code).Send(nil)
			}

			// Return statuscode with error message
			return c.Status(ferr.Code).JSON(fiber.Map{
				"message": ferr.Message,
			})
		}

		logger.Errorf("unknown error was occured: %v", err)

		// Return statuscode
		return c.Status(http.StatusInternalServerError).Send(nil)
	}
}

// setCustomConverters sets http query params converters for custom data types.
func setCustomConverters() {
	var timeConverter = func(value string) reflect.Value {
		if v, err := time.Parse("2006-01-02 15:04:05.000", value); err == nil {
			return reflect.ValueOf(v)
		}
		return reflect.Value{}
	}

	customTime := fiber.ParserType{
		Customtype: time.Time{},
		Converter:  timeConverter,
	}

	// Add setting to the Decoder
	fiber.SetParserDecoder(fiber.ParserConfig{
		IgnoreUnknownKeys: true,
		ParserType: []fiber.ParserType{
			customTime,
		},
		ZeroEmpty: true,
	})
}
