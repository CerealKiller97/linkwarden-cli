package tags

import (
	"strconv"
	"strings"
)

// parseIDs parses a comma-separated list of integer IDs, e.g. "1,2,3".
func parseIDs(csv string) ([]int, error) {
	parts := strings.Split(csv, ",")
	ids := make([]int, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		id, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}
