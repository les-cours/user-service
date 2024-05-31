package service

import (
	"github.com/les-cours/user-service/api/learning"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
	"os"
	"runtime"

	"github.com/les-cours/user-service/api/auth"
	"github.com/les-cours/user-service/api/payment"
	"github.com/les-cours/user-service/api/users"
	"github.com/les-cours/user-service/resolvers"

	"github.com/les-cours/user-service/database"
	"github.com/les-cours/user-service/env"
	"google.golang.org/grpc"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	registry       = prometheus.NewRegistry()
	requestCounter = prometheus.NewGauge(prometheus.GaugeOpts{Name: "request_counter", Help: "request counter"})
	memoryUsage    = prometheus.NewGauge(prometheus.GaugeOpts{Name: "memory_usage", Help: "memory usage"})
	goRoutineNum   = prometheus.NewGauge(prometheus.GaugeOpts{Name: "go_routines_num", Help: "the number of go routine "})
)

func monitoringMiddleware(originalHandler http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		memoryUsage.Set(float64(m.Alloc))
		goRoutineNum.Set(float64(runtime.NumGoroutine()))
		requestCounter.Inc()
		originalHandler.ServeHTTP(w, r)
	}
}

func loggerInit() *zap.Logger {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(zap.InfoLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0))
	return logger
}

func Start() {
	logger := loggerInit()
	defer logger.Sync()
	registry.MustRegister(requestCounter, memoryUsage, goRoutineNum)
	promHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	http.HandleFunc("/metrics", monitoringMiddleware(promHandler))
	logger.Info("Starting http server on port " + env.Settings.HttpPort)
	go func() {
		err := http.ListenAndServe(":"+env.Settings.HttpPort, nil)
		if err != nil {
			logger.Info("Error http server on port :"+env.Settings.HttpPort, zap.Error(err))
		}
	}()

	lis, err := net.Listen("tcp", ":"+env.Settings.GrpcPort)
	if err != nil {
		logger.Fatal("Failed to listen on port : "+env.Settings.GrpcPort, zap.Error(err))
	}

	db, err := database.StartDatabase()
	if err != nil {
		logger.Error(err.Error())
	}

	//defer db.Close()
	//mongoDB, err := database.StartMongoDB()
	//if err != nil {
	//	logger.Fatalln(err)
	//}
	//
	//defer mongoDB.Disconnect(context.Background())

	logger.Info("auth service connection: " + env.Settings.AuthService.Host + ":" + env.Settings.AuthService.Port)
	authConnectionService, err := grpc.Dial(env.Settings.AuthService.Host+":"+env.Settings.AuthService.Port, grpc.WithInsecure())
	if err != nil {
		logger.Error(err.Error())
	}
	defer authConnectionService.Close()
	authServiceClient := auth.NewAuthServiceClient(authConnectionService)

	learningConnectionService, err := grpc.Dial(env.Settings.LearningService.Host+":"+env.Settings.LearningService.Port, grpc.WithInsecure())
	if err != nil {
		logger.Error(err.Error())
	}
	defer learningConnectionService.Close()
	learningServiceClient := learning.NewLearningServiceClient(learningConnectionService)

	paymentConnectionService, err := grpc.Dial("payment-api:8080", grpc.WithInsecure())
	if err != nil {
		logger.Error(err.Error())
	}
	defer paymentConnectionService.Close()
	paymentServiceClient := payment.NewPaymentServiceClient(paymentConnectionService)

	var s = resolvers.GetInstance(
		db,
		authServiceClient,
		learningServiceClient,
		paymentServiceClient,
		logger,
	)

	grpcServer := grpc.NewServer()
	users.RegisterUserServiceServer(grpcServer, s)
	logger.Info("Starting grpc server on port " + env.Settings.GrpcPort)
	err = grpcServer.Serve(lis)
	if err != nil {
		logger.Error("Failed to start gRPC server on port: "+env.Settings.GrpcPort, zap.Error(err))
	}

}
