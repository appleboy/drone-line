package random

import "testing"

func TestStringWithCharset(t *testing.T) {
	type args struct {
		length  int
		charset Charset
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Rand",
			args: args{
				length:  5,
				charset: "a",
			},
			want: "aaaaa",
		},
	}
	for _, tt := range tests {
		if got := StringWithCharset(tt.args.length, tt.args.charset); got != tt.want {
			t.Errorf("%q. StringWithCharset() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestString(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test Rand",
			args: args{
				length: 5,
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		if got := String(tt.args.length); len(got) != tt.want {
			t.Errorf("%q. String() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
