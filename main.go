package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/yunomu/send/internal/handler"
)

var (
	logger = log.New(os.Stderr, "", log.LstdFlags)
	mu     = &sync.Mutex{}

	bind = flag.String("bind", ":8080", "Bind address")
)

func init() {
	flag.Parse()
}

func SetLogger(l *log.Logger) {
	mu.Lock()
	defer mu.Unlock()

	if l == nil {
		logger = log.New(os.Stderr, "", log.LstdFlags)
	} else {
		logger = l
	}
}

func main() {
	defer logger.Printf("Stopped")
	logger.Printf("Start")

	http.Handle("/", handler.NewHandler())

	if err := http.ListenAndServe(*bind, nil); err != nil {
		logger.Fatalf("http.ListenAndServe: %v", err)
	}
}
