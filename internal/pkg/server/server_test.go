package server

import (
	"net/http"
	"net/http/httptest"
	"proj1/internal/pkg/storage"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandlerSetSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stor, _ := storage.NewSliceStorage()
	s := New("localhost:8090", &stor)
	router := s.newAPI()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/scalar/set/testkey", strings.NewReader(`{"value":"testvalue"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	res, _ := stor.Get("testkey")
	assert.Equal(t, "testvalue", res)
}

func TestHandlerSetBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stor, _ := storage.NewSliceStorage()
	s := New("localhost:8090", &stor)
	router := s.newAPI()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/scalar/set/testkey", strings.NewReader(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerGetSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stor, _ := storage.NewSliceStorage()
	stor.Set("testkey", "testvalue")
	s := New("localhost:8090", &stor)
	router := s.newAPI()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/scalar/get/testkey", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedBody := `{"value":"testvalue"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestHandlerGetNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stor, _ := storage.NewSliceStorage()
	s := New("localhost:8090", &stor)
	router := s.newAPI()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/scalar/get/nonexistent", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
