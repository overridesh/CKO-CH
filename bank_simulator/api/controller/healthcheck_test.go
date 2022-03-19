package controller

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthcheckRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Router with gin
	router := gin.Default()
	router.GET("/healthcheck", HealthcheckController)

	tests := []testCase{
		{
			name:       "healthcheck_success",
			path:       "/healthcheck",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			expect: gin.H{
				"ok": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := testRequest(router, tt.method, tt.path, tt.input, nil)
			statusCode := w.Result().StatusCode

			if statusCode != tt.statusCode {
				t.Fatalf("got: %v output: %v", statusCode, tt.statusCode)
			}

			if tt.expect != nil {
				var response map[string]interface{}
				err := json.Unmarshal([]byte(w.Body.Bytes()), &response)
				if err != nil {
					t.Fatal(err)
				}

				for key, _ := range tt.expect {
					value, exists := response[key]
					if !exists || value != tt.expect[key] {
						t.Fatalf("got: %v output: %v", value, tt.expect[key])
					}
				}
			}
		})
	}

}
