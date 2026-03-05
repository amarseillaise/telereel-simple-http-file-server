package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amarseillaise/simple-http-file-server/internal/handlers"
	"github.com/amarseillaise/simple-http-file-server/internal/service"
	"github.com/amarseillaise/simple-http-file-server/internal/storage"
	"github.com/amarseillaise/simple-http-file-server/pkg/config"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.Load()

	msg := "Starting server with config: port=%d, contentDir=%s"
	if cfg.TLSEnabled() {
		log.Printf(msg, cfg.ServerPort, cfg.ContentDir)
	} else {
		log.Printf(msg, cfg.ServerPort, cfg.ContentDir)
	}

	fs, err := storage.NewFileSystem(cfg.ContentDir)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	videoService := service.NewVideoService(fs, &service.Downloader{})
	videoHandler := handlers.NewVideoHandler(videoService)

	router := mux.NewRouter()
	videoHandler.RegisterRoutes(router)

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods(http.MethodGet)

	router.Use(loggingMiddleware)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 5 * time.Minute,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if cfg.TLSEnabled() {
			log.Printf("Server listening on port %d (HTTPS)", cfg.ServerPort)
			if err := server.ListenAndServeTLS(cfg.TLSCertFile, cfg.TLSKeyFile); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Server failed: %v", err)
			}
		} else {
			log.Printf("Server listening on port %d (HTTP)", cfg.ServerPort)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Server failed: %v", err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		log.Printf(
			"%s %s %d %s",
			r.Method,
			r.RequestURI,
			wrapped.statusCode,
			time.Since(start),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
