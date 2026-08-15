package utils

import (
	"strconv"

	"github.com/manifoldco/promptui"
)

// PickerItem is a single choice in an interactive ID picker.
type PickerItem struct {
	Label string
	Value int
}

// SelectID shows an interactive promptui select over items and returns the chosen
// item's Value.
func SelectID(label string, items []PickerItem) (int, error) {
	names := make([]string, len(items))
	for i, item := range items {
		names[i] = item.Label
	}

	prompt := promptui.Select{
		Label: label,
		Items: names,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	return items[idx].Value, nil
}

// ResolveID returns flagValue parsed as an int when non-empty; otherwise it calls
// fetch to list the available items and shows an interactive picker for label.
func ResolveID(flagValue string, label string, fetch func() ([]PickerItem, error)) (int, error) {
	if flagValue != "" {
		return strconv.Atoi(flagValue)
	}

	items, err := fetch()
	if err != nil {
		return 0, err
	}

	return SelectID(label, items)
}
