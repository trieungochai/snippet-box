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