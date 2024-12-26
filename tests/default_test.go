package tests

import (
	"cat-app/controllers"
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
	// Get the absolute path of the root of your Beego application
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
	beego.Router("/", &controllers.MainController{})
}

func TestMainControllerGetDef(t *testing.T) {
	// Create a new HTTP request to test the MainController Get method
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record the response
	rr := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200")

}
