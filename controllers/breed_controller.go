package controllers

import (
	"cat-app/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
)

type BreedController struct {
	beego.Controller
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// GetBreeds handles the GET request for retrieving all cat breeds
func (c *BreedController) GetBreeds() {
	// Make HTTP request to the cat API
	resp, err := http.Get("https://api.thecatapi.com/v1/breeds")
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = ErrorResponse{Message: "Failed to fetch breeds"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = ErrorResponse{Message: "Failed to read response"}
		c.ServeJSON()
		return
	}

	// Check if the response status is not OK
	if resp.StatusCode != http.StatusOK {
		c.Ctx.Output.SetStatus(resp.StatusCode)
		c.Data["json"] = ErrorResponse{Message: fmt.Sprintf("API returned status: %d", resp.StatusCode)}
		c.ServeJSON()
		return
	}

	// Try to parse the response
	var breeds []models.Breed
	if err := json.Unmarshal(body, &breeds); err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = ErrorResponse{Message: "Failed to parse response"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = breeds
	c.ServeJSON()
}

// GetBreedDetails retrieves details and images for a specific breed
func (c *BreedController) GetBreedDetails() {
	// Get the breedID from the URL parameter
	breedID := c.Ctx.Input.Param(":id")

	// Fetch the API URL and API key from the Beego config
	apiUrl, err := beego.AppConfig.String("cat_api_url")
	if err != nil || apiUrl == "" {
		c.Data["json"] = map[string]interface{}{
			"status":  "error",
			"message": "API URL is missing or not configured",
		}
		c.ServeJSON()
		return
	}

	apiKey, err := beego.AppConfig.String("api_key")
	if err != nil || apiKey == "" {
		c.Data["json"] = map[string]interface{}{
			"status":  "error",
			"message": "API key is missing or not configured",
		}
		c.ServeJSON()
		return
	}

	// Construct the URLs for the breed details and breed images
	breedUrl := fmt.Sprintf("%s/breeds/%s", apiUrl, breedID)
	imagesUrl := fmt.Sprintf("%s/images/search?limit=5&breed_ids=%s&api_key=%s", apiUrl, breedID, apiKey)

	// Create channels for the breed data, images, and errors
	breedChan := make(chan models.Breed)
	imagesChan := make(chan []models.BreedImage)
	errChan := make(chan error)

	// Fetch breed details asynchronously
	go func() {
		req, err := http.NewRequest("GET", breedUrl, nil)
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

	// Fetch breed images asynchronously
	go func() {
		req, err := http.NewRequest("GET", imagesUrl, nil)
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

	// Wait for both goroutines to finish and handle the result
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