package testutils

import (
	"aviation-service/internal/dto"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"
)

type ExpectedResult struct {
	Status  int
	Data    interface{}
	Message string
	Error   string
}

func AssertHandlerResponse(t *testing.T, rr *httptest.ResponseRecorder, expected ExpectedResult) error {
	t.Helper()
	
	if rr.Code != expected.Status {
		return fmt.Errorf("Expected status %d, got %d", expected.Status, rr.Code)
	}

	var body dto.Response
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		return fmt.Errorf("Failed to decode response body: %v", err)
	}

	if expected.Data != nil {
		var expectedMap map[string]interface{}
		expectedBytes, _ := json.Marshal(expected.Data)
		json.Unmarshal(expectedBytes, &expectedMap)

		if !reflect.DeepEqual(body.Data, expectedMap) {
			return fmt.Errorf("Expected data %+v, got %+v", expectedMap, body.Data)
		}
	}

	if expected.Message != "" && body.Message != expected.Message {
		return fmt.Errorf("Expected message %q, got %q", expected.Message, body.Message)
	}

	if expected.Error != "" && body.Error != expected.Error {
		return fmt.Errorf("Expected error %q, got %q", expected.Error, body.Error)
	}
	return nil
}
