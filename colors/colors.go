package colors

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Colors []string

// region Formatters

func (clrs *Colors) Fprintf(writer io.Writer, format string, args ...any) (n int, err error) {
	var buf strings.Builder
	for _, clr := range *clrs {
		buf.WriteString(clr)
	}

	buf.WriteString(fmt.Sprintf(format, args...))
	buf.WriteString(Default)
	return fmt.Fprint(writer, buf.String())
}

func (clrs *Colors) Printf(format string, args ...any) (n int, err error) {
	return clrs.Fprintf(os.Stdout, format, args...)
}

// endregion

// region Other printers

func (clrs *Colors) Fprint(writer io.Writer, args ...any) (n int, err error) {
	var buf strings.Builder
	buf.WriteString(strings.Join(*clrs, ""))
	buf.WriteString(fmt.Sprint(args...))
	buf.WriteString(Default)

	return fmt.Fprint(writer, buf.String())
}

func (clrs *Colors) Fprintln(writer io.Writer, args ...any) (n int, err error) {
	var buf strings.Builder
	buf.WriteString(strings.Join(*clrs, ""))
	buf.WriteString(fmt.Sprintln(args...))
	buf.WriteString(Default)

	return fmt.Fprint(writer, buf)
}

// endregion
