package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/HeadGardener/TaxiApp/order-service/internal/config"
	grpc_client "github.com/HeadGardener/TaxiApp/order-service/internal/grpc-client"
	grpc_handler "github.com/HeadGardener/TaxiApp/order-service/internal/handlers/grpc-handler"
	http_handler "github.com/HeadGardener/TaxiApp/order-service/internal/handlers/http-handler"
	"github.com/HeadGardener/TaxiApp/order-service/internal/models"
	grpc_server "github.com/HeadGardener/TaxiApp/order-service/internal/server/grpc-server"
	http_server "github.com/HeadGardener/TaxiApp/order-service/internal/server/http-sever"
	"github.com/HeadGardener/TaxiApp/order-service/internal/services"
	"github.com/HeadGardener/TaxiApp/order-service/internal/storage"
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

	db, err := storage.NewMongoCollection(ctx, conf.DBConfig)
	if err != nil {
		stop()
		log.Fatalf("[FATAL] error while establishing db connection: %e", err)
	}

	var (
		orderStorage = storage.NewOrderStorage(db)
	)

	userConn, err := grpc.Dial("localhost:"+conf.GRPCUserClientConfig.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		stop()
		log.Fatalf("fail to dial: %v", err)
	}

	userClient := grpc_client.NewUserServiceClient(userConn)

	driverConn, err := grpc.Dial("localhost:"+conf.GRPCDriverClientConfig.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		stop()
		userConn.Close()
		log.Fatalf("fail to dial: %v", err)
	}

	driverClient := grpc_client.NewDriverServiceClient(driverConn)

	var (
		orderService  = services.NewOrderService(userClient, orderStorage)
		orderNotifier = services.NewOrderNotifier(driverClient, models.NewUsersQueues(), models.NewDriversQueues())
	)

	handler := http_handler.NewHandler(orderService)

	srv := &http_server.Server{}
	go func() {
		if err = srv.Run(conf.ServerConfig, handler.InitRoutes()); err != nil {
			log.Printf("[ERROR] failed to run server: %e", err)
		}
	}()
	log.Println("[INFO] server start working")

	grpcHandler := grpc_handler.NewProcessOrderHandler(orderService, orderNotifier)

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

	if err = db.Database().Client().Disconnect(ctx); err != nil {
		log.Printf("[INFO] db connection forced to shutdown: %e", err)
	}

	userConn.Close()
	driverConn.Close()

	log.Println("[INFO] server exiting")
}
