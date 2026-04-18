package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

type proxyHandler struct {
	proxy      *httputil.ReverseProxy
	proxyToken string
	apiKey     string
	enclave    string
	upstream   string
	methods    []string
}

type errorResponse struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

func newProxyHandler(secureTransport http.RoundTripper, enclave, apiKey, proxyToken, upstream string, methods []string) *proxyHandler {
	h := &proxyHandler{
		proxyToken: proxyToken,
		apiKey:     apiKey,
		enclave:    enclave,
		upstream:   upstream,
		methods:    methods,
	}

	h.proxy = &httputil.ReverseProxy{
		Director:     h.director,
		Transport:    secureTransport,
		ErrorHandler: h.errorHandler,
	}

	return h
}

func (h *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	methodAllowed := false
	for _, m := range h.methods {
		if r.Method == m {
			methodAllowed = true
			break
		}
	}
	if !methodAllowed {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", fmt.Sprintf("Allowed methods: %s", strings.Join(h.methods, ", ")))
		return
	}

	if !h.authenticate(r) {
		writeError(w, http.StatusUnauthorized, "authentication_error", "Invalid or missing API key")
		return
	}

	h.proxy.ServeHTTP(w, r)
}

func (h *proxyHandler) authenticate(r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return false
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	return token == h.proxyToken
}

func (h *proxyHandler) director(req *http.Request) {
	req.URL.Scheme = "https"
	req.URL.Host = h.enclave
	req.URL.Path = h.upstream
	req.URL.RawPath = ""
	req.URL.RawQuery = ""
	req.Host = h.enclave

	req.Header.Set("Authorization", "Bearer "+h.apiKey)
	req.Header.Set("Content-Type", "application/json")
}

func (h *proxyHandler) errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("proxy error: %v", err)
	writeError(w, http.StatusBadGateway, "upstream_error", fmt.Sprintf("Failed to reach upstream: %v", err))
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{
		Error: errorDetail{
			Message: message,
			Type:    "proxy_error",
			Code:    code,
		},
	})
}
