In this section, we’re going to focus on adding an HTML form for creating new snippets. The form will look a bit like this:
![snippetbox-create](snippetbox-create.png)

The high-level flow for processing this form will follow a standard `Post-Redirect-Get` pattern and will work like so:
- The user is shown the blank form when they make a `GET` request to `/snippet/create`.
- The user completes the form and it’s submitted to the server via a `POST` request to `/snippet/create`.
- The form data will be validated by our `snippetCreatePost` handler. If there are any validation failures the form will be re-displayed with the appropriate form fields highlighted. If it passes our validation checks, the data for the new snippet will be added to the database and then we’ll redirect the user to `GET` `/snippet/view/{id}`.
