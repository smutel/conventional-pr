package constants

import "errors"

// Metadata error
var (
	ErrMalformedMetadata = errors.New("[Meta] Malformed repository metadata")
)

// Config error
var (
	ErrMissingToken       = errors.New("[Config] Access token is empty")
	ErrNegativeFileChange = errors.New("[Config] Maximum file change must not be a negative number")
	ErrInvalidPattern     = errors.New("[Config] Invalid pull request title pattern")
)

// Event error
var (
	ErrEventFileRead  = errors.New("[Event] Failed to read event file")
	ErrEventFileParse = errors.New("[Event] Failed to parse event file")
)

// validator error
var (
	ErrInvalidTitle = errors.New("pull request title does not follow the conventional commit style")
	ErrNoIssue      = errors.New("pull request title does not mention any issues")
)
