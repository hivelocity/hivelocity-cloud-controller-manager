// Package errutils provides some helper functions around error handling.
package errutils

import (
	"fmt"
	"runtime"
	"strings"
)

// Wrap takes an error, a message an optionally pairs of variable names and variables.
// Wrap returns a wrapped error with a unified format.
// Example: errutils.Wrap(err, "no space left on device", "requestedSize", requestedSize).
// returns err wrapped, with this messages "[callerMethod] no space left on device. requestedSize 1234: %w".
func Wrap(err error, msg string, args ...any) error {
	if len(args)%2 != 0 {
		args = append(args, "#### missing argument in call to Wrap()###")
	}
	pc, _, _, _ := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	funcName := "unknownMethod"
	if details != nil {
		funcName = details.Name()

		// split github.com/../errutils.TestWrap.func1 to "errutils.TestWrap.func1"
		funcName = funcName[strings.LastIndex(funcName, "/")+1:]
	}
	nameValuePairs := make([]string, 0, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		varName := args[i]
		value := args[i+1]
		nameValuePairs = append(nameValuePairs, fmt.Sprintf("%s %s", varName, fmt.Sprint(value)))
	}
	return fmt.Errorf("[%s] %s. %s: %w", funcName, msg, strings.Join(nameValuePairs, ", "), err)
}
