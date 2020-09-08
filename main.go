package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
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

var defaultAddr = "127.0.0.1"

func main() {

	var address, port, path string
	var v bool

	path = getEnv("HOME")
	flag.StringVar(&address, "addr", defaultAddr, "IP address of fileserver where it runs")
	flag.StringVar(&port, "port", "8080", "Port where fileserver runs on")
	flag.StringVar(&path, "path", path, "Directory Path which you want to share using fileserver")
	flag.BoolVar(&v, "v", false, "display version of fileserver")
	flag.Parse()

	if v {
		fmt.Println("Version: " + version.VERSION)
		fmt.Println("Git Commit: " + version.GITCOMMIT)
		os.Exit(0)
	}

	if address == defaultAddr {
		ip, err := externalIP()
		if err != nil {
			log.Fatal(err)
		}
		address = ip
	}

	if path == "." {
		cwd, err := os.Getwd()
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

	signal.Notify(c, os.Interrupt)

	server := http.Server{
		Handler:      fileServerHandler,
		Addr:         address + ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		time.Sleep(time.Millisecond * 300)
		fmt.Printf("Server running on %v%v:%v serving %v\n", color("http://"), color(address), color(port), color(path))
	}()

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	<-c
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	server.Shutdown(ctx)
	log.Println("Recevied SIGINT signal")
	log.Println("shutting down server")
	os.Exit(0)
}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return defaultAddr, nil
}
