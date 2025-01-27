Let’s follow HTTP good practice and restrict our application so that it only responds to requests with an appropriate HTTP method.

As we progress through our application build, our `home`, `snippetView` and `snippetCreate` handlers will merely retrieve information and display pages to the user, so it makes sense that these handlers should be restricted to acting on GET requests.

---
### Method precedence
It’s important to be aware that a route pattern which doesn’t include a method — like `"/article/{id}"` — will match incoming HTTP requests with any method. In contrast, a route like `"POST /article/{id}"` will only match requests which have the method POST. So if you declare the overlapping routes `"/article/{id}"` and `"POST /article/{id}"` in your application, then the `"POST /article/{id}"` route will take precedence.