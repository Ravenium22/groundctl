package recovery

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// Wrap returns a function that recovers from panics and prints a friendly error.
// Use as: defer recovery.Wrap("command-name")()
func Wrap(context string) func() {
	return func() {
		if r := recover(); r != nil {
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("\ngroundctl encountered an unexpected error in '%s'.\n", context))
			sb.WriteString(fmt.Sprintf("Error: %v\n\n", r))
			sb.WriteString("Stack trace:\n")

			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			sb.Write(buf[:n])

			sb.WriteString("\n\nPlease report this at: https://github.com/groundctl/groundctl/issues\n")
			sb.WriteString("Include the above output and the output of 'ground doctor'.\n")

			fmt.Fprint(os.Stderr, sb.String())
			os.Exit(3)
		}
	}
}
