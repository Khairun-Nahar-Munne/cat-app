package controllers

import (
	"cat-app/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

type FavController struct {
	beego.Controller
}

// SubmitFavorite handles favorite 
func (c *FavController) SubmitFavorite() {
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
	favoriteAPIURL, err := beego.AppConfig.String("cat_api_url")
	if err != nil || favoriteAPIURL == "" {
		resultChannel <- models.FavoriteResponse{Status: "error", Message: "API URL is missing or not configured"}
		return
	}

	apiKey, err := beego.AppConfig.String("api_key")
	if err != nil || apiKey == "" {
		resultChannel <- models.FavoriteResponse{Status: "error", Message: "API key is missing or not configured"}
		return
	}

	// Append /favorites to the API URL
	favoriteAPIURL = favoriteAPIURL + "/favourites"

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


func (c *FavController) GetFavorites() {
	// Create a result channel
	resultChannel := make(chan models.FavoriteResponse)

	go func() {
		// Fetch the API URL and API key from the Beego config
		favoriteAPIURL, err := beego.AppConfig.String("cat_api_url")
		if err != nil || favoriteAPIURL == "" {
			resultChannel <- models.FavoriteResponse{
				Status:  "error",
				Message: "API URL is missing or not configured",
			}
			return
		}

		apiKey, err := beego.AppConfig.String("api_key")
		if err != nil || apiKey == "" {
			resultChannel <- models.FavoriteResponse{
				Status:  "error",
				Message: "API key is missing or not configured",
			}
			return
		}

		// Append /favourites to the base URL
		favoriteAPIURL = favoriteAPIURL + "/favourites"

		// Create request to fetch favorites
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

	// Retrieve the result from the channel and send it as a JSON response
	result := <-resultChannel
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *FavController) DeleteFavorite() {
	// Get the favorite ID from the URL parameter
	favoriteID := c.Ctx.Input.Param(":id")
	if favoriteID == "" {
		c.Data["json"] = models.FavoriteResponse{
			Status:  "error",
			Message: "Favorite ID is required",
		}
		c.ServeJSON()
		return
	}

	// Create a result channel to handle the response
	resultChannel := make(chan models.FavoriteResponse)

	go func() {
		// Fetch the API URL and API key from the Beego config
		favoriteAPIURL, err := beego.AppConfig.String("cat_api_url")
		if err != nil || favoriteAPIURL == "" {
			resultChannel <- models.FavoriteResponse{
				Status:  "error",
				Message: "API URL is missing or not configured",
			}
			return
		}

		apiKey, err := beego.AppConfig.String("api_key")
		if err != nil || apiKey == "" {
			resultChannel <- models.FavoriteResponse{
				Status:  "error",
				Message: "API key is missing or not configured",
			}
			return
		}

		// Append /favourites to the base URL to form the full URL for deleting the favorite
		url := fmt.Sprintf("%s/favourites/%s", favoriteAPIURL, favoriteID)

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
