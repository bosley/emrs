package api

import (
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
)

var ErrMalformedRoute = errors.New("invalid route")

// UUIDs and the chunks that compose an action route (chunk.chunk.chunk)
// all follow the same naming convention: alphanum, with any `-` or `_`
// All other symbols are not allowed
func ValidateCommonChunk(chunk string) bool {
	if len(strings.TrimSpace(chunk)) == 0 {
		return false
	}
	// This is verbose for a reason -
	// reading the code prefixed with `!`
	// seems unintuitive
	if regexp.MustCompile(`[^a-zA-Z0-9]+`).MatchString(chunk) {
		return false
	}
	return true
}

// Transform a '.' encoded action route into a string list, whils
// validating each piece of the route. Returns ErrMalformedRoute if
// any chunk of the given route is invalid
func DecomposeRoute(route string) ([]string, error) {

	slog.Debug("decompose route", "source", route)

	r := strings.Split(route, `.`)

	if len(r) < 1 {
		slog.Error("improper route, expected '.' composed string of at least size 1", "given", route, "split-length", len(r))
		return r, ErrMalformedRoute
	}

	for i, v := range r {
		if !ValidateCommonChunk(v) {
			slog.Error("invalid route due to bad chunk", "chunk", i)
			return r, ErrMalformedRoute
		}
	}

	return r, nil
}

// Transform a list of strings into an EMRS action route
// that is validated as its encoded. Returns ErrMalformedRoute
// if any chink of the given route is invalid
func ComposeRoute(route []string) (string, error) {
	result := ""
	for i, v := range route {
		if !ValidateCommonChunk(v) {
			slog.Error("invalid route due to bad chunk", "chunk", i)
			return result, ErrMalformedRoute
		}
		if i == 0 {
			result = v
		} else {
			result = fmt.Sprintf("%s.%s", result, v)
		}
	}
	return result, nil
}
