Our web application’s `main.go` file currently contains a couple of hard-coded configuration settings:
- The network address for the server to listen on (currently `":4000"`)
- The file path for the static files directory (currently `"./ui/static"`)

Having these hard-coded isn’t ideal. There’s no separation between our configuration settings and code, and we can’t change the settings at runtime (which is important if you need different settings for development, testing and production environments).

In this chapter we’ll start to improve that, beginning by making the network address for our server configurable at runtime.

---
### Command-line flags
In Go, a common and idiomatic way to manage configuration settings is to use command-line flags when starting an application. For example:
```
    $ go run ./cmd/web -addr=":80”
```
The easiest way to accept and parse a command-line flag in your application is with a line of code like this:
```go
addr := flag.String("addr", ":4000", "HTTP network address")
```
This essentially defines a new command-line flag with the name `addr`, a default value of `":4000"` and some short help text explaining what the flag controls. The value of the flag will be stored in the `addr` variable at runtime.

>Note: Ports `0-1023` are restricted and (typically) can only be used by services which have root privileges. If you try to use one of these ports you should get a bind: permission denied error message on start-up.

---
### Default values
Command-line flags are completely optional. For instance, if you run the application with no `-addr` flag the server will fall back to listening on address `":4000"` (which is the default value we specified).
```
$ go run ./cmd/web
2024/03/18 11:29:23 starting server on :4000
```

There are no rules about what to use as the default values for your command-line flags.

---
### Type conversions
In the code above we’ve used the `flag.String()` function to define the command-line flag. This has the benefit of converting whatever value the user provides at runtime to a string type. If the value can’t be converted to a string then the application will print an error message and exit.

Go also has a range of other functions including [flag.Int()](https://pkg.go.dev/flag#Int), [flag.Bool()](https://pkg.go.dev/flag#Bool), [flag.Float64()](https://pkg.go.dev/flag#Float64) and [flag.Duration](https://pkg.go.dev/flag#Duration) for defining flags. These work in exactly the same way as `flag.String()`, except they automatically convert the command-line flag value to the appropriate type.

---
### Automated help
Another great feature is that you can use the `-help` flag to list all the available command-line flags for an application and their accompanying help text. Give it a try:

```
$ go run ./cmd/web -help
Usage of /tmp/go-build3672328037/b001/exe/web:
  -addr string
        HTTP network address (default ":4000")
```
