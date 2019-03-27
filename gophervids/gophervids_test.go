package gophervids

import (
	"testing"
)

func TestVideo_URL(t *testing.T) {
	type fields struct {
		ID    string
		Date  string
		Title string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Testing URL", fields{"test01", "date", "TestingURL"}, "https://www.youtube.com/watch?v=test01"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Video{
				ID:    tt.fields.ID,
				Date:  tt.fields.Date,
				Title: tt.fields.Title,
			}
			if got := v.URL(); got != tt.want {
				t.Errorf("Video.URL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVideo_Filename(t *testing.T) {
	type fields struct {
		ID    string
		Date  string
		Title string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"with blanks", fields{"", "", "Testing Title"}, "testing-title"},
		{"with non-alpha chars", fields{"", "", "Testing!@# Title"}, "testing-title"},
		{"with extra spaces", fields{"", "", " Testing!@# Title     "}, "testing-title"},
		{"with random casing", fields{"", "", " TeStiNg!@# TitlE%@ !#>     "}, "testing-title"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Video{
				ID:    tt.fields.ID,
				Date:  tt.fields.Date,
				Title: tt.fields.Title,
			}
			if got := v.Filename(); got != tt.want {
				t.Errorf("Video.Filename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVideo_FullPath(t *testing.T) {
	type fields struct {
		ID    string
		Date  string
		Title string
	}
	type args struct {
		p string
		a string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"With Author", fields{"", "01-01-2001", "testing"}, args{"output", "author"}, "output/author/01-01-2001-testing"},
		{"Sanitize Author", fields{"", "01-01-2001", "testing"}, args{"output", " !@#$ #!@#Author "}, "output/author/01-01-2001-testing"},
		{"Missing Author", fields{"", "01-01-2001", "testing"}, args{"output", ""}, "output/01-01-2001-testing"},
		{"Missing Date", fields{"", "", "testing"}, args{"output", "author"}, "output/author/00-testing"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Video{
				ID:    tt.fields.ID,
				Date:  tt.fields.Date,
				Title: tt.fields.Title,
			}
			if got := v.FullPath(tt.args.p, tt.args.a); got != tt.want {
				t.Errorf("Video.FullPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
