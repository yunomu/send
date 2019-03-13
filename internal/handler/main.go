package handler

import (
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	logger = log.New(os.Stderr, "", log.LstdFlags)
	mu     = &sync.Mutex{}
)

func SetLogger(l *log.Logger) {
	mu.Lock()
	defer mu.Unlock()

	if l == nil {
		logger = log.New(os.Stderr, "", log.LstdFlags)
	} else {
		logger = l
	}
}

type handler struct {
	pipe map[string]io.WriteCloser
}

func NewHandler() *handler {
	return &handler{
		pipe: make(map[string]io.WriteCloser),
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger.Printf("%v: %v", req.Method, req.URL)
	path := req.URL.Path
	pipe, ok := h.pipe[path]

	switch req.Method {
	case http.MethodGet:
		if ok {
			w.WriteHeader(http.StatusConflict)
			logger.Printf("Resource busy")
			return
		}

		reader, writer := io.Pipe()
		defer reader.Close()

		h.pipe[path] = writer
		defer delete(h.pipe, path)

		w.Header().Add("Content-Type", "application/octet-stream")

		if _, err := io.Copy(w, reader); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Printf("io.Copy: %v", err)
			return
		}
	case http.MethodPut:
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			logger.Printf("Resource not found")
			return
		}

		defer req.Body.Close()
		defer pipe.Close()

		if _, err := io.Copy(pipe, req.Body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Printf("io.Copy: %v", err)
			return
		}
	}
}
