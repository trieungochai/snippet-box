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
