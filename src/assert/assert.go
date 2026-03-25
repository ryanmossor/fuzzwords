package assert

import (
	"fmt"
	"log/slog"
)

func Assert(truth bool, msg string, v ...any) {
	if !truth {
		slog.Error(fmt.Sprintf("[assert] %s", msg), v...)
		panic(msg)
	}
}
