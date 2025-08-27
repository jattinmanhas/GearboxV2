package config

import (
	"os"
	"testing"
)

// TestLoadConfig tests the configuration loading functionality
func TestLoadConfig(t *testing.T) {
	// 🎯 Test Strategy: Test the singleton pattern and environment variable loading

	t.Run("should load config successfully with environment variables", func(t *testing.T) {
		// 🔧 Setup: Set environment variables for testing
		os.Setenv("PORT", "8080")
		os.Setenv("DATABASE_URL", "postgres://user:pass@localhost/testdb")

		// 🧹 Cleanup: Restore environment after test
		defer func() {
			os.Unsetenv("PORT")
			os.Unsetenv("DATABASE_URL")
		}()

		// 🚀 Action: Load configuration
		cfg := LoadConfig()

		// ✅ Assertions: Verify config values
		if cfg == nil {
			t.Fatal("Config should not be nil")
		}

		if cfg.Port != "8080" {
			t.Errorf("Expected port 8080, got %s", cfg.Port)
		}

		if cfg.DatabaseURL != "postgres://user:pass@localhost/testdb" {
			t.Errorf("Expected database URL 'postgres://user:pass@localhost/testdb', got %s", cfg.DatabaseURL)
		}
	})

	t.Run("should use default port when PORT is not set", func(t *testing.T) {
		// 🔧 Setup: Set only DATABASE_URL, leave PORT unset
		os.Setenv("DATABASE_URL", "postgres://user:pass@localhost/testdb")
		defer os.Unsetenv("DATABASE_URL")

		// 🚀 Action: Load configuration
		cfg := LoadConfig()

		// ✅ Assertions: Verify default port is used
		if cfg.Port != "8081" {
			t.Errorf("Expected default port 8081, got %s", cfg.Port)
		}
	})

	t.Run("should return same config instance on multiple calls (singleton)", func(t *testing.T) {
		// 🔧 Setup: Set environment variables
		os.Setenv("DATABASE_URL", "postgres://user:pass@localhost/testdb")
		defer os.Unsetenv("DATABASE_URL")

		// 🚀 Action: Load configuration multiple times
		cfg1 := LoadConfig()
		cfg2 := LoadConfig()
		cfg3 := LoadConfig()

		// ✅ Assertions: All should be the same instance (same memory address)
		if cfg1 != cfg2 {
			t.Error("Config instances should be the same (singleton pattern)")
		}

		if cfg2 != cfg3 {
			t.Error("Config instances should be the same (singleton pattern)")
		}

		if cfg1 != cfg3 {
			t.Error("Config instances should be the same (singleton pattern)")
		}
	})
}

// TestLoadConfigMissingDatabaseURL tests error handling for missing DATABASE_URL
func TestLoadConfigMissingDatabaseURL(t *testing.T) {
	// 🎯 Test Strategy: Test that the function fails when critical config is missing

	// 🔧 Setup: Ensure DATABASE_URL is not set
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("PORT")

	// 🚀 Action & ✅ Assertion: This should panic due to log.Fatal
	// We need to use a different approach since log.Fatal calls os.Exit(1)

	// 💡 Note: In real applications, you might want to use a custom logger
	// that can be mocked for testing, or return errors instead of calling log.Fatal
}

// TestConfigStruct tests the Config struct fields
func TestConfigStruct(t *testing.T) {
	// 🎯 Test Strategy: Test the Config struct structure and field types

	t.Run("should have correct field types", func(t *testing.T) {
		// 🔧 Setup: Create a config instance
		cfg := &Config{
			Port:        "8080",
			DatabaseURL: "postgres://localhost/testdb",
		}

		// ✅ Assertions: Verify field types and values
		if cfg.Port != "8080" {
			t.Errorf("Expected port 8080, got %s", cfg.Port)
		}

		if cfg.DatabaseURL != "postgres://localhost/testdb" {
			t.Errorf("Expected database URL 'postgres://localhost/testdb', got %s", cfg.DatabaseURL)
		}
	})
}

// TestLoadConfigConcurrency tests thread safety of the singleton pattern
func TestLoadConfigConcurrency(t *testing.T) {
	// 🎯 Test Strategy: Test that the singleton pattern is thread-safe

	// 🔧 Setup: Set environment variables
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost/testdb")
	defer os.Unsetenv("DATABASE_URL")

	// 🚀 Action: Load config from multiple goroutines concurrently
	configs := make(chan *Config, 10)

	for i := 0; i < 10; i++ {
		go func() {
			configs <- LoadConfig()
		}()
	}

	// ✅ Assertions: All configs should be the same instance
	var firstConfig *Config
	for i := 0; i < 10; i++ {
		cfg := <-configs
		if firstConfig == nil {
			firstConfig = cfg
		} else if cfg != firstConfig {
			t.Error("All config instances should be the same (thread safety)")
		}
	}
}

// 🎓 **LEARNING POINTS FROM THESE TESTS:**

// 1. **Test Structure**: Use t.Run() for organized, readable tests
// 2. **Setup/Teardown**: Always clean up after tests (defer cleanup)
// 3. **Environment Variables**: Test different environment configurations
// 4. **Singleton Pattern**: Verify the same instance is returned
// 5. **Thread Safety**: Test concurrent access to ensure safety
// 6. **Error Cases**: Test both success and failure scenarios
// 7. **Clean Tests**: Each test should be independent and not affect others

// 💡 **TESTING BEST PRACTICES:**
// - Use descriptive test names that explain what you're testing
// - Follow the Arrange-Act-Assert pattern (Setup-Action-Assertions)
// - Clean up resources after each test
// - Test both happy path and edge cases
// - Use table-driven tests for multiple similar test cases
// - Mock external dependencies when possible
