package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/beego/beego/v2/server/web"
)

type MainController struct {
	web.Controller
}
type FavoriteResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	ImageID  string `json:"image_id"`
	ImageURL string `json:"image_url"`
}
type FavoriteItem struct {
    ID        int       `json:"id"`
    ImageID   string    `json:"image_id"`
    SubID     string    `json:"sub_id"`
    CreatedAt string    `json:"created_at"`
    Image     CatImage  `json:"image"`
}
type Breed struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Origin        string `json:"origin"`
	Temperament   string `json:"temperament"`
	Description   string `json:"description"`
	Wikipedia_URL string `json:"wikipedia_url"`
}
type CatImage struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}
type BreedImage struct {
	URL string `json:"url"`
}

type VoteResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	ImageID  string `json:"image_id"`
	ImageURL string `json:"image_url"`
	Value    int    `json:"value,omitempty"`
}

// Get handles the main page
func (c *MainController) Get() {
	// Get initial cat image
	resultChannel := make(chan VoteResponse)
	go func() {
		resultChannel <- fetchNewImage()
	}()
	result := <-resultChannel

	if result.Status == "success" {
		c.Data["ImageID"] = result.ImageID
		c.Data["ImageURL"] = result.ImageURL
	}

	c.TplName = "index.tpl"
}
func (c *MainController) FetchNewImage() {
	resultChannel := make(chan VoteResponse)
	go func() {
		resultChannel <- fetchNewImage() // This function should fetch and return the new image
	}()
	result := <-resultChannel

	// Return the image or an error response
	c.Data["json"] = result
	c.ServeJSON()
}

func fetchNewImage() VoteResponse {
	apiUrl := "https://api.thecatapi.com/v1/images/search?api_key=live_i5ikdgttExQhChEMpt7UuGateqcwS8IVYlMGGBwSO3Mp7CatoYAhl9VAUgZ76Pqa"
	resp, err := http.Get(apiUrl)
	if err != nil {
		return VoteResponse{Status: "error", Message: "Error fetching new image"}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return VoteResponse{Status: "error", Message: "Error reading new image response"}
	}

	var images []CatImage
	err = json.Unmarshal(body, &images)
	if err != nil {
		return VoteResponse{Status: "error", Message: "Error parsing new image JSON"}
	}

	if len(images) > 0 {
		return VoteResponse{
			Status:   "success",
			Message:  "Vote submitted successfully",
			ImageID:  images[0].ID,
			ImageURL: images[0].URL,
		}
	} else {
		return VoteResponse{Status: "error", Message: "No new cat image found"}
	}
}

// SubmitVote handles vote submission
func (c *MainController) SubmitVote() {
	voteValue, _ := c.GetInt("value", 0)
	imageID := c.GetString("image_id")

	resultChannel := make(chan VoteResponse)
	go submitVote(imageID, voteValue, resultChannel)
	result := <-resultChannel

	c.Data["json"] = result
	c.ServeJSON()
}

func (c *MainController) GetVotes() {
	// Create request to fetch votes
	req, err := http.NewRequest("GET", "https://api.thecatapi.com/v1/votes?limit=10&order=DESC", nil)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// Add API key header
	req.Header.Set("x-api-key", "live_i5ikdgttExQhChEMpt7UuGateqcwS8IVYlMGGBwSO3Mp7CatoYAhl9VAUgZ76Pqa")

	// Add query parameters if needed
	q := req.URL.Query()
	q.Add("sub_id", "123") // Filter by sub_id used when voting
	req.URL.RawQuery = q.Encode()

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// Parse JSON response
	var votes interface{}
	err = json.Unmarshal(body, &votes)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		}
		c.ServeJSON()
		return
	}

	// Return the raw votes data
	c.Data["json"] = votes
	c.ServeJSON()
}
func submitVote(imageID string, value int, resultChannel chan VoteResponse) {
	// Construct the API request body
	payload := fmt.Sprintf(`{
        "image_id": "%s",
        "sub_id": "123", 
        "value": %d
    }`, imageID, value)

	// Send POST request to Cat API
	req, err := http.NewRequest("POST", "https://api.thecatapi.com/v1/votes",
		strings.NewReader(payload))
	if err != nil {
		resultChannel <- VoteResponse{Status: "error",
			Message: fmt.Sprintf("Error creating request: %s", err)}
		return
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "live_i5ikdgttExQhChEMpt7UuGateqcwS8IVYlMGGBwSO3Mp7CatoYAhl9VAUgZ76Pqa")

	// Send the request and get response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		resultChannel <- VoteResponse{Status: "error",
			Message: fmt.Sprintf("Error sending request: %s", err)}
		return
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			resultChannel <- VoteResponse{Status: "error",
				Message: fmt.Sprintf("Error reading response body: %s", err)}
			return
		}

		// Parse the response JSON
		var responseBody map[string]interface{}
		err = json.Unmarshal(body, &responseBody)
		if err != nil {
			resultChannel <- VoteResponse{Status: "error",
				Message: fmt.Sprintf("Error parsing response JSON: %s", err)}
			return
		}

		if message, exists := responseBody["message"].(string); exists && message == "SUCCESS" {
			// Fetch a new image asynchronously
			resultChannel <- fetchNewImage()
		} else {
			resultChannel <- VoteResponse{Status: "error",
				Message: fmt.Sprintf("Failed to submit vote: %v", responseBody)}
		}
	} else {
		resultChannel <- VoteResponse{Status: "error",
			Message: fmt.Sprintf("Request failed with status code: %d", resp.StatusCode)}
	}
}
func (c *MainController) SubmitFavorite() {
	imageID := c.GetString("image_id")

	resultChannel := make(chan FavoriteResponse)
	go submitFavorite(imageID, resultChannel)
	result := <-resultChannel

	c.Data["json"] = result
	c.ServeJSON()
}

func submitFavorite(imageID string, resultChannel chan FavoriteResponse) {
	payload := fmt.Sprintf(`{
        "image_id": "%s",
        "sub_id": "user123"
    }`, imageID)

	req, err := http.NewRequest("POST", "https://api.thecatapi.com/v1/favourites",
		strings.NewReader(payload))
	if err != nil {
		resultChannel <- FavoriteResponse{Status: "error",
			Message: fmt.Sprintf("Error creating request: %s", err)}
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "live_i5ikdgttExQhChEMpt7UuGateqcwS8IVYlMGGBwSO3Mp7CatoYAhl9VAUgZ76Pqa")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		resultChannel <- FavoriteResponse{Status: "error",
			Message: fmt.Sprintf("Error sending request: %s", err)}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		// Fetch a new image after successful favorite
		newImage := fetchNewImage()
		resultChannel <- FavoriteResponse{
			Status:   "success",
			Message:  "Image favorited successfully",
			ImageID:  newImage.ImageID,
			ImageURL: newImage.ImageURL,
		}
	} else {
		resultChannel <- FavoriteResponse{Status: "error",
			Message: fmt.Sprintf("Request failed with status code: %d", resp.StatusCode)}
	}
}
func (c *MainController) GetFavorites() {
    resultChannel := make(chan struct {
        Status    string         `json:"status"`
        Message   string         `json:"message"`
        Favorites []FavoriteItem `json:"favorites,omitempty"`
    })

    go func() {
        // Create request to Cat API
        req, err := http.NewRequest("GET", "https://api.thecatapi.com/v1/favourites", nil)
        if err != nil {
            resultChannel <- struct {
                Status    string         `json:"status"`
                Message   string         `json:"message"`
                Favorites []FavoriteItem `json:"favorites,omitempty"`
            }{
                Status:  "error",
                Message: fmt.Sprintf("Error creating request: %s", err),
            }
            return
        }

        // Add API key header
        req.Header.Set("x-api-key", "live_i5ikdgttExQhChEMpt7UuGateqcwS8IVYlMGGBwSO3Mp7CatoYAhl9VAUgZ76Pqa")

        // Send request
        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            resultChannel <- struct {
                Status    string         `json:"status"`
                Message   string         `json:"message"`
                Favorites []FavoriteItem `json:"favorites,omitempty"`
            }{
                Status:  "error",
                Message: fmt.Sprintf("Error sending request: %s", err),
            }
            return
        }
        defer resp.Body.Close()

        // Read response body
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            resultChannel <- struct {
                Status    string         `json:"status"`
                Message   string         `json:"message"`
                Favorites []FavoriteItem `json:"favorites,omitempty"`
            }{
                Status:  "error",
                Message: fmt.Sprintf("Error reading response: %s", err),
            }
            return
        }

        // Parse JSON response
        var favorites []FavoriteItem
        err = json.Unmarshal(body, &favorites)
        if err != nil {
            resultChannel <- struct {
                Status    string         `json:"status"`
                Message   string         `json:"message"`
                Favorites []FavoriteItem `json:"favorites,omitempty"`
            }{
                Status:  "error",
                Message: fmt.Sprintf("Error parsing JSON: %s", err),
            }
            return
        }

        // Send successful response
        resultChannel <- struct {
            Status    string         `json:"status"`
            Message   string         `json:"message"`
            Favorites []FavoriteItem `json:"favorites,omitempty"`
        }{
            Status:    "success",
            Message:   "Favorites retrieved successfully",
            Favorites: favorites,
        }
    }()

    result := <-resultChannel
    c.Data["json"] = result
    c.ServeJSON()
}
func (c *MainController) DeleteFavorite() {
    favoriteID := c.Ctx.Input.Param(":id")
    if favoriteID == "" {
        c.Data["json"] = struct {
            Status  string `json:"status"`
            Message string `json:"message"`
        }{
            Status:  "error",
            Message: "Favorite ID is required",
        }
        c.ServeJSON()
        return
    }

    resultChannel := make(chan struct {
        Status  string `json:"status"`
        Message string `json:"message"`
    })

    go func() {
        url := fmt.Sprintf("https://api.thecatapi.com/v1/favourites/%s", favoriteID)
        req, err := http.NewRequest("DELETE", url, nil)
        if err != nil {
            resultChannel <- struct {
                Status  string `json:"status"`
                Message string `json:"message"`
            }{
                Status:  "error",
                Message: fmt.Sprintf("Error creating request: %s", err),
            }
            return
        }

        req.Header.Set("x-api-key", "live_i5ikdgttExQhChEMpt7UuGateqcwS8IVYlMGGBwSO3Mp7CatoYAhl9VAUgZ76Pqa")

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            resultChannel <- struct {
                Status  string `json:"status"`
                Message string `json:"message"`
            }{
                Status:  "error",
                Message: fmt.Sprintf("Error sending request: %s", err),
            }
            return
        }
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            resultChannel <- struct {
                Status  string `json:"status"`
                Message string `json:"message"`
            }{
                Status:  "error",
                Message: fmt.Sprintf("Error reading response: %s", err),
            }
            return
        }

        if resp.StatusCode != http.StatusOK {
            resultChannel <- struct {
                Status  string `json:"status"`
                Message string `json:"message"`
            }{
                Status:  "error",
                Message: fmt.Sprintf("API error: %s", string(body)),
            }
            return
        }

        resultChannel <- struct {
            Status  string `json:"status"`
            Message string `json:"message"`
        }{
            Status:  "success",
            Message: "Favorite deleted successfully",
        }
    }()

    result := <-resultChannel
    c.Data["json"] = result
    c.ServeJSON()
}
// GetBreeds handles fetching all breeds
func (c *MainController) GetBreeds() {
	apiUrl := "https://api.thecatapi.com/v1/breeds"
	apiKey := "live_i5ikdgttExQhChEMpt7UuGateqcwS8IVYlMGGBwSO3Mp7CatoYAhl9VAUgZ76Pqa"

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		}
		c.ServeJSON()
		return
	}

	req.Header.Set("x-api-key", apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		}
		c.ServeJSON()
		return
	}

	var breeds []Breed
	err = json.Unmarshal(body, &breeds)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = breeds
	c.ServeJSON()
}

// GetBreedDetails handles fetching details for a specific breed
func (c *MainController) GetBreedDetails() {
	breedID := c.Ctx.Input.Param(":id")
	apiKey := "live_i5ikdgttExQhChEMpt7UuGateqcwS8IVYlMGGBwSO3Mp7CatoYAhl9VAUgZ76Pqa"

	// Fetch breed details and images concurrently
	breedChan := make(chan Breed)
	imagesChan := make(chan []BreedImage)
	errChan := make(chan error)

	// Fetch breed details
	go func() {
		url := fmt.Sprintf("https://api.thecatapi.com/v1/breeds/%s", breedID)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			errChan <- err
			return
		}
		req.Header.Set("x-api-key", apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errChan <- err
			return
		}

		var breed Breed
		if err := json.Unmarshal(body, &breed); err != nil {
			errChan <- err
			return
		}

		breedChan <- breed
	}()

	// Fetch breed images
	go func() {
		url := fmt.Sprintf("https://api.thecatapi.com/v1/images/search?limit=5&breed_ids=%s&api_key=%s", breedID, apiKey)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			errChan <- err
			return
		}
		req.Header.Set("x-api-key", apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errChan <- err
			return
		}

		var images []BreedImage
		if err := json.Unmarshal(body, &images); err != nil {
			errChan <- err
			return
		}

		imagesChan <- images
	}()

	// Wait for both goroutines
	select {
	case err := <-errChan:
		c.Data["json"] = map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		}
	case breed := <-breedChan:
		images := <-imagesChan
		c.Data["json"] = map[string]interface{}{
			"status": "success",
			"breed":  breed,
			"images": images,
		}
	}

	c.ServeJSON()
}
