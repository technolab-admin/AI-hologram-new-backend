package utils

import "errors"

// Responsible for handling errors

var (
	// Meshy lifecycle errors
	ErrMeshyJobFailed = errors.New("meshy job failed")
	ErrMeshyTimeout   = errors.New("meshy job timed out")
	ErrMeshyInvalid   = errors.New("invalid meshy response")

	// Job runner errors
	ErrJobCancelled    = errors.New("job cancelled")
	ErrInvalidJobState = errors.New("invalid job state")
)
