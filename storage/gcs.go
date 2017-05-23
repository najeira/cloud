package storage

import (
	"errors"
	"fmt"
	"net/http"
	"unicode/utf8"

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
	if r.Bucket == "" {
		return nil, errors.New("storage: bucket name is empty")
	} else if r.Key == "" {
		return nil, errors.New("storage: object name is empty")
	} else if r.Body == nil {
		return nil, errors.New("storage: object body is nil")
	} else if !utf8.ValidString(r.Key) {
		return nil, fmt.Errorf("storage: object name %q is not valid UTF-8", r.Key)
	}

	object := &storage.Object{
		Bucket:             r.Bucket,
		Name:               r.Key,
		ContentEncoding:    r.ContentEncoding,
		ContentLanguage:    r.ContentLanguage,
		CacheControl:       r.CacheControl,
		ContentDisposition: r.ContentDisposition,
		//StorageClass:       r.StorageClass,
	}

	mediaOpts := []googleapi.MediaOption{
		googleapi.ChunkSize(googleapi.DefaultUploadChunkSize),
	}
	if r.ContentType != "" {
		mediaOpts = append(mediaOpts, googleapi.ContentType(r.ContentType))
	}

	call := s.Service.Objects.Insert(r.Bucket, object)
	//call.Context(ctx)

	switch r.ACL {
	case AclPrivate:
		call.PredefinedAcl("private")
	case AclPublicRead:
		call.PredefinedAcl("publicRead")
	}

	call.Media(r.Body, mediaOpts...)
	call.Projection("full")

	_, err := call.Do()
	if err != nil {
		return nil, err
	}
	return &PutResponse{}, nil
}

func (s *GCS) Delete(r *DeleteRequest) (*DeleteResponse, error) {
	if r.Bucket == "" {
		return nil, errors.New("storage: bucket name is empty")
	} else if r.Key == "" {
		return nil, errors.New("storage: object name is empty")
	} else if !utf8.ValidString(r.Key) {
		return nil, fmt.Errorf("storage: object name %q is not valid UTF-8", r.Key)
	}

	call := s.Service.Objects.Delete(r.Bucket, r.Key)
	//call.Context(ctx)
	err := call.Do()
	if err != nil {
		return nil, err
	}
	return &DeleteResponse{}, nil
}

func (s *GCS) DeleteMulti(r *DeleteMultiRequest) (*DeleteMultiResponse, error) {
	panic("not implemented")
	return nil, nil
}

func (s *GCS) List(r *ListRequest) (*ListResponse, error) {
	if r.Bucket == "" {
		return nil, errors.New("storage: bucket name is empty")
	}

	call := s.Service.Objects.List(r.Bucket)
	//call.Context(ctx)
	call.Projection("full")
	call.Delimiter("/")
	if len(r.Cursor) > 0 {
		call.PageToken(r.Cursor)
	}
	if len(r.Prefix) > 0 {
		call.Prefix(r.Prefix)
	}
	if r.Size > 0 {
		call.MaxResults(int64(r.Size))
	}

	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	objects := make([]*Object, 0, len(resp.Prefixes)+len(resp.Items))
	for _, prefix := range resp.Prefixes {
		object := &Object{
			Bucket:    r.Bucket,
			Name:      prefix,
			Directory: true,
		}
		objects = append(objects, object)
	}
	for _, item := range resp.Items {
		object := &Object{
			Bucket:             item.Bucket,
			CacheControl:       item.CacheControl,
			ComponentCount:     item.ComponentCount,
			ContentDisposition: item.ContentDisposition,
			ContentEncoding:    item.ContentEncoding,
			ContentLanguage:    item.ContentLanguage,
			ContentType:        item.ContentType,
			Etag:               item.Etag,
			Generation:         item.Generation,
			Metadata:           item.Metadata,
			Name:               item.Name,
			Size:               item.Size,
			StorageClass:       item.StorageClass,
			Updated:            item.Updated,
		}
		objects = append(objects, object)
	}
	return &ListResponse{
		Objects: objects,
		Cursor:  resp.NextPageToken,
	}, nil
}
