# Meal ORG

A food recipes database, with features for meal planning, shopping list.

## Motivations

- Keep hand written recipes as well as ones from the internet (allrecipes, youtube, instagram, etc.) in one centralized location.
- Plan meals weekly or daily and groceries shopping with these recipes through search and tags.
- Share and collaborate with family and friends.

## 🚀 Quick Start

### Prerequisites

- [Go 1.22 and above](https://go.dev/doc/install)
- [Postgresql](https://www.postgresql.org/download/)

### Setup as web server

Clone and cd into the project to install dependencies.

```sh
git clone https://github.com/quangd42/meal-org.git
cd meal-org
go mod tidy
```

Create a .env file. Here's an example:

```sh
PORT=8080

# Generate a random string for jwt-secret
JWT_SECRET=[jwt-secret]

# Replace [db-user] and [db-name] with your local postgres info
# Note the "" and ${} on DATABASE_URL
DB_USER=[db-user]
DB_NAME=[db-name]
DATABASE_URL="postgres://${DB_USER}:@localhost:5432/${DB_NAME}?sslmode=disable"
```

You can generate your own JWT_SECRET with a command like this:

```sh
openssl rand -base64 64
```

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

### APIs

Exposed APIs can be found in the [tests](tests/integration) in form of [hurl files](https://hurl.dev/docs/hurl-file.html).

## 🛠️ Local development

### Live reloading

After setting up the local environment per [quick start](#-quick-start), to develop locally with live reloading on the server:

```sh
make live/server
```

When extending the UI with Templ, additionally you can get browser live reloading by running:

```sh
# See details of what this does in the Makefile.
# make live/server in included in make live
make live
```

- [Makefile for your Go project](https://www.alexedwards.net/blog/a-time-saving-makefile-for-your-go-projects)
- [Templ live reload](https://templ.guide/commands-and-tools/live-reload-with-other-tools)

### Run the tests

I'm currently testing the APIs with [hurl](https://hurl.dev/docs/installation.html). Make sure to have it installed.
Test setup and tear down are done in [test_integration.sh](scripts/test_integration.sh)

To start, create a new .env in the [scripts](scripts) dir:

```sh
# Replace [db-user] with your postgres user
DB_USER=[db-user]
```

To run the test script:

```sh
make test
```

## 🤝 Contributing

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.
