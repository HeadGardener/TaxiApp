package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/HeadGardener/TaxiApp/driver-service/internal/config"
	grpc_client "github.com/HeadGardener/TaxiApp/driver-service/internal/grpc-client"
	grpc_handler "github.com/HeadGardener/TaxiApp/driver-service/internal/handlers/grpc-handler"
	http_handler "github.com/HeadGardener/TaxiApp/driver-service/internal/handlers/http-handler"
	grpc_server "github.com/HeadGardener/TaxiApp/driver-service/internal/server/grpc-server"
	http_server "github.com/HeadGardener/TaxiApp/driver-service/internal/server/http-server"
	"github.com/HeadGardener/TaxiApp/driver-service/internal/services"
	"github.com/HeadGardener/TaxiApp/driver-service/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const shutdownTimeout = 5 * time.Second

var confPath = flag.String("conf-path", "./config/.env", "path to config env")

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	conf, err := config.Init(*confPath)
	if err != nil {
		stop()
		log.Fatalf("[FATAL] error while initializing config: %e", err)
	}

	db, err := storage.NewDB(conf.DBConfig)
	if err != nil {
		stop()
		log.Fatalf("[FATAL] error while establishing db connection: %e", err)
	}

	rdb := storage.NewRedisDB(conf.RedisConfig)
	if err != nil {
		stop()
		log.Fatalf("[FATAL] error while establishing db connection: %e", err)
	}

	var (
		driverStorage    = storage.NewDriverStorage(db)
		tripStorage      = storage.NewTripStorage(db)
		balanceProcessor = storage.NewBalanceProcessor(db)

		orderStorage = storage.NewOrderStorage()

		tokenStorage = storage.NewTokenStorage(rdb)
	)

	conn, err := grpc.Dial("localhost:"+conf.GRPCClientConfig.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		stop()
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := grpc_client.NewOrderServiceClient(conn)

	var (
		driverService   = services.NewDriverService(driverStorage)
		orderService    = services.NewOrderService(client, orderStorage)
		tripService     = services.NewTripService(tripStorage)
		balanceService  = services.NewBalanceService(balanceProcessor, driverStorage)
		identityService = services.NewIdentityService(driverStorage, tokenStorage)
	)

	handler := http_handler.NewHandler(identityService, driverService, tripService, balanceService, orderService)

	srv := &http_server.Server{}
	go func() {
		if err = srv.Run(conf.ServerConfig, handler.InitRoutes()); err != nil {
			log.Printf("[ERROR] failed to run server: %e", err)
		}
	}()
	log.Println("[INFO] server start working")

	grpcHandler := grpc_handler.NewProcessOrderHandler(orderService)

	grpcsrv := &grpc_server.GRPCServer{}
	go func() {
		if err = grpcsrv.Init(conf.GRPCServerConfig, grpcHandler); err != nil {
			log.Printf("[ERROR] failed to run grpc server: %e", err)
		}
	}()
	log.Println("[INFO] grpc server start working")

	<-ctx.Done()
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Printf("[INFO] server forced to shutdown: %e", err)
	}

	if err = db.Close(); err != nil {
		log.Printf("[INFO] db connection forced to shutdown: %e", err)
	}

	if err = rdb.Close(); err != nil {
		log.Printf("[INFO] redis db connection forced to shutdown: %e", err)
	}

	log.Println("[INFO] server exiting")
}
