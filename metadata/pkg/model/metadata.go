package model

// Metadata defines the movie metadata.
type Metadata struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Director    string `json:"director,omitempty"`
}
