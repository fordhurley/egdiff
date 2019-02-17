package main

import (
	"bytes"
	"strconv"
	"testing"
)

func TestNthFromLastIndex(t *testing.T) {
	tts := []struct {
		data     string
		sep      string
		n        int
		expected int
	}{
		{
			"hello",
			"l",
			0,
			3,
		},
		{
			"hello",
			"l",
			1,
			2,
		},
		{
			"hello\nmy\ndarling\n",
			"\n",
			2,
			5,
		},
	}

	for i, tt := range tts {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			index := NthFromLastIndex([]byte(tt.data), []byte(tt.sep), tt.n)
			if index != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, index)
			}

			lastIndexReal := bytes.LastIndex([]byte(tt.data), []byte(tt.sep))
			lastIndex := NthFromLastIndex([]byte(tt.data), []byte(tt.sep), 0)
			if lastIndex != lastIndexReal {
				t.Errorf("expected last index to be %d, got %d", lastIndexReal, lastIndex)
			}
		})
	}
}
