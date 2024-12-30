package http

import (
	"fmt"
	"io"
)

type Response struct {
	HttpVersion   string
	StatusCode    int
	ContentLength int
	ContentType   string
	Body          *string
}

func NewResponse() (response Response) {
	return Response{
		HttpVersion:   "HTTP/1.1",
		StatusCode:    200,
		ContentLength: -1,
	}
}

func (r *Response) Write(writer io.Writer) (err error) {
	if _, err = fmt.Fprintf(writer, "%s %d %s\r\n", r.HttpVersion, r.StatusCode, StatusText(r.StatusCode)); err != nil {
		return err
	}

	if r.ContentType != "" {
		if _, err = fmt.Fprintf(writer, "Content-Type: %s\r\n", r.ContentType); err != nil {
			return err
		}
	}

	if r.Body == nil {
		if _, err = fmt.Fprint(writer, "\r\n"); err != nil {
			return err
		}
		return nil
	}

	if r.ContentLength == -1 {
		r.ContentLength = len(*r.Body)
	}

	if _, err = fmt.Fprintf(writer, "Content-Length: %d\r\n", r.ContentLength); err != nil {
		return err
	}

	if _, err = fmt.Fprintf(writer, "\r\n%s", *r.Body); err != nil {
		return err
	}

	return nil
}

func StatusText(statusCode int) string {
	switch statusCode {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 404:
		return "Not Found"
	case 500:
		return "Server Error"
	default:
		return ""
	}

}
