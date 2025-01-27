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

---
### Content sniffing
In order to automatically set the Content-Type header, Go content sniffs the response body with the [http.DetectContentType()](https://pkg.go.dev/net/http#DetectContentType) function. If this function can’t guess the content type, Go will fall back to setting the header `Content-Type: application/octet-stream` instead.

The `http.DetectContentType()` function generally works quite well, but a common gotcha for web developers is that it can’t distinguish JSON from plain text. So, by default, JSON responses will be sent with a `Content-Type: text/plain; charset=utf-8` header. You can prevent this from happening by setting the correct header manually in your handler like so:

```go
w.Header().Set("Content-Type", "application/json")
w.Write([]byte(`{"name":"Alex"}`))”
```

---
### Manipulating the header map

```
// Set a new cache-control header. If an existing "Cache-Control" header exists
// it will be overwritten.
w.Header().Set("Cache-Control", "public, max-age=31536000")

// In contrast, the Add() method appends a new "Cache-Control" header and can
// be called multiple times.
w.Header().Add("Cache-Control", "public")
w.Header().Add("Cache-Control", "max-age=31536000")

// Delete all values for the "Cache-Control" header.
w.Header().Del("Cache-Control")

// Retrieve the first value for the "Cache-Control" header.
w.Header().Get("Cache-Control")

// Retrieve a slice of all values for the "Cache-Control" header.
w.Header().Values("Cache-Control")
```

---
### Header canonicalization
When you’re using the `Set()`, `Add()`, `Del()`, `Get()` and `Values()` methods on the header map, the header name will always be canonicalized using the [textproto.CanonicalMIMEHeaderKey()](https://pkg.go.dev/net/textproto#CanonicalMIMEHeaderKey) function. This converts the first letter and any letter following a hyphen to upper case, and the rest of the letters to lowercase. This has the practical implication that when calling these methods the header name is case-insensitive.

If you need to avoid this canonicalization behavior, you can edit the underlying header map directly. It has the type `map[string][]string` behind the scenes. For example:

```go
w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}
```

> Note: If a `HTTP/2` connection is being used, Go will always automatically convert the header names and values to lowercase for you when writing the response, as per the [HTTP/2 specifications](https://datatracker.ietf.org/doc/html/rfc7540#section-8.1.2).
