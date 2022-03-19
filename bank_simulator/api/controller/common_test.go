package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
)

// testCase structure for test routes
type testCase struct {
	name       string
	path       string
	method     string
	statusCode int
	headers    map[string]string
	input      string
	expect     gin.H
}

// testRequest func for record requests
func testRequest(r http.Handler, method, path string, body string, headers map[string]string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(body))

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
