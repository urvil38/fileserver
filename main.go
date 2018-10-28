package main

import (
	"context"
	"os/signal"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func getEnv(env string) string {
	val := os.Getenv(env)
	return val
}

var version string

func main() {

	var address, port, path string
	var v bool
	defultAddr := "127.0.0.1" 

	path = getEnv("HOME")
	flag.StringVar(&address, "addr", defultAddr, "IP address of fileserver where it runs")
	flag.StringVar(&port, "port", "8080", "Port where fileserver runs on")
	flag.StringVar(&path, "path", path, "Directory Path which you want to share using fileserver")
	flag.BoolVar(&v,"v",false,"display version of fileserver")
	flag.Parse()

	if v {
		fmt.Println("Version: "+version)
		os.Exit(0)
	}

	var ipv4Addr []string
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		log.Fatal(err)
	}
	for _, addr := range addrs {
		if !strings.Contains(addr, "::") {
			ipv4Addr = append(ipv4Addr, addr)
		}
	}
	
	if len(ipv4Addr) > 0 {
		address = ipv4Addr[0]
	}

	if path == "." {
		cwd,err := os.Getwd()
		if err != nil {
			log.Fatalln(err)
		}
		path = cwd
	}

	color := func(s string) string {
		return fmt.Sprintf("\x1b[1;33m%v\x1b[0m", s)
	}

	fileServerHandler := http.FileServer(http.Dir(path))
	c := make(chan os.Signal, 1)

	signal.Notify(c,os.Interrupt)

	server := http.Server{
		Handler: fileServerHandler,
		Addr: address+":"+port,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		time.Sleep(time.Millisecond * 300)
		fmt.Printf("Server running on %v%v:%v serving %v\n",color("http://") ,color(address), color(port), color(path))
	}()

	go func() {
		if err := server.ListenAndServe() ; err != nil {
			log.Println(err)
		}
	}()

	<-c
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	server.Shutdown(ctx)
	log.Println("Recevied SIGINT signal")
	log.Println("shutting down server")
    os.Exit(0)
}