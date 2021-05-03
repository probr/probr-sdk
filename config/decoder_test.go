package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewConfigDecoder(t *testing.T) {
	realDir, _ := os.Getwd()
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "Validate that a directory for config path results in an error",
			path:    realDir,
			wantErr: true,
		},
		{
			name:    "Validate that a directory for config path results in an error",
			path:    filepath.Join(realDir, "decoder_test.go"), // targeting this file should be sufficient
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, file, err := NewConfigDecoder(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfigDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			file.Close()
		})
	}
}
