package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"text/tabwriter"
)

const (
	correctEmoji = "✔︎"
	wrongEmoji   = "✗"
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
		s.Handler = loggingHandler(s.Handler, c)
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

func yellow(s string) string {
	return fmt.Sprintf("\x1b[1;33m%v\x1b[0m", s)
}

func red(s string) string {
	return fmt.Sprintf("\x1b[1;31m%v\x1b[0m", s)
}

func yellowUnderline(s string) string {
	return fmt.Sprintf("\x1b[4;33m%v\x1b[0m", s)
}

func redUnderline(s string) string {
	return fmt.Sprintf("\x1b[4;31m%v\x1b[0m", s)
}

func green(s string) string {
	return fmt.Sprintf("\x1b[1;32m%v\x1b[0m", s)
}

func (fs fileServer) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Starting up File Server, Serving %v\n", yellow(fs.c.rootDir)))

	conf := fs.getConf()

	if conf != "" {
		sb.WriteString(yellowUnderline("Configuration:\n"))
		sb.WriteString(conf)
	}

	sb.WriteString(yellowUnderline("Available on:\n"))
	var addrs []string
	if fs.c.host == defaultAddr {
		addrs = append(addrs, green(fmt.Sprintf("\t%v%v:%v", fs.scheme, "127.0.0.1", fs.c.port)))
		extIP, _ := externalIP()
		addrs = append(addrs, green(fmt.Sprintf("\t%v%v:%v", fs.scheme, extIP, fs.c.port)))
	} else {
		addrs = append(addrs, green(fmt.Sprintf("\t%v%v:%v", fs.scheme, fs.c.host, fs.c.port)))
	}
	sb.WriteString(strings.Join(addrs, "\n"))

	sb.WriteString(fmt.Sprintf("\nHit %v to stop the server\n", redUnderline("CTRL+C")))

	return sb.String()
}

func (fs fileServer) getConf() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 5, 0, 1, ' ', tabwriter.AlignRight)

	if fs.BasicAuthEnable {
		fmt.Fprintln(w, "Basic Auth:\t "+green(correctEmoji))
	} else {
		fmt.Fprintln(w, "Basic Auth:\t "+red(wrongEmoji))
	}

	if fs.c.gzipEnable {
		fmt.Fprintln(w, "gzip:\t "+green(correctEmoji))
	} else {
		fmt.Fprintln(w, "gzip:\t "+red(wrongEmoji))
	}

	if fs.TLSEnable {
		fmt.Fprintln(w, "TLS:\t "+green(correctEmoji))
	} else {
		fmt.Fprintln(w, "TLS:\t "+red(wrongEmoji))
	}

	if fs.c.silent {
		fmt.Fprintln(w, "silent mode:\t "+green(correctEmoji))
	} else {
		fmt.Fprintln(w, "silent mode:\t "+red(wrongEmoji))
	}

	if fs.c.logIP {
		fmt.Fprintln(w, "log IP Addr:\t "+green(correctEmoji))
	} else {
		fmt.Fprintln(w, "log IP Addr:\t "+red(wrongEmoji))
	}

	if fs.c.hideDotFiles {
		fmt.Fprintln(w, "hide dot files:\t "+green(correctEmoji))
	} else {
		fmt.Fprintln(w, "hide dot files:\t "+red(wrongEmoji))
	}

	fmt.Fprintln(w, fmt.Sprintf("Read/Write Timeout:\t %v", fs.c.timeout))

	err := w.Flush()
	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
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
