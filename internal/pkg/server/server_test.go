package server

import (
	"net/http"
	"net/http/httptest"
	"proj1/internal/pkg/storage"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	file = "slice_storage.json"
)

func setupTestServer(stor *storage.SliceStorage) *gin.Engine {
	s := New("localhost:8090", stor)
	return s.engine
}

func TestHandlerSetSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	stor2, _ := storage.NewSliceStorage(file)
	router := setupTestServer(&stor2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/scalar/set/testkey/testvalue", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandlerSetBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stor2, _ := storage.NewSliceStorage(file)
	router := setupTestServer(&stor2)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodPost, "/scalar/set/testkey/testval?exp=uuu", nil)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerGetSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	stor2, _ := storage.NewSliceStorage(file)
	stor2.Set("testkey", "42")
	router := setupTestServer(&stor2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/scalar/get/testkey", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedBody := `{"value":"42"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestHandlerGetNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	stor2, _ := storage.NewSliceStorage(file)
	router := setupTestServer(&stor2)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/scalar/get/nonexistent", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
