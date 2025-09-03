package httpx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAPIResponse tests the APIResponse struct
func TestAPIResponse(t *testing.T) {
	// 🎯 Test Strategy: Test the response struct structure and JSON marshaling

	t.Run("should create APIResponse with all fields", func(t *testing.T) {
		// 🔧 Setup: Create a response with all fields
		now := time.Now().UTC()
		response := APIResponse{
			Timestamp: now,
			Status:    200,
			Success:   true,
			Message:   "Operation successful",
			Data:      map[string]string{"key": "value"},
			Error:     nil,
		}

		// ✅ Assertions: Verify all fields are set correctly
		assert.Equal(t, now, response.Timestamp)
		assert.Equal(t, 200, response.Status)
		assert.True(t, response.Success)
		assert.Equal(t, "Operation successful", response.Message)
		assert.Equal(t, map[string]string{"key": "value"}, response.Data)
		assert.Nil(t, response.Error)
	})

	t.Run("should marshal APIResponse to JSON correctly", func(t *testing.T) {
		// 🔧 Setup: Create a response
		response := APIResponse{
			Timestamp: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			Status:    201,
			Success:   true,
			Message:   "User created",
			Data:      map[string]int{"id": 123},
		}

		// 🚀 Action: Marshal to JSON
		jsonData, err := json.Marshal(response)

		// ✅ Assertions: Should marshal successfully
		require.NoError(t, err)

		// Verify JSON structure
		var unmarshaled APIResponse
		err = json.Unmarshal(jsonData, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, response.Status, unmarshaled.Status)
		assert.Equal(t, response.Success, unmarshaled.Success)
		assert.Equal(t, response.Message, unmarshaled.Message)

		// When unmarshaling, the Data field becomes map[string]interface{}
		// So we need to check the values individually
		expectedData := response.Data.(map[string]int)
		actualData := unmarshaled.Data.(map[string]interface{})
		assert.Equal(t, expectedData["id"], int(actualData["id"].(float64)))
	})
}

// TestWriteJSON tests the WriteJSON function
func TestWriteJSON(t *testing.T) {
	// 🎯 Test Strategy: Test HTTP response writing with different scenarios

	t.Run("should write success response correctly", func(t *testing.T) {
		// 🔧 Setup: Create response recorder
		rr := httptest.NewRecorder()

		// 🚀 Action: Write JSON response
		WriteJSON(rr, http.StatusOK, true, "Success message", map[string]string{"key": "value"}, nil)

		// ✅ Assertions: Verify response
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		// Parse response body
		var response APIResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response.Success)
		assert.Equal(t, "Success message", response.Message)

		// Data becomes map[string]interface{} when unmarshaled
		actualData := response.Data.(map[string]interface{})
		assert.Equal(t, "value", actualData["key"])
		assert.Nil(t, response.Error)
	})

	t.Run("should write error response correctly", func(t *testing.T) {
		// 🔧 Setup: Create response recorder
		rr := httptest.NewRecorder()

		// 🚀 Action: Write error response
		WriteJSON(rr, http.StatusBadRequest, false, "Error message", nil, map[string]string{"detail": "error detail"})

		// ✅ Assertions: Verify response
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		// Parse response body
		var response APIResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.False(t, response.Success)
		assert.Equal(t, "Error message", response.Message)
		assert.Nil(t, response.Data)

		// Error becomes map[string]interface{} when unmarshaled
		actualError := response.Error.(map[string]interface{})
		assert.Equal(t, "error detail", actualError["detail"])
	})

	t.Run("should handle nil data and error gracefully", func(t *testing.T) {
		// 🔧 Setup: Create response recorder
		rr := httptest.NewRecorder()

		// 🚀 Action: Write response with nil values
		WriteJSON(rr, http.StatusOK, true, "Message", nil, nil)

		// ✅ Assertions: Should not panic and write valid JSON
		assert.Equal(t, http.StatusOK, rr.Code)

		var response APIResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response.Success)
		assert.Equal(t, "Message", response.Message)
	})
}

// TestOK tests the OK helper function
func TestOK(t *testing.T) {
	// 🎯 Test Strategy: Test the OK helper function

	t.Run("should write OK response correctly", func(t *testing.T) {
		// 🔧 Setup: Create response recorder
		rr := httptest.NewRecorder()

		// 🚀 Action: Write OK response
		OK(rr, "User found", map[string]string{"username": "john_doe"})

		// ✅ Assertions: Verify response
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		var response APIResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, response.Status)
		assert.True(t, response.Success)
		assert.Equal(t, "User found", response.Message)

		// Data becomes map[string]interface{} when unmarshaled
		actualData := response.Data.(map[string]interface{})
		assert.Equal(t, "john_doe", actualData["username"])
		assert.Nil(t, response.Error)
	})
}

// TestCreated tests the Created helper function
func TestCreated(t *testing.T) {
	// 🎯 Test Strategy: Test the Created helper function

	t.Run("should write Created response correctly", func(t *testing.T) {
		// 🔧 Setup: Create response recorder
		rr := httptest.NewRecorder()

		// 🚀 Action: Write Created response
		Created(rr, "User registered", map[string]int{"id": 123})

		// ✅ Assertions: Verify response
		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		var response APIResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, response.Status)
		assert.True(t, response.Success)
		assert.Equal(t, "User registered", response.Message)

		// Data becomes map[string]interface{} when unmarshaled
		// Numbers become float64 in JSON
		actualData := response.Data.(map[string]interface{})
		assert.Equal(t, float64(123), actualData["id"])
		assert.Nil(t, response.Error)
	})
}

// TestError tests the Error helper function
func TestError(t *testing.T) {
	// 🎯 Test Strategy: Test the Error helper function with different scenarios

	t.Run("should write error response with error detail", func(t *testing.T) {
		// 🔧 Setup: Create response recorder
		rr := httptest.NewRecorder()

		// Create a test error
		testError := assert.AnError

		// 🚀 Action: Write error response
		Error(rr, http.StatusBadRequest, "Validation failed", testError)

		// ✅ Assertions: Verify response
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		var response APIResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, response.Status)
		assert.False(t, response.Success)
		assert.Equal(t, "Validation failed", response.Message)
		assert.Nil(t, response.Data)

		// Verify error payload - Error becomes map[string]interface{} when unmarshaled
		errorPayload, ok := response.Error.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "Validation failed", errorPayload["message"])
		assert.Equal(t, testError.Error(), errorPayload["detail"])
	})

	t.Run("should write error response without error detail", func(t *testing.T) {
		// 🔧 Setup: Create response recorder
		rr := httptest.NewRecorder()

		// 🚀 Action: Write error response without error
		Error(rr, http.StatusInternalServerError, "Server error", nil)

		// ✅ Assertions: Verify response
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		var response APIResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.False(t, response.Success)
		assert.Equal(t, "Server error", response.Message)

		// Verify error payload - Error becomes map[string]interface{} when unmarshaled
		errorPayload, ok := response.Error.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "Server error", errorPayload["message"])
		// When err is nil, the "detail" field is not set in the payload
		_, hasDetail := errorPayload["detail"]
		assert.False(t, hasDetail, "Detail field should not be present when err is nil")
	})

	t.Run("should handle different HTTP status codes", func(t *testing.T) {
		// 🔧 Setup: Test different status codes
		testCases := []struct {
			statusCode int
			statusText string
		}{
			{http.StatusBadRequest, "Bad Request"},
			{http.StatusUnauthorized, "Unauthorized"},
			{http.StatusForbidden, "Forbidden"},
			{http.StatusNotFound, "Not Found"},
			{http.StatusInternalServerError, "Internal Server Error"},
		}

		for _, tc := range testCases {
			t.Run(tc.statusText, func(t *testing.T) {
				// 🔧 Setup: Create response recorder
				rr := httptest.NewRecorder()

				// 🚀 Action: Write error response
				Error(rr, tc.statusCode, "Error message", nil)

				// ✅ Assertions: Verify status code
				assert.Equal(t, tc.statusCode, rr.Code)

				var response APIResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Equal(t, tc.statusCode, response.Status)
				assert.False(t, response.Success)
			})
		}
	})
}

// TestResponseHeaders tests response headers
func TestResponseHeaders(t *testing.T) {
	// 🎯 Test Strategy: Test that proper headers are set

	t.Run("should set correct content type header", func(t *testing.T) {
		// 🔧 Setup: Create response recorder
		rr := httptest.NewRecorder()

		// 🚀 Action: Write any response
		OK(rr, "Test", nil)

		// ✅ Assertions: Verify headers
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	})

	t.Run("should not override existing headers", func(t *testing.T) {
		// 🔧 Setup: Create response recorder with existing headers
		rr := httptest.NewRecorder()
		rr.Header().Set("X-Custom-Header", "custom-value")
		rr.Header().Set("Content-Type", "text/plain")

		// 🚀 Action: Write JSON response
		OK(rr, "Test", nil)

		// ✅ Assertions: Verify headers
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type")) // Should override
		assert.Equal(t, "custom-value", rr.Header().Get("X-Custom-Header"))  // Should preserve
	})
}

// TestResponseTimestamp tests timestamp handling
func TestResponseTimestamp(t *testing.T) {
	// 🎯 Test Strategy: Test that timestamps are set correctly

	t.Run("should set timestamp in response", func(t *testing.T) {
		// 🔧 Setup: Create response recorder
		rr := httptest.NewRecorder()

		// 🚀 Action: Write response
		OK(rr, "Test", nil)

		// ✅ Assertions: Verify timestamp
		var response APIResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		// Timestamp should be recent (within last second)
		now := time.Now().UTC()
		diff := now.Sub(response.Timestamp)
		assert.True(t, diff >= 0 && diff < time.Second, "Timestamp should be recent")
	})
}

// 🎓 **LEARNING POINTS FROM THESE TESTS:**

// 1. **HTTP Testing**: Use httptest.NewRecorder() to test HTTP handlers
// 2. **JSON Marshaling**: Test both marshaling and unmarshaling
// 3. **Header Testing**: Verify HTTP headers are set correctly
// 4. **Status Code Testing**: Test different HTTP status codes
// 5. **Error Handling**: Test error scenarios and edge cases
// 6. **Timestamp Testing**: Test time-based functionality
// 7. **Response Structure**: Verify response format and content

// 💡 **TESTING BEST PRACTICES:**
// - Use httptest for HTTP testing (no real server needed)
// - Test both success and error scenarios
// - Verify HTTP status codes and headers
// - Test JSON serialization/deserialization
// - Test edge cases (nil values, empty data)
// - Use table-driven tests for similar scenarios
// - Verify response structure and content

// 🔧 **TOOLS USED:**
// - httptest: HTTP testing utilities
// - testify/assert: Cleaner assertions
// - testify/require: Fail fast assertions
// - encoding/json: JSON marshaling testing
