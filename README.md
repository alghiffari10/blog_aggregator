# gator

A CLI RSS feed aggregator written in Go. Subscribe to RSS feeds, collect posts, and browse them from your terminal.

## Prerequisites

- [PostgreSQL](https://www.postgresql.org/) (running)
- [Go](https://go.dev/) 1.26+

## Installation

```bash
go install github.com/alghiffari10/blog_aggregator@latest
```

This compiles a static binary to `$GOPATH/bin`. Make sure that directory is on your `$PATH`.

## Quick Start

### 1. Create the database

```bash
createdb gator
```

### 2. Run migrations

Use [goose](https://github.com/pressly/goose) to apply the schema:

```bash
goose -dir sql/schema postgres "postgres://youruser:@localhost:5432/gator?sslmode=disable" up
```

### 3. Configure gator

Create `~/.gatorconfig.json`:

```json
{
  "db_url": "postgres://youruser:@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

## Commands

| Command     | Usage                              | Auth required | Description                             |
|-------------|------------------------------------|:-------------:|-----------------------------------------|
| `register`  | `gator register <name>`            | No            | Create a new user                       |
| `login`     | `gator login <name>`               | No            | Switch to an existing user              |
| `users`     | `gator users`                      | No            | List all registered users               |
| `reset`     | `gator reset`                      | No            | Delete all users                        |
| `addfeed`   | `gator addfeed <name> <url>`       | Yes           | Add an RSS feed (automatically follows) |
| `feeds`     | `gator feeds`                      | No            | List all feeds in the system            |
| `follow`    | `gator follow <url>`               | Yes           | Follow an existing feed                 |
| `following` | `gator following`                  | Yes           | List feeds you follow                   |
| `unfollow`  | `gator unfollow <url>`             | Yes           | Unfollow a feed                         |
| `agg`       | `gator agg <duration>`             | No            | Start the RSS aggregation loop          |
| `browse`    | `gator browse [limit]`             | Yes           | Browse recent posts (default: 2)        |

## Example workflow

```bash
# Register a user
blog_aggregator register alice

# Add a feed (you automatically follow it)
blog_aggregator addfeed HackerNews https://news.ycombinator.com/rss

# Start the aggregator in a separate terminal
blog_aggregator agg 60s

# Follow another feed
blog_aggregator follow https://blog.golang.org/feed.atom

# See what you follow
blog_aggregator following

# Browse the latest posts
blog_aggregator browse 5
```
