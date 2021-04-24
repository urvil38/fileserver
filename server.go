package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
)

type fileServer struct {
	config          Config
	server          *http.Server
	scheme          string
	TLSEnable       bool
	BasicAuthEnable bool
}

func NewFileServer(c Config) *fileServer {
	fs := &fileServer{
		config: c,
		server: &http.Server{
			ReadTimeout:  c.timeout,
			WriteTimeout: c.timeout,
			Addr:         net.JoinHostPort(c.host, c.port),
		},
	}

	var fSys http.FileSystem

	if c.hideDotFiles {
		fSys = dotFileHidingFileSystem{http.Dir(c.rootDir)}
	} else {
		fSys = http.Dir(c.rootDir)
	}

	h := http.FileServer(fSys)

	if !c.quiet {
		h = loggingHandler(h, c.logIP)
	}

	if c.gzipEnable {
		h = gzipHandler(h)
	}

	if c.username != "" && c.password != "" {
		fs.BasicAuthEnable = true
		a := auth{
			username: c.username,
			password: c.password,
			relm:     "Please enter your username and password for this site",
		}

		h = a.basicAuthHandler(h)
	}

	fs.server.Handler = h

	if c.keyFile != "" && c.certFile != "" {
		fs.TLSEnable = true
		fs.scheme = "https://"
	} else {
		fs.scheme = "http://"
	}

	return fs
}

func (fs *fileServer) Start() {

	fmt.Print(fs)

	if fs.TLSEnable {
		err := fs.server.ListenAndServeTLS(fs.config.certFile, fs.config.keyFile)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	} else {
		err := fs.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}
}

func (fs *fileServer) Stop(ctx context.Context) {
	err := fs.server.Shutdown(ctx)
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
