# Gator 🐊

> A multi-user CLI RSS feed aggregator written in Go, backed by PostgreSQL.

Register users, follow RSS feeds, scrape them continuously with a background worker, and read the latest posts without leaving your terminal.

<p>
  <img alt="Go" src="https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white">
  <img alt="PostgreSQL" src="https://img.shields.io/badge/PostgreSQL-14+-4169E1?logo=postgresql&logoColor=white">
</p>

---

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Database Setup](#database-setup)
- [Quick Start](#quick-start)
- [Command Reference](#command-reference)
- [How the Aggregator Works](#how-the-aggregator-works)
- [Project Structure](#project-structure)

---

## Features

- **Multi-user** — each user keeps their own set of followed feeds.
- **Persistent sessions** — the current user is stored in a local config file.
- **Background scraping** — a long-running worker fetches feeds on a fixed interval.
- **Deduplicated posts** — the same article is never stored twice.
- **Terminal reader** — browse the most recent posts from the feeds you follow.

---

## Prerequisites

| Requirement | Version / Notes |
| --- | --- |
| [Go](https://go.dev/dl/) | 1.26 or higher |
| [PostgreSQL](https://www.postgresql.org/download/) | A running instance you can connect to |
| [goose](https://github.com/pressly/goose) | Used to run the SQL migrations |
| WSL 2 | Recommended if you are on Windows |

Install `goose` if you don't have it yet:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

---

## Installation

Install the binary with `go install`:

```bash
go install github.com/dariodesalvo/gator@latest
```

Make sure `$GOPATH/bin` (usually `~/go/bin`) is in your `PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Verify the installation:

```bash
gator users
```

---

## Configuration

Gator reads its configuration from `~/.gatorconfig.json`. Create it with your connection string:

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"
}
```

| Field | Description |
| --- | --- |
| `db_url` | PostgreSQL connection string. |
| `current_user_name` | Managed by Gator — set automatically by `register` and `login`. |

> **Note:** don't commit this file. It contains your database credentials.

---

## Database Setup

Create the database:

```bash
createdb gator
```

Then run the migrations from the schema directory:

```bash
cd sql/schema
goose postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" up
```

---

## Quick Start

```bash
# 1. Create a user (also logs you in)
gator register alice

# 2. Add a feed — this follows it automatically
gator addfeed "Boot.dev Blog" https://blog.boot.dev/index.xml

# 3. Start the scraper in a separate terminal (Ctrl+C to stop)
gator agg 1m

# 4. Read what came in
gator browse 10
```

---

## Command Reference

### User Management

| Command | Description |
| --- | --- |
| `gator register <username>` | Register a new user and set them as the current user. |
| `gator login <username>` | Log in as an existing user. |
| `gator users` | List all users and mark the one currently logged in. |

### Feed Management

| Command | Description |
| --- | --- |
| `gator addfeed <name> <url>` | Add a new RSS feed and follow it automatically. |
| `gator feeds` | List every registered feed along with the user who created it. |
| `gator follow <url>` | Follow a feed that already exists. |
| `gator following` | List the feeds the current user follows. |
| `gator unfollow <url>` | Stop following a feed. |

### Scraper & Reader

| Command | Description |
| --- | --- |
| `gator agg <time_between_reqs>` | Run the background scraping service. Accepts Go durations: `10s`, `1m`, `1h`. |
| `gator browse [limit]` | Show the most recent posts from followed feeds. Defaults to `2`. |

### Maintenance

| Command | Description |
| --- | --- |
| `gator reset` | Delete all data from the database. |

> ⚠️ `gator reset` is destructive and irreversible. It exists for local development only.

---

## How the Aggregator Works

`gator agg` runs in a loop until you stop it with `Ctrl+C`:

1. Pick the feed that was fetched least recently.
2. Mark it as fetched.
3. Download and parse the RSS XML.
4. Save any new posts, skipping ones already stored.
5. Wait for the configured interval and repeat.

Because the worker always picks the oldest feed, every feed gets refreshed fairly without any extra scheduling logic. Keep the interval reasonable (`1m` or more) so you don't hammer the servers you're scraping.

---

## Project Structure

```
gator/
├── main.go              # Entry point and command dispatch
├── internal/
│   ├── config/          # Reads and writes ~/.gatorconfig.json
│   └── database/        # Generated query code
└── sql/
    ├── schema/          # goose migrations
    └── queries/         # SQL used to generate the database package
```

---

## Acknowledgements

Built as the capstone project for the [Boot.dev](https://boot.dev) Backend Developer Path.