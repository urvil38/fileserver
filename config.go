package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/tabwriter"
	"time"
)

type Config struct {
	host         string
	port         string
	rootDir      string
	certFile     string
	keyFile      string
	timeout      time.Duration
	gzipEnable   bool
	hideDotFiles bool
	silent       bool
	logIP        bool
	username     string
	password     string
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
