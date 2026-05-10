package logic

import "testing"

func TestNormalizeResourceApiPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "colon param unchanged", path: "/api/users/:id", want: "/api/users/:id"},
		{name: "brace param normalized", path: "/api/users/{id}", want: "/api/users/:id"},
		{name: "multiple brace params normalized", path: "/api/users/{userId}/orders/{orderId}", want: "/api/users/:userId/orders/:orderId"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeResourceApiPath(tt.path); got != tt.want {
				t.Fatalf("normalizeResourceApiPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValidateResourceApiPathTemplate(t *testing.T) {
	validPaths := []string{
		"/api/users/:id",
		"/api/users/:userId/orders/:orderId",
		"/api/users/{id}",
	}
	for _, path := range validPaths {
		t.Run("valid "+path, func(t *testing.T) {
			if err := validateResourceApiPathTemplate(normalizeResourceApiPath(path)); err != nil {
				t.Fatalf("validateResourceApiPathTemplate(%q) returned error: %v", path, err)
			}
		})
	}

	invalidPaths := []string{
		"/api/users/:",
		"/api/users/:1id",
		"/api/users/:user-id",
		"/api/users/{bad-name}",
		"/api/users/{id",
	}
	for _, path := range invalidPaths {
		t.Run("invalid "+path, func(t *testing.T) {
			if err := validateResourceApiPathTemplate(normalizeResourceApiPath(path)); err == nil {
				t.Fatalf("validateResourceApiPathTemplate(%q) returned nil", path)
			}
		})
	}
}
