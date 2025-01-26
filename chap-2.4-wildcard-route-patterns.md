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

---
### Precedence and conflicts
When defining route patterns with wildcard segments, it’s possible that some of your patterns will ‘overlap’.

For example, if you define routes with the patterns `"/post/edit"` and `"/post/{id}"` they overlap because an incoming HTTP request with the path `/post/edit` is a valid match for both patterns.

When route patterns overlap, Go’s servemux needs to decide which pattern takes precedent so it can dispatch the request to the appropriate handler.

The rule for this is very neat and succinct: the most specific route pattern wins. Formally, Go defines a pattern as more specific than another if it matches only a subset of requests that the other pattern matches.

Continuing with the example above, the route pattern `"/post/edit"` only matches requests with the exact path `/post/edit`, whereas the pattern `"/post/{id}"` matches requests with the path `/post/edit`, `/post/123`, `/post/abc` and many more. Therefore `"/post/edit"` is the more specific route pattern and will take precedent.

While we’re on this topic, there are a few other things worth mentioning:
- A nice side-effect of the most specific pattern wins rule is that you can register patterns in any order and it won’t change how the servemux behaves.
- There is a potential edge case where you have 2 overlapping route patterns but neither one is obviously more specific than the other. For example, the patterns `"/post/new/{id}"` and `"/post/{author}/latest"` overlap because they both match the request path `/post/new/latest`, but it’s not clear which one should take precedence. In this scenario, Go’s servemux considers the patterns to conflict, and will panic at runtime when initializing the routes.
- Just because Go’s servemux supports overlapping routes, it doesn’t mean that you should use them! Having overlapping route patterns increases the risk of bugs and unintended behavior in your application, and if you have the freedom to design the URL structure for your application it’s generally good practice to keep overlaps to a minimum or avoid them completely.

---
### Subtree path patterns with wildcards
In particular, if your route pattern ends in a trailing slash and has no `{$}` at the end, then it is treated as a subtree path pattern and it only requires the start of a request URL path to match.

So, if you have a subtree path pattern like `"/user/{id}/"` in your routes (note the trailing slash), this pattern will match requests like `/user/1/`, `/user/2/a`, `/user/2/a/b/c` and so on.

Again, if you don’t want that behavior, stick a `{$}` at the end — like `"/user/{id}/{$}"`.

---
### Remainder wildcards
If a route pattern ends with a wildcard, and this final wildcard identifier ends in ..., then the wildcard will match any and all remaining segments of a request path.

For example, if you declare a route pattern like `"/post/{path...}"` it will match requests like `/post/a`, `/post/a/b`, `/post/a/b/c` and so on — very much like a subtree path pattern does. But the difference is that you can access the entire wildcard part via the `r.PathValue()` method in your handlers. In this example, you could get the wildcard value for `{path...}` by calling `r.PathValue("path")`.