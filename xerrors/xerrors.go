package xerrors

import "fmt"

func Wrap(err *error, format string, args ...interface{}) {
	if err == nil {
		return
	}

	*err = fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), *err)
}
