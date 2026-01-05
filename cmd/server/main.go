package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/exaring/otelpgx"
	"github.com/gorilla/mux"
	"github.com/grafana/pyroscope-go"
	"github.com/jackc/pgx/v5/multitracer"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"gopkg.in/mail.v2"

	"github.com/OJOMB/fightpicker/internal/config"
	"github.com/OJOMB/fightpicker/internal/consumers/email"
	mediaconsumer "github.com/OJOMB/fightpicker/internal/consumers/media"
	repoauth "github.com/OJOMB/fightpicker/internal/repo/auth"
	repofighters "github.com/OJOMB/fightpicker/internal/repo/fighters"
	repousers "github.com/OJOMB/fightpicker/internal/repo/users"
	"github.com/OJOMB/fightpicker/internal/server"
	handlersv1auth "github.com/OJOMB/fightpicker/internal/server/handlers/v1/auth"
	handlersv1fighters "github.com/OJOMB/fightpicker/internal/server/handlers/v1/fighters"
	handlersv1users "github.com/OJOMB/fightpicker/internal/server/handlers/v1/users"
	serviceauth "github.com/OJOMB/fightpicker/internal/service/auth"
	servicefighters "github.com/OJOMB/fightpicker/internal/service/fighters"
	serviceusers "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/auth"
	postgresclient "github.com/OJOMB/fightpicker/pkg/clients/postgres"
	"github.com/OJOMB/fightpicker/pkg/datetimes"
	"github.com/OJOMB/fightpicker/pkg/id"
	"github.com/OJOMB/fightpicker/pkg/jsonwebtokens"
	"github.com/OJOMB/fightpicker/pkg/logs"
	mediaprocessor "github.com/OJOMB/fightpicker/pkg/media"
	"github.com/OJOMB/fightpicker/pkg/otel"
)

const appName = "fightpicker"

var otelShutdown func(ctx context.Context) error

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		log.Fatalf("env environment variable missing")
	}

	//////////////////////////////
	// setup config with viper //
	////////////////////////////

	viper.SetConfigName(env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/") // path to look for the config file

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	// enable env var for secrets and overrides
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	//////////////////////
	// PYROSCOPE SETUP //
	////////////////////

	if !cfg.Observability.Pyroscope.Enable {
		log.Printf("pyroscope disabled")
	} else {
		log.Printf("pyroscope enabled")
		pyroProf, err := pyroscope.Start(pyroscope.Config{
			ApplicationName: appName,
			// the OTel collector is actually Grafana Alloy, which has a built-in Pyroscope receiver
			ServerAddress: cfg.Observability.Pyroscope.Endpoint,
			// you can disable logging by setting this to nil
			Logger:   nil, //pyroscope.StandardLogger, -- disabled for less noisy logs
			TenantID: "tenant_123",
			ProfileTypes: []pyroscope.ProfileType{
				pyroscope.ProfileCPU,
				pyroscope.ProfileAllocObjects,
				pyroscope.ProfileAllocSpace,
				pyroscope.ProfileInuseObjects,
				pyroscope.ProfileInuseSpace,
				// these profile types are optional:
				pyroscope.ProfileGoroutines,
				pyroscope.ProfileMutexCount,
				pyroscope.ProfileMutexDuration,
				pyroscope.ProfileBlockCount,
				pyroscope.ProfileBlockDuration,
			},
		})
		if err != nil {
			log.Fatalf("failed to start pyroscope: %v", err)
		}

		defer pyroProf.Stop()

		log.Print("pyroscope started")
	}

	/////////////////
	// OTel SETUP //
	///////////////

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// set up OpenTelemetry.
	var err error
	loggerHandlers := []slog.Handler{slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(cfg.LogLevel)})}
	if !cfg.Observability.OTel.Enable {
		log.Printf("OTel disabled")
	} else {
		log.Print("OTel enabled")
		otelShutdown, err = otel.SetupOTelSDK(ctx, cfg.Observability.OTel.Endpoint, appName)
		if err != nil {
			log.Fatalf("failed to setup OTel: %v", err)
		}

		// Handle shutdown properly so nothing leaks.
		defer func() {
			log.Fatalf("shutdown error: %v", errors.Join(err, otelShutdown(context.Background())))
		}()

		loggerHandlers = append(loggerHandlers, otelslog.NewHandler("fightpicker"))

		log.Print("OTel setup complete")
	}

	///////////////////////////////
	// OTel instrumented logger //
	/////////////////////////////

	// TODO: slogmulti can be removed once we have go 1.26 is released as it includes built-in support for multiple handlers https://tip.golang.org/doc/go1.26
	baseLogger := logs.NewMultiSlogger(loggerHandlers...)

	baseLogger = baseLogger.With("app", appName, "env", env)
	baseLogger.InfoContext(ctx, "configuration loaded successfully")

	/////////////////////////////////////////////
	// create OTel instrumented DB connection //
	///////////////////////////////////////////

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Auth.SSLMode,
	)

	// DB Client init
	pool, dbClient, err := newDB(ctx, dbURL, baseLogger)
	if err != nil {
		baseLogger.FatalContext(ctx, "failed to create database client", "error", err)
	}

	///////////////////
	// utils init  	//
	/////////////////

	// we strictly use UUIDv7 for the added benefit of non-random sortable IDs
	idGen := id.NewUUIDV7Generator()
	dateTimeTool := datetimes.NewUTCNow()
	authTool := auth.NewAuthTool(cfg.Auth.HashingCost)
	jwtTool := jsonwebtokens.NewJWTTool[serviceauth.AuthClaims]()

	///////////////////
	// Clients init //
	/////////////////

	// Load the Shared AWS Configuration (~/.aws/config)
	awscfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(cfg.AWS.Region),
	)

	s3Client := s3.NewFromConfig(awscfg, func(o *s3.Options) {
		if cfg.AWS.S3Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.AWS.S3Endpoint)
			o.UsePathStyle = true
			creds, err := o.Credentials.Retrieve(ctx)
			if err != nil {
				baseLogger.ErrorContext(ctx, "failed to load AWS credentials", "access_key", creds.AccessKeyID)

				return
			}
		}
	})

	// Redis Client init
	rdbClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Cache.Host, cfg.Cache.Port),
		Password: cfg.Cache.Password,
		DB:       0, // use default DB
		Protocol: 2,
	})

	// Enable tracing instrumentation.
	if err := redisotel.InstrumentTracing(rdbClient); err != nil {
		baseLogger.FatalContext(ctx, "failed to instrument redis client for tracing", "error", err)
	}

	// Enable metrics instrumentation.
	if err := redisotel.InstrumentMetrics(rdbClient); err != nil {
		baseLogger.FatalContext(ctx, "failed to instrument redis client for metrics", "error", err)
	}

	kafkaMediaProfilePicClient, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.EventBroker.SeedBrokers...),
		kgo.ConsumeTopics(cfg.EventBroker.TopicProfilePictureUpload),
		// consumer group ID is is required for CommitRecords to work
		kgo.ConsumerGroup(cfg.EventBroker.GroupID),
		// start from the earliest offset if there is no committed offset
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	if err != nil {
		baseLogger.FatalContext(ctx, "failed to create kafka client for media profile picture uploads", "error", err)
	}

	kafkaPostUserCreateClient, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.EventBroker.SeedBrokers...),
		kgo.ConsumeTopics(cfg.EventBroker.TopicPostUserCreated),
		// consumer group ID is is required for CommitRecords to work
		kgo.ConsumerGroup(cfg.EventBroker.GroupID),
		// start from the earliest offset if there is no committed offset
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	if err != nil {
		baseLogger.FatalContext(ctx, "failed to create kafka client for post user create events", "error", err)
	}

	emailDialer := mail.NewDialer(cfg.Email.SMTPHost, cfg.Email.SMTPPort, cfg.Email.SMTPUser, cfg.Email.SMTPPassword)
	if cfg.Email.SkipTLS {
		emailDialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	/////////////////
	// REPOS INIT //
	///////////////

	presignedGetURLTTLDuration := time.Duration(cfg.AWS.PresignedGetURLTTL) * time.Minute
	presignedPutURLTTLDuration := time.Duration(cfg.AWS.PresignedPutURLTTL) * time.Minute

	usersRepo, err := repousers.New(
		pool,
		cfg.Domain,
		dbClient,
		s3Client,
		dateTimeTool,
		cfg.AWS.S3MediaBucket,
		cfg.AWS.S3MediaBucketRaw,
		presignedGetURLTTLDuration,
		presignedPutURLTTLDuration,
		kafkaPostUserCreateClient,
		cfg.EventBroker.TopicPostUserCreated,
		cfg.EventBroker.TopicPostUserVerified,
		emailDialer,
		cfg.Email.AddressNoReply,
		baseLogger,
	)
	if err != nil {
		baseLogger.FatalContext(ctx, "failed to create users repo", "error", err)
	}

	authRepo := repoauth.New(pool, dbClient, idGen, baseLogger)

	/////////////////////////////////
	// SERVICES AND HANDLERS INIT //
	///////////////////////////////

	var handlers = make([]server.RouteRegistrar, 0)

	// auth Service + Handler
	accessTokenTTL := time.Hour * time.Duration(cfg.Auth.AccessTTLHours)
	refreshTokenTTL := time.Hour * time.Duration(cfg.Auth.RefreshTTLHours)

	authService, err := serviceauth.New(
		usersRepo,
		authRepo,
		authTool,
		idGen,
		jwtTool,
		accessTokenTTL,
		refreshTokenTTL,
		cfg.Auth.PrivateKey,
		cfg.Auth.TokenAudience,
		cfg.Auth.TokenIssuer,
		baseLogger,
	)
	if err != nil {
		baseLogger.FatalContext(ctx, "failed to create auth service", "error", err)
	}

	authHandler := handlersv1auth.New(authService)
	handlers = append(handlers, authHandler)

	// Users Service + Handler
	if cfg.API.Users {
		baseLogger.InfoContext(ctx, "users API enabled, initializing users service and handler")

		imageProcessor := mediaprocessor.NewImageProcessor()
		usersService, err := serviceusers.NewService(usersRepo, idGen, dateTimeTool, authTool, imageProcessor, cfg.Domain, cfg.Email.AddressNoReply, baseLogger)
		if err != nil {
			baseLogger.FatalContext(ctx, "failed to create user service", "error", err)
		}

		imageConsumer, err := mediaconsumer.NewUserProfilePictureConsumer(
			kafkaMediaProfilePicClient,
			usersService,
			baseLogger,
		)
		if err != nil {
			baseLogger.FatalContext(ctx, "failed to create user profile picture consumer", "error", err)
		}

		emailConsumer, err := email.NewUserCreationEmailConsumer(
			kafkaPostUserCreateClient,
			usersService,
			baseLogger,
		)
		if err != nil {
			baseLogger.FatalContext(ctx, "failed to create user creation email consumer", "error", err)
		}

		// kick off the consumers in separate goroutines
		go imageConsumer.Run(ctx)
		go emailConsumer.Run(ctx)

		usersHandler := handlersv1users.New(usersService)

		handlers = append(handlers, usersHandler)
	}

	if cfg.API.Fighters {
		fightersrepo, err := repofighters.New(pool, dbClient, rdbClient, dateTimeTool, baseLogger)
		if err != nil {
			baseLogger.FatalContext(ctx, "failed to create fighters repo", "error", err)
		}

		baseLogger.InfoContext(ctx, "fighters API enabled, initializing fighters service and handler")
		fightersService, err := servicefighters.New(fightersrepo, idGen, dateTimeTool, baseLogger)
		if err != nil {
			baseLogger.FatalContext(ctx, "failed to create fighters service", "error", err)
		}

		fightersHandler := handlersv1fighters.New(fightersService)
		handlers = append(handlers, fightersHandler)
	}

	srv, err := server.New(cfg.Domain, cfg.Port, mux.NewRouter(), jwtTool, cfg.Auth.PrivateKey, cfg.Observability.OTel.Enable, env, baseLogger)
	if err != nil {
		baseLogger.FatalContext(ctx, "failed to create server", "error", err)
	}

	srv.WithHandlers(handlers)

	////////////////////
	// RUN THE SERVER //
	///////////////////

	if err := srv.Run(); err != nil {
		baseLogger.FatalContext(ctx, "server encountered an error", "error", err)
	}
}

func newDB(ctx context.Context, dsn string, logger logs.Logger) (*pgxpool.Pool, *postgresclient.Queries, error) {
	// 1. Create pool config
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, nil, err
	}

	// create multitracer for otelsql
	tracer := multitracer.New(otelpgx.NewTracer(), &tracelog.TraceLog{
		Logger:   SlogTracer{logger},
		LogLevel: tracelog.LogLevelInfo,
	})

	cfg.ConnConfig.Tracer = tracer

	// 3. Create the pool (thread-safe)
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, nil, err
	}

	// 4. Inject pool into sqlc generated client
	q := postgresclient.New(pool)

	return pool, q, nil
}

// SlogTracer implements pgx tracelog.Logger using our slog.Logger
type SlogTracer struct {
	logger logs.Logger
}

// Log implements pgx tracelog.Logger converting it to a regular slog log
func (s SlogTracer) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	// Convert map[string]any → slog style attributes
	attrs := make([]slog.Attr, 0, len(data))
	for k, v := range data {
		attrs = append(attrs, slog.Any(k, v))
	}

	s.logger.Log(ctx, logs.Level(level), msg, attrs)
}

type s3EndpointResolver struct {
	endpoint *url.URL
}

func (r *s3EndpointResolver) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (smithyendpoints.Endpoint, error) {
	if r.endpoint != nil {
		return smithyendpoints.Endpoint{
			URI: *r.endpoint,
		}, nil
	}

	return s3.NewDefaultEndpointResolverV2().
		ResolveEndpoint(ctx, params)
}
