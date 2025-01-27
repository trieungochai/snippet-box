For this project, we’ll implement an outline structure that follows a [popular and tried-and-tested](https://go.dev/doc/modules/layout#server-project) approach. It’s a solid starting point, and you should be able to reuse the general structure in a wide variety of projects.

- `cmd` directory will contain the application-specific code for the executable applications in the project. For now our project will have just one executable application — the web application — which will live under the `cmd/web` directory.
- `internal` directory will contain the ancillary non-application-specific code used in the project. We’ll use it to hold potentially reusable code like validation helpers and the SQL database models for the project.
- `ui` directory will contain the user-interface assets used by the web application. Specifically, the `ui/html` directory will contain HTML templates, and the `ui/static` directory will contain static files (like CSS and images).

---
### why are we using this structure?
There are 2 big benefits:
1. It gives a clean separation between Go and non-Go assets. All the Go code we write will live exclusively under the `cmd` and `internal` directories, leaving the project root free to hold non-Go assets like UI files, makefiles and module definitions (including our go.mod file).
2. It scales really nicely if you want to add another executable application to your project. For example, you might want to add a CLI to automate some administrative tasks in the future. With this structure, you could create this CLI application under `cmd/cli` and it will be able to import and reuse all the code you’ve written under the internal directory.
