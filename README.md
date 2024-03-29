# Fileserver

![fileserver-0](./docs/img/fileserver-0.png)
![fileserver-1](./docs/img/fileserver-1.png)

# Download

- Download appropriate pre-compiled binary from the [release](https://github.com/urvil38/fileserver/releases) page.

```sh
# download binary using cURL
curl -L https://github.com/urvil38/fileserver/releases/download/3.0.2/fileserver-darwin-amd64 -o fileserver

# make binary executable
chmod +x ./fileserver

# move it to bin dir (user need to has root privileges. run following command as root user using sudo.
sudo mv ./fileserver /usr/local/bin
```

- Download using `go get`

```sh
go get -u github.com/urvil38/fileserver
```

# Docker
- Build Docker image:
```sh
make image
```

- Run fileserver using Docker:
```sh
docker run -d -p 8080:8080 -v $(PWD):/$(PWD) urvil38/fileserver:3.0.2 fileserver -port 8080 -path /
```
# Usage

```
Usage of fileserver:
  -cert string
    	path to the public cert file
  -gzip
    	enable gzip
  -host string
    	IP address of fileserver on which it listen on (default "0.0.0.0")
  -key string
    	path to the private key file
  -log-ip
    	Log ip address of incoming request
  -no-dotfiles
    	weather to show file starting with dot i.e. hidden files
  -password string
    	Password for basic authentication
  -path string
    	path to the directory you want to share using fileserver (default "/Users/urvilpatel")
  -port string
    	port on which fileserver runs on (default "8080")
  -quiet
    	Suppress log messages from output
  -timeout duration
    	read/write timeout (default 10s)
  -username string
    	Username for basic authentication
  -v	display version of fileserver
```

# Build

- If you want to build fileserver right away, you need a working [Go environment](https://golang.org/doc/install). It requires Go version 1.12 and above.

```sh
git clone https://github.com/urvil38/fileserver.git
cd fileserver
make build
```
