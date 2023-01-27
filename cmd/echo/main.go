package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[Echo] %v: %v", r.Method, r.URL.Path)
		io.WriteString(w, "OK")
	})

	server := &http.Server{Addr: ":8890", Handler: mux}

	log.Println("[Echo] Starting...")

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("[Echo] Statup error: %s", err.Error())
		}
	}()

	log.Println("[Echo] Serving simple echo")

	// shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-stop

	log.Println("[Echo] Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[Echo] Shutdown error: %s", err.Error())
	} else {
		log.Println("[Echo] Ready to exit")
	}
	cancel()
	os.Exit(0)
}
