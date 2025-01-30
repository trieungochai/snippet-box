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

---
### Environment variables
If you’ve built and deployed web applications before, then you’re probably thinking what about environment variables? Surely it’s [good-practice](http://12factor.net/config) to store configuration settings there?

You can store your configuration settings in environment variables and access them directly from your application by using the [os.Getenv()](https://pkg.go.dev/os/#Getenv) function like so:
```go
addr := os.Getenv("SNIPPETBOX_ADDR")
```

But this has some drawbacks compared to using command-line flags. You can’t specify a default setting (the return value from `os.Getenv()` is the empty string if the environment variable doesn’t exist), you don’t get the `-help` functionality that you do with command-line flags, and the return value from `os.Getenv()` is always a string — you don’t get automatic type conversions like you do with `flag.Int()`, `flag.Bool()` and the other command line flag functions.

Instead, you can get the best of both worlds by passing the environment variable as a command-line flag when starting the application. For example:

```
$ export SNIPPETBOX_ADDR=":9999"
$ go run ./cmd/web -addr=$SNIPPETBOX_ADDR
2024/03/18 11:29:23 starting server on :9999”
```

---
### Boolean flags
For flags defined with `flag.Bool()`, omitting a value when starting the application is the same as writing `-flag=true`. The following 2 commands are equivalent:
```
$ go run ./example -flag=true
$ go run ./example -flag
```

You must explicitly use `-flag=false` if you want to set a boolean flag value to false.

---
### Pre-existing variables
It’s possible to parse command-line flag values into the memory addresses of pre-existing variables, using [flag.StringVar()](https://pkg.go.dev/flag/#FlagSet.StringVar), [flag.IntVar()](https://pkg.go.dev/flag#FlagSet.IntVar), [flag.BoolVar()](https://pkg.go.dev/flag#FlagSet.BoolVar), and similar functions for other types.

These functions are particularly useful if you want to store all your configuration settings in a single struct. As a rough example:
```go
type config struct {
    addr      string
    staticDir string
}

...

var cfg config
flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
flag.Parse()
```