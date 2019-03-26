package utils

import "testing"

func TestSanitize(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{"Front Trim", " test", "test"},
		{"Back Trim", "test ", "test"},
		{"Remove non-alphnumeric", "!$!%!@#!#test!$!%!$>?<!@ !@#", "test"},
		{"Replace spaces between words with '-'", "test test", "test-test"},
		{"Replace spaces between words with '-' including special chars", "test!@#$#%$ !@#!@#test", "test-test"},
		{"LowerCase everthing", "TestTestTesT", "testtesttest"},
		{"LowerCase everthing with '-'", "Test Test TesT", "test-test-test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sanitize(tt.arg); got != tt.want {
				t.Errorf("Sanitize() = %v, want %v", got, tt.want)
			}
		})
	}
}
