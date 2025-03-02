When you’re building a web application there’s probably some shared functionality that you want to use for many (or even all) HTTP requests. For example, you might want to log every request, compress every response, or check a cache before passing the request to your handlers.

A common way of organizing this shared functionality is to set it up as middleware. This is essentially some self-contained code which independently acts on a request before or after your normal application handlers.

---
### 6.1 How middleware works
> You can think of a Go web application as a chain of `ServeHTTP()` methods being called one after another.

Currently, in our application, when our server receives a new HTTP request it calls the servemux’s `ServeHTTP()` method. This looks up the relevant handler based on the request method and URL path, and in turn calls that handler’s `ServeHTTP()` method.

The basic idea of middleware is to insert another handler into this chain. The middleware handler executes some logic, like logging a request, and then calls the `ServeHTTP()` method of the next handler in the chain.

In fact, we’re actually already using some middleware in our application — the `http.StripPrefix()` function from serving static files, which removes a specific prefix from the request’s URL path before passing the request on to the file server.

### The pattern
The standard pattern for creating your own middleware looks like this:
```go
func myMiddleware(next http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        // TODO: Execute our middleware logic here...
        next.ServeHTTP(w, r)
    }

    return http.HandlerFunc(fn)
}
```

- The `myMiddleware()` function is essentially a wrapper around the next handler, which we pass to it as a parameter.
- It establishes a function fn which closes over the next handler to form a closure. When fn is run it executes our middleware logic and then transfers control to the next handler by calling it’s `ServeHTTP()` method.
- Regardless of what you do with a closure it will always be able to access the variables that are local to the scope it was created in — which in this case means that fn will always have access to the next variable.
- In the final line of code, we then convert this closure to a `http.Handler` and return it using the `http.HandlerFunc()` adapter.

If this feels confusing, you can think of it more simply: `myMiddleware()` is a function that accepts the next handler in a chain as a parameter. It returns a handler which executes some logic and then calls the next handler.

### Simplifying the middleware
A tweak to this pattern is to use an anonymous function inside `myMiddleware()` middleware:
```go   
func myMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // TODO: Execute our middleware logic here...
        next.ServeHTTP(w, r)
    })
}
```
This pattern is very common in the wild, and the one that you’ll probably see most often if you’re reading the source code of other applications or third-party packages.

### Positioning the middleware
It’s important to explain that where you position the middleware in the chain of handlers will affect the behavior of your application.

If you position your middleware before the `servemux` in the chain then it will act on every request that your application receives.
```
myMiddleware → servemux → application handler
```

A good example of where this would be useful is middleware to log requests — as that’s typically something you would want to do for all requests.

Alternatively, you can position the middleware after the `servemux` in the chain — by wrapping a specific application handler. This would cause your middleware to only be executed for a specific route.
```
servemux → myMiddleware → application handler
```

An example of this would be something like authorization middleware, which you may only want to run on specific routes.
