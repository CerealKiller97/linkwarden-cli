# 🔗 linkwarden-cli

A command-line interface for the [Linkwarden](https://linkwarden.app) platform. Manage links, collections, tokens, and migrations directly from your terminal.

## ✨ Prerequisites

- **Go 1.25+** (for building from source)
- A Linkwarden account and API token

## 🚀 Installation

### From source

```bash
git clone https://github.com/CerealKiller97/linkwarden-cli.git
cd linkwarden-cli
go build -o linkwarden-cli .
```

### ⚙️ Configuration

The CLI requires two environment variables:

| Variable           | Required | Default                        | Description                                                                            |
| ------------------ | -------- | ------------------------------ | -------------------------------------------------------------------------------------- |
| `LINKWARDEN_TOKEN` | Yes      | —                              | Your Linkwarden API token (generate via `tokens generate` or your Linkwarden settings) |
| `LINKWARDEN_URL`   | No       | `https://cloud.linkwarden.app` | Your Linkwarden instance URL (use for self-hosted instances)                           |

Configure before running any command:

```bash
export LINKWARDEN_TOKEN="your-api-token"
export LINKWARDEN_URL="https://your-instance.linkwarden.app"  # optional
```

## 📋 Commands

### `links`

Manage bookmarks and links.

#### `links create`

Bulk import links from a JSON file into a collection.

**Usage:**

```bash
linkwarden-cli links create --file <path-to-json>
```

**File format**

The JSON file must be an array of objects with `url` and `title` fields:

```json
[
  {
    "url": "https://example.com",
    "title": "Example Website"
  },
  {
    "url": "https://another-site.com",
    "title": "Another Site"
  }
]
```

**Flow**

1. Prompts you to select a target collection
2. Reads links from the specified file
3. Imports each link with a progress bar
4. Displays a summary table (total, success, failed)

---

### `tokens`

Manage API tokens for authentication.

#### `tokens generate`

Create a new API token interactively.

**Usage:**

```bash
linkwarden-cli tokens generate
```

**Flow**

1. Prompts for a token name
2. Prompts for expiration (7, 30, 60, 90 days, or no expiration)
3. Creates the token and prints it
4. On Linux/macOS: optionally adds the token to your shell config (`~/.zshrc`, `~/.bashrc`, etc.)

#### `tokens list`

List all available tokens.

**Usage:**

```bash
linkwarden-cli tokens list
```

#### `tokens revoke`

Revoke an existing token via an interactive prompt.

**Usage:**

```bash
linkwarden-cli tokens revoke
```

---

### `collections`

Manage collections (groups of links).

#### `collections create`

Create a new collection.

**Usage:**

```bash
linkwarden-cli collections create
```

---

### `migrations`

Import and export Linkwarden data.

#### `migrations export`

Export your Linkwarden data to a JSON backup file in the current directory.

**Usage:**

```bash
linkwarden-cli migrations export
```

The backup filename is derived from the API response. The file is written to your current working directory.

#### `migrations import`

Import data from external sources (planned).

**Usage:**

```bash
linkwarden-cli migrations import --format <format>
```

**Supported formats (planned):**

- `linkwarden` — Linkwarden backup format
- `htmlFile` — HTML bookmarks file
- `wallbag` — Wallabag export
- `omnivore` — Omnivore export

---

## 💡 Examples

**Bulk import links from a JSON file:**

```bash
export LINKWARDEN_TOKEN="your-token"
linkwarden-cli links create --file ./bookmarks.json
```

**Generate a token and optionally persist it:**

```bash
linkwarden-cli tokens generate
```

**Export your Linkwarden data:**

```bash
linkwarden-cli migrations export
```

---

## 📄 License

See repository for license details.
