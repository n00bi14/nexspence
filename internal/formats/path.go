package formats

import (
	"errors"
	"fmt"
	"path"
	"strings"
)

// ErrNonCanonicalPath marks a repository path that names an artifact by a
// second spelling. See CanonicalPath.
var ErrNonCanonicalPath = errors.New("non-canonical path")

// nonCanonicalError is an ErrNonCanonicalPath whose text is the reason alone,
// ready to hand to the client.
type nonCanonicalError string

func (e nonCanonicalError) Error() string { return string(e) }

func (nonCanonicalError) Is(target error) bool { return target == ErrNonCanonicalPath }

// CanonicalPath returns the path a format handler serves for p, or
// ErrNonCanonicalPath when p reaches that artifact by a second spelling.
//
// Format handlers serve path.Clean of the request path. A check judged on the
// path as requested therefore misses every spelling Clean rewrites: an
// anchored ^/private/ rule does not match "//private/x" or "/./private/x", and
// the handler serves "/private/x" anyway. Those spellings are refused rather
// than rewritten, so no legitimate path changes meaning on the way through: an
// artifact path never has an empty, "." or ".." segment.
//
// A trailing slash is an ordinary spelling (PyPI's simple index needs it). It
// passes, and is dropped from the result as path.Clean drops it for the
// handler; see RuleAllows for how a rule judges it.
func CanonicalPath(p string) (string, error) {
	segs := strings.Split(strings.TrimPrefix(p, "/"), "/")
	for i, seg := range segs {
		switch {
		case seg == "" && i == len(segs)-1:
			// trailing slash, or the root itself
		case seg == "":
			return "", nonCanonicalError("path must not contain an empty segment")
		case seg == "." || seg == "..":
			return "", nonCanonicalError(fmt.Sprintf("path must not contain a %q segment", seg))
		}
	}
	return path.Clean("/" + strings.TrimPrefix(p, "/")), nil
}
