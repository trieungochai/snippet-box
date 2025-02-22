Let’s neaten up our application by moving some of the error handling code into helper methods. This will help [separate our concerns](https://deviq.com/principles/separation-of-concerns) and stop us repeating code as we progress through the build.

The [http.StatusText()](https://pkg.go.dev/net/http/#StatusText) function. This returns a human-friendly text representation of a given HTTP status code.
For example
- `http.StatusText(400)` will return the string `Bad Request`.
- `http.StatusText(500)` will return the string `Internal Server Error`.

---
### Stack traces
You can use the [debug.Stack()](https://pkg.go.dev/runtime/debug#Stack) function to get a stack trace outlining the execution path of the application for the current goroutine. Including this as an attribute in your log entries can be helpful for debugging errors.