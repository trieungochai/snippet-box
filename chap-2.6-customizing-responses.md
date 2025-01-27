By default, every response that your handlers send has the HTTP status code 200 OK (which indicates to the user that their request was received and processed successfully), plus 3 automatic system-generated headers: a `Date header`, and the `Content-Length` and  `Content-Type` of the response body.

```
$ curl -i localhost:4000/
HTTP/1.1 200 OK
Date: Wed, 18 Mar 2024 11:29:23 GMT
Content-Length: 21
Content-Type: text/plain; charset=utf-8

Hello from Snippetbox
```
