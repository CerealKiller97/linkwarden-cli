# Records interactive.gif: tables, pickers and a delete confirmation.
#
# See README.md in this folder for how to record.

Output .github/assets/interactive.gif
Source .github/assets/demo/setup.tape

Sleep 500ms
Type "# readable tables in a terminal, plain text when piped"
Enter
Type "linkwarden collections list"
Sleep 400ms
Enter
Sleep 2.5s

Type "# no --id? pick from a list: type to filter, Enter to select"
Enter
Type "linkwarden links get"
Sleep 400ms
Enter
Sleep 1.2s
Down@300ms 2
Sleep 500ms
Type@120ms "pipe"
Sleep 1s
Enter
Sleep 3s

Hide
Type "clear"
Enter
Sleep 200ms
Show

Type "# deletes ask first, with No selected"
Enter
Type "linkwarden links delete"
Sleep 400ms
Enter
Sleep 1.2s
Type@120ms "broken"
Sleep 800ms
Enter
Sleep 1.5s
Right
Sleep 600ms
Left
Sleep 800ms
Enter
Sleep 3s
