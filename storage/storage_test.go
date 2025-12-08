package storage

import (
	"os"
	"testing"
)

func TestBuildPostgresDSN(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected string
	}{
		{
			name:     "default values",
			envVars:  map[string]string{},
			expected: "host=localhost port=5432 user=postgres password=postgres dbname=microservice_overview sslmode=disable",
		},
		{
			name: "custom values",
			envVars: map[string]string{
				"DB_HOST":     "db.example.com",
				"DB_PORT":     "5433",
				"DB_USER":     "admin",
				"DB_PASSWORD": "secret123",
				"DB_NAME":     "test_db",
			},
			expected: "host=db.example.com port=5433 user=admin password=secret123 dbname=test_db sslmode=disable",
		},
		{
			name: "partial override",
			envVars: map[string]string{
				"DB_HOST": "custom-host",
				"DB_PORT": "9999",
			},
			expected: "host=custom-host port=9999 user=postgres password=postgres dbname=microservice_overview sslmode=disable",
		},
		{
			name: "empty values use defaults",
			envVars: map[string]string{
				"DB_HOST":     "",
				"DB_PORT":     "",
				"DB_USER":     "",
				"DB_PASSWORD": "",
				"DB_NAME":     "",
			},
			expected: "host=localhost port=5432 user=postgres password=postgres dbname=microservice_overview sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Zapisz oryginalne wartości
			originalVars := make(map[string]string)
			for key := range tt.envVars {
				originalVars[key] = os.Getenv(key)
			}

			// Ustaw zmienne środowiskowe dla testu
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Wywołaj funkcję
			result := buildPostgresDSN()

			// Przywróć oryginalne wartości
			for key, value := range originalVars {
				if value == "" {
					os.Unsetenv(key)
				} else {
					os.Setenv(key, value)
				}
			}

			// Sprawdź wynik
			if result != tt.expected {
				t.Errorf("buildPostgresDSN() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		envValue     string
		defaultValue string
		expected     string
	}{
		{
			name:         "env var set",
			key:          "TEST_VAR",
			envValue:     "custom_value",
			defaultValue: "default",
			expected:     "custom_value",
		},
		{
			name:         "env var not set",
			key:          "TEST_VAR_NOT_SET",
			envValue:     "",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "env var empty string",
			key:          "TEST_VAR_EMPTY",
			envValue:     "",
			defaultValue: "default",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Zapisz oryginalną wartość
			originalValue := os.Getenv(tt.key)

			// Ustaw wartość dla testu
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			} else {
				os.Unsetenv(tt.key)
			}

			// Wywołaj funkcję
			result := getEnv(tt.key, tt.defaultValue)

			// Przywróć oryginalną wartość
			if originalValue == "" {
				os.Unsetenv(tt.key)
			} else {
				os.Setenv(tt.key, originalValue)
			}

			// Sprawdź wynik
			if result != tt.expected {
				t.Errorf("getEnv(%q, %q) = %q, want %q", tt.key, tt.defaultValue, result, tt.expected)
			}
		})
	}
}

