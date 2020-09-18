package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/urvil38/fileserver/version"
)

func getEnv(env string) string {
	val := os.Getenv(env)
	return val
}

const (
	defaultAddr    = "localhost"
	defaultTimeout = 10 * time.Second
)

func main() {

	homePath := getEnv("HOME")
	host := flag.String("host", defaultAddr, "IP address of fileserver where it runs")
	port := flag.String("port", "8080", "Port where fileserver runs on")
	rootDir := flag.String("path", homePath, "Directory Path which you want to share using fileserver")
	certFile := flag.String("cert", "", "path to the public cert")
	keyFile := flag.String("key", "", "path to the private key")
	gzip := flag.Bool("gzip", false, "enable gzip")
	v := flag.Bool("v", false, "display version of fileserver")
	flag.Parse()

	if *v {
		fmt.Println("Version: " + version.VERSION)
		fmt.Println("Git Commit: " + version.GITCOMMIT)
		os.Exit(0)
	}

	if *host == defaultAddr {
		ip, err := externalIP()
		if err != nil {
			log.Fatal(err)
		}
		*host = ip
	}

	if *rootDir == "." {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalln(err)
		}
		*rootDir = cwd
	}

	var handle http.Handler
	if *gzip {
		handle = &GzHandler{path: *rootDir}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	fs := NewFileServer(*host, *port, *rootDir, *certFile, *keyFile, defaultTimeout, handle)

	go fs.Start()

	<-c
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Recevied SIGINT signal")
	log.Println("shutting down server")

	fs.Stop(ctx)
	os.Exit(0)
}
