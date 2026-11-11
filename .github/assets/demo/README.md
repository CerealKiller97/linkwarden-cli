# Demo recordings

The GIFs in the main README are recorded with [VHS](https://github.com/charmbracelet/vhs) from the tapes here:

| Tape | GIF | Shows |
| --- | --- | --- |
| [`interactive.tape`](interactive.tape) | [`interactive.gif`](../interactive.gif) | Tables, pickers and the delete confirmation |
| [`import.tape`](import.tape) | [`import.gif`](../import.gif) | `links import` with a collection picker and progress bar |
| [`scripting.tape`](scripting.tape) | [`scripting.gif`](../scripting.gif) | `--json`, `-q` and failing fast without a terminal |
| [`../install.tape`](../install.tape) | [`install.gif`](../install.gif) | The install script (see the tape for its own setup) |

[`setup.tape`](setup.tape) holds the shared look, and [`bashrc`](bashrc) the hidden shell setup.

## Recording

1. Point `LINKWARDEN_URL` and `LINKWARDEN_TOKEN` at an **empty** throwaway account. The demos create and
   delete data, so don't use your real one. [`scripts/dev-linkwarden.sh`](../../../scripts/dev-linkwarden.sh)
   starts a local instance:

   ```bash
   scripts/dev-linkwarden.sh up
   eval "$(scripts/dev-linkwarden.sh env)"
   ```

2. Seed it, while it's still empty, so the IDs match the tapes:

   ```bash
   .github/assets/demo/seed.sh
   ```

3. Build a Linux binary into `bin/` (git-ignored), where the tapes look for it:

   ```bash
   CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/linkwarden-linux-amd64 ./cmd/linkwarden
   ```

4. Build a VHS image with `curl` and `jq`, and record from the repo root. Inside the container, a local
   instance is reached through `host.docker.internal` instead of `localhost`:

   ```bash
   printf 'FROM ghcr.io/charmbracelet/vhs\nRUN apt-get update && apt-get install -y --no-install-recommends curl jq\n' | docker build -t lw-vhs -
   for tape in interactive import scripting; do
     docker run --rm --shm-size=1g -v "$PWD:/vhs" \
       -e LINKWARDEN_URL="${LINKWARDEN_URL/localhost/host.docker.internal}" -e LINKWARDEN_TOKEN \
       lw-vhs ".github/assets/demo/$tape.tape"
   done
   ```

Each tape cleans up what it creates, so you can re-record without re-seeding.
