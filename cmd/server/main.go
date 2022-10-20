package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Boutit/api-gateway/internal/config"
	"github.com/Boutit/api-gateway/pkg/middleware"
	userService "github.com/Boutit/user/api"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const deadline = 20 * time.Millisecond

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	cfg := config.GetConfig(env)
	
	// use net package's Listen method to create a server
	h := cfg.AppConfig.Host

	gatewayPort := cfg.AppConfig.GatewayPort
	gatewayValues := []interface{}{h, gatewayPort}
	gatewayConnStr := fmt.Sprintf("%s:%d", gatewayValues...)


	// Create a context w/ a deadline to pass to grpc DialContext
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2 * time.Second))

	defer cancel()

	// Create client connection to GRPC server
	// GRPC server must be running at url specified in config
	conn, err := grpc.DialContext(
		ctx,
		cfg.AppConfig.UserServiceUrl,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	
	if err != nil {
		log.Fatalln("Failed to connect to user service with err:", err)
	}


	r := mux.NewRouter()
	r.Host(cfg.AppConfig.UserServiceUrl, middleware.CreateMiddlewareHandler)







	

	// Create a multiplexer to match http requests to patterns and invoke corresponding handlers
	gwmux := runtime.NewServeMux()

	// Register the user service handlers with the the multiplexer and user service client connection
	err = userService.RegisterUserServiceHandler(context.Background(), gwmux, conn)
	
	if err != nil {
		log.Fatalln("Failed to register user service handlers with err:", err)
	}

	// Create a gateway server to expose the multiplexer
	gwServer := &http.Server{
		Addr:    gatewayConnStr,
		Handler: gwmux,
	}

	log.Printf("Serving gRPC-Gateway on http:%s", gatewayConnStr)

	// Serve the gateway server
	log.Fatalln(gwServer.ListenAndServe())
}