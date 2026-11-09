# Records install.gif, the installer demo in the README, with VHS
# (https://github.com/charmbracelet/vhs).
#
# The hidden setup serves the release from a local mirror instead of GitHub
# (LINKWARDEN_RELEASES_URL) and makes `curl …/install.sh` print the local
# install.sh, so the demo works before a release is public. The output is
# what the real one-liner shows.
#
# Record from the repo root, with a mirror of a release on the host at :8765
# (download/vX.Y.Z/* plus a /latest → /tag/vX.Y.Z redirect) and an image
# that adds curl to ghcr.io/charmbracelet/vhs:
#
#   docker run --rm --shm-size=1g -v "$PWD:/vhs" lw-vhs .github/assets/install.tape

Output .github/assets/install.gif

Set Shell "bash"
Set FontSize 14
Set Width 1100
Set Height 1120
Set Padding 24
Set Theme "Catppuccin Mocha"
Set WindowBar Colorful
Set TypingSpeed 35ms

Hide
Type `export LINKWARDEN_RELEASES_URL=http://host.docker.internal:8765 COLORTERM=truecolor SHELL=/bin/bash PS1='$ '`
Enter
Type `curl() { if [[ "$*" == *install.sh ]]; then cat /vhs/install.sh; else command curl "$@"; fi; }`
Enter
Type "clear"
Enter
Sleep 500ms
Show

Sleep 500ms
Type "curl -fsSL https://raw.githubusercontent.com/CerealKiller97/linkwarden-cli/main/install.sh | bash"
Sleep 800ms
Enter

Wait@30s /Pick one/
Sleep 1.5s
Type "1"
Sleep 500ms
Enter

Wait@30s /tab completions/
Sleep 1s
Type "y"
Sleep 300ms
Enter

Wait@15s /Add it to/
Sleep 1s
Type "y"
Sleep 300ms
Enter

Wait+Screen@30s /Happy bookmarking/
Sleep 4s
