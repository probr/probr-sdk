package config

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestGlobalOpts_OutputDir(t *testing.T) {
	tests := []struct {
		name string
		base string
	}{
		{
			name: "Test outputDir",
			base: filepath.Join("imaginary", "dir"),
		},
		{
			name: "Test outputDir",
			base: filepath.Join("other", "imaginary", "dir"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gc := GlobalOpts{
				InstallDir: tt.base,
			}
			got := gc.outputDir()
			if !strings.Contains(got, tt.base) {
				t.Errorf("Expected output to contain '%s' but found '%s'", tt.base, got)
			}
		})
	}
}
