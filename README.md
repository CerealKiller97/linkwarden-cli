# 🔗 linkwarden-cli

A command-line interface and Go SDK for the [Linkwarden](https://linkwarden.app) platform. Manage links, collections, tags, tokens, RSS subscriptions, highlights, archives, and more — directly from your terminal.

## ✨ Prerequisites

- **Go 1.25+** (for building from source)
- A Linkwarden account and API token

## 🚀 Installation

### From source

```bash
git clone https://github.com/CerealKiller97/linkwarden-cli.git
cd linkwarden-cli
go build -o linkwarden .
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

## 🧰 Go SDK

Everything the CLI does goes through `pkg/linkwarden`, a standalone Go client for the Linkwarden API:

```go
client := linkwarden.New("https://your-instance.linkwarden.app", "your-api-token")

links, err := client.Links.List(ctx, linkwarden.ListLinksOptions{CollectionID: 3})
```

Every resource is exposed as a service on `*linkwarden.Client` (`Links`, `Collections`, `Tags`, `Tokens`, `RSS`, `Search`, `Dashboard`, `Highlights`, `Archives`, `Config`, `Favicon`, `Migrations`). Errors from non-2xx responses are returned as `*linkwarden.APIError` (see `linkwarden.IsNotFound`/`IsConflict`).

## 📋 Commands

Commands not listed with a full walkthrough below follow the same shape: flags map directly to the underlying API parameters, and any command that takes a resource ID shows an interactive picker (via `--flag ""`) when the ID isn't passed explicitly.

| Group        | Subcommands                                                                                     |
| ------------ | ------------------------------------------------------------------------------------------------ |
| `links`      | `create`, `list`, `get`, `update`, `delete`, `archive`, `bulk-update`, `bulk-delete`, `delete-archive`, `get-highlights` |
| `collections`| `create`, `list`, `get`, `update`, `delete`                                                     |
| `tags`       | `list`, `get`, `create`, `update`, `delete`, `bulk-delete`                                       |
| `tokens`     | `generate`, `list`, `revoke`                                                                     |
| `rss`        | `list`, `create`, `delete`                                                                       |
| `highlights` | `create`, `delete`                                                                                |
| `dashboard`  | `get`, `get-v2`, `update-v2`                                                                      |
| `archives`   | `get`, `upload`, `upload-standalone`                                                             |
| `migrations` | `export`, `import`                                                                                |
| `search`     | (single command)                                                                                  |
| `config`     | `get`                                                                                              |
| `favicon`    | `get`                                                                                              |

Run `linkwarden <command> --help` or `linkwarden <command> <subcommand> --help` for the full flag list of any command.

### `links`

Manage bookmarks and links.

#### `links create`

Bulk import links from a JSON file into a collection.

**Usage:**

```bash
linkwarden links create --file <path-to-json>
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

1. Prompts you to select a target collection (unless `--collection` is set)
2. Reads links from the specified file
3. Imports each link with a progress bar
4. Displays a summary table (total, success, failed)

#### Other `links` subcommands

```bash
linkwarden links list [--collection-id ID] [--tag-id ID] [--pinned-only] [--search QUERY] [--cursor N]
linkwarden links get [--id ID]
linkwarden links update [--id ID] [--name NAME] [--url URL] ...
linkwarden links delete [--id ID]
linkwarden links archive [--id ID]
linkwarden links delete-archive --ids 1,2,3
linkwarden links bulk-update --ids 1,2,3 [--collection-id ID] [--remove-previous-tags]
linkwarden links bulk-delete --ids 1,2,3
linkwarden links get-highlights [--id ID]
```

---

### `tokens`

Manage API tokens for authentication.

#### `tokens generate`

Create a new API token interactively.

**Usage:**

```bash
linkwarden tokens generate
```

**Flow**

1. Prompts for a token name
2. Prompts for expiration (7, 30, 60, 90 days, or no expiration)
3. Creates the token and prints it
4. On Linux/macOS: optionally adds the token to your shell config (`~/.zshenv`, `~/.bashrc`, `~/.config/fish/config.fish`)

#### `tokens list`

```bash
linkwarden tokens list
```

#### `tokens revoke`

```bash
linkwarden tokens revoke
```

---

### `collections`

Manage collections (groups of links).

```bash
linkwarden collections create [--name NAME] ...
linkwarden collections list
linkwarden collections get [--id ID]
linkwarden collections update [--id ID] ...
linkwarden collections delete [--id ID]
```

---

### `tags`

Manage tags.

```bash
linkwarden tags list [--search QUERY] [--cursor N]
linkwarden tags get [--id ID]
linkwarden tags create --label LABEL [--label LABEL ...]
linkwarden tags update [--id ID] --name NAME
linkwarden tags delete [--id ID]
linkwarden tags bulk-delete --ids 1,2,3
```

---

### `rss`

Manage RSS subscriptions that feed links into a collection.

```bash
linkwarden rss list
linkwarden rss create --name NAME --url FEED_URL [--collection-id ID | --collection-name NAME]
linkwarden rss delete [--id ID]
```

---

### `highlights`

Manage highlighted passages on a link. There's no list-all endpoint, so `delete` walks
you through a link picker, then a highlight picker, when `--id` isn't given.

```bash
linkwarden highlights create --link-id ID --text "..." --color yellow --start-offset 0 --end-offset 42
linkwarden highlights delete [--id ID] [--link-id ID]
```

---

### `dashboard`

View dashboard data and configure the v2 dashboard layout.

```bash
linkwarden dashboard get
linkwarden dashboard get-v2
linkwarden dashboard update-v2 --file layout.json
```

`layout.json` is an array of widgets, e.g.:

```json
[
  { "type": "STATS", "enabled": true, "order": 0 },
  { "type": "RECENT_LINKS", "enabled": true, "order": 1 },
  { "type": "COLLECTION", "collectionId": 3, "enabled": true, "order": 2 }
]
```

---

### `archives`

Retrieve or upload a link's archived formats (screenshot, PDF, readability text, monolith HTML).

```bash
linkwarden archives get --link-id ID --format 2 --output archive.pdf
linkwarden archives upload --link-id ID --file screenshot.png --format 0
linkwarden archives upload-standalone --file page.html --url https://example.com --format 4
```

Format values: `0`=PNG, `1`=JPEG, `2`=PDF, `3`=Readability JSON (get only), `4`=Monolith HTML.

---

### `migrations`

Import and export Linkwarden data.

#### `migrations export`

```bash
linkwarden migrations export [--output path.json]
```

Defaults to `./linkwarden-export-<timestamp>.json` in the current directory.

#### `migrations import`

```bash
linkwarden migrations import --file backup.json --format 0
```

Format values: `0`=Linkwarden JSON, `1`=Bookmarks HTML, `2`=Wallabag JSON, `3`=Omnivore JSON, `4`=Pocket JSON.

---

### `search`

```bash
linkwarden search --query "some text" [--sort N] [--cursor N]
```

---

### `config`

```bash
linkwarden config get
```

Prints the instance's public runtime configuration.

---

### `favicon`

```bash
linkwarden favicon get --url https://example.com [--output favicon.ico]
```

---

## 💡 Examples

**Bulk import links from a JSON file:**

```bash
export LINKWARDEN_TOKEN="your-token"
linkwarden links create --file ./bookmarks.json
```

**Generate a token and optionally persist it:**

```bash
linkwarden tokens generate
```

**Export your Linkwarden data:**

```bash
linkwarden migrations export
```

**Search and download a link's PDF archive:**

```bash
linkwarden search --query "kubernetes"
linkwarden archives get --link-id 42 --format 2 --output kubernetes.pdf
```

---

## 📄 License

See repository for license details.
