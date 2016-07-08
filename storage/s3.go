package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	_ Service = (*S3)(nil)
)

type S3 struct {
	Service *s3.S3
}

func NewS3(crds *credentials.Credentials, region string) *S3 {
	cfg := &aws.Config{
		Credentials: crds,
		Region:      aws.String(region),
	}
	sess := session.New(cfg)
	svc := s3.New(sess)
	return &S3{Service: svc}
}

func (s *S3) Head(r *HeadRequest) (*HeadResponse, error) {
	res, err := s.Service.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(r.Bucket),
		Key:    aws.String(r.Key),
	})
	if err != nil {
		return nil, err
	}
	return &HeadResponse{
		Headers: Headers{
			AcceptRanges:       aws.StringValue(res.AcceptRanges),
			CacheControl:       aws.StringValue(res.CacheControl),
			ContentDisposition: aws.StringValue(res.ContentDisposition),
			ContentEncoding:    aws.StringValue(res.ContentEncoding),
			ContentLanguage:    aws.StringValue(res.ContentLanguage),
			ContentLength:      aws.Int64Value(res.ContentLength),
			ContentType:        aws.StringValue(res.ContentType),
			ETag:               aws.StringValue(res.ETag),
			Expires:            aws.StringValue(res.Expires),
			LastModified:       aws.TimeValue(res.LastModified),
		},
	}, nil
}

func (s *S3) Get(r *GetRequest) (*GetResponse, error) {
	res, err := s.Service.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(r.Bucket),
		Key:    aws.String(r.Key),
	})
	if err != nil {
		return nil, err
	}
	return &GetResponse{
		Body: res.Body,
		Headers: Headers{
			AcceptRanges:       aws.StringValue(res.AcceptRanges),
			CacheControl:       aws.StringValue(res.CacheControl),
			ContentDisposition: aws.StringValue(res.ContentDisposition),
			ContentEncoding:    aws.StringValue(res.ContentEncoding),
			ContentLanguage:    aws.StringValue(res.ContentLanguage),
			ContentLength:      aws.Int64Value(res.ContentLength),
			ContentType:        aws.StringValue(res.ContentType),
			ETag:               aws.StringValue(res.ETag),
			Expires:            aws.StringValue(res.Expires),
			LastModified:       aws.TimeValue(res.LastModified),
		},
	}, nil
}

func (s *S3) Put(r *PutRequest) (*PutResponse, error) {
	params := &s3.PutObjectInput{
		Bucket: aws.String(r.Bucket), // Required
		Key:    aws.String(r.Key),    // Required
		Body:   r.Body,               // Required
	}

	if r.ContentType != "" {
		params.ContentType = aws.String(r.ContentType)
	} else {
		// detect content type
		ct, err := detectContentType(r.Body)
		if err != nil {
			return nil, err
		}
		if ct != "" {
			params.ContentType = aws.String(ct)
		}
	}

	params.ACL = aws.String(r.ACL)

	if r.CacheControl != "" {
		params.CacheControl = aws.String(r.CacheControl)
	}

	_, err := s.Service.PutObject(params)
	if err != nil {
		return nil, err
	}
	return &PutResponse{}, nil
}

func (s *S3) Delete(r *DeleteRequest) (*DeleteResponse, error) {
	params := &s3.DeleteObjectInput{
		Bucket: aws.String(r.Bucket), // Required
		Key:    aws.String(r.Key),    // Required
	}
	_, err := s.Service.DeleteObject(params)
	if err != nil {
		return nil, err
	}
	return &DeleteResponse{}, nil
}

func (s *S3) DeleteMulti(r *DeleteMultiRequest) (*DeleteMultiResponse, error) {
	objects := make([]*s3.ObjectIdentifier, len(r.Keys))
	for i, key := range r.Keys {
		objects[i] = &s3.ObjectIdentifier{
			Key: aws.String(key),
		}
	}
	s3params := &s3.DeleteObjectsInput{
		Bucket: aws.String(r.Bucket), // Required
		Delete: &s3.Delete{ // Required
			Objects: objects,
			Quiet:   aws.Bool(r.Quiet),
		},
	}
	resp, err := s.Service.DeleteObjects(s3params)
	out := &DeleteMultiResponse{}
	if resp != nil {
		out.Errors = make([]error, len(resp.Errors))
		for i, s3err := range resp.Errors {
			out.Errors[i] = &S3Error{s3err}
		}
		out.Keys = make([]string, len(resp.Deleted))
		for i, key := range resp.Deleted {
			out.Keys[i] = aws.StringValue(key.Key)
		}
	}
	return out, err
}

type S3Error struct {
	s3err *s3.Error
}

func (e *S3Error) Error() string {
	return e.s3err.String()
}
