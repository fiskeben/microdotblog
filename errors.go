package microdotblog

import (
	"fmt"
	"io"
	"io/ioutil"
)

// HTTPError is used to indicate error responses from the API.
type httpError interface {
	error
	withServerResponse(reason string) httpError
}

// NotFound is returned when the API returns status 404.
type NotFound struct {
	msg            string
	ServerResponse string
}

func (e NotFound) Error() string {
	return fmt.Sprintf("%s (%s)", e.msg, e.ServerResponse)
}

func (e NotFound) withServerResponse(response string) httpError {
	return NotFound{
		msg:            e.msg,
		ServerResponse: response,
	}
}

// NotAuthorized is returned when the API returns 401.
type NotAuthorized struct {
	msg            string
	ServerResponse string
}

func (e NotAuthorized) Error() string {
	return fmt.Sprintf("%s (%s)", e.msg, e.ServerResponse)
}

func (e NotAuthorized) withServerResponse(response string) httpError {
	return NotAuthorized{
		msg:            e.msg,
		ServerResponse: response,
	}
}

// Forbidden is returned when the API returns 403.
type Forbidden struct {
	msg            string
	ServerResponse string
}

func (e Forbidden) Error() string {
	return fmt.Sprintf("%s (%s)", e.msg, e.ServerResponse)
}

func (e Forbidden) withServerResponse(response string) httpError {
	return Forbidden{
		msg:            e.msg,
		ServerResponse: response,
	}
}

// ServerError is a default error that is used when the API
// returns an 5xx status code that isn't already mapped to a
// more specific error type.
type ServerError struct {
	msg            string
	ServerResponse string
	StatusCode     int
}

func (e ServerError) Error() string {
	return fmt.Sprintf("Server error: %s (%d %s)", e.msg, e.StatusCode, e.ServerResponse)
}

func (e ServerError) withServerResponse(response string) httpError {
	return ServerError{
		msg:            e.msg,
		ServerResponse: response,
		StatusCode:     e.StatusCode,
	}
}

// ClientError is a default error that is used when the API
// returns a 3xx or 4xx status code that isn't already mapped to a
// more specific error type.
type ClientError struct {
	msg            string
	ServerResponse string
	StatusCode     int
}

func (e ClientError) Error() string {
	return fmt.Sprintf("Server error: %s (%d %s)", e.msg, e.StatusCode, e.ServerResponse)
}

func (e ClientError) withServerResponse(response string) httpError {
	return ClientError{
		msg:            e.msg,
		StatusCode:     e.StatusCode,
		ServerResponse: response,
	}
}

func newAPIError(status int, body io.ReadCloser) error {
	if status < 300 {
		return nil
	}

	var err httpError

	switch status {
	case 401:
		err = NotAuthorized{msg: "you are not authorized to access this resource"}
	case 403:
		err = Forbidden{msg: "insufficient privileges"}
	case 404:
		err = NotFound{msg: "the resource was not found"}
	default:
		if status >= 500 {
			err = ServerError{msg: "the server returned an error", StatusCode: status}
		} else if status >= 300 {
			err = ClientError{msg: "client error", StatusCode: status}
		}
	}

	if err != nil {
		defer body.Close()
		reason, readErr := ioutil.ReadAll(body)
		if readErr == nil {
			return err.withServerResponse(string(reason))
		}
		return err
	}
	return nil
}
