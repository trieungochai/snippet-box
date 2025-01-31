Both the `log.Printf()` and `log.Fatal()` functions output log entries using Go’s standard logger from the log package, which — by default — prefixes a message with the local date and time and writes out to the standard error stream (which should display in your terminal window).

For many applications, using the standard logger will be good enough, and there’s no need to do anything more complex.

But for applications which do a lot of logging, you may want to make the log entries easier to filter and work with. For example, you might want to distinguish between different severities of log entries (like informational and error entries), or to enforce a consistent structure for log entries so that they are easy for external programs or services to parse.

To support this, the Go standard library includes the [log/slog](https://pkg.go.dev/log/slog) package which lets you create custom structured loggers that output log entries in a set format. Each log entry includes the following things:
- A timestamp with millisecond precision.
- The severity level of the log entry (Debug, Info, Warn or Error).
- The log message (an arbitrary string value).
- Optionally, any number of key-value pairs (known as attributes) containing additional information.

---
### Creating a structured logger
The key thing to understand is that all structured loggers have a structured logging handler associated with them (not to be confused with a HTTP handler), and it’s actually this handler that controls how log entries are formatted and where they are written to.

The code for creating a logger looks like this:
```go
loggerHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{...})
logger := slog.New(loggerHandler)
```

In the 1st line of code we first use the [slog.NewTextHandler()](https://pkg.go.dev/log/slog#NewTextHandler) function to create the structured logging handler. This function accepts 2 arguments:
- The 1st argument is the write destination for the log entries. In the example above we’ve set it to `os.Stdout`, which means it will write log entries to the standard out stream.
- The 2nd argument is a pointer to a [slog.HandlerOptions](https://pkg.go.dev/log/slog#HandlerOptions) struct , which you can use to customize the behavior of the handler.

Then in the 2nd line of code, we actually create the structured logger by passing the handler to the [slog.New()](https://pkg.go.dev/log/slog#New) function.

In practice, it’s more common to do all this  in a single line of code:
```go
logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{...}))
```

---
### Using a structured logger
Once you’ve created a structured logger, you can then write a log entry at a specific severity level by calling the `Debug()`, `Info()`, `Warn()` or `Error()` methods on the logger. As an example, the following line of code:
```go
logger.Info("request received")
```

Would result in a log entry that looks like this:
```
time=2024-03-18T11:29:23.000+00:00 level=INFO msg="request received”
```

The `Debug()`, `Info()`, `Warn()` or `Error()` methods are variadic methods which accept an arbitrary number of additional attributes (key-value pairs). Like so:
```go
logger.Info("request received", "method", "GET", "path", "/")
```
In this example, we’ve added 2 extra attributes to the log entry: the key `"method"` and value `"GET"`, and the key `"path"` and value `"/"`.  Attribute keys must always be strings, but the values can be of any type. In this example, the log entry will look like this:
```
time=2024-03-18T11:29:23.000+00:00 level=INFO msg="request received" method=GET path=/”
```

>Note: If your attribute keys, values, or log message contain " or = characters or any whitespace, they will be wrapped in double quotes in the log output. We can see this behavior in the example above, where the log message msg="request received" is quoted.

---
### JSON formatted logs
The `slog.NewTextHandler()` function that we’ve used in this chapter creates a handler that writes plaintext log entries. But it’s possible to create a handler that writes log entries as JSON objects instead using the `slog.NewJSONHandler()` function. Like so:
```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
```

When using the JSON handler, the log output will look similar to this:
```
{"time":"2024-03-18T11:29:23.00000000+00:00","level":"INFO","msg":"starting server","addr":":4000"}
{"time":"2024-03-18T11:29:23.00000000+00:00","level":"ERROR","msg":"listen tcp :4000: bind: address already in use"}
```

---
### Minimum log level
The `log/slog` package supports 4 severity levels: `Debug`, `Info`, `Warn` and `Error` in that order. Debug is the least severe level, and Error is the most severe.

By default, the minimum log level for a structured logger is Info. That means that any log entries with a severity less than Info  — i.e. Debug level entries — will be silently discarded.

You can use the `slog.HandlerOptions` struct to override this and set the minimum level to Debug (or any other level) if you want:
```go
logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))
```

---
### Caller location
You can also customize the handler so that it includes the filename and line number of the calling source code in the log entries, like so:
```go
logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
    AddSource: true,
}))
```

The log entries will look similar to this, with the caller location recorded under the source key:
```
time=2024-03-18T11:29:23.000+00:00 level=INFO source=/home/alex/code/snippetbox/cmd/web/main.go:32 msg="starting server" addr=:4000
```

---
### Decoupled logging
In this chapter we’ve set up our structured logger to write entries to `os.Stdout` — the standard out stream.

The big benefit of writing log entries to `os.Stdout` is that your application and logging are decoupled. Your application itself isn’t concerned with the routing or storage of the logs, and that can make it easier to manage the logs differently depending on the environment.

During development, it’s easy to view the log output because the standard out stream is displayed in the terminal.

In staging or production environments, you can redirect the stream to a final destination for viewing and archival. This destination could be on-disk files, or a logging service such as Splunk. Either way, the final destination of the logs can be managed by your execution environment independently of the application.

For example, we could redirect the standard out stream to a on-disk file when starting the application like so:
```
$ go run ./cmd/web >>/tmp/web.log
```
>Note: Using the double arrow >> will append to an existing file, instead of truncating it when starting the application.

---
### Concurrent logging
Custom loggers created by `slog.New()` are concurrency-safe. You can share a single logger and use it across multiple goroutines and in your HTTP handlers without needing to worry about race conditions.

That said, if you have multiple structured loggers writing to the same destination then you need to be careful and ensure that the destination’s underlying `Write()` method is also safe for concurrent use.
