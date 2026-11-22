package chronary

import (
	"net/url"
	"strings"
)

// invalidPathSegment stands in for a refused path parameter. It holds a NUL
// byte, which url.PathEscape never produces, so Client.do can recognise it and
// refuse the request before anything is sent.
const invalidPathSegment = "\x00invalid-path-segment\x00"

// pathSegment escapes one caller-supplied ID as a single URL path segment
// (#2519). url.PathEscape keeps "/", "?", "#" and "%" inside the segment, so
// an ID like "../agents/agt_1" or "evt_1?expand=all" can't reach another
// resource or add a query string. Escaping leaves dots alone, and a "." or
// ".." segment would be resolved by the edge, so those and "" are refused:
// pathSegment returns a marker and Client.do returns a validation error.
func pathSegment(id string) string {
	if id == "" || id == "." || id == ".." {
		return invalidPathSegment
	}
	return url.PathEscape(id)
}

// invalidPathError reports whether path holds a refused segment, and the
// error Client.do returns for it.
func invalidPathError(path string) error {
	if !strings.Contains(path, invalidPathSegment) {
		return nil
	}
	return &Error{
		Type:    "validation_error",
		Message: `chronary: invalid ID: a path parameter must not be empty, "." or "..". No request was sent.`,
	}
}
