package binary

import "testing"

func TestPathWasReturned(t *testing.T) {
	tests := []struct {
		name string
		line string
		want bool
	}{
		{
			name: "which no result",
			line: "which: no ahfafka in (/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin)",
			want: false,
		},
		{
			name: "absolute path result",
			line: "/usr/bin/base64",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PathWasReturned(tt.line); got != tt.want {
				t.Fatalf("PathWasReturned(%q) = %t, want %t", tt.line, got, tt.want)
			}
		})
	}
}
