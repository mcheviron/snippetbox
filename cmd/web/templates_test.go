package main

import (
	"snippetbox/internal/assert"
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		args time.Time
		want string
	}{
		// TODO: Add test cases.
		{
			name: "UTC",
			args: time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2022 at 10:15",
		},
		{
			name: "Empty",
			args: time.Time{},
			want: "",
		},
		{
			name: "CET",
			args: time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2022 at 09:15",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := humanDate(tt.args)
			assert.Equal(t, got, tt.want)
		})
	}
}
