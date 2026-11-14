// Package modelid validates declared builder and reviewer model attribution.
package modelid

import "regexp"

const Pattern = `^[a-z0-9][a-z0-9._-]{0,63}/[a-z0-9][a-z0-9._-]{0,63}$`

var pattern = regexp.MustCompile(Pattern)

func Valid(model string) bool { return pattern.MatchString(model) }
