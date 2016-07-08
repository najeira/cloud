package storage

import (
	"io"
	"time"
)

const (
	// Owner gets FULL_CONTROL. No one else has access rights (default).
	// Object owner gets OWNER access.
	AclPrivate = "private"

	// Owner gets FULL_CONTROL. The AllUsers group gets READ access.
	// Object owner gets OWNER access, and allUsers get READER access.
	AclPublicRead = "public-read"
)

type Service interface {
	Copy(*CopyRequest) (*CopyResponse, error)
	Head(*HeadRequest) (*HeadResponse, error)
	Get(*GetRequest) (*GetResponse, error)
	Put(*PutRequest) (*PutResponse, error)
	Delete(*DeleteRequest) (*DeleteResponse, error)
	DeleteMulti(*DeleteMultiRequest) (*DeleteMultiResponse, error)
}

type Object interface {
}

type Headers struct {
	AcceptRanges string

	// Specifies caching behavior along the request/reply chain.
	CacheControl string

	// Specifies presentational information for the object.
	ContentDisposition string

	// Specifies what content encodings have been applied to the object and thus
	// what decoding mechanisms must be applied to obtain the media-type referenced
	// by the Content-Type header field.
	ContentEncoding string

	// The language the content is in.
	ContentLanguage string

	// Size of the body in bytes.
	ContentLength int64

	// The portion of the object returned in the response.
	ContentRange string

	// A standard MIME type describing the format of the object data.
	ContentType string

	// An ETag is an opaque identifier assigned by a web server to a specific version
	// of a resource found at a URL
	ETag string

	// The date and time at which the object is no longer cacheable.
	Expires string

	// Last modified date of the object
	LastModified time.Time
}

type CopyRequest struct {
	SourceBucket string
	SourceKey    string
	Bucket       string
	Key          string
}

type CopyResponse struct {
}

type HeadRequest struct {
	Bucket string
	Key    string
}

type HeadResponse struct {
	Headers Headers
}

type GetRequest struct {
	Bucket string
	Key    string
}

type GetResponse struct {
	Headers Headers

	// Object data.
	Body io.ReadCloser
}

type PutRequest struct {
	Bucket             string
	Key                string
	Body               io.ReadSeeker
	ACL                string
	CacheControl       string
	ContentType        string
	ContentEncoding    string
	ContentLanguage    string
	ContentDisposition string
}

type PutResponse struct {
}

type DeleteRequest struct {
	Bucket string
	Key    string
}

type DeleteResponse struct {
}

type DeleteMultiRequest struct {
	Bucket string
	Keys   []string
	Quiet  bool
}

type DeleteMultiResponse struct {
	Keys   []string
	Errors []error
}
