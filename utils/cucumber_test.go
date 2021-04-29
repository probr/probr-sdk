package utils

import (
	"testing"
)

// TestCucumberTagsListToString ...
func TestCucumberTagsListToString(t *testing.T) {
	tests := []struct {
		name string
		tags []string
		want string
	}{
		{
			name: "Ensure a single tag is converted properly",
			tags: []string{"tag1"},
			want: "@tag1",
		},
		{
			name: "Ensure a list of several tags is converted properly",
			tags: []string{"a1", "b1", "c2", "d3", "e5", "f8"},
			want: "@a1,@b1,@c2,@d3,@e5,@f8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CucumberTagsListToString(tt.tags); got != tt.want {
				t.Errorf("CucumberTagsListToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCucumberTagExclusionsListToString ...
func TestCucumberTagExclusionsListToString(t *testing.T) {
	tests := []struct {
		name string
		tags []string
		want string
	}{
		{
			name: "Ensure a single tag is converted properly",
			tags: []string{"tag1"},
			want: "~@tag1",
		},
		{
			name: "Ensure a list of several tags is converted properly",
			tags: []string{"a1", "b1", "c2", "d3", "e5", "f8"},
			want: "~@a1 && ~@b1 && ~@c2 && ~@d3 && ~@e5 && ~@f8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CucumberTagExclusionsListToString(tt.tags); got != tt.want {
				t.Errorf("ConfigTagExclusionsListToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
