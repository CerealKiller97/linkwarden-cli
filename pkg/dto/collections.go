package dto

import (
	"time"
)

type (
	CreateCollectionRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Color       string `json:"color"`
		Icon        string `json:"icon"`
		IconWeight  string `json:"iconWeight"`
		ParentID    int    `json:"parentId"`
	}

	Member struct {
		UserID       int       `json:"userId"`
		CollectionID int       `json:"collectionId"`
		CanCreate    bool      `json:"canCreate"`
		CanUpdate    bool      `json:"canUpdate"`
		CanDelete    bool      `json:"canDelete"`
		CreatedAt    time.Time `json:"createdAt"`
		UpdatedAt    time.Time `json:"updatedAt"`
		User         struct {
			Username string  `json:"username"`
			Name     string  `json:"name"`
			Image    *string `json:"image"`
		} `json:"user"`
	}

	CreateCollectionResponse struct {
		Response struct {
			ID          int       `json:"id"`
			Name        string    `json:"name"`
			Description string    `json:"description"`
			Icon        string    `json:"icon"`
			IconWeight  string    `json:"iconWeight"`
			Color       string    `json:"color"`
			ParentID    int       `json:"parentId"`
			IsPublic    bool      `json:"isPublic"`
			OwnerID     int       `json:"ownerId"`
			CreatedByID int       `json:"createdById"`
			CreatedAt   time.Time `json:"createdAt"`
			UpdatedAt   time.Time `json:"updatedAt"`
			Count       struct {
				Links int `json:"links"`
			} `json:"_count"`
			Members []Member `json:"members"`
		} `json:"response"`
	}
)
