# Meal planner

## Description

Current goal - A web app that helps with meal planning with the following features:

- Recipe database: users can save their favorite recipes from different sources on the internet, and tag them appropriately (cuisines, meal types, etc)
- Meal plan generation:
  - users can specify their preferences for the week: how many meals with each cuisines
  - the app randomly generates the meal plan based on the preferences
  - users can update the generated meal plan
- Groceries list: groceries for the week will be generated based on the ingredients required in the meal plan.

Built with Go, Templ, HTMX, Tailwind CSS.

## Installation

### Prerequisites

- [Go 1.22 and above](https://go.dev/doc/install)
- [Postgresql](https://www.postgresql.org/download/)
- [Goose](https://github.com/pressly/goose)
- [sqlc](https://github.com/sqlc-dev/sqlc)
- [templ](https://templ.guide/quick-start/installation)

### Setup

Clone and cd into the project to install dependencies.

```sh
git clone https://github.com/quangd42/meal-planner.git
cd meal-planner
go mod tidy
```

Create a .env file. Here's an example:

```sh
PORT=8080
DATABASE_URL=postgres://[db-user]:@localhost:5432/[db-name]?sslmode=disable
JWT_SECRET=IpRoF6GpEewWJcHr8QqI5g4nj6RkKvaVYMNMJFa6svOgDbLWwyDg1jDictjfBIzY
```

You can generate JWT_SECRET with a command like this:

```sh
openssl rand -base64 64
```

### Local development

To start the server locally:

```sh
# setup the database
make db/reset
# generate code and build server binary
make build
# run the binary
make run
```

`make help` for more details.
