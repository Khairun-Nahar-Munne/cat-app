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

type VoteController struct {
	beego.Controller
	
}


// FetchNewImage handles fetching a new image
func (c *VoteController) FetchNewImage() {
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
	apiUrl:= "https://api.thecatapi.com/v1"

	apiKey, err := beego.AppConfig.String("api_key")
	if err != nil || apiKey == "" {
		return models.VoteResponse{Status: "error", Message: "API key is missing or not configured"}
	}

	// Construct the URL with the API key
	url := apiUrl + "/images/search?api_key=" + apiKey

	// Making the API request
	resp, err := http.Get(url)
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
func (c *VoteController) SubmitVote() {
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


	apiURL, err := beego.AppConfig.String("cat_api_url")
	if err != nil || apiURL == "" {
		resultChannel <- models.VoteResponse{Status: "error", Message: "API URL is missing or not configured"}
		return
	}

	apiKey, err :=beego.AppConfig.String("api_key")
	if err != nil || apiKey == "" {
		resultChannel <- models.VoteResponse{Status: "error", Message: "API key is missing or not configured"}
		return
	}

	// Construct the API request body
	payload := fmt.Sprintf(`{
		"image_id": "%s",
		"sub_id": "123", 
		"value": %d
	}`, imageID, value)

	// Send POST request to the Cat API (with /votes in the URL)
	req, err := http.NewRequest("POST", apiURL+"/votes", strings.NewReader(payload))
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

func (c *VoteController) GetVotes() {
	// Fetch the API URL and API Key from the config
	apiURL, err := beego.AppConfig.String("cat_api_url")
	if err != nil || apiURL == "" {
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

	// Create a result channel
	resultChannel := make(chan map[string]interface{}, 1)

	// Launch a Goroutine to handle the HTTP request
	go func() {
		// Create request to fetch votes
		req, err := http.NewRequest("GET", apiURL+"/votes?limit=10&order=DESC", nil)
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

