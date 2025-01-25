It’s also possible to define route patterns that contain wildcard segments. You can use these to create more flexible routing rules, and also to pass variables to your Go application via a request URL.

Wildcard segments in a route pattern are denoted by an wildcard identifier inside `{}` brackets.

```go
mux.HandleFunc("/products/{category}/item/{itemID}", exampleHandler)
```

In this example, the route pattern contains 2 wildcard segments. 
- The 1st segment has the identifier category.
- The 2nd has the identifier itemID.

<b>Important</b>: When defining a route pattern, each path segment (the bit between forward slash characters) can only contain one wildcard and the wildcard needs to fill the whole path segment. Patterns like `"/products/c_{category}"`, `/date/{y}-{m}-{d}` or `/{slug}.html` are not valid.

Inside your handler, you can retrieve the corresponding value for a wildcard segment using its identifier and the `r.PathValue()` method.

```go
func exampleHandler(w http.ResponseWriter, r *http.Request) {
    category := r.PathValue("category")
    itemID := r.PathValue("itemID")

    ...
}
```

The `r.PathValue()` method always returns a string value.