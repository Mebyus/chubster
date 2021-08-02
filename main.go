package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	timeout       = 10 * time.Second
	maxHeaderSize = 1 << 20

	infoPrefix    = "[info]    "
	errorPrefix   = "[error]   "
	requestPrefix = "[request] "
)

type LogHandler struct {
	logger  *log.Logger
	handler http.Handler
}

func NewLogHandler(handler http.Handler) *LogHandler {
	return &LogHandler{
		handler: handler,
		logger:  log.New(os.Stdout, requestPrefix, log.Lmsgprefix|log.Ltime),
	}
}

func (h *LogHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	h.handler.ServeHTTP(rw, req)
	h.logger.Printf("%-8s %s", req.Method, req.URL.Path)
}

func main() {
	loggerInfo := log.New(os.Stdout, infoPrefix, log.Lmsgprefix|log.Ltime)
	loggerError := log.New(os.Stdout, errorPrefix, log.Lmsgprefix|log.Ltime)

	var host, dir string
	var port uint

	flag.StringVar(&host, "host", "localhost", "HTTP server's host")
	flag.UintVar(&port, "port", 80, "HTTP server's port")
	flag.StringVar(&dir, "dir", ".", "HTTP server's serve directory")

	flag.Parse()

	absDirPath, err := filepath.Abs(dir)
	if err != nil {
		loggerError.Fatalln(err)
	}

	hostAndPort := net.JoinHostPort(host, strconv.FormatUint(uint64(port), 10))

	server := http.Server{
		Addr:           hostAndPort,
		Handler:        NewLogHandler(http.FileServer(http.Dir(absDirPath))),
		ReadTimeout:    timeout,
		WriteTimeout:   timeout,
		MaxHeaderBytes: maxHeaderSize,
	}

	loggerInfo.Printf("Serving:   %s", absDirPath)
	loggerInfo.Printf("Listening: %s", hostAndPort)
	loggerError.Fatalln(server.ListenAndServe())
}
