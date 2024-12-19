package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/beego/beego/v2/server/web"
)

// BreedController handles breed-related actions.
type BreedController struct {
	web.Controller
}

type Breed struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Origin      string `json:"origin"`
	Temperament string `json:"temperament"`
	Description string `json:"description"`
	Wikipedia_URL string `json:"wikipedia_url"`
}

type BreedImage struct {
	URL string `json:"url"`
}

// DisplayBreeds handles fetching and displaying breeds.
func (c *BreedController) DisplayBreeds() {
	apiUrl := "https://api.thecatapi.com/v1/breeds"
	apiKey, _ := web.AppConfig.String("apiKey")

	// Helper function for making API requests
	makeRequest := func(url string, ch chan<- []byte, errCh chan<- error) {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			errCh <- fmt.Errorf("error creating request: %w", err)
			return
		}
		req.Header.Set("x-api-key", apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			errCh <- fmt.Errorf("error sending request: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			errCh <- fmt.Errorf("received non-200 response: %d", resp.StatusCode)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errCh <- fmt.Errorf("error reading response body: %w", err)
			return
		}

		ch <- body
	}

	// Fetch all breeds
	breedsCh := make(chan []byte)
	errCh := make(chan error)
	go makeRequest(apiUrl, breedsCh, errCh)

	select {
	case body := <-breedsCh:
		var breeds []Breed
		if err := json.Unmarshal(body, &breeds); err != nil {
			fmt.Println("Error unmarshalling breeds JSON:", err)
			c.TplName = "breeds.tpl"
			return
		}
		c.Data["Breeds"] = breeds

		// Process selected breed concurrently
		selectedBreedID := c.GetString("breed_id")
		if selectedBreedID == "" && len(breeds) > 0 {
			selectedBreedID = breeds[0].ID
		}

		if selectedBreedID != "" {
			breedDetailsCh := make(chan []byte)
			breedImagesCh := make(chan []byte)
			api_key:="live_TDpPM0HUwR1qIh40x1rYJwnQXhFjJN00iDCU387Jy6CBoCnbZCkbIPBzIkBktFXc"
			go makeRequest(fmt.Sprintf("https://api.thecatapi.com/v1/breeds/%s", selectedBreedID), breedDetailsCh, errCh)
			go makeRequest(fmt.Sprintf("https://api.thecatapi.com/v1/images/search?limit=5&breed_ids=%s&api_key=%s", selectedBreedID, api_key), breedImagesCh, errCh)


			select {
			case breedBody := <-breedDetailsCh:
				var breed Breed
				if err := json.Unmarshal(breedBody, &breed); err == nil {
					c.Data["Breed"] = breed
				} else {
					fmt.Println("Error unmarshalling breed details:", err)
				}
			case err := <-errCh:
				fmt.Println("Error fetching breed details:", err)
			}

			select {
			case imageBody := <-breedImagesCh:
				var images []BreedImage
				if err := json.Unmarshal(imageBody, &images); err == nil {
					c.Data["BreedImages"] = images
				} else {
					fmt.Println("Error unmarshalling breed images:", err)
				}
			case err := <-errCh:
				fmt.Println("Error fetching breed images:", err)
			}
		}
	case err := <-errCh:
		fmt.Println("Error fetching breeds:", err)
		c.TplName = "breeds.tpl"
		return
	}

	c.TplName = "breeds.tpl"
}
