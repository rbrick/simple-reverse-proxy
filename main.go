package main

import (
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func init() {

}

func main() {
	instance := NewProxiedHandler("hosts.json")
	go http.ListenAndServe(":80", instance)

	if os.Getenv("SRP_TLS") == "true" {
		go http.ListenAndServeTLS(":443", os.Getenv("SRP_CERT_FILE"), os.Getenv("SRP_KEY_FILE"), instance)
	}

	signals := make(chan os.Signal, 1) // allocate a channel with a size of one

	// SIGINT = Ctrl+C
	// SIGTERM = Termination request
	// Listen for these
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	{
		wg.Add(1)
	}

	go func() {
		select {
		case <-signals:
			wg.Done()
		}
	}()

	wg.Wait()
}
