### The http.Fileserver handler

Go’s `net/http` package ships with a built-in [http.FileServer](https://pkg.go.dev/net/http#FileServer) handler which you can use to serve files over HTTP from a specific directory.

---
### File server features and functions
Go’s `http.FileServer` handler has a few really nice features that are worth mentioning:
- It sanitizes all request paths by running them through the [path.Clean()](https://pkg.go.dev/path#Clean) function before searching for a file. This removes any . and .. elements from the URL path, which helps to stop directory traversal attacks. This feature is particularly useful if you’re using the fileserver in conjunction with a router that doesn’t automatically sanitize URL paths.
- [Range requests](https://benramsey.com/blog/2008/05/206-partial-content-and-range-requests) are fully supported. This is great if your application is serving large files and you want to support resumable downloads. You can see this functionality in action if you use curl to request bytes 100-199 of the logo.png file, like so:
    ```
    $ curl -i -H "Range: bytes=100-199" --output - http://localhost:4000/static/img/logo.png
    HTTP/1.1 206 Partial Content
    Accept-Ranges: bytes
    Content-Length: 100
    Content-Range: bytes 100-199/1075
    Content-Type: image/png
    Last-Modified: Wed, 18 Mar 2024 11:29:23 GMT
    Date: Wed, 18 Mar 2024 11:29:23 GMT
    [binary data]
    ```
- The Last-Modified and If-Modified-Since headers are transparently supported. If a file hasn’t changed since the user last requested it, then http.FileServer will send a 304 Not Modified status code instead of the file itself. This helps reduce latency and processing overhead for both the client and server.
- The Content-Type is automatically set from the file extension using the [mime.TypeByExtension()](https://pkg.go.dev/mime#TypeByExtension) function. You can add your own custom extensions and content types using the [mime.AddExtensionType()](https://pkg.go.dev/mime#AddExtensionType) function if necessary.

---
### Performance
In this chapter we set up the file server so that it serves files out of the `./ui/static` directory on your hard disk.

But it’s important to note that http.FileServer probably won’t be reading these files from disk once the application is up-and-running. Both [Windows](https://learn.microsoft.com/en-us/windows/win32/fileio/file-caching) and [Unix-based](https://tldp.org/LDP/sag/html/buffer-cache.html) operating systems cache recently-used files in RAM, so (for frequently-served files at least) it’s likely that `http.FileServer` will be serving them from RAM rather than making the [relatively slow](https://gist.github.com/jboner/2841832) round-trip to your hard disk.

---
### Serving single file

Sometimes you might want to serve a single file from within a handler. For this there’s the [http.ServeFile()](https://pkg.go.dev/net/http#ServeFile) function, which you can use like so:

```go
func downloadHandler(w http.ResponseWriter, r“*http.Request) {
    http.ServeFile(w, r, "./ui/static/file.zip")
}
```

> Warning: http.ServeFile() does not automatically sanitize the file path. If you’re constructing a file path from untrusted user input, to avoid directory traversal attacks you must sanitize the input with filepath.Clean() before using it.

---
### Disabling directory listings

If you want to disable directory listings there are a few different approaches you can take.

The simplest way? Add a blank index.html file to the specific directory that you want to disable listings for. This will then be served instead of the directory listing, and the user will get a 200 OK response with no body. If you want to do this for all directories under `./ui/static` you can use the command:

```
$ find ./ui/static -type d -exec touch {}/index.html \;
```

A more complicated (but arguably better) solution is to create a custom implementation of [http.FileSystem](https://pkg.go.dev/net/http#FileSystem), and have it return an os.ErrNotExist error for any directories. A full explanation and sample code can be found in [this blog post](https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings).
