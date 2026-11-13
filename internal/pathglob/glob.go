// Package pathglob matches root-anchored repository paths. Only * (within
// a folder) and ** (across folders) are operators; **/ also matches no folders.
package pathglob

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Glob struct{ expression *regexp.Regexp }

func Compile(pattern string) (Glob, error) {
	if !validPath(pattern) || strings.ContainsAny(pattern, `\?[]{}`) || strings.Contains(pattern, "***") {
		return Glob{}, fmt.Errorf("invalid repository glob %q: use a relative path with * or **", pattern)
	}
	var expression strings.Builder
	expression.WriteString("^")
	for index := 0; index < len(pattern); index++ {
		switch {
		case strings.HasPrefix(pattern[index:], "**/"):
			expression.WriteString("(?:[^/]+/)*")
			index += 2
		case strings.HasPrefix(pattern[index:], "**"):
			expression.WriteString(".*")
			index++
		case pattern[index] == '*':
			expression.WriteString("[^/]*")
		default:
			expression.WriteString(regexp.QuoteMeta(pattern[index : index+1]))
		}
	}
	expression.WriteString("$")
	compiled, err := regexp.Compile(expression.String())
	if err != nil {
		return Glob{}, fmt.Errorf("compile repository glob %q: %w", pattern, err)
	}
	return Glob{expression: compiled}, nil
}

func (glob Glob) Match(path string) bool {
	return glob.expression != nil && validPath(path) && glob.expression.MatchString(path)
}

func validPath(path string) bool {
	if !utf8.ValidString(path) || strings.ContainsRune(path, '\\') || strings.ContainsFunc(path, unicode.IsControl) {
		return false
	}
	for _, component := range strings.Split(path, "/") {
		if component == "" || component == "." || component == ".." {
			return false
		}
	}
	return true
}
