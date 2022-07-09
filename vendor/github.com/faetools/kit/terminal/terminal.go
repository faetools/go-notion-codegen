package terminal

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

// Transform transforms the terminal output in some way.
type Transform func(interface{}) aurora.Value

// Println prints something in color or after having gone through some other auora transformation.
func Println(transform Transform, a ...interface{}) (int, error) {
	return fmt.Fprint(os.Stdout, transform(fmt.Sprintln(a...)))
}

// Print prints something in color or after having gone through some other aurora transformation.
func Print(transform Transform, a ...interface{}) (int, error) {
	return fmt.Fprint(os.Stdout, transform(fmt.Sprint(a...)))
}

// Printf prints something in color or after having gone through some other aurora transformation.
func Printf(transform Transform, format string, a ...interface{}) (int, error) {
	return fmt.Fprint(os.Stdout, transform(fmt.Sprintf(format, a...)))
}
