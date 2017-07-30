package api

type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Language    string `json:"language"`
}
