package routes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/magiconair/properties/assert"
	"github.com/phazyy/golang-rest-api/middleware"
	"github.com/phazyy/golang-rest-api/routes"
	"github.com/phazyy/sqlboiler-gin/models"
	"gopkg.in/inconshreveable/log15.v2"
)

var log = log15.New()
var pilotID int

// SetupRouter :
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Logger(log))
	r.Use(middleware.Database(log))
	gin.SetMode(gin.TestMode)

	v1 := r.Group("/v1")
	{
		pilot := new(routes.PilotRoutes)

		v1.GET("/pilots", pilot.GetAll)
		v1.GET("/pilots/:id", pilot.Get)
		v1.POST("/pilots", pilot.Create)
		v1.PUT("/pilots/:id", pilot.Update)
		v1.DELETE("/pilots/:id", pilot.Delete)
	}
	return r
}

func main() {
	r := SetupRouter()
	r.Run()
}

// TestCreatePilot : Assert pilot creation - must return 201
func TestCreatePilot(t *testing.T) {
	testRouter := SetupRouter()
	testPilot := &models.Pilot{Name: "Adam"}

	data, _ := json.Marshal(testPilot)
	req, err := http.NewRequest("POST", "/v1/pilots", bytes.NewBufferString(string(data)))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}

	res := httptest.NewRecorder()
	testRouter.ServeHTTP(res, req)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	resp := struct {
		Message string
		ID      int
	}{}

	json.Unmarshal(body, &resp)
	pilotID = resp.ID

	assert.Equal(t, res.Code, 201)
}

// TestCreateInvalidPilot : Assert invalid pilot create - must return 400
func TestCreateInvalidPilot(t *testing.T) {
	testRouter := SetupRouter()

	req, err := http.NewRequest("POST", "/v1/pilots", bytes.NewBufferString("Test"))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}

	res := httptest.NewRecorder()

	testRouter.ServeHTTP(res, req)
	assert.Equal(t, res.Code, 400)
}

// TestGetPilot : Assert pilot fetch - must return 200
func TestGetPilot(t *testing.T) {
	testRouter := SetupRouter()

	url := fmt.Sprintf("/v1/pilots/%d", pilotID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}

	res := httptest.NewRecorder()
	testRouter.ServeHTTP(res, req)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	resp := struct {
		ID   int
		Name string
	}{}

	json.Unmarshal(body, &resp)
	pilotID = resp.ID

	assert.Equal(t, res.Code, 200)
	assert.Equal(t, resp.Name, "Adam")
}

// TestGetInvalidPilot : Assert negative pilot fetch - must return 404
func TestGetInvalidPilot(t *testing.T) {
	testRouter := SetupRouter()

	url := fmt.Sprintf("/v1/pilots/%d", pilotID+1)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}

	res := httptest.NewRecorder()

	testRouter.ServeHTTP(res, req)
	assert.Equal(t, res.Code, 404)
}

// TestUpdatePilot : Assert pilot update - must return 200
func TestUpdatePilot(t *testing.T) {
	testRouter := SetupRouter()
	testPilot := &models.Pilot{Name: "updateAdam"}

	data, _ := json.Marshal(testPilot)
	url := fmt.Sprintf("/v1/pilots/%d", pilotID)
	req, err := http.NewRequest("PUT", url, bytes.NewBufferString(string(data)))

	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}

	res := httptest.NewRecorder()
	testRouter.ServeHTTP(res, req)

	assert.Equal(t, res.Code, 200)
}

// TestDeletePilot : Assert pilot deletion - must return 204
func TestDeletePilot(t *testing.T) {
	testRouter := SetupRouter()

	url := fmt.Sprintf("/v1/pilots/%d", pilotID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println(err)
	}

	res := httptest.NewRecorder()

	testRouter.ServeHTTP(res, req)
	assert.Equal(t, res.Code, 204)
}
