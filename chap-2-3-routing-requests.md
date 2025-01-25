![routing-requests](routing-requests.png)

---

### Trailing slashes in route patterns

It’s important to know that Go’s servemux has different matching rules depending on whether a route pattern ends with a trailing slash or not.

- When a pattern doesn’t have a trailing slash, it will only be matched (and the corresponding handler called) when the request URL path exactly matches the pattern in full.
- When a route pattern ends with a trailing slash — like `"/"` or `"/static/"` — it is known as a subtree path pattern. Subtree path patterns are matched (and the corresponding handler called) whenever the start of a request URL path matches the subtree path. If it helps your understanding, you can think of subtree paths as acting a bit like they have a wildcard at the end, like `"/**"` or `"/static/**"`.

This helps explain why the `"/"` route pattern acts like a catch-all.

---
### Restricting subtree paths

To prevent subtree path patterns from acting like they have a wildcard at the end, you can append the special character sequence `{$}` to the end of the pattern  — like `"/{$}"` or `"/static/{$}`.

So if you have the route pattern `"/{$}"`, it effectively means match a single slash, followed by nothing else. It will only match requests where the URL path is exactly `/`.

---
### Additional servemux features

Request URL paths are automatically sanitized. If the request path contains any `.` or `..` elements or repeated slashes, the user will automatically be redirected to an equivalent clean URL. For example, if a user makes a request to `/foo/bar/..//baz` they will automatically be sent a 301 Permanent Redirect to `/foo/baz` instead.

If a subtree path has been registered and a request is received for that subtree path without a trailing slash, then the user will automatically be sent a 301 Permanent Redirect to the subtree path with the slash added. For example, if you have registered the subtree path `/foo/`, then any request to `/foo` will be redirected to `/foo/`.

---
### Host name matching

It’s possible to include host names in your route patterns. This can be useful when you want to redirect all HTTP requests to a canonical URL, or if your application is acting as the back end for multiple sites or services.

For example:
```go
mux := http.NewServeMux()
mux.HandleFunc("foo.example.org/", fooHandler)
mux.HandleFunc("bar.example.org/", barHandler)
mux.HandleFunc("/baz", bazHandler)
```

When it comes to pattern matching, any host-specific patterns will be checked first and if there is a match the request will be dispatched to the corresponding handler. Only when there isn’t a host-specific match found will the non-host specific patterns also be checked.

---
### The default servemux

The `http.Handle()` and `http.HandleFunc()` functions - these allow you to register routes without explicitly declaring a servemux:

```go
func main() {
    http.HandleFunc("/", home)
    http.HandleFunc("/snippet/view", snippetView)
    http.HandleFunc("/snippet/create", snippetCreate)

    log.Print("starting server on :4000")
    
    err := http.ListenAndServe(":4000", nil)
    log.Fatal(err)
}
```

Behind the scenes, these functions register their routes with something called the default `servemux`. This is just a regular servemux like we’ve already been using, but which is initialized automatically by Go and stored in the `http.DefaultServeMux` global variable.

If you pass nil as the second argument to `http.ListenAndServe()`, the server will use `http.DefaultServeMux` for routing.

Although this approach can make your code slightly shorter, we don’t recommend it for 2 reasons:
- It’s less explicit and feels more ‘magic’ than declaring and using your own locally-scoped servemux.
- Because `http.DefaultServeMux` is a global variable in the standard library, it means any Go code in your project can access it and potentially register a route. If you have a large project codebase (especially one that is being worked on by many people), that can make it harder to ensure all route declarations for your application are easily discoverable in one central place.
