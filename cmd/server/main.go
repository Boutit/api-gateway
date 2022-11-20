package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/Boutit/api-gateway/internal/config"
	"github.com/Boutit/api-gateway/pkg/middleware"
	userService "github.com/Boutit/user/api"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const deadline = 20 * time.Millisecond

	// NewProxy takes target host and creates a reverse proxy
	func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
		url, err := url.Parse(targetHost)
		if err != nil {
				return nil, err
		}
	
		return httputil.NewSingleHostReverseProxy(url), nil
	}
	
	// ProxyRequestHandler handles the http request using proxy
	func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
				proxy.ServeHTTP(w, r)
		}
	}

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

	// Create a multiplexer to match http requests to patterns and invoke corresponding handlers
	gwmux := runtime.NewServeMux()

	// Register the user service handlers with the the multiplexer and user service client connection
	err = userService.RegisterUserServiceHandler(context.Background(), gwmux, conn)
	
	if err != nil {
		log.Fatalln("Failed to register user service handlers with err:", err)
	}

	r := mux.NewRouter()

	/* Handlers */
	userServiceHandler := middleware.CreateHandler([]middleware.Middleware{
		middleware.AuthenticateToken(),
	}, gwmux)

	/* Reverse Proxy for GraphQL path */
	fmt.Println("before")
	rp := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = cfg.GraphQLConfig.Scheme
			r.URL.Host = cfg.GraphQLConfig.Host
			r.URL.Path = cfg.GraphQLConfig.Path
		},
	}

	fmt.Println("after")

	graphQLHandler := middleware.CreateHandler([]middleware.Middleware{
		middleware.AuthenticateToken(),
	}, rp)
	

	/* Routes */
	r.PathPrefix("/v1/user").Handler(userServiceHandler)
	r.PathPrefix("/graphql").Handler(graphQLHandler)
	
	r.Use(cors.New(cors.Options{
		AllowedHeaders: 	[]string{"Access-Control-Allow-Headers", "Access-Control-Allow-Methods", "Access-Control-Allow-Origin", "Authorization", "Content-Type"},
		AllowedOrigins:   []string{"http://localhost:8080"},
	}).Handler)

	// Create a gateway server to expose the multiplexer
	gwServer := &http.Server{
		Addr:    gatewayConnStr,
		Handler: r,
	}

	log.Printf("API gateway listening on http:%s", gatewayConnStr)

	// Serve the gateway server
	log.Fatalln(gwServer.ListenAndServe())
}