[![Go](https://github.com/eliben/static-server/actions/workflows/go.yml/badge.svg)](https://github.com/eliben/static-server/actions/workflows/go.yml)

----

# static-server

Simple static file server, with support of HTTP and HTTPS. There are no configuration files and no dependencies (except one for testing). Serving the current directory on HTTP port 8080 is as simple as invoking:

```
$ go run github.com/eliben/static-server@latest
<timestamp> Serving directory "." on http://127.0.0.1:8080
```

If you want to install `static-server` locally first, you can run:

```
$ go install github.com/eliben/static-server@latest
```

And then invoke `static-server` as needed. Command-line flags can be used to
configure the behavior of the server:

```
$ static-server -h
Usage: ./static-server [dir]

  [dir] is optional; if not passed, '.' is used.

  By default, the server listens on localhost:8080. Both the
  host and the port are configurable with flags. Set the host
  to something else if you want the server to listen on a
  specific network interface. Setting the port to 0 will
  instruct the server to pick a random available port.

  -addr string
    	full address (host:port) to listen on; don't use this if 'port' or 'host' are set (default "localhost:8080")
  -certfile string
    	TLS certificate file to use with -tls (default "cert.pem")
  -cors
    	enable CORS by returning Access-Control-Allow-Origin header
  -host string
    	specific host to listen on (default "localhost")
  -keyfile string
    	TLS key file to use with -tls (default "key.pem")
  -port string
    	port to listen on; if 0, a random available port will be used (default "8080")
  -silent
    	suppress messages from output (reporting only errors)
  -tls
    	enable HTTPS serving with TLS
  -version
    	print version and exit
```
