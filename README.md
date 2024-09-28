---

# Gator

Gator is a command-line application (CLI) built in Go that allows users to manage RSS feed subscriptions and aggregate posts from various feeds. It uses PostgreSQL as the backend database to store users, feeds, feed follows, and posts.

## Features

- User registration and login.
- Ability to follow and unfollow RSS feeds.
- Automatic fetching of RSS feeds and storing of posts in the database.
- Browse RSS posts in the terminal.
- Feeds are aggregated continuously in a long-running process.
- Middleware to handle logged-in users for specific commands.

## Prerequisites

Before running Gator, you'll need to have the following installed on your system:

1. **Go**: You can install it via [this link](https://golang.org/doc/install).
2. **PostgreSQL**: Install using your system's package manager (e.g. Homebrew for macOS or apt for Linux).

### Installing PostgreSQL

On macOS, you can install PostgreSQL using Homebrew:
```bash
brew install postgresql
brew services start postgresql
```

## Installation

To install the Gator CLI, use the `go install` command:

```bash
go install github.com/Romasav/gator@latest
```

## Database Setup

Gator requires a PostgreSQL database to function. You can set it up by creating a database:

```bash
createdb gator
```

After creating the database, you can apply the migrations using **Goose**:

```bash
goose postgres "postgres://username:password@localhost:5432/gator?sslmode=disable" up
```

Make sure to update the database connection string to match your PostgreSQL setup.

## Config File

The config file, `.gatorconfig.json`, needs to be located in your project’s root directory. This file holds important configuration, including the currently logged-in user and database connection URL.

### Example `.gatorconfig.json`

```json
{
  "current_user_name": "your-username",
  "db_url": "postgres://your-username:your-password@localhost:5432/gator?sslmode=disable"
}
```

### Setting Up the Config

To initialize the configuration, create the `.gatorconfig.json` file in your project root directory with your PostgreSQL details.

## Running Gator

You can run Gator by using the following command:

```bash
./gator <command> [arguments...]
```

### Example Commands

- **Login**: Login as an existing user.
  
  ```bash
  ./gator login <username>
  ```

- **Register**: Create a new user.

  ```bash
  ./gator register <username>
  ```

- **Reset**: Delete all users from the database.

  ```bash
  ./gator reset
  ```

- **Add Feed**: Add a new feed to follow.

  ```bash
  ./gator addfeed <feed-name> <feed-url>
  ```

- **Browse**: Browse posts from the feeds you follow.

  ```bash
  ./gator browse [limit]
  ```

- **Aggregator**: Continuously fetch new posts from all followed feeds.

  ```bash
  ./gator agg <time_between_reqs>
  ```

- **Follow a Feed**: Follow an existing feed.

  ```bash
  ./gator follow <feed-url>
  ```

- **Unfollow a Feed**: Unfollow a feed.

  ```bash
  ./gator unfollow <feed-url>
  ```

- **Show Followed Feeds**: Display all feeds the user is following.

  ```bash
  ./gator following
  ```

## Development

### SQL Migrations

Migrations are managed using **Goose**. Migrations can be found in the `sql/schema` directory.

To apply migrations, run:

```bash
goose postgres "postgres://username:password@localhost:5432/gator?sslmode=disable" up
```

### SQL Queries with SQLC

SQLC is used to generate type-safe Go code from SQL queries. Queries are stored in the `sql/queries` directory. To regenerate Go code from SQL queries, run:

```bash
sqlc generate
```

### Install the Required Go Packages

For development, you will need to install **Goose** and **SQLC** for managing database migrations and generating Go code from SQL queries:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
```

--- 
