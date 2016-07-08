package storage

import (
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

func ErrAWS(err error) bool {
	_, ok := err.(awserr.Error)
	return ok
}

func Err4xx(err error) bool {
	if aerr, ok := err.(awserr.RequestFailure); ok {
		sc := aerr.StatusCode()
		if sc >= 400 && sc <= 499 {
			return true
		}
	}
	return false
}

func Err404(err error) bool {
	if aerr, ok := err.(awserr.RequestFailure); ok {
		return aerr.StatusCode() == 404
	}
	return false
}

func detectContentType(rs io.ReadSeeker) (string, error) {
	// read header
	buf := make([]byte, 64)
	n, err := rs.Read(buf)
	if err != nil {
		return "", err
	}

	// restore original position
	_, err = rs.Seek(int64(0-n), 1)
	if err != nil {
		return "", err
	}

	// detect content type
	ct := http.DetectContentType(buf[:n])
	return ct, nil
}
