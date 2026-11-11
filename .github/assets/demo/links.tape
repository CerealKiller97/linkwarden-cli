# Records links.gif: save, list, search, bulk-retag and rename links.
#
# Run it through record.sh: it expects freshly seeded data (see seed.sh). It
# puts the links it changes back the way they were and deletes the one it adds.

Output .github/assets/links.gif
Source .github/assets/demo/setup.tape
Set FontSize 15
Set Width 1500
Set Height 820

Sleep 500ms
Type "# save a link; its name comes from the page's title"
Enter
Type "linkwarden links create https://doc.rust-lang.org/book/ --collection-id 2 --tag rust --tag books"
Sleep 400ms
Enter
Sleep 2s
Type "linkwarden links list --collection-id 2"
Sleep 400ms
Enter
Sleep 3s

Hide
Type "clear"
Enter
Sleep 200ms
Show

Type "# full-text search across every collection"
Enter
Type "linkwarden search --query kubernetes"
Sleep 400ms
Enter
Sleep 2.5s

Type "# tag (or move) many links in one call"
Enter
Type "linkwarden links bulk-update --ids 6,7 --tag weekend"
Sleep 400ms
Enter
Sleep 1.5s
Type "linkwarden links list --collection-id 3"
Sleep 400ms
Enter
Sleep 3s

Hide
Type "clear"
Enter
Sleep 200ms
Show

Type "# no --id? pick the link to rename"
Enter
Type `linkwarden links update --name "Kubernetes: the overview"`
Sleep 400ms
Enter
Sleep 1.2s
Type@120ms "overview"
Sleep 900ms
Enter
Sleep 1.5s
Type "linkwarden links get --id 7"
Sleep 400ms
Enter
Sleep 3.5s

# Put the seed data back.
Hide
Sleep 500ms
Type `linkwarden links list --collection-id 2 --json | jq '.[] | select(.url | contains("rust-lang")) | .id' | xargs -n1 linkwarden-demo links delete --yes --id`
Enter
Type "linkwarden links bulk-update --ids 6 --tag k8s --tag books --remove-previous-tags"
Enter
Type "linkwarden links bulk-update --ids 7 --tag k8s --tag docs --remove-previous-tags"
Enter
Type `linkwarden links update --id 7 --name "Kubernetes Overview"`
Enter
Sleep 3s
