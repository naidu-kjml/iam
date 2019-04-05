package shared

// APIError used for errors that will be returned to the client
type APIError struct {
	Message string
	Code    int
}

func (err APIError) Error() string {
	return err.Message
}
