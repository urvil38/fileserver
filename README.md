# Download

- Download appropriate pre-compiled binary from the [release](https://github.com/urvil38/fileserver/releases) page.

# Usage

```
Usage of ./fileserver:
  -cert string
        path to the public cert
  -gzip
        enable gzip
  -host string
        IP address of fileserver where it runs (default "localhost")
  -key string
        path to the private key
  -path string
        Directory Path which you want to share using fileserver (default "/Users/urvilpatel")
  -port string
        Port where fileserver runs on (default "8080")
  -v    display version of fileserver
```


# Fileserver

![fileserver-0](./docs/img/fileserver-0.png)
![fileserver-1](./docs/img/fileserver-1.png)

# Build

- If you want to build fileserver right away, you need a working [Go environment](https://golang.org/doc/install). It requires Go version 1.12 and above.

```
$ git clone https://github.com/urvil38/fileserver.git
$ cd fileserver
$ make
```
