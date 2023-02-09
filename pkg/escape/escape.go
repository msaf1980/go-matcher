package escape

import (
	"bytes"
)

// Path escapes the string so it can be safely used as a URL path.
func Path(s string) string {
	return escape(s, encodePath)
}

func PathTo(s string, sb *bytes.Buffer) {
	escapeTo(s, encodePath, sb)
}

// Query escapes the string so it can be safely placed inside a URL query.
func Query(s string) string {
	return escape(s, encodeQueryComponent)
}

func QueryTo(s string, sb *bytes.Buffer) {
	escapeTo(s, encodeQueryComponent, sb)
}
