Strictly speaking, what we mean by handler is an object which satisfies the `http.Handler` interface:
```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

In simple terms, this basically means that to be a handler an object must have a `ServeHTTP()` method with the exact signature:
```go
ServeHTTP(http.ResponseWriter, *http.Request)
```

So in its simplest form a handler might look something like this:
```go
type home struct {}

func (h *home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("This is my home page"))
}
```

Here we have an object (in this case it’s an empty home struct, but it could equally be a string or function or anything else), and we’ve implemented a method with the signature `ServeHTTP(http.ResponseWriter, *http.Request)` on it. That’s all we need to make a handler.

You could then register this with a servemux using the Handle method like so:
```go
mux := http.NewServeMux()
mux.Handle("/", &home{})
```

When this servemux receives a HTTP request for `"/"`, it will then call the `ServeHTTP()` method of the home struct — which in turn writes the HTTP response.

---
### Chaining handlers
The [http.ListenAndServe()](https://pkg.go.dev/net/http#ListenAndServe) function takes a `http.Handler` object as the second parameter:
```go
func ListenAndServe(addr string, handler Handler) error
```

… but we’ve been passing in a servemux.

We were able to do this because the servemux also has a [ServeHTTP()](https://pkg.go.dev/net/http#ServeMux.ServeHTTP) method, meaning that it too satisfies the `http.Handler` interface.

The servemux as just being a special kind of handler, which instead of providing a response itself passes the request on to a second handler. Chaining handlers together is a very common idiom in Go.

In fact, what exactly is happening is this: When our server receives a new HTTP request, it calls the servemux’s `ServeHTTP()` method. This looks up the relevant handler based on the request method and URL path, and in turn calls that handler’s `ServeHTTP()` method. You can think of a Go web application as a chain of `ServeHTTP()` methods being called one after another.

---
### Requests are handled concurrently
There is one more thing that’s really important to point out: all incoming HTTP requests are served in their own goroutine. For busy servers, this means it’s very likely that the code in or called by your handlers will be running concurrently. While this helps make Go blazingly fast, the downside is that you need to be aware of (and protect against) [race conditions](https://www.alexedwards.net/blog/understanding-mutexes) when accessing shared resources from your handlers.
