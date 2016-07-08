package storage

import (
	"net/http"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/storage/v1"
)

var (
	_ Service = (*GCS)(nil)
)

type GCS struct {
	Service *storage.Service
}

func NewGCS(client *http.Client) (*GCS, error) {
	svc, err := storage.New(client)
	if err != nil {
		return nil, err
	}
	return &GCS{Service: svc}, nil
}

func (s *GCS) Copy(r *CopyRequest) (*CopyResponse, error) {
	panic("not implemented")
	return nil, nil
}

func (s *GCS) Head(r *HeadRequest) (*HeadResponse, error) {
	panic("not implemented")
	return nil, nil
}

func (s *GCS) Get(r *GetRequest) (*GetResponse, error) {
	panic("not implemented")
	return nil, nil
}

func (s *GCS) Put(r *PutRequest) (*PutResponse, error) {
	object := &storage.Object{Name: r.Key}
	object.CacheControl = r.CacheControl
	object.ContentEncoding = r.ContentEncoding
	object.ContentLanguage = r.ContentLanguage
	object.ContentDisposition = r.ContentDisposition

	call := s.Service.Objects.Insert(r.Bucket, object)

	var acl string
	switch r.ACL {
	case AclPrivate:
		acl = "private"
	case AclPublicRead:
		acl = "publicRead"
	}
	call.PredefinedAcl(acl)

	if r.ContentType != "" {
		call.Media(r.Body, googleapi.ContentType(r.ContentType))
	} else {
		call.Media(r.Body)
	}

	_, err := call.Do()
	if err != nil {
		return nil, err
	}
	return &PutResponse{}, nil
}

func (s *GCS) Delete(r *DeleteRequest) (*DeleteResponse, error) {
	panic("not implemented")
	return nil, nil
}

func (s *GCS) DeleteMulti(r *DeleteMultiRequest) (*DeleteMultiResponse, error) {
	panic("not implemented")
	return nil, nil
}
