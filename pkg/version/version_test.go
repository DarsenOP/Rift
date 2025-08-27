package version

import "testing"

func TestVersionIsSet(t *testing.T) {
	if Version == "" {
		t.Fatal("Version constant should be set")
	}
}

func TestValidateVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{"valid dev version", "1.2.3-dev", false},
		{"valid release version", "2.0.0-release", false},
		{"invalid missing suffix", "1.2.3", true},
		{"invalid wrong suffix", "1.2.3-alpha", true},
		{"invalid format", "1.2-dev", true},
		{"invalid characters", "a.b.c-dev", true},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVersion(%q) error = %v, wantErr %v", tt.version, err, tt.wantErr)
			}
		})
	}
}

func TestCurrentVersionFormat(t *testing.T) {
	// Test that our actual Version constant is properly formatted
	if err := ValidateVersion(Version); err != nil {
		t.Errorf("Current Version constant %q is invalid: %v", Version, err)
	}
}
