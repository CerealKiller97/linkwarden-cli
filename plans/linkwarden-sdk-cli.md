# Linkwarden Go SDK + Full CLI Command Coverage

## Context

The CLI (`linkwarden-cli`, Go/cobra) currently covers only a handful of Linkwarden API
endpoints (`links create`, `collections create`/list-for-picker, `tokens
generate`/`revoke`, `migrations export`), and each resource re-implements its own
`net/http` request/response boilerplate with inconsistent error handling
(`pkg/services/link`, `pkg/services/collection`, `pkg/services/token`,
`pkg/services/migration`, `pkg/services/tag`). Types are duplicated across
`pkg/dto/*` and each service package (e.g. `dto.Collection` vs.
`collection.Collection`). Several stubs exist but do nothing (`tokens list`,
`migrations import`), and `pkg/services/tag` has no CLI command wired to it at all.

The user wants the CLI to act as a proper **SDK-backed client for the full
Linkwarden API**, so every action available via the API is available via the CLI.
Decisions locked in with the user:

- **Scope**: Links, Collections, Tags, Tokens, RSS, Search, Dashboard, Highlights,
  Archives, Config, Favicon. Excluded (web-app/admin/login-only, not meaningful for
  a token-authed personal CLI): Users CRUD, Auth (session/password-reset/verify-email),
  Worker stats, Avatar, Preserved/Public.
- **Architecture**: one shared Go SDK client package replaces the duplicated
  per-resource `net/http` code, with centralized types and consistent error handling.
- **Delivery**: full design now, implemented in 3 reviewable stages.

Endpoint list and request/response shapes were pulled from Linkwarden's own OpenAPI
spec (`github.com/linkwarden/docs`, `openapi/linkwarden.yaml`), not guessed.

Two spec quirks worth knowing going in:
- `POST /api/v1/collections/{id}` is the documented *create* operation (confirmed via
  `operationId: createCollection`), even though creating shouldn't need an id. Current
  working code posts to `/api/v1/collections` (no id) and that's what actually works
  against a live instance today. Implement `CollectionsService.Create` to POST to
  `/api/v1/collections` (matching proven behavior) and only fall back to the
  `{id}`-suffixed path if a live smoke test shows the plain path stops working.
- `POST /api/v1/migration` (import) takes a **JSON** body `{format: int, data: string}`
  where `format` is `0`=Linkwarden JSON, `1`=Bookmarks HTML, `2`=Wallabag JSON,
  `3`=Omnivore JSON, `4`=Pocket JSON, and `data` is the raw import payload as a string
  — not a multipart file upload. `GET /api/v1/migration` (export) returns the full
  user/collections/links export object (confirmed shape, includes nested
  `collections[].links[].tags[]`).

## SDK client package: `pkg/linkwarden`

New package, one file per resource plus shared plumbing:

**`pkg/linkwarden/client.go`**
```go
type Client struct {
    baseURL    string
    token      string
    httpClient *http.Client
    logger     zerolog.Logger

    Links, Collections, Tags, Tokens, RSS, Search,
    Dashboard, Highlights, Archives, Config, Favicon, Migrations *XService
}

type Option func(*Client)
func WithHTTPClient(hc *http.Client) Option
func WithLogger(logger zerolog.Logger) Option
func New(baseURL, token string, opts ...Option) *Client
```
`New` trims a trailing `/` from `baseURL`, applies options, and wires each
`*XService{c: c}`.

Internal request helper used by every resource method:
```go
type requestOptions struct {
    method, path string
    query        url.Values
    jsonBody     any       // marshaled JSON body
    rawBody      io.Reader // multipart/binary bodies instead of jsonBody
    contentType  string    // overrides "application/json" when rawBody is set
}
func (c *Client) request(ctx context.Context, opt requestOptions) (*http.Response, []byte, error)
```
Sets `Authorization: Bearer <token>`, builds the query string, does the request,
fully reads + closes the body, and on non-2xx status decodes a best-effort message
into a typed `*APIError` (checks `{"message":...}`, falls back to `{"response":"..."}`,
then `http.StatusText`). This is the single place status-code/error handling happens —
no more hand-rolled `errors.New(resp.Status)` scattered per service.

Two generics collapse the two response envelope shapes seen across the API:
```go
func decodeResponseEnvelope[T any](body []byte) (T, error) // unwraps {"response": T}
func decodeDataEnvelope[T any](body []byte) (T, error)     // unwraps {"data": T, "message":..., "success":...}
```

**`pkg/linkwarden/errors.go`**
```go
type APIError struct{ StatusCode int; Message string; RawBody []byte }
func (e *APIError) Error() string
func IsNotFound(err error) bool
func IsConflict(err error) bool
```

**`pkg/linkwarden/types.go`** — only genuinely shared types (`Member`, `PinnedBy`) and
small query-building helpers (`setIfNonEmpty`, `setIfPositive`). Resource-specific
types (`Link`, `Collection`, `Tag`, `Token`, `RSSSubscription`, `Highlight`, ...) live
next to their service file, as the single source of truth (replacing `pkg/dto/*`).

Multipart bodies (archive upload / upload-standalone) share a helper in `client.go`:
`buildMultipartBody(fields map[string]string, fileField, filename string, r io.Reader) (io.Reader, string, error)`.
Binary GETs (`Archives.Get`, `Favicon.Get`) return `Data []byte`, `ContentType string`,
and (archives) `Filename string` parsed from `Content-Disposition`.

### Resource services

| File | Service | Methods |
|---|---|---|
| `links.go` | `LinksService` | `List(ctx, ListLinksOptions)`, `Get(ctx, id)`, `Create(ctx, CreateLinkRequest)`, `Update(ctx, id, UpdateLinkRequest)`, `Delete(ctx, id)`, `BulkUpdate(ctx, BulkUpdateLinksRequest)`, `BulkDelete(ctx, ids)`, `Archive(ctx, id)`, `DeleteArchive(ctx, ids)`, `Highlights(ctx, linkID)` |
| `collections.go` | `CollectionsService` | `List(ctx)`, `Get(ctx, id)`, `Create(ctx, CreateCollectionRequest)`, `Update(ctx, id, UpdateCollectionRequest)`, `Delete(ctx, id)` |
| `tags.go` | `TagsService` | `List(ctx, ListTagsOptions)`, `Get(ctx, id)`, `CreateOrUpdateBulk(ctx, []Tag)`, `Update(ctx, id, name)`, `Delete(ctx, id)`, `BulkDelete(ctx, ids)` |
| `tokens.go` | `TokensService` | `List(ctx)`, `Create(ctx, name, expiresDays)`, `Revoke(ctx, id)` |
| `rss.go` | `RSSService` | `List(ctx)`, `Create(ctx, CreateRSSRequest)`, `Delete(ctx, id)` |
| `search.go` | `SearchService` | `Links(ctx, SearchOptions)` |
| `dashboard.go` | `DashboardService` | `GetV1(ctx)`, `GetV2(ctx)`, `UpdateV2Layout(ctx, UpdateDashboardLayoutRequest)` |
| `highlights.go` | `HighlightsService` | `Create(ctx, CreateHighlightRequest)`, `Delete(ctx, id)` |
| `archives.go` | `ArchivesService` | `Get(ctx, linkID, format, preview)`, `Upload(ctx, linkID, format, preview, filename, r)`, `UploadStandalone(ctx, format, preview, url, filename, r)` |
| `config.go` | `ConfigService` | `Get(ctx)` |
| `favicon.go` | `FaviconService` | `Get(ctx, url)` |
| `migrations.go` | `MigrationsService` | `Export(ctx)`, `Import(ctx, format int, data string)` |

`DashboardData`/`DashboardV2Data`/`PublicConfig` should carry enough fields to be
useful but tolerate unknowns (e.g. embed `Raw json.RawMessage` or use `map[string]any`
for less-critical nested bits) since the plan doesn't have their full schemas
confirmed — verify against a live instance during Stage B/C smoke testing.

## Migration of existing code

Once each command is rewired and verified, delete: `pkg/services/link/`,
`pkg/services/collection/`, `pkg/services/token/`, `pkg/services/tag/`,
`pkg/services/migration/`, `pkg/dto/links.go`, `pkg/dto/collections.go` (and the
directories if left empty).

New helpers:
- **`pkg/utils/client.go`**: `func NewLinkwardenClient() *linkwarden.Client` — builds
  the client from `utils.GetInstanceURL()`/`GetLinkwardenToken()`, used by every
  command instead of constructing a per-resource service.
- **`pkg/utils/prompt.go`**: generalizes the picker pattern already in
  `cmd/commands/links/create.go:91-106` (promptui select when a flag is empty):
  ```go
  type PickerItem struct{ Label string; Value int }
  func SelectID(label string, items []PickerItem) (int, error)
  func ResolveID(flagValue, label string, fetch func() ([]PickerItem, error)) (int, error)
  ```

Rewiring (behavior preserved, output formatting via `go-pretty/table` unchanged):
- `cmd/commands/links/create.go`: `collection.New`/`link.New` → `utils.NewLinkwardenClient()`; picker via `utils.ResolveID` over `client.Collections.List`.
- `cmd/commands/collections/create.go`: same swap; `Get`→`client.Collections.List`, `Create`→`client.Collections.Create`.
- `cmd/commands/tokens/{generate,revoke,list}.go`: swap to the shared client; `list.go`'s empty `executeList` gets a real body (fetch + table render).
- `cmd/commands/migrations/migrations.go` (Stage C): rewire `Export`; implement `Import` with new `--file` (path to raw import payload) and existing `--format` flags, reading the file into a string for `client.Migrations.Import`.

## Full CLI command tree

Same cobra pattern throughout: `cmd/commands/<group>/<group>.go` exposes
`Command() *cobra.Command`, one file per subcommand, each building its own client via
`utils.NewLinkwardenClient()`.

- **`links`**: keep `create`; add `list`, `get`, `update`, `delete`, `archive`,
  `bulk-update`, `bulk-delete`, `delete-archive`, `get-highlights`. ID-taking
  subcommands accept `--id`/positional or fall back to `utils.ResolveID` over
  `client.Links.List` titles.
- **`collections`**: keep `create`; add `list`, `get`, `update`, `delete`.
- **`tags`** (new package): `list`, `get`, `create`, `update`, `delete`, `bulk-delete`.
  Registered in `cmd/root.go`.
- **`tokens`**: `generate` (unchanged), `list` (filled in), `revoke` (unchanged).
- **`rss`** (new package): `list`, `create`, `delete`.
- **`search`** (new package, single leaf command): `search --query ... --sort ... --cursor ...`. No picker (free-text).
- **`dashboard`** (new package): `get` (v1), `get-v2`, `update-v2` (`--file` JSON layout). No picker (singleton resource).
- **`highlights`** (new package): `create` (`--link-id`, `--text`, `--start-offset`, `--end-offset`, `--color`), `delete` (two-step picker: pick link, then pick highlight, since there's no list-all-highlights endpoint).
- **`archives`** (new package): `get` (`--link-id`, `--output <path>` to save; without it, prints content-type/size), `upload` (`--link-id`, `--file`, `--format`, `--preview`), `upload-standalone` (`--file`, `--url`, `--format`, `--preview`).
- **`config`** (new package): `get`. No picker.
- **`favicon`** (new package): `get` (`--url` required, `--output` optional path). No picker.
- **`migrations`**: `export` (rewired), `import` (filled in).

Picker does **not** apply to: `bulk-update`/`bulk-delete` (use `--ids` comma-separated
or repeated args — promptui has no multi-select), `search`, `config get`,
`favicon get`, `dashboard *`, `archives upload-standalone`.

## Testing requirement (every endpoint, every command)

Per the user's explicit instruction, **every SDK method (one per API endpoint) and
every CLI command gets a test**, added in the same stage that introduces it — not
deferred to a later cleanup pass.

- **SDK layer** (`pkg/linkwarden/*_test.go`): one `httptest.Server`-backed test per
  resource method, asserting the outgoing request (method, path, query params,
  JSON/multipart body) and the decoded result for a success response, plus at least
  one non-2xx case asserting `*linkwarden.APIError` is returned with the right status
  code. This covers all ~35 endpoints across the 11 in-scope resources.
- **CLI layer** (`cmd/commands/<group>/*_test.go`): one test per subcommand file.
  Since commands build their client via `utils.NewLinkwardenClient()` (reading
  `LINKWARDEN_URL`/`LINKWARDEN_TOKEN`), command tests point that env at an
  `httptest.Server` (set via `t.Setenv`) and drive the `*cobra.Command` directly
  (`cmd.SetArgs(...)`, `cmd.Execute()`), asserting on flags parsed correctly, the
  request the fake server received, and rendered output (table content / error
  message) for both success and API-error paths. Commands with an interactive
  promptui picker get a test for the flag-provided path (no prompt) plus a test
  proving the picker triggers when the flag is omitted (`utils.ResolveID`'s `fetch`
  callback is what actually needs coverage — test it directly in
  `pkg/utils/prompt_test.go` rather than trying to drive promptui's TTY prompt from a
  command test).
- Existing rewired commands (`links create`, `collections create`, `tokens
  generate/revoke`) get their current behavior covered too, since they're moving onto
  the new client and currently have zero tests.

This adds `pkg/linkwarden/client_test.go` (shared `httptest` setup helpers reused by
every resource test file) and `pkg/utils/{client,prompt}_test.go` alongside the
production files, plus a `*_test.go` next to every `cmd/commands/**/*.go` command
file listed below.

## Staged delivery

**Stage A — SDK core + Links/Collections/Tags/Tokens**
New: `pkg/linkwarden/{client,errors,types,links,collections,tags,tokens}.go` +
matching `*_test.go` for each (`httptest.Server`-backed, per the testing requirement
above); `pkg/utils/{client,prompt}.go` + `*_test.go`; new command files under
`cmd/commands/links/`, `cmd/commands/collections/` each with a sibling `_test.go`;
new package `cmd/commands/tags/` (commands + tests).
Changed: `cmd/commands/links/create.go`, `cmd/commands/collections/create.go`,
`cmd/commands/tokens/{generate,list,revoke}.go`, `cmd/root.go` — each gets a new or
updated `_test.go` since none exist today.
Removed (after verified rewiring): `pkg/services/{link,collection,token,tag}/`,
`pkg/dto/{links,collections}.go`.

**Stage B — Search/RSS/Highlights/Dashboard**
New: `pkg/linkwarden/{search,rss,highlights,dashboard}.go` + tests; new packages
`cmd/commands/{search,rss,highlights,dashboard}/` (commands + tests).
Changed: `cmd/commands/links/` gains `get-highlights` + test (depends on `Highlight`
type from this stage); `cmd/root.go`.

**Stage C — Archives/Config/Favicon + migrations import + cleanup**
New: `pkg/linkwarden/{archives,config,favicon,migrations}.go` + tests (archive/favicon
tests cover binary response handling and `--output` file-writing); new packages
`cmd/commands/{archives,config,favicon}/` (commands + tests).
Changed: `cmd/commands/migrations/migrations.go` (implement import) + test covering
both export and the new import flag combination.
Removed: `pkg/services/migration/`, now-empty `pkg/dto/`/`pkg/services/` dirs.
Also: update `README.md` command tables to document the full command set.

## Verification

- `go build ./...` and `go vet ./...` after each stage.
- `go test ./...` — every new SDK method and every new/changed CLI command has a
  test per the requirement above; CI-style full run at the end of each stage.
- Manual smoke test against a real/self-hosted instance
  (`export LINKWARDEN_TOKEN=... LINKWARDEN_URL=https://your-instance`), exercising
  each new subcommand at least once per stage, confirming table output, the picker
  flow when flags are omitted, and that non-2xx responses (e.g. `links get <bad-id>`)
  surface a clean `*linkwarden.APIError` rather than a raw error/panic.
- Stage C binary endpoints: verify `archives get --output` produces a valid file for
  the requested format, and `favicon get --url` produces a valid image file.

### Critical files
- [pkg/linkwarden/client.go](pkg/linkwarden/client.go) (new)
- [pkg/linkwarden/errors.go](pkg/linkwarden/errors.go) (new)
- [pkg/utils/client.go](pkg/utils/client.go) (new)
- [pkg/utils/prompt.go](pkg/utils/prompt.go) (new)
- [cmd/root.go](cmd/root.go)
- [cmd/commands/links/create.go](cmd/commands/links/create.go)
