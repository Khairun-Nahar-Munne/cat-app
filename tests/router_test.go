package tests

import (
	_ "cat-app/routers" // Import your routers so Beego initializes them
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/stretchr/testify/assert" // Use testify for easier assertion
)

func init() {
	// Initialize the Beego app for testing
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)

}

func TestMainControllerGet(t *testing.T) {
	// Create a test request for the root route ("/")
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a recorder to capture the response
	recorder := httptest.NewRecorder()

	// Serve the request using Beego
	beego.BeeApp.Handlers.ServeHTTP(recorder, req)

	// Assert the status code is 200 OK
	assert.Equal(t, 200, recorder.Code, "Expected status code 200")

	// You can also assert that certain elements are in the response (e.g., ImageID, ImageURL)
	assert.Contains(t, recorder.Body.String(), "imageId", "Expected the response to contain 'imageId'")
	assert.Contains(t, recorder.Body.String(), "catImage", "Expected the response to contain 'catImage'")
}

// Test GET /api/breeds
func TestGetBreedsRou(t *testing.T) {
	// Create a test request for the /api/breeds endpoint
	req, err := http.NewRequest("GET", "/api/breeds", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a recorder to capture the response
	recorder := httptest.NewRecorder()

	// Serve the request using Beego
	beego.BeeApp.Handlers.ServeHTTP(recorder, req)

	// Assert the status code is 200 OK
	assert.Equal(t, 200, recorder.Code)
}

// Test POST /api/vote
func TestSubmitVote(t *testing.T) {
	// Create a test request for the /api/vote endpoint
	req, err := http.NewRequest("POST", "/api/vote", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a recorder to capture the response
	recorder := httptest.NewRecorder()

	// Serve the request using Beego
	beego.BeeApp.Handlers.ServeHTTP(recorder, req)

	// Assert the status code is 200 OK (adjust based on your controller logic)
	assert.Equal(t, 200, recorder.Code)
}

// Test GET /api/favourite
func TestGetFavorites(t *testing.T) {
	// Create a test request for the /api/favourite endpoint
	req, err := http.NewRequest("GET", "/api/favourite", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a recorder to capture the response
	recorder := httptest.NewRecorder()

	// Serve the request using Beego
	beego.BeeApp.Handlers.ServeHTTP(recorder, req)

	// Assert the status code is 200 OK
	assert.Equal(t, 200, recorder.Code)
}

// Test DELETE /api/favourite/:id
func TestDeleteFavorite(t *testing.T) {
	// Create a test request for the /api/favourite/:id endpoint
	req, err := http.NewRequest("DELETE", "/api/favourite/123", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a recorder to capture the response
	recorder := httptest.NewRecorder()

	// Serve the request using Beego
	beego.BeeApp.Handlers.ServeHTTP(recorder, req)

	// Assert the status code is 200 OK (adjust based on your controller logic)
	assert.Equal(t, 200, recorder.Code)
}

// Test GET /api/breed/:id
func TestGetBreedDetailsRou(t *testing.T) {
	// Create a test request for the /api/breed/:id endpoint
	req, err := http.NewRequest("GET", "/api/breed/abys", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a recorder to capture the response
	recorder := httptest.NewRecorder()

	// Serve the request using Beego
	beego.BeeApp.Handlers.ServeHTTP(recorder, req)

	// Assert the status code is 200 OK
	assert.Equal(t, 200, recorder.Code)
}

// Test GET /api/cat/fetch
// Test GET /api/cat/fetch
func TestFetchNewImage(t *testing.T) {
	// Create a test request for the /api/cat/fetch endpoint
	req, err := http.NewRequest("GET", "/api/cat/fetch", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a recorder to capture the response
	recorder := httptest.NewRecorder()

	// Serve the request using Beego
	beego.BeeApp.Handlers.ServeHTTP(recorder, req)

	// Assert the status code is 200 OK
	assert.Equal(t, 200, recorder.Code, "Expected status code 200")

	// Correct the key names to match the actual response structure
	assert.Contains(t, recorder.Body.String(), "image_id", "Expected response to contain 'image_id'")
	assert.Contains(t, recorder.Body.String(), "image_url", "Expected response to contain 'image_url'")
}

// Test GET /api/vote
func TestGetVotes(t *testing.T) {
	// Create a test request for the /api/vote endpoint
	req, err := http.NewRequest("GET", "/api/vote", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a recorder to capture the response
	recorder := httptest.NewRecorder()

	// Serve the request using Beego
	beego.BeeApp.Handlers.ServeHTTP(recorder, req)

	// Assert the status code is 200 OK
	assert.Equal(t, 200, recorder.Code, "Expected status code 200")

	// Correct the key name to 'data' for the /api/vote response
	assert.Contains(t, recorder.Body.String(), "data", "Expected response to contain 'data'")
}
