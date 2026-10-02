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

	"github.com/hextec/spec-auditor/internal/adapter/llm"
	"github.com/hextec/spec-auditor/internal/adapter/repository"
	deliveryHttp "github.com/hextec/spec-auditor/internal/delivery/http"
	"github.com/hextec/spec-auditor/internal/service"
	"github.com/hextec/spec-auditor/web"
)

func main() {
	// Parse CLI flags and Environment variables
	port := flag.String("port", getEnv("PORT", "8080"), "HTTP server listening port")
	ollamaURL := flag.String("ollama-url", getEnv("OLLAMA_URL", "http://localhost:11434"), "Ollama Base URL")
	ollamaModel := flag.String("model", getEnv("OLLAMA_MODEL", "qwen2.5:7b"), "Ollama LLM model name")
	flag.Parse()

	log.Println("==================================================================")
	log.Println("  SpecAuditor Gov-Edition (Oʻz DSt 1987:2018 Ekspertizasi)  ")
	log.Println("==================================================================")

	// 1. Initialize Adapters (Repository & LLM)
	repo := repository.NewMemoryRepository()

	llmCfg := llm.DefaultOllamaConfig()
	llmCfg.BaseURL = *ollamaURL
	llmCfg.Model = *ollamaModel
	llmClient := llm.NewOllamaClient(llmCfg)

	// 2. Initialize Service Layer
	auditorService := service.NewAuditorService(repo, llmClient)

	// 3. Initialize Transport / HTTP Handlers
	apiHandler := deliveryHttp.NewHandler(auditorService)

	mux := http.NewServeMux()
	apiHandler.RegisterRoutes(mux)

	// 4. Mount Embedded Static Web UI
	indexHTML, err := web.ReadIndexHTML()
	if err != nil {
		log.Fatalf("[FATAL] index.html yuklanmadi: %v", err)
	}

	fileServer := http.FileServer(web.GetFileSystem())
	mux.Handle("/static/", fileServer)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			fileServer.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(indexHTML)
	})

	// 5. Wrap Middleware Stack
	var handler http.Handler = mux
	handler = deliveryHttp.CORSMiddleware(handler)
	handler = deliveryHttp.LoggingMiddleware(handler)
	handler = deliveryHttp.RecoveryMiddleware(handler)

	// 6. Check Initial LLM Status for terminal diagnostics
	checkCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	llmOnline := auditorService.IsLLMAvailable(checkCtx)
	cancel()

	if llmOnline {
		log.Printf("[STATUS] Lokal LLM (%s) MUVOFAQIYATLI ULANGAN (Faol)\n", *ollamaModel)
	} else {
		log.Printf("[STATUS] Lokal LLM ulanmadi (Ollama: %s). Avtomatik 'Offline Demo Koʻrgazma' rejimi faollashdi.\n", *ollamaURL)
	}

	addr := fmt.Sprintf(":%s", *port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 240 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 7. Start Server in Background
	go func() {
		log.Printf("[SERVER] SpecAuditor ishga tushdi: http://localhost:%s\n", *port)
		log.Printf("[SERVER] Brauzerda http://localhost:%s manzilini oching\n", *port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[FATAL] Server toʻxtadi: %v", err)
		}
	}()

	// 8. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("[SHUTDOWN] Tizim toʻxtatilmoqda...")
	ctx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[ERROR] Majburiy toʻxtatish: %v", err)
	}

	log.Println("[DONE] Server muvaffaqiyatli toʻxtatildi")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
