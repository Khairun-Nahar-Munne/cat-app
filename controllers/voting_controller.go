package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/beego/beego/v2/server/web"
)

type VotingController struct {
	web.Controller
}

type CatImage struct {
	URL string `json:"url"`
}

func (c *VotingController) Get() {
	apiUrl, _ := web.AppConfig.String("cat_api_url")

	imageChan := make(chan string)

	// Go routine to fetch cat image
	go func() {
		fmt.Println("Fetching URL:", apiUrl) // Log the URL
		resp, err := http.Get(apiUrl)
		if err != nil {
			fmt.Println("Error making HTTP request:", err)
			imageChan <- ""
			return
		}
		defer resp.Body.Close()

		// Log the HTTP status
		fmt.Println("Response Status:", resp.Status)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			imageChan <- ""
			return
		}

		// Log the response body
		fmt.Println("Response Body:", string(body))

		var images []CatImage
		err = json.Unmarshal(body, &images)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			imageChan <- ""
			return
		}

		if len(images) > 0 {
			imageChan <- images[0].URL
		} else {
			imageChan <- ""
		}
	}()

	// Set the fetched image URL to the template
	c.Data["ImageURL"] = <-imageChan
	fmt.Println("Fetched Image URL:", c.Data["ImageURL"])
	c.TplName = "voting.tpl"
}
