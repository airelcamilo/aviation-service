package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeWriter struct{}

func (f *fakeWriter) Header() http.Header {
	return http.Header{}
}
func (f *fakeWriter) Write([]byte) (int, error) {
	return 0, fmt.Errorf("forced write error")
}
func (f *fakeWriter) WriteHeader(statusCode int) {}

func TestRespondWithJSON(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		writer     http.ResponseWriter
		payload    any
	}{
		{
			name:       "Success create response",
			statusCode: http.StatusOK,
			writer:     httptest.NewRecorder(),
			payload:    map[string]string{"Hello": "World"},
		},
		{
			name:       "Error marshal payload",
			statusCode: http.StatusOK,
			writer:     httptest.NewRecorder(),
			payload:    map[string]interface{}{"bad": make(chan int)},
		},
		{
			name:       "Error write payload",
			statusCode: http.StatusOK,
			writer:     &fakeWriter{},
			payload:    map[string]string{"Hello": "World"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respondWithJSON(tt.writer, tt.statusCode, tt.payload)
		})
	}
}
