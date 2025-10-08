package dto_test

import (
	"reflect"
	"testing"

	"aviation-service/internal/dto"
)

func TestNewSuccessResponse(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		message  string
		expected dto.Response
	}{
		{
			name:    "Success to return success response",
			data:    map[string]string{"Hello": "World"},
			message: "Operation successful",
			expected: dto.Response{
				Success: true,
				Data:    map[string]string{"Hello": "World"},
				Message: "Operation successful",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dto.NewSuccessResponse(tt.data, tt.message)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Expected %+v, got %+v", tt.expected, got)
			}
		})
	}
}

func TestNewErrorResponse(t *testing.T) {
	tests := []struct {
		name     string
		err      interface{}
		message  string
		expected dto.Response
	}{
		{
			name:    "Success to return error response",
			err:     "Something went wrong",
			message: "Error occurred",
			expected: dto.Response{
				Success: false,
				Error:   "Something went wrong",
				Message: "Error occurred",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dto.NewErrorResponse(tt.err, tt.message)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Expected %+v, got %+v", tt.expected, got)
			}
		})
	}
}
