package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Boutit/api-gateway/internal/config"
	userService "github.com/Boutit/user/api"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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


		// Create client connection to GRPC server that was started
	// Proxy requests using grpc-gateway
	conn, err := grpc.DialContext(
		context.Background(),
		cfg.AppConfig.UserServiceUrl,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()

	err = userService.RegisterUserServiceHandler(context.Background(), gwmux, conn)
	
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    gatewayConnStr,
		Handler: gwmux,
	}

	log.Printf("Serving gRPC-Gateway on http:%s", gatewayConnStr)

	log.Fatalln(gwServer.ListenAndServe())
}