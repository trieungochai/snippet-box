Let’s neaten up our application by moving some of the error handling code into helper methods. This will help [separate our concerns](https://deviq.com/principles/separation-of-concerns) and stop us repeating code as we progress through the build.

The [http.StatusText()](https://pkg.go.dev/net/http/#StatusText) function. This returns a human-friendly text representation of a given HTTP status code.
For example
- `http.StatusText(400)` will return the string `Bad Request`.
- `http.StatusText(500)` will return the string `Internal Server Error`.
