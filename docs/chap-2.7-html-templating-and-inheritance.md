Use Go’s [html/template](https://pkg.go.dev/html/template) package, which provides a family of functions for safely parsing and rendering HTML templates. We can use the functions in this package to parse the template file and then execute the template.

---
### Template composition
As we add more pages to our web application, there will be some shared, boilerplate, HTML markup that we want to include on every page — like the header, navigation and metadata inside the `<head>` HTML element.

To prevent duplication and save typing, it’s a good idea to create a `base` (or `master`) template which contains this shared content, which we can then compose with the page-specific markup for the individual pages.

---
### Embedding partials
For some applications you might want to break out certain bits of HTML into partials that can be reused in different pages or layouts.

