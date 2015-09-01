package cli

type Exit struct {
	Err        error
	StatusCode int
}

func (e Exit) IsError() bool {
	return e.Err != nil
}

func (e Exit) Error() string {
	if e.IsError() {
		return e.Err.Error()
	}
	return ""
}
