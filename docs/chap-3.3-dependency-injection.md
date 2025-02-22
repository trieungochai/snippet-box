DI is a recommended design pattern and a software development technique used to implement the `Inversion of Control (IoC)` principle in software applications. IoC is the transfer of control of objects, object dependencies, and general application flow to a different component or part of the application.

While IoC and DI are sometimes used interchangeably, it's important to note that they are not identical concepts. Dependency injection is a specific technique for achieving IoC, but IoC encompasses broader techniques and patterns.

In DI, the dependencies of an object (i.e. the objects it relies on) are provided externally rather than created internally by the object itself. As an example, suppose Service X requires a function from Service Y to perform its operations effectively. Instead of Service X internally creating a new instance of Service Y, the recommended approach with DI is to have a separate component responsible for creating an instance of Service Y and then "inject" that instance into Service X.

---
[Dependency Injection](https://dev.to/dsysd_dev/dependency-injection-like-a-pro-in-golang-43ao) is just a fancy way of saying `"passing stuff into a function"`. It's about giving a function or object the things it needs to work.

It's like baking pizza. You don't bake the pizza at home; you call a pizza place and tell them what you want. The pizza delivery is like a function that takes your order (dependencies) and delivers the result (pizza).
```go
func OrderPizza() {
    oven := NewOven()
    pizza := oven.BakePizza()
    // ...
}
```

This function creates its own dependencies (the oven) instead of having them provided.
```go
func OrderPizza(oven Oven) {
    pizza := oven.BakePizza()
    // ...
}
```

The function is more flexible because you can now choose which oven to use (e.g., regular oven or wood-fired oven).

---
Dependency Injection is widely used in Go, especially in building HTTP servers. Explain that you can pass different routers, middleware, and database connections to the server as dependencies.
```go
package main

import (
 "fmt"
 "net/http"
)

// Router defines the interface for a router.
type Router interface {
 HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

// MyRouter is an example implementation of the Router interface.
type MyRouter struct{}

func (r *MyRouter) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
 http.HandleFunc(pattern, handler)
}

// LoggerMiddleware is an example middleware.
func LoggerMiddleware(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Logging:", r.URL.Path)
  next.ServeHTTP(w, r)
 })
}

func main() {
 // Create a router instance
 router := &MyRouter{}

 // Attach middleware
 http.Handle("/", LoggerMiddleware(router))

 // Define routes
 router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Welcome to our website!")
 })

 router.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "About Us")
 })

 router.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Contact Us")
 })

 // Start the server
 fmt.Println("Server is listening on :8080")
 http.ListenAndServe(":8080", nil)
}
```

In this example:
- We define a `Router interface` that has a HandleFunc method, which mimics the behavior of Go's `http.HandleFunc`.
- `MyRouter` is an implementation of the `Router interface`.
- We create a `LoggerMiddleware` function, which is a middleware that logs incoming requests.
- In the `main` function, we create an instance of `MyRouter`.
- We attach the `LoggerMiddleware` to the router to log requests.
- We define routes using the `HandleFunc` method provided by our router, which is the equivalent of http.HandleFunc.
- Finally, we start the server using `http.ListenAndServe`.

This example demonstrates how you can use dependency injection to pass in routers and middleware to create a flexible and modular web server in Go.

You can easily swap out different routers or middleware components to customize your server's behavior.

---

The home handler function is still writing error messages using Go’s standard logger, not the structured logger that we now want to be using.
```go
func home(w http.ResponseWriter, r *http.Request) {
    ...

    ts, err := template.ParseFiles(files...)
    if err != nil {
        log.Print(err.Error()) // This isn't using our new structured logger.
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    err = ts.ExecuteTemplate(w, "base", nil)
    if err != nil {
        log.Print(err.Error()) // This isn't using our new structured logger.
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}
```

This raises a good question: `how can we make our new structured logger available to our home function from main()`?

Most web applications will have multiple dependencies that their handlers need to access, such as a database connection pool, centralized error handlers, and template caches. What we really want to answer is: `how can we make any dependency available to our handlers?`

There are a few different ways to do this, the simplest being to just put the dependencies in global variables. But in general, it is good practice to inject dependencies into your handlers. It makes your code more explicit, less error-prone, and easier to unit test than if you use global variables.

For applications where all your handlers are in the same package, like ours, a neat way to inject dependencies is to put them into a custom application struct, and then define your handler functions as methods against application.

---
### Closures for DI
The pattern that we’re using to inject dependencies won’t work if your handlers are spread across multiple packages. In that case, an alternative approach is to create a standalone config package which exports an Application struct, and have your handler functions close over this to form a closure.

```go
// package config

type Application struct {
    Logger *slog.Logger
}
```

```go
// package foo

func ExampleHandler(app *config.Application) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ...
        ts, err := template.ParseFiles(files...)
        if err != nil {
            app.Logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
        ...
    }
}
```

```go
// package main

func main() {
    app := &config.Application{
        Logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
    }
    ...
    mux.Handle("/", foo.ExampleHandler(app))
    ...
}
```