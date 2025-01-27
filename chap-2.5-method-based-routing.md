Let’s follow HTTP good practice and restrict our application so that it only responds to requests with an appropriate HTTP method.

As we progress through our application build, our `home`, `snippetView` and `snippetCreate` handlers will merely retrieve information and display pages to the user, so it makes sense that these handlers should be restricted to acting on GET requests.
