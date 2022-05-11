package strutil

import (
	"fmt"
	"testing"
)

func TestClen(t *testing.T) {
	cases := []string{
		"/case",
		"//case",
		"case//",
	}

	for _, c := range cases {
		s := clean(c)
		if s != "case" {
			fmt.Printf("error: result %s\n", s)
		}
	}
}

func TestReverse(t *testing.T) {
	cases := []string{
		"a-b-c",
		// "a/b/c",
		"a-b",
	}

	for _, c := range cases {
		s := Reverse(c, "-")
		fmt.Printf("%s\n", s)
	}
}
