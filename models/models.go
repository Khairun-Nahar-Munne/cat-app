package models

// VoteResponse represents the response when a vote is submitted or when fetching an image
type VoteResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	ImageID  string `json:"image_id"`
	ImageURL string `json:"image_url"`
	Value    int    `json:"value,omitempty"`
}
// CatImage represents the structure of the cat image data returned by the API
type CatImage struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}
type FavoriteItem struct {
    ID        int       `json:"id"`
    ImageID   string    `json:"image_id"`
    SubID     string    `json:"sub_id"`
    CreatedAt string    `json:"created_at"`
    Image     CatImage  `json:"image"`
}
type FavoriteResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	ImageID  string `json:"image_id"`
	ImageURL string `json:"image_url"`
	Favorites []FavoriteItem `json:"favorites,omitempty"`
}

// Breed represents a breed of cat
type Breed struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Origin        string `json:"origin"`
	Temperament   string `json:"temperament"`
	Description   string `json:"description"`
	Wikipedia_URL string `json:"wikipedia_url"`
}

// BreedImage represents an image of a breed
type BreedImage struct {
	URL string `json:"url"`
}