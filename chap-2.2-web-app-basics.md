3 absolute essentials:
- The 1st thing we need is a handler. If you’ve previously built web applications using a MVC pattern, you can think of handlers as being a bit like controllers. They’re responsible for executing your application logic and for writing HTTP response headers and bodies.
- The 2nd component is a router (or servemux in Go terminology). This stores a mapping between the URL routing patterns for your application and the corresponding handlers. Usually you have one servemux for your application containing all your routes.
- The last thing we need is a web server. One of the great things about Go is that you can establish a web server and listen for incoming requests as part of your application itself. You don’t need an external third-party server like Nginx, Apache or Caddy.

---
### Network addresses

```go
err := http.ListenAndServe(":4000", mux)
```

The TCP network address that you pass to `http.ListenAndServe()` should be in the format `host:port`.

If you omit the host (like we did with `:4000`) then the server will listen on all your computer’s available network interfaces. Generally, you only need to specify a host in the address if your computer has multiple network interfaces and you want to listen on just one of them.

In other Go projects or documentation you might sometimes see network addresses written using named ports like `:http` or `:http-alt` instead of a number. If you use a named port then the `http.ListenAndServe()` function will attempt to look up the relevant port number from your `/etc/services` file when starting the server, returning an error if a match can’t be found.

---
### Using `go run`

During development the go run command is a convenient way to try out your code. It’s essentially a shortcut that compiles your code, creates an executable binary in your `/tmp` directory, and then runs this binary in one step.