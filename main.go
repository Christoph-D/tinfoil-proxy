package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/openai/openai-go/v3/option"
	"github.com/tinfoilsh/tinfoil-go"
)

func main() {
	listenAddr := flag.String("listen", "127.0.0.1:17349", "listen address")
	flag.Parse()

	tinfoilAPIKey := os.Getenv("TINFOIL_API_KEY")
	if tinfoilAPIKey == "" {
		log.Fatal("TINFOIL_API_KEY environment variable is required")
	}

	proxyAPIKey := os.Getenv("TINFOIL_PROXY_API_KEY")
	if proxyAPIKey == "" {
		log.Fatal("TINFOIL_PROXY_API_KEY environment variable is required")
	}

	log.Println("initializing tinfoil client and verifying enclave...")
	client, err := tinfoil.NewClient(option.WithAPIKey(tinfoilAPIKey))
	if err != nil {
		log.Fatalf("failed to create tinfoil client: %v", err)
	}

	enclave := client.Enclave()
	secureTransport := client.HTTPClient().Transport
	log.Printf("enclave verified: %s", enclave)

	chatHandler := newProxyHandler(secureTransport, enclave, tinfoilAPIKey, proxyAPIKey, "/v1/chat/completions", []string{http.MethodPost})
	modelsHandler := newProxyHandler(secureTransport, enclave, tinfoilAPIKey, proxyAPIKey, "/v1/models", []string{http.MethodGet})

	mux := http.NewServeMux()
	mux.Handle("/v1/chat/completions", chatHandler)
	mux.Handle("/v1/models", modelsHandler)

	srv := &http.Server{
		Addr:    *listenAddr,
		Handler: mux,
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("shutdown error: %v", err)
		}
	}()

	fmt.Fprintf(os.Stderr, "tinfoil-proxy listening on %s\n", *listenAddr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
