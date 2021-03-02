package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type LogHandler struct {
	logger  *log.Logger
	handler http.Handler
}

func (h *LogHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	h.handler.ServeHTTP(rw, req)
	h.logger.Printf("%s  %s\n", req.Method, req.URL.Path)
}

func main() {
	logger := log.New(os.Stdout, "", 0)
	var host, port, dir *string

	host = flag.String("host", "", "HTTP server's host. Default: \"localhost\"")
	port = flag.String("port", "80", "HTTP server's port. Default: \"80\"")
	dir = flag.String("dir", ".", "HTTP server's dir. Default is CWD")

	flag.Parse()
	dirpath, err := filepath.Abs(*dir)
	if err != nil {
		logger.Printf("WARN: %v\n", err)
		err = nil
		dirpath = *dir
	}
	hostNiceStr := *host
	if hostNiceStr == "" {
		hostNiceStr = "127.0.0.1"
	}

	server := http.Server{
		Addr:           net.JoinHostPort(*host, *port),
		Handler:        &LogHandler{handler: http.FileServer(http.Dir(*dir)), logger: logger},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logger.Printf("Listening: [ %s:%s ]\nServing dir: [ %s ]\n", hostNiceStr, *port, dirpath)
	log.Fatal(server.ListenAndServe())
	return
}
