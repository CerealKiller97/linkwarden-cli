# Records scripting.gif: --json, -q and fail-fast behavior for scripts.
#
# See README.md in this folder for how to record.

Output .github/assets/scripting.gif
Source .github/assets/demo/setup.tape
Set Height 800

Sleep 500ms
Type "# --json prints the API's own fields, ready for jq"
Enter
Type `linkwarden links list --collection-id 1 --json | jq -r '.[] | "\(.id)  \(.url)"'`
Sleep 400ms
Enter
Sleep 2.5s

Type "# -q prints bare IDs, so results feed the next command"
Enter
Type `COL=$(linkwarden collections create --name "Rust" --color "#f74c00" -q)`
Sleep 300ms
Enter
Sleep 1s
Type `linkwarden links create https://doc.rust-lang.org/book/ --collection-id "$COL" --tag rust -q`
Sleep 300ms
Enter
Sleep 2.5s

Type "# no terminal? it fails fast instead of waiting on a prompt"
Enter
Type `linkwarden collections get < /dev/null; echo "exit=$?"`
Sleep 300ms
Enter
Sleep 2.5s
Type `linkwarden --non-interactive links delete --id 2; echo "exit=$?"`
Sleep 300ms
Enter
Sleep 4s

Hide
Type `linkwarden collections delete --id "$COL" --yes`
Enter
Sleep 2s
