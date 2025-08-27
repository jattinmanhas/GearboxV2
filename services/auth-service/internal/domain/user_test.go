package domain

import (
	"encoding/json"
	"testing"
	"time"
)

// TestUserStruct tests the User struct and its behavior
func TestUserStruct(t *testing.T) {
	// 🎯 Test Strategy: Test struct creation, field access, and JSON marshaling

	t.Run("should create user with all fields", func(t *testing.T) {
		// 🔧 Setup: Create a user with all fields populated
		now := time.Now().UTC()
		user := &User{
			ID:          1,
			Username:    "john_doe",
			Password:    "hashedpassword123",
			Email:       "john@example.com",
			FirstName:   "John",
			MiddleName:  "Michael",
			LastName:    "Doe",
			Avatar:      "https://example.com/avatar.jpg",
			Gender:      "male",
			DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			CreatedAt:   now,
			UpdatedAt:   now,
			IsDeleted:   false,
		}

		// ✅ Assertions: Verify all fields are set correctly
		if user.ID != 1 {
			t.Errorf("Expected ID 1, got %d", user.ID)
		}

		if user.Username != "john_doe" {
			t.Errorf("Expected username 'john_doe', got %s", user.Username)
		}

		if user.Password != "hashedpassword123" {
			t.Errorf("Expected password 'hashedpassword123', got %s", user.Password)
		}

		if user.Email != "john@example.com" {
			t.Errorf("Expected email 'john@example.com', got %s", user.Email)
		}

		if user.FirstName != "John" {
			t.Errorf("Expected first name 'John', got %s", user.FirstName)
		}

		if user.MiddleName != "Michael" {
			t.Errorf("Expected middle name 'Michael', got %s", user.MiddleName)
		}

		if user.LastName != "Doe" {
			t.Errorf("Expected last name 'Doe', got %s", user.LastName)
		}

		if user.Avatar != "https://example.com/avatar.jpg" {
			t.Errorf("Expected avatar 'https://example.com/avatar.jpg', got %s", user.Avatar)
		}

		if user.Gender != "male" {
			t.Errorf("Expected gender 'male', got %s", user.Gender)
		}

		if user.DateOfBirth.Year() != 1990 {
			t.Errorf("Expected year of birth 1990, got %d", user.DateOfBirth.Year())
		}

		if user.IsDeleted != false {
			t.Errorf("Expected IsDeleted false, got %t", user.IsDeleted)
		}
	})

	t.Run("should create user with minimal required fields", func(t *testing.T) {
		// 🔧 Setup: Create a user with only required fields
		user := &User{
			Username: "jane_doe",
			Email:    "jane@example.com",
		}

		// ✅ Assertions: Verify required fields are set
		if user.Username != "jane_doe" {
			t.Errorf("Expected username 'jane_doe', got %s", user.Username)
		}

		if user.Email != "jane@example.com" {
			t.Errorf("Expected email 'jane@example.com', got %s", user.Email)
		}

		// Verify zero values for unset fields
		if user.ID != 0 {
			t.Errorf("Expected ID 0 for unset field, got %d", user.ID)
		}

		if user.Password != "" {
			t.Errorf("Expected empty password for unset field, got %s", user.Password)
		}

		if !user.CreatedAt.IsZero() {
			t.Errorf("Expected zero time for unset CreatedAt, got %v", user.CreatedAt)
		}
	})
}

// TestUserJSONMarshaling tests JSON serialization and deserialization
func TestUserJSONMarshaling(t *testing.T) {
	// 🎯 Test Strategy: Test JSON struct tags and marshaling behavior

	t.Run("should marshal user to JSON correctly", func(t *testing.T) {
		// 🔧 Setup: Create a user with all fields
		user := &User{
			ID:          1,
			Username:    "john_doe",
			Password:    "secretpassword", // This should be hidden in JSON
			Email:       "john@example.com",
			FirstName:   "John",
			LastName:    "Doe",
			DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			IsDeleted:   false,
		}

		// 🚀 Action: Marshal to JSON
		jsonData, err := json.Marshal(user)
		if err != nil {
			t.Fatalf("Failed to marshal user to JSON: %v", err)
		}

		// ✅ Assertions: Verify JSON structure
		jsonString := string(jsonData)

		// Password should NOT be in JSON (json:"-" tag)
		if contains(jsonString, "secretpassword") {
			t.Error("Password should not appear in JSON output")
		}

		// Other fields should be present
		if !contains(jsonString, "john_doe") {
			t.Error("Username should appear in JSON output")
		}

		if !contains(jsonString, "john@example.com") {
			t.Error("Email should appear in JSON output")
		}

		if !contains(jsonString, "John") {
			t.Error("FirstName should appear in JSON output")
		}

		if !contains(jsonString, "Doe") {
			t.Error("LastName should appear in JSON output")
		}

		// Verify JSON field names (snake_case)
		if !contains(jsonString, "first_name") {
			t.Error("JSON should use snake_case field names")
		}

		if !contains(jsonString, "last_name") {
			t.Error("JSON should use snake_case field names")
		}

		if !contains(jsonString, "date_of_birth") {
			t.Error("JSON should use snake_case field names")
		}
	})

	t.Run("should unmarshal JSON to user correctly", func(t *testing.T) {
		// 🔧 Setup: JSON data representing a user
		jsonData := `{
			"id": 2,
			"username": "jane_doe",
			"email": "jane@example.com",
			"first_name": "Jane",
			"last_name": "Smith",
			"gender": "female",
			"is_deleted": false
		}`

		// 🚀 Action: Unmarshal JSON to User struct
		var user User
		err := json.Unmarshal([]byte(jsonData), &user)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON to user: %v", err)
		}

		// ✅ Assertions: Verify fields are set correctly
		if user.ID != 2 {
			t.Errorf("Expected ID 2, got %d", user.ID)
		}

		if user.Username != "jane_doe" {
			t.Errorf("Expected username 'jane_doe', got %s", user.Username)
		}

		if user.Email != "jane@example.com" {
			t.Errorf("Expected email 'jane@example.com', got %s", user.Email)
		}

		if user.FirstName != "Jane" {
			t.Errorf("Expected first name 'Jane', got %s", user.FirstName)
		}

		if user.LastName != "Smith" {
			t.Errorf("Expected last name 'Smith', got %s", user.LastName)
		}

		if user.Gender != "female" {
			t.Errorf("Expected gender 'female', got %s", user.Gender)
		}

		if user.IsDeleted != false {
			t.Errorf("Expected IsDeleted false, got %t", user.IsDeleted)
		}

		// Fields not in JSON should have zero values
		if user.Password != "" {
			t.Errorf("Expected empty password for unmarshaled field, got %s", user.Password)
		}

		if user.MiddleName != "" {
			t.Errorf("Expected empty middle name for unmarshaled field, got %s", user.MiddleName)
		}

		if !user.CreatedAt.IsZero() {
			t.Errorf("Expected zero time for unmarshaled field, got %v", user.CreatedAt)
		}
	})
}

// TestUserStructTags tests the struct tags for database and JSON
func TestUserStructTags(t *testing.T) {
	// 🎯 Test Strategy: Test that struct tags are correctly defined

	t.Run("should have correct JSON tags", func(t *testing.T) {
		// 🔧 Setup: Create a user
		user := &User{
			Username: "test_user",
			Email:    "test@example.com",
		}

		// 🚀 Action: Marshal to JSON
		jsonData, err := json.Marshal(user)
		if err != nil {
			t.Fatalf("Failed to marshal user: %v", err)
		}

		jsonString := string(jsonData)

		// ✅ Assertions: Verify JSON field names
		expectedFields := map[string]string{
			"username": "test_user",
			"email":    "test@example.com",
		}

		for fieldName, expectedValue := range expectedFields {
			if !contains(jsonString, fieldName) {
				t.Errorf("JSON should contain field '%s'", fieldName)
			}

			if !contains(jsonString, expectedValue) {
				t.Errorf("JSON should contain value '%s' for field '%s'", expectedValue, fieldName)
			}
		}
	})

	t.Run("should hide password in JSON", func(t *testing.T) {
		// 🔧 Setup: Create a user with password
		user := &User{
			Username: "test_user",
			Password: "secret123",
			Email:    "test@example.com",
		}

		// 🚀 Action: Marshal to JSON
		jsonData, err := json.Marshal(user)
		if err != nil {
			t.Fatalf("Failed to marshal user: %v", err)
		}

		jsonString := string(jsonData)

		// ✅ Assertions: Password should not appear in JSON
		if contains(jsonString, "secret123") {
			t.Error("Password should be hidden in JSON output")
		}

		if contains(jsonString, "password") {
			t.Error("Password field name should not appear in JSON output")
		}
	})
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 1; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}())))
}

// 🎓 **LEARNING POINTS FROM THESE TESTS:**

// 1. **Struct Testing**: Test both creation and field access
// 2. **JSON Marshaling**: Test struct tags and serialization behavior
// 3. **Zero Values**: Understand Go's zero value behavior
// 4. **Struct Tags**: Test json:"-" for hiding sensitive fields
// 5. **Field Mapping**: Test JSON field name mapping (snake_case)
// 6. **Unmarshaling**: Test JSON to struct conversion
// 7. **Security**: Ensure sensitive data (password) is hidden

// 💡 **TESTING BEST PRACTICES:**
// - Test both happy path and edge cases
// - Verify struct tags work as expected
// - Test JSON serialization/deserialization
// - Ensure sensitive data is properly hidden
// - Test zero values and default behavior
// - Use descriptive test names and assertions
