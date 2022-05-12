package strutil

import (
	"fmt"
	"strings"
)

func JoinGIT(path ...string) string {
	return JoinPath("", ".git", path...)
}

func JoinImage(version string, path ...string) string {
	return JoinPath("", fmt.Sprintf(":%s", version), path...)
}

func JoinPath(prefix, suffix string, path ...string) string {
	p := path[:0]
	for _, s := range path {
		p = append(p, clean(s))
	}
	s := strings.Join(path, "/")
	return fmt.Sprintf("%s%s%s", prefix, s, suffix)
}

func GenName(s ...string) string {
	for i, v := range s {
		s[i] = strings.ToLower(v)
	}
	return strings.Join(s, "-")
}

func clean(path string) string {
	for i, c := range path {
		if c != '/' {
			path = path[i:]
			break
		}
	}

	for i := len(path) - 1; i >= 0; i-- {
		if path[i] != '/' {
			path = path[:i+1]
			break
		}
	}
	return path
}

func Reverse(s string, sep ...string) string {
	defaultSep := "-"
	if len(sep) != 0 {
		defaultSep = sep[0]
	}
	set := strings.Split(s, defaultSep)
	for i := 0; i < len(set)/2; i++ {
		temp := set[i]
		set[i] = set[len(set)-1-i]
		set[len(set)-1-i] = temp
	}
	return strings.Join(set, defaultSep)
}
