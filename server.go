package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type fileServer struct {
	s         *http.Server
	TLSEnable bool
	scheme    string
	host   string
	port      string
	rootDir   string
	certFile  string
	keyFile   string
}

func NewFileServer(host, port, rootDir, certFile, keyFile string, timeout time.Duration, handler http.Handler) *fileServer {
	s := http.Server{
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		Addr:         net.JoinHostPort(host,port),
	}

	if handler == nil {
		s.Handler = http.FileServer(http.Dir(rootDir))
	} else {
		s.Handler = handler
	}

	fs := &fileServer{
		host:  host,
		port:     port,
		rootDir:  rootDir,
		certFile: certFile,
		keyFile:  keyFile,
		s:        &s,
	}

	if keyFile != "" && certFile != "" {
		fs.TLSEnable = true
		fs.scheme = "https://"
	} else {
		fs.scheme = "http://"
	}

	return fs
}

func (fs *fileServer) Start() {

	fmt.Printf("Server is running on %v%v:%v serving %v\n", color(fs.scheme), color(fs.host), color(fs.port), color(fs.rootDir))

	if fs.certFile != "" && fs.keyFile != "" {
		err := fs.s.ListenAndServeTLS(fs.certFile, fs.keyFile)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := fs.s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func color(s string) string {
	return fmt.Sprintf("\x1b[1;33m%v\x1b[0m", s)
}

func (fs *fileServer) Stop(ctx context.Context) {
	err := fs.s.Shutdown(ctx)
	if err != nil {
		if err == context.DeadlineExceeded {
			log.Println("finish: shutdown timeout")
		} else {
			log.Println("finish: error while shutting down, ", err)
		}
	} else {
		log.Println("finish: closed")
	}
}
