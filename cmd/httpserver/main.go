package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/louispy/miniloan/internal/api"
)

func main() {
	r := api.NewAPIRouter()
	r.Register()

	port := ":8989"

	srv := &http.Server{Addr: port, Handler: r.GetRouter()}
	go func() {
		log.Println("Listening on port", port)
		log.Fatal(srv.ListenAndServe())
	}()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}
