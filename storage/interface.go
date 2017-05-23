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

// Object: An object.
type Object struct {
	// Bucket: The name of the bucket containing this object.
	Bucket string `json:"bucket,omitempty"`

	// CacheControl: Cache-Control directive for the object data. If
	// omitted, and the object is accessible to all anonymous users, the
	// default will be public, max-age=3600.
	CacheControl string `json:"cacheControl,omitempty"`

	// ComponentCount: Number of underlying components that make up this
	// object. Components are accumulated by compose operations.
	ComponentCount int64 `json:"componentCount,omitempty"`
	// S3: PartsCount *int64 `location:"header" locationName:"x-amz-mp-parts-count" type:"integer"`

	// ContentDisposition: Content-Disposition of the object data.
	ContentDisposition string `json:"contentDisposition,omitempty"`

	// ContentEncoding: Content-Encoding of the object data.
	ContentEncoding string `json:"contentEncoding,omitempty"`

	// ContentLanguage: Content-Language of the object data.
	ContentLanguage string `json:"contentLanguage,omitempty"`

	// ContentType: Content-Type of the object data. If contentType is not
	// specified, object downloads will be served as
	// application/octet-stream.
	ContentType string `json:"contentType,omitempty"`

	// Etag: HTTP 1.1 Entity tag for the object.
	Etag string `json:"etag,omitempty"`

	// Generation: The content generation of this object. Used for object
	// versioning.
	Generation int64 `json:"generation,omitempty,string"`
	// S3: VersionId *string `location:"header" locationName:"x-amz-version-id" type:"string"`

	// Metadata: User-provided metadata, in key/value pairs.
	Metadata map[string]string `json:"metadata,omitempty"`
	// S3: Metadata map[string]*string `location:"headers" locationName:"x-amz-meta-" type:"map"`

	// Name: The name of the object. Required if not specified by URL
	// parameter.
	Name string `json:"name,omitempty"`

	// Size: Content-Length of the data in bytes.
	Size uint64 `json:"size,omitempty,string"`
	// S3: ContentLength *int64 `location:"header" locationName:"Content-Length" type:"long"`

	// StorageClass: Storage class of the object.
	StorageClass string `json:"storageClass,omitempty"`
	// S3: StorageClass *string `location:"header" locationName:"x-amz-storage-class" type:"string" enum:"StorageClass"`

	// Updated: The modification time of the object metadata in RFC 3339
	// format.
	Updated string `json:"updated,omitempty"`
	// S3: LastModified *time.Time `location:"header" locationName:"Last-Modified" type:"timestamp" timestampFormat:"rfc822"`

	// Directory: Whether the object is a directory or not.
	Directory bool
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

type ListRequest struct {
	Bucket string
	Prefix string
	Size   int
	Cursor string
}

type ListResponse struct {
	Objects []*Object
	Cursor  string
}
