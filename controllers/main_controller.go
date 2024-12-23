package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"cat-app/models"
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}


type Breed struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Origin        string `json:"origin"`
	Temperament   string `json:"temperament"`
	Description   string `json:"description"`
	Wikipedia_URL string `json:"wikipedia_url"`
}

type BreedImage struct {
	URL string `json:"url"`
}


// Get handles the main page
func (c *MainController) Get() {
	// Get initial cat image
	resultChannel := make(chan models.VoteResponse)
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

// FetchNewImage handles fetching a new image
func (c *MainController) FetchNewImage() {
	resultChannel := make(chan models.VoteResponse)
	go func() {
		resultChannel <- fetchNewImage() // This function should fetch and return the new image
	}()
	result := <-resultChannel

	// Return the image or an error response
	c.Data["json"] = result
	c.ServeJSON()
}

func fetchNewImage() models.VoteResponse {
	// Fetch API URL and API key from the config file
	apiUrl := beego.AppConfig.String("cat_api_url")
	apiKey := beego.AppConfig.String("api_key")

	// Construct the URL with the API key
	if apiUrl == "" || apiKey == "" {
		return models.VoteResponse{Status: "error", Message: "API URL or API key is missing"}
	}

	// Making the API request
	resp, err := http.Get(apiUrl + "&api_key=" + apiKey)
	if err != nil {
		return models.VoteResponse{Status: "error", Message: "Error fetching new image"}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.VoteResponse{Status: "error", Message: "Error reading new image response"}
	}

	var images []models.CatImage
	err = json.Unmarshal(body, &images)
	if err != nil {
		return models.VoteResponse{Status: "error", Message: "Error parsing new image JSON"}
	}

	if len(images) > 0 {
		return models.VoteResponse{
			Status:   "success",
			Message:  "Vote submitted successfully",
			ImageID:  images[0].ID,
			ImageURL: images[0].URL,
		}
	} else {
		return models.VoteResponse{Status: "error", Message: "No new cat image found"}
	}
}
// SubmitVote handles vote submission
func (c *MainController) SubmitVote() {
	voteValue, _ := c.GetInt("value", 0)
	imageID := c.GetString("image_id")

	// Channel should use models.VoteResponse
	resultChannel := make(chan models.VoteResponse)
	go submitVote(imageID, voteValue, resultChannel)
	result := <-resultChannel

	c.Data["json"] = result
	c.ServeJSON()
}

func submitVote(imageID string, value int, resultChannel chan models.VoteResponse) {
	// Fetch the API URL and API key from the Beego config
	apiURL := beego.AppConfig.String("vote_api_url")
	apiKey := beego.AppConfig.String("api_key")

	// Construct the API request body
	payload := fmt.Sprintf(`{
		"image_id": "%s",
		"sub_id": "123", 
		"value": %d
	}`, imageID, value)

	// Send POST request to the Cat API
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(payload))
	if err != nil {
		resultChannel <- models.VoteResponse{Status: "error", Message: fmt.Sprintf("Error creating request: %s", err)}
		return
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	// Send the request and get the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		resultChannel <- models.VoteResponse{Status: "error", Message: fmt.Sprintf("Error sending request: %s", err)}
		return
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			resultChannel <- models.VoteResponse{Status: "error", Message: fmt.Sprintf("Error reading response body: %s", err)}
			return
		}

		// Parse the response JSON
		var responseBody map[string]interface{}
		err = json.Unmarshal(body, &responseBody)
		if err != nil {
			resultChannel <- models.VoteResponse{Status: "error", Message: fmt.Sprintf("Error parsing response JSON: %s", err)}
			return
		}

		if message, exists := responseBody["message"].(string); exists && message == "SUCCESS" {
			// Fetch a new image asynchronously after vote
			resultChannel <- fetchNewImage()
		} else {
			resultChannel <- models.VoteResponse{Status: "error", Message: fmt.Sprintf("Failed to submit vote: %v", responseBody)}
		}
	} else {
		resultChannel <- models.VoteResponse{Status: "error", Message: fmt.Sprintf("Request failed with status code: %d", resp.StatusCode)}
	}
}
// GetVotes handles fetching all votes for a specific user

func (c *MainController) GetVotes() {
	// Get the API URL and API Key from the config
	apiURL := beego.AppConfig.String("vote_api_url")  // Fetches the API URL from the config
	apiKey := beego.AppConfig.String("api_key")  // Fetches the API key from the config

	// Create a result channel
	resultChannel := make(chan map[string]interface{}, 1)

	// Launch a Goroutine to handle the HTTP request
	go func() {
		// Create request to fetch votes
		req, err := http.NewRequest("GET", apiURL+"?limit=10&order=DESC", nil)
		if err != nil {
			resultChannel <- map[string]interface{}{
				"status":  "error",
				"message": err.Error(),
			}
			return
		}

		// Add API key header
		req.Header.Set("x-api-key", apiKey)

		// Add query parameters if needed
		q := req.URL.Query()
		q.Add("sub_id", "123") // Filter by sub_id used when voting
		req.URL.RawQuery = q.Encode()

		// Send request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			resultChannel <- map[string]interface{}{
				"status":  "error",
				"message": err.Error(),
			}
			return
		}
		defer resp.Body.Close()

		// Read response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			resultChannel <- map[string]interface{}{
				"status":  "error",
				"message": err.Error(),
			}
			return
		}

		// Parse JSON response
		var votes interface{}
		err = json.Unmarshal(body, &votes)
		if err != nil {
			resultChannel <- map[string]interface{}{
				"status":  "error",
				"message": err.Error(),
			}
			return
		}

		// Send the result back via the channel
		resultChannel <- map[string]interface{}{
			"status": "success",
			"data":   votes,
		}
	}()

	// Wait for the result to come back from the channel
	result := <-resultChannel

	// Return the result as a JSON response
	c.Data["json"] = result
	c.ServeJSON()
}


// SubmitFavorite handles favorite 
func (c *MainController) SubmitFavorite() {
	imageID := c.GetString("image_id")

	resultChannel := make(chan models.FavoriteResponse)
	newImageChannel := make(chan models.VoteResponse) // Channel to receive the new image result

	// Submit the favorite in a goroutine
	go submitFavorite(imageID, resultChannel, newImageChannel)

	// Wait for both the favorite submission result and the new image
	result := <-resultChannel
	newImageResult := <-newImageChannel

	// Combine the results (you can handle logic based on your requirements here)
	if result.Status == "success" {
		// If the favorite submission was successful, include the new image information
		result.ImageID = newImageResult.ImageID
		result.ImageURL = newImageResult.ImageURL
	}

	// Return the response as JSON
	c.Data["json"] = result
	c.ServeJSON()
}

func submitFavorite(imageID string, resultChannel chan models.FavoriteResponse, newImageChannel chan models.VoteResponse) {
	// Fetch the API URL and API key from the Beego config
	favoriteAPIURL := beego.AppConfig.String("favorite_api_url")
	apiKey := beego.AppConfig.String("api_key")

	// Construct the API request payload
	payload := fmt.Sprintf(`{
		"image_id": "%s",
		"sub_id": "user123"
	}`, imageID)

	// Send POST request to the Cat API
	req, err := http.NewRequest("POST", favoriteAPIURL, strings.NewReader(payload))
	if err != nil {
		resultChannel <- models.FavoriteResponse{Status: "error", Message: fmt.Sprintf("Error creating request: %s", err)}
		return
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	// Send the request and get the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		resultChannel <- models.FavoriteResponse{Status: "error", Message: fmt.Sprintf("Error sending request: %s", err)}
		return
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		// Send a successful result for favorite submission
		resultChannel <- models.FavoriteResponse{
			Status:  "success",
			Message: "Image favorited successfully",
		}

		// Fetch a new image concurrently
		go func() {
			newImage := fetchNewImage()
			newImageChannel <- newImage // Send the new image response through the channel
		}()
	} else {
		resultChannel <- models.FavoriteResponse{Status: "error", Message: fmt.Sprintf("Request failed with status code: %d", resp.StatusCode)}
	}
}

func (c *MainController) GetFavorites() {
	resultChannel := make(chan models.FavoriteResponse)

	go func() {
		// Fetch the API URL and API key from the Beego config
		favoriteAPIURL := beego.AppConfig.String("favorite_api_url")
		apiKey := beego.AppConfig.String("api_key")

		// Create request to Cat API
		req, err := http.NewRequest("GET", favoriteAPIURL, nil)
		if err != nil {
			resultChannel <- models.FavoriteResponse{
				Status:  "error",
				Message: fmt.Sprintf("Error creating request: %s", err),
			}
			return
		}

		// Add API key header
		req.Header.Set("x-api-key", apiKey)

		// Send request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			resultChannel <- models.FavoriteResponse{
				Status:  "error",
				Message: fmt.Sprintf("Error sending request: %s", err),
			}
			return
		}
		defer resp.Body.Close()

		// Read response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			resultChannel <- models.FavoriteResponse{
				Status:  "error",
				Message: fmt.Sprintf("Error reading response: %s", err),
			}
			return
		}

		// Parse JSON response
		var favorites []models.FavoriteItem
		err = json.Unmarshal(body, &favorites)
		if err != nil {
			resultChannel <- models.FavoriteResponse{
				Status:  "error",
				Message: fmt.Sprintf("Error parsing JSON: %s", err),
			}
			return
		}

		// Send successful response
		resultChannel <- models.FavoriteResponse{
			Status:    "success",
			Message:   "Favorites retrieved successfully",
			Favorites: favorites,
		}
	}()

	// Retrieve the result from the channel and send it as JSON response
	result := <-resultChannel
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *MainController) DeleteFavorite() {
	favoriteID := c.Ctx.Input.Param(":id")
	if favoriteID == "" {
		c.Data["json"] = models.FavoriteResponse{
			Status:  "error",
			Message: "Favorite ID is required",
		}
		c.ServeJSON()
		return
	}

	resultChannel := make(chan models.FavoriteResponse)

	go func() {
		// Fetch the API URL and API key from the Beego config
		favoriteAPIURL := beego.AppConfig.String("favorite_api_url")
		apiKey := beego.AppConfig.String("api_key")

		// Construct the URL for deleting the favorite
		url := fmt.Sprintf("%s/%s", favoriteAPIURL, favoriteID)

		// Create the DELETE request
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			resultChannel <- models.FavoriteResponse{
				Status:  "error",
				Message: fmt.Sprintf("Error creating request: %s", err),
			}
			return
		}

		// Add API key header
		req.Header.Set("x-api-key", apiKey)

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			resultChannel <- models.FavoriteResponse{
				Status:  "error",
				Message: fmt.Sprintf("Error sending request: %s", err),
			}
			return
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			resultChannel <- models.FavoriteResponse{
				Status:  "error",
				Message: fmt.Sprintf("Error reading response: %s", err),
			}
			return
		}

		// Check the response status
		if resp.StatusCode != http.StatusOK {
			resultChannel <- models.FavoriteResponse{
				Status:  "error",
				Message: fmt.Sprintf("API error: %s", string(body)),
			}
			return
		}

		// Send the success response
		resultChannel <- models.FavoriteResponse{
			Status:  "success",
			Message: "Favorite deleted successfully",
		}
	}()

	// Retrieve the result from the channel and send it as JSON response
	result := <-resultChannel
	c.Data["json"] = result
	c.ServeJSON()
}
// GetBreeds handles fetching all breeds
func (c *MainController) GetBreeds() {
	apiUrl := beego.AppConfig.String("breed_api_url")
	apiKey := beego.AppConfig.String("api_key")

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

	var breeds []models.Breed
	if err := json.Unmarshal(body, &breeds); err != nil {
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

// GetBreedDetails retrieves details and images for a specific breed
func (c *MainController) GetBreedDetails() {
	breedID := c.Ctx.Input.Param(":id")
	apiUrl := beego.AppConfig.String("breed_api_url")
	apiKey := beego.AppConfig.String("api_key")
	imagesUrl := beego.AppConfig.String("breed_images_url")

	breedChan := make(chan models.Breed)
	imagesChan := make(chan []models.BreedImage)
	errChan := make(chan error)

	// Fetch breed details
	go func() {
		url := fmt.Sprintf("%s/%s", apiUrl, breedID)
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

		var breed models.Breed
		if err := json.Unmarshal(body, &breed); err != nil {
			errChan <- err
			return
		}

		breedChan <- breed
	}()

	// Fetch breed images
	go func() {
		url := fmt.Sprintf(imagesUrl, breedID, apiKey)
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

		var images []models.BreedImage
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