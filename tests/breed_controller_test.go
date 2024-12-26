package tests

import (
	"bytes"
	"cat-app/controllers"
	"cat-app/models"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/stretchr/testify/assert"
)

type roundTripperFunc func(req *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/breeds", &controllers.BreedController{}, "get:GetBreeds")
	beego.Router("/api/breed/:id", &controllers.BreedController{}, "get:GetBreedDetails")
}

func TestBreedController_GetBreeds(t *testing.T) {
	// Test cases for GetBreeds
	tests := []struct {
		name            string
		mockResponse    string
		mockStatusCode  int
		mockError       error
		expectedCode    int
		apiUrl          string
		expectedMessage string
	}{
		{
			name:            "Valid Response",
			mockResponse:    `[{"id": "abys", "name": "Abyssinian"}]`,
			mockStatusCode:  http.StatusOK,
			mockError:       nil,
			expectedCode:    http.StatusOK,
			apiUrl:          "https://api.example.com",
			expectedMessage: "",
		},
		{
			name:            "Empty Response",
			mockResponse:    `[]`,
			mockStatusCode:  http.StatusOK,
			mockError:       nil,
			expectedCode:    http.StatusOK,
			apiUrl:          "https://api.example.com",
			expectedMessage: "",
		},
		{
			name:            "API Error",
			mockResponse:    `{"message": "Internal Server Error"}`,
			mockStatusCode:  http.StatusInternalServerError,
			mockError:       nil,
			expectedCode:    http.StatusInternalServerError,
			apiUrl:          "https://api.example.com",
			expectedMessage: "API returned status: 500", // Add the specific error message you are expecting
		},
		{
			name:            "Network Error",
			mockResponse:    "",
			mockStatusCode:  0,
			mockError:       errors.New("network error"),
			expectedCode:    http.StatusInternalServerError,
			apiUrl:          "https://api.example.com",
			expectedMessage: "Failed to fetch breeds", // Expect this message for network error
		},
		{
			name:            "Invalid JSON Response",
			mockResponse:    `{invalid json}`,
			mockStatusCode:  http.StatusOK,
			mockError:       nil,
			expectedCode:    http.StatusInternalServerError,
			apiUrl:          "https://api.example.com",
			expectedMessage: "Failed to parse response", // Expect this message for invalid JSON
		},
		
		{
			name:            "Missing API URL",
			mockResponse:    "",
			mockStatusCode:  0,
			mockError:       nil,
			expectedCode:    http.StatusOK,
			apiUrl:          "",
			expectedMessage: "API returned status: 0", // Handle empty API URL
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original config and restore after test
			originalAPIUrl, _ := beego.AppConfig.String("apiurl")
			defer func() {
				beego.AppConfig.Set("apiurl", originalAPIUrl)
			}()

			// Set test-specific API URL
			beego.AppConfig.Set("apiurl", tt.apiUrl)

			r, _ := http.NewRequest("GET", "/api/breeds", nil)
			w := httptest.NewRecorder()

			originalHTTPClient := http.DefaultClient
			defer func() { http.DefaultClient = originalHTTPClient }()

			// Mocking HTTP Client
			http.DefaultClient = &http.Client{
				Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					// Simulate body read error for the "Body Read Error" case
					if tt.name == "Body Read Error" {
						return &http.Response{
							StatusCode: tt.mockStatusCode,
							Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))), // Empty body to simulate error
						}, nil
					}

					return &http.Response{
						StatusCode: tt.mockStatusCode,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(tt.mockResponse))),
					}, nil
				}),
			}

			beego.BeeApp.Handlers.ServeHTTP(w, r)

			// Check the status code
			assert.Equal(t, tt.expectedCode, w.Code)

			// If we expect a successful response (Status OK)
			if tt.expectedCode == http.StatusOK && tt.expectedMessage == "" {
				var response []models.Breed
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
			} else if tt.expectedMessage != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMessage, response["message"])
			}
		})
	}
}
