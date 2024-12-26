package tests

import (
	"bytes"
	"cat-app/controllers"
	_ "cat-app/routers" // Import your routers so Beego initializes them
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"io/ioutil"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/stretchr/testify/assert" // Use testify for easier assertion
	"github.com/stretchr/testify/mock"   // For mocking HTTP client requests
)

// Initialize Beego app for testing
func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)

	beego.Router("/api/cat/fetch", &controllers.VoteController{}, "get:FetchNewImage")
	beego.Router("/api/vote", &controllers.VoteController{}, "post:SubmitVote")
	beego.Router("/api/vote", &controllers.VoteController{}, "get:GetVotes")
}

// Mock HTTP client to simulate the external API response
type MockRoundTripper struct {
	mock.Mock
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

// Test for FetchNewImage (GET request to /api/cat/fetch)
func TestFetchNewImageDef(t *testing.T) {
	// Mock the HTTP client (RoundTripper)
	mockRoundTripper := new(MockRoundTripper)

	// Simulate a successful response from the external API
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`[{"id":"MTgzMzAyOA","url":"https://cdn2.thecatapi.com/images/MTgzMzAyOA.jpg"}]`)),
	}

	// Mock the RoundTrip method to return the simulated response
	mockRoundTripper.On("RoundTrip", mock.Anything).Return(mockResponse, nil)

	// Set configuration values for the test
	beego.AppConfig.Set("cat_api_url", "https://api.thecatapi.com")
	beego.AppConfig.Set("api_key", "fake-api-key")

	// Override the default HTTP client with the mock client
	http.DefaultClient = &http.Client{Transport: mockRoundTripper}

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

	// Assert the response contains the expected keys (ImageID and ImageURL)
	assert.Contains(t, recorder.Body.String(), "image_id", "Expected the response to contain 'image_id'")
	assert.Contains(t, recorder.Body.String(), "image_url", "Expected the response to contain 'image_url'")

	// Assert that the mocked RoundTripper was called
	mockRoundTripper.AssertExpectations(t)
}

// Test for FetchNewImage when API returns an error
func TestFetchNewImageError(t *testing.T) {
	// Mock the HTTP client (RoundTripper)
	mockRoundTripper := new(MockRoundTripper)

	// Simulate an error response from the external API
	mockResponse := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`{"status": "error", "message": "Error fetching image"}`)),
	}

	// Mock the RoundTrip method to return the simulated response
	mockRoundTripper.On("RoundTrip", mock.Anything).Return(mockResponse, nil)

	// Set configuration values for the test
	beego.AppConfig.Set("cat_api_url", "https://api.thecatapi.com")
	beego.AppConfig.Set("api_key", "fake-api-key")

	// Override the default HTTP client with the mock client
	http.DefaultClient = &http.Client{Transport: mockRoundTripper}

	// Create a test request for the /api/cat/fetch endpoint
	req, err := http.NewRequest("GET", "/api/cat/fetch", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a recorder to capture the response
	recorder := httptest.NewRecorder()

	// Serve the request using Beego
	beego.BeeApp.Handlers.ServeHTTP(recorder, req)

	// Assert the response contains the error message
	assert.Contains(t, recorder.Body.String(), "Error parsing new image JSON", "Expected response to contain 'Error parsing new image JSON'")
	// Assert that the mocked RoundTripper was called
	mockRoundTripper.AssertExpectations(t)
}

