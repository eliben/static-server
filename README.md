# static-server

Simple static file server, written in Go. There are no configuration files and no dependencies (except one for testing). Serving the current directory on HTTP port 8080 is as simple as invoking:

```
$ go run github.com/eliben/static-server@latest
<timestamp> Serving directory "." on http://127.0.0.1:8080
```

If you want to install `static-server` locally first, you can run:

```
$ go install github.com/eliben/static-server@latest
```

And then invoke `static-server` as needed. The `-h` flag will show usage. That's pretty much it!
