package conf

import "fmt"

type GenericConfigError struct {
	Err  error
	File string
}

func newGenericConfigError(err error, file string) *GenericConfigError {
	return &GenericConfigError{
		Err:  err,
		File: file,
	}
}

func (e *GenericConfigError) Error() string {
	return fmt.Sprintf("error loading configuration from file %s: %v", e.File, e.Err)
}

// ---

type DuplicateKeyError struct {
	Key  string
	File string
}

func newDuplicateKeyError(key, file string) *DuplicateKeyError {
	return &DuplicateKeyError{
		Key:  key,
		File: file,
	}
}

func (e *DuplicateKeyError) Error() string {
	return fmt.Sprintf("error loading configuration from file %s: duplicate key %s", e.File, e.Key)
}
