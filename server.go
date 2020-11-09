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
	c               Config
	s               *http.Server
	scheme          string
	TLSEnable       bool
	BasicAuthEnable bool
}

func NewFileServer(c Config) *fileServer {
	s := http.Server{
		ReadTimeout:  c.timeout,
		WriteTimeout: c.timeout,
		Addr:         net.JoinHostPort(c.host, c.port),
	}

	fs := &fileServer{
		c: c,
		s: &s,
	}

	var fSys http.FileSystem

	if c.hideDotFiles {
		fSys = dotFileHidingFileSystem{http.Dir(c.rootDir)}
	} else {
		fSys = http.Dir(c.rootDir)
	}

	h := http.FileServer(fSys)

	if c.gzipEnable {
		s.Handler = gzipHandler(h)
	} else {
		s.Handler = h
	}

	if c.username != "" && c.password != "" {
		fs.BasicAuthEnable = true
		a := auth{
			username: c.username,
			password: c.password,
			relm:     "Please enter your username and password for this site",
		}

		s.Handler = a.basicAuthHandler(s.Handler)
	}

	if !c.silent {
		s.Handler = loggingHandler(s.Handler, c.logIP)
	}

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
		err := fs.s.ListenAndServeTLS(fs.c.certFile, fs.c.keyFile)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	} else {
		err := fs.s.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}
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
