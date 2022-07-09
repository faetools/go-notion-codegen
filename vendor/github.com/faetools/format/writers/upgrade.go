package writers

import "io"

// A Writer is an expanded io.Writer.
type Writer interface {
	io.Writer
	ByteWriter
	StringWriter
}

type allWriter struct {
	io.Writer
	ByteWriter
	StringWriter
}

// Upgrade expands the functionality of a writer.
func Upgrade(w io.Writer) Writer {
	if w, ok := w.(Writer); ok {
		return w
	}

	upgrade := &allWriter{Writer: w}

	if bw, ok := w.(ByteWriter); ok {
		upgrade.ByteWriter = bw
	} else {
		upgrade.ByteWriter = &byteWriter{w}
	}

	if sw, ok := w.(StringWriter); ok {
		upgrade.StringWriter = sw
	} else {
		upgrade.StringWriter = &stringWriter{w}
	}

	return upgrade
}

type (
	// ByteWriter can write a single byte.
	ByteWriter interface{ WriteByte(byte) error }

	// StringWriter is a writer that can write strings.
	StringWriter interface{ WriteString(s string) (int, error) }

	byteWriter   struct{ w io.Writer }
	stringWriter struct{ w io.Writer }
)

func (bw *byteWriter) WriteByte(c byte) error {
	_, err := bw.w.Write([]byte{c})
	return err
}

func (sw *stringWriter) WriteString(s string) (int, error) {
	return sw.w.Write([]byte(s))
}
