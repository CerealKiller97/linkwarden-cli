package dto

type (
	Tag struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	}

	Collection struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	CreateLinkRequest struct {
		Name        string     `json:"name"`
		URL         string     `json:"url"`
		Type        string     `json:"type"`
		Description string     `json:"description"`
		Tags        []Tag      `json:"tags"`
		Collection  Collection `json:"collection"`
	}

	CreateLinkResponse struct {
		Name        string     `json:"name"`
		URL         string     `json:"url"`
		Type        string     `json:"type"`
		Description string     `json:"description"`
		Tags        []Tag      `json:"tags"`
		Collection  Collection `json:"collection"`
	}
)
