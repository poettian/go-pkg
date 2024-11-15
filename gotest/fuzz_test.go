package gotest

import (
	"testing"
	"unicode/utf8"
)

func TestReverse(t *testing.T) {
	testcases := []struct {
		in, want string
	}{
		{in: "Hello, world", want: "dlrow ,olleH"},
		{in: " ", want: " "},
		{in: "!12345", want: "54321!"},
	}
	for _, tc := range testcases {
		rev, err := Reverse(tc.in)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if rev != tc.want {
			t.Errorf("Reverse: %q, want %q", rev, tc.want)
		}
	}
}

func FuzzReverse(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!123456"}
	for _, tc := range testcases {
		f.Add(tc) // 输入测试种子
	}
	f.Fuzz(func(t *testing.T, orig string) {
		if !utf8.ValidString(orig) {
			return
		}
		rev, err := ReverseV2(orig)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !utf8.ValidString(rev) {
			t.Fatalf("Reverse produced invalid UTF-8 string %q", rev)
		}
		doubleRev, err := ReverseV2(rev)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if orig != doubleRev {
			t.Errorf("Before:%q, After:%q", orig, doubleRev)
		}
	})
}
