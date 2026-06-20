package util

import (
	"fmt"
	"runtime"
	"strings"
)

// GetFunctionName returns the name of the calling function.
// skip: number of stack frames to skip (0 = current function, 1 = caller, etc.)
func GetFunctionName(skip int) string {
	pc, _, _, _ := runtime.Caller(skip + 1)
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "[Unknown]"
	}

	fullName := fn.Name()
	parts := strings.Split(fullName, ".")
	if len(parts) > 0 {
		return fmt.Sprintf("[%s]", parts[len(parts)-1])
	}
	return "[Unknown]"
}
