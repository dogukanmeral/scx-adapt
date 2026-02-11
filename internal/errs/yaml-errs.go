package errs

// NO nested errors: if an error occured from a function call; return the error as it is

type ConflictPrioritiesError struct {
	Msg string
}

func (e *ConflictPrioritiesError) Error() string {
	return e.Msg
}

type InvalidSchedulerError struct {
	Msg string
}

func (e *InvalidSchedulerError) Error() string {
	return e.Msg
}

type ConflictCriteriasError struct {
	Msg string
}

func (e *ConflictCriteriasError) Error() string {
	return e.Msg
}

type InvalidValueNameError struct {
	Msg string
}

func (e *InvalidValueNameError) Error() string {
	return e.Msg
}

type MissingParameterError struct {
	Msg string
}

func (e *MissingParameterError) Error() string {
	return e.Msg
}

type ConflictParametersError struct {
	Msg string
}

func (e *ConflictParametersError) Error() string {
	return e.Msg
}
