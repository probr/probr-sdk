package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestTmpDir(t *testing.T) {
	// Testing "create" and "delete" together to avoid making a mess
	tests := []struct {
		name string
		path string
	}{
		{
			name: "Validate SetTmpDir (1)",
			path: "TESTtmp1",
		},
		{
			name: "Validate SetTmpDir (2)",
			path: "TEST_tmp_2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetTmpDir(tt.path)
			if GlobalConfig.TmpDir != tt.path {
				t.Errorf("Expected GlobalConfig.TmpDir to be '%s' but found'%s'", tt.path, GlobalConfig.TmpDir)
			}
			fileInfo, err := os.Stat(tt.path)
			if err != nil {
				t.Errorf("Expected '%s' to be created, but found error: %s", tt.path, err)
			} else if !fileInfo.IsDir() {
				t.Errorf("Expected '%s' to be a directory, but it was a regular file", tt.path)
			}
			GlobalConfig.CleanupTmp()
			_, err = os.Stat(tt.path)
			if err == nil {
				t.Errorf("Expected file to be removed by GlobalConfig.CleanupTmp, but it was not: %s", tt.path)
			}
		})
	}
}

func TestGlobalOpts_OutputDir(t *testing.T) {
	time1 := time.Now()
	time2 := time.Now().Local().AddDate(1, 2, 3)
	tests := []struct {
		name string
		time time.Time
		base string
	}{
		{
			name: "Test OutputDir",
			time: time1,
			base: filepath.Join("imaginary", "dir"),
		},
		{
			name: "Test OutputDir",
			time: time2,
			base: filepath.Join("other", "imaginary", "dir"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			year, month, day := tt.time.Date()
			hour, min, sec := tt.time.Clock()
			gc := GlobalOpts{
				InstallDir: tt.base,
				StartTime:  tt.time,
			}
			got := gc.OutputDir()
			if !strings.Contains(got, tt.base) {
				t.Errorf("Expected output to contain '%s' but found '%s'", tt.base, got)
			}
			for _, time := range []interface{}{year, day, month, hour, min, sec} {
				if !strings.Contains(got, fmt.Sprint(time)) {
					t.Errorf("Expected output to contain '%d' but found '%s'", time, got)
				}
			}
		})
	}
}
