package controllers

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"

    "github.com/beego/beego/v2/server/web"
)

type VotingController struct {
    web.Controller
}

type CatImage struct {
    ID  string `json:"id"`
    URL string `json:"url"`
}

type VoteResponse struct {
    Status   string `json:"status"`
    Message  string `json:"message"`
    ImageID  string `json:"image_id"`
    ImageURL string `json:"image_url"`
    Value    int    `json:"value,omitempty"`
}


// Function to handle the API call and vote submission asynchronously
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
    req.Header.Set("x-api-key", "live_TDpPM0HUwR1qIh40x1rYJwnQXhFjJN00iDCU387Jy6CBoCnbZCkbIPBzIkBktFXc")

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

// Fetch a new cat image asynchronously
func fetchNewImage() VoteResponse {
    apiUrl := "https://api.thecatapi.com/v1/images/search?api_key=live_TDpPM0HUwR1qIh40x1rYJwnQXhFjJN00iDCU387Jy6CBoCnbZCkbIPBzIkBktFXc"
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

// Handle GET request for voting page
func (c *VotingController) Get() {
    // Create a channel for asynchronous image fetching
    resultChannel := make(chan VoteResponse)

    // Start a goroutine to fetch the image
    go func() {
        resultChannel <- fetchNewImage()
    }()

    // Wait for the result from the goroutine
    result := <-resultChannel

    if result.Status == "success" {
        c.Data["ImageID"] = result.ImageID
        c.Data["ImageURL"] = result.ImageURL
        c.TplName = "voting.tpl"
    } else {
        c.Ctx.WriteString("Error fetching cat image")
    }
}

// Handle POST request for submitting vote
func (c *VotingController) Post() {
    voteValue, _ := c.GetInt("value", 0)
    imageID := c.GetString("image_id")
    
    // Check if this is an AJAX request
    isAjax := c.Ctx.Input.IsAjax()

    // Create a channel to receive the result of the vote submission
    resultChannel := make(chan VoteResponse)

    // Start the goroutine for submitting the vote
    go submitVote(imageID, voteValue, resultChannel)

    // Wait for the result from the goroutine
    result := <-resultChannel

    if isAjax {
        // If it's an AJAX request, return JSON
        c.Data["json"] = result
        c.ServeJSON()
    } else {
        // For regular form submissions, redirect back to the voting page
        if result.Status == "success" {
            c.Redirect("/voting", 302)
        } else {
            // Handle error case
            c.Data["Error"] = result.Message
            c.TplName = "voting.tpl"
        }
    }
}

func (c *VotingController) GetVotes() {
    // Create request to fetch votes
    req, err := http.NewRequest("GET", "https://api.thecatapi.com/v1/votes", nil)
    if err != nil {
        c.Data["json"] = map[string]interface{}{
            "status": "error",
            "message": err.Error(),
        }
        c.ServeJSON()
        return
    }

    // Add API key header
    req.Header.Set("x-api-key", "live_TDpPM0HUwR1qIh40x1rYJwnQXhFjJN00iDCU387Jy6CBoCnbZCkbIPBzIkBktFXc")

    // Add query parameters if needed
    q := req.URL.Query()
    q.Add("sub_id", "123")  // Filter by sub_id used when voting
    req.URL.RawQuery = q.Encode()

    // Send request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        c.Data["json"] = map[string]interface{}{
            "status": "error",
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
            "status": "error",
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
            "status": "error",
            "message": err.Error(),
        }
        c.ServeJSON()
        return
    }

    // Return the raw votes data
    c.Data["json"] = votes
    c.ServeJSON()
}

