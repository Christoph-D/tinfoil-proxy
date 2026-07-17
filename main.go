package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/openai/openai-go/v3/option"
	"github.com/tinfoilsh/tinfoil-go"
)

func resolveKey(directEnvName, pathEnvName string) (string, error) {
	direct := os.Getenv(directEnvName)
	path := os.Getenv(pathEnvName)

	if direct != "" && path != "" {
		return "", fmt.Errorf("both %s and %s are set; provide exactly one", directEnvName, pathEnvName)
	}
	if direct != "" {
		return direct, nil
	}
	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("failed to read %s from %s: %w", directEnvName, path, err)
		}
		return strings.TrimSpace(string(data)), nil
	}
	return "", fmt.Errorf("missing key: %s or %s must be set", directEnvName, pathEnvName)
}

func main() {
	listenAddr := flag.String("listen", "127.0.0.1:17349", "listen address")
	flag.Parse()

	tinfoilAPIKey, err := resolveKey("TINFOIL_API_KEY", "TINFOIL_API_KEY_PATH")
	if err != nil {
		log.Fatal(err)
	}

	proxyAPIKey, err := resolveKey("TINFOIL_PROXY_API_KEY", "TINFOIL_PROXY_API_KEY_PATH")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("initializing tinfoil client and verifying enclave...")
	cacheSecret := make([]byte, 32)
	if _, err := rand.Read(cacheSecret); err != nil {
		log.Fatalf("failed to generate user cache secret: %v", err)
	}
	client, err := tinfoil.NewClientWithOptions(
		tinfoil.WithOpenAIOptions(option.WithAPIKey(tinfoilAPIKey)),
		tinfoil.WithUserCacheSecret(hex.EncodeToString(cacheSecret)),
	)
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
