By default, every response that your handlers send has the HTTP status code 200 OK (which indicates to the user that their request was received and processed successfully), plus 3 automatic system-generated headers: a `Date header`, and the `Content-Length` and  `Content-Type` of the response body.

```
$ curl -i localhost:4000/
HTTP/1.1 200 OK
Date: Wed, 18 Mar 2024 11:29:23 GMT
Content-Length: 21
Content-Type: text/plain; charset=utf-8

Hello from Snippetbox
```

---
### Status code constants
The `net/http` package provides [constants for HTTP status codes](https://pkg.go.dev/net/http#pkg-constants), which we can use instead of writing the status code number ourselves. Using these constants is good practice because it helps prevent mistakes due to typos, and it can also help make your code clearer and self-documenting especially when dealing with less-commonly-used status codes.

---
### Customizing headers
You can also customize the HTTP headers sent to a user by changing the response header map. Probably the most common thing you’ll want to do is include an additional header in the map, which you can do using the `w.Header().Add()` method.

---
### Writing response bodies
So far we’ve been using `w.Write()` to send specific  HTTP response bodies to a user. And while this is the simplest and most fundamental way to send a response, in practice it’s far more common to pass your `http.ResponseWriter` value to another function that writes the response for you.

In fact, there are a lot of functions that you can use to write a response!

The key thing to understand is this… because the `http.ResponseWriter` value in your handlers has a `Write()` method, it satisfies the [io.Writer](https://pkg.go.dev/io#Writer) interface.

Any functions where you see an `io.Writer` parameter, you can pass in your `http.ResponseWriter` value and whatever is being written will subsequently be sent as the body of the HTTP response.

That means you can use standard library functions like [io.WriteString()](https://pkg.go.dev/io#WriteString) and the [fmt.Fprint*()](https://pkg.go.dev/fmt#Fprint) family (all of which accept an `io.Writer` parameter) to write plain-text response bodies too.

```go
// Instead of this...
w.Write([]byte("Hello world"))

// You can do this...
io.WriteString(w, "Hello world")
fmt.Fprint(w, "Hello world")”
```
