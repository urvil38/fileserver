package main

import (
	"context"
	"flag"
	"fmt"
	"log"
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
	defaultAddr    = "0.0.0.0"
	defaultTimeout = 10 * time.Second
)

func main() {

	homePath := getEnv("HOME")
	host := flag.String("host", defaultAddr, "IP address of fileserver on which it listen on")
	port := flag.String("port", "8080", "port on which fileserver runs on")
	rootDir := flag.String("path", homePath, "path to the directory you want to share using fileserver")
	certFile := flag.String("cert", "", "path to the public cert file")
	keyFile := flag.String("key", "", "path to the private key file")
	timeout := flag.Duration("timeout", defaultTimeout, "read/write timeout")
	hideDotFiles := flag.Bool("no-dotfiles", false, "weather to show file starting with dot i.e. hidden files")
	silent := flag.Bool("silent", false, "Suppress log messages from output")
	logIP := flag.Bool("log-ip", false, "Log ip address of incoming request")
	gzip := flag.Bool("gzip", false, "enable gzip")
	username := flag.String("username", "", "Username for basic authentication")
	password := flag.String("password", "", "Password for basic authentication")
	v := flag.Bool("v", false, "display version of fileserver")
	flag.Parse()

	if *v {
		fmt.Println("Version: " + version.VERSION)
		fmt.Println("Git Commit: " + version.GITCOMMIT)
		os.Exit(0)
	}

	if *rootDir == "." {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalln(err)
		}
		*rootDir = cwd
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	fs := NewFileServer(Config{
		rootDir:      *rootDir,
		host:         *host,
		port:         *port,
		certFile:     *certFile,
		keyFile:      *keyFile,
		gzipEnable:   *gzip,
		timeout:      *timeout,
		hideDotFiles: *hideDotFiles,
		silent:       *silent,
		logIP:        *logIP,
		username:     *username,
		password:     *password,
	})

	go fs.Start()

	<-c
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Recevied SIGINT signal")
	log.Println("shutting down server")

	fs.Stop(ctx)
	os.Exit(0)
}
