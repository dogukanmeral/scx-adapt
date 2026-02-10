package errs

type NotObjFileError struct {
	Msg string
}

func (e *NotObjFileError) Error() string {
	return e.Msg
}

type NotBPFFileError struct {
	Msg string
}

func (e *NotBPFFileError) Error() string {
	return e.Msg
}

type NoStructOpsError struct {
	Msg string
}

func (e *NoStructOpsError) Error() string {
	return e.Msg
}
