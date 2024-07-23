package colors

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Colors []string

func (clrs *Colors) String() string {
	if clrs == nil {
		return ""
	}
	return strings.Join(*clrs, "")
}

// region String Printers

func (clrs *Colors) Sprintf(format string, args ...any) string {
	var buf strings.Builder
	buf.WriteString(clrs.String())
	buf.WriteString(fmt.Sprintf(format, args...))
	buf.WriteString(Default)
	return buf.String()
}

func (clrs *Colors) Sprint(args ...any) string {
	var buf strings.Builder
	buf.WriteString(clrs.String())
	buf.WriteString(fmt.Sprint(args...))
	buf.WriteString(Default)
	return buf.String()
}

func (clrs *Colors) Sprintln(args ...any) string {
	var buf strings.Builder
	buf.WriteString(clrs.String())
	buf.WriteString(fmt.Sprintln(args...))
	buf.WriteString(Default)
	return buf.String()
}

// endregion

// region File Printers

func (clrs *Colors) Fprintf(writer io.Writer, format string, args ...any) (n int, err error) {
	return writer.Write([]byte(clrs.Sprintf(format, args...)))
}

func (clrs *Colors) Fprint(writer io.Writer, args ...any) (n int, err error) {
	return writer.Write([]byte(clrs.Sprint(args...)))
}

func (clrs *Colors) Fprintln(writer io.Writer, args ...any) (n int, err error) {
	return writer.Write([]byte(clrs.Sprintln(args...)))
}

// endregion

// region Stdout Printers

func (clrs *Colors) Printf(format string, args ...any) (n int, err error) {
	return clrs.Fprintf(os.Stdout, format, args...)
}

func (clrs *Colors) Print(args ...any) (n int, err error) {
	return clrs.Fprint(os.Stdout, args...)
}

func (clrs *Colors) Println(args ...any) (n int, err error) {
	return clrs.Fprintln(os.Stdout, args...)
}

// endregion

func (clrs *Colors) Errorf(format string, args ...any) error {
	return errors.New(clrs.Sprintf(format, args...))
}
