# Records import.gif: a bulk import into a collection picked from a list.
#
# See README.md in this folder for how to record.

Output .github/assets/import.gif
Source .github/assets/demo/setup.tape

# A fresh target collection, removed again at the end so the tape can be re-run.
Hide
Type `WEEKEND=$(linkwarden collections create --name "Weekend Reads" --description "Essays to read slowly" --color "#22c55e" -q)`
Enter
Sleep 1s
Type "clear"
Enter
Sleep 200ms
Show

Sleep 500ms
Type "# bulk import from JSON; duplicates are skipped"
Enter
Type "jq -r '.[].title' reading-list.json"
Enter
Sleep 2s

Type "linkwarden links import --file reading-list.json"
Sleep 400ms
Enter
Sleep 1.2s
Type@120ms "week"
Sleep 900ms
Enter
Wait+Screen@60s /Failed/
Sleep 4s

Hide
Sleep 500ms
Type `linkwarden collections delete --id "$WEEKEND" --yes`
Enter
Sleep 2s
