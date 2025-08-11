package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/RomanRyabinkin/antipidorandel/internal/config"
	"github.com/RomanRyabinkin/antipidorandel/internal/hub"
	"github.com/RomanRyabinkin/antipidorandel/internal/janitor"
	"github.com/RomanRyabinkin/antipidorandel/internal/store/postgres"
	"github.com/RomanRyabinkin/antipidorandel/internal/transport/ws"
)

func main() {
	cfg, err := config.Load()
	if err != nil { log.Fatal(err) }

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// storage
	st, err := postgres.New(ctx, cfg.DatabaseURL)
	if err != nil { log.Fatal(err) }
	defer st.Close()
	if err := st.Init(ctx); err != nil { log.Fatal(err) }

	// hub + ws
	h := hub.New()
	wsSrv := ws.New(h, st, cfg)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(200); _, _ = w.Write([]byte("ok")) })
	wsSrv.Router(mux)

	// janitor
	janitor.Start(ctx, st, cfg)

	httpSrv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: cfg.HandshakeTimeout,
	}

	go func() {
		log.Printf("listen %s", cfg.Addr)
		if err := httpSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")
	shCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(shCtx)
}