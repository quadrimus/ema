package file

import (
	"regexp"
)

// IsCommandName checks if string is valid command name.
//
// Command name regular expression:
//
//	^[A-Za-z][-0-9A-Za-z]*$
func IsCommandName(s string) bool {
	return commandNameRegexp.MatchString(s)
}

var commandNameRegexp = regexp.MustCompile(`^[A-Za-z][-0-9A-Za-z]*$`)
