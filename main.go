package main

import (
	"fmt"
	"time"
	"log"
	"flag"
	"net/http"
	"os"
)

func getEnv(env string) string {
	val := os.Getenv(env)
	return val
}

func main() {

	var address, port, path string
	path = getEnv("HOME")
	flag.StringVar(&address, "addr", "127.0.0.1", "IP address of fileserver where it runs")
	flag.StringVar(&port, "port", "8080", "Port where fileserver runs on")
	flag.StringVar(&path, "path", path, "Directory Path which you want to share using fileserver")
	flag.Parse()

	color := func(s string) string {
		return fmt.Sprintf("\x1b[1;33m%v\x1b[0m",s)
	}

	fileServerHandler := http.FileServer(http.Dir(path))
	go func(){
		time.Sleep(time.Second*1)
		fmt.Printf("Server running on %v and port %v serving %v",color(address),color(port),color(path))
	}()
	log.Fatalln(http.ListenAndServe(address+":"+port, fileServerHandler))
}