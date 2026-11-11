# Records collections.gif: create, nest, update through a picker, list and
# delete with confirmation.
#
# Run it through record.sh: it expects freshly seeded data, so the first new
# collection gets ID 5. It deletes what it creates.

Output .github/assets/collections.gif
Source .github/assets/demo/setup.tape
Set Height 760

Sleep 500ms
Type "# create a collection"
Enter
Type `linkwarden collections create --name "Rust" --description "The book, std and crates" --color "#f74c00"`
Sleep 400ms
Enter
Sleep 1.5s

Type "# nest another one inside it"
Enter
Type `linkwarden collections create --name "Async" --parent-id 5 --color "#fb923c"`
Sleep 400ms
Enter
Sleep 1.5s

Type "# no --id? pick the collection to change"
Enter
Type `linkwarden collections update --name "Async Rust" --description "Tokio, futures and friends"`
Sleep 400ms
Enter
Sleep 1.2s
Type@120ms "async"
Sleep 900ms
Enter
Sleep 1.5s

Type "linkwarden collections list"
Sleep 400ms
Enter
Sleep 3s

Hide
Type "clear"
Enter
Sleep 200ms
Show

Type "# deleting asks first, and says what goes with it"
Enter
Type "linkwarden collections delete --id 5"
Sleep 400ms
Enter
Sleep 2s
Right
Sleep 800ms
Enter
Sleep 1.5s
Type "linkwarden collections list"
Sleep 400ms
Enter
Sleep 3.5s
