package photoservice

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
	"github.com/lestrrat-go/backoff"

	"github.com/gphotosuploader/google-photos-api-client-go/v2/internal/log"
)

// GooglePhotosService represents a Google Photos client using the `gphotosuploader/googlemirror`
// photos library with exponential back-off implementation.
type GooglePhotosService struct {
	service *photoslibrary.Service
	log     log.Logger
}

// NewGooglePhotosService returns a Google Photos client using the `gphotosuploader/googlemirror`
// photos library with exponential back-off implementation. It can ve configured with Options.
func NewGooglePhotosService(httpClient *http.Client, options ...Option) (*GooglePhotosService, error) {
	logger := defaultLogger()

	s, err := photoslibrary.New(httpClient)
	if err != nil {
		return nil, err
	}

	for _, o := range options {
		switch o.Name() {
		case optkeyLogger:
			logger = o.Value().(log.Logger)
		}
	}

	return &GooglePhotosService{
		service: s,
		log:     logger,
	}, nil
}

// ListAlbums makes a `Albums.List` call to Google Photos API and returns the response.
// It implements retries using exponential back-off.
func (s *GooglePhotosService) ListAlbums(ctx context.Context, pageSize int64, pageToken string) (*photoslibrary.ListAlbumsResponse, error) {
	list := s.service.Albums.List().PageSize(pageSize).PageToken(pageToken)
	b, cancel := defaultRetryPolicy.Start(ctx)
	defer cancel()
	for backoff.Continue(b) {
		res, err := list.Do()
		switch {
		case err == nil:
			return res, nil
		case IsRetryableError(err):
			s.log.Debugf("Error while listing albums: %s", err)
		case IsRateLimitError(err):
			s.log.Warnf("Rate limit reached.")
			return nil, errors.New("rate limit reached. wait ~30 seconds before trying again")
		default:
			s.log.Errorf("Error while listing albums: %s", err)
			return nil, fmt.Errorf("error listing albums: err=%v", err)
		}
	}

	return nil, errors.New("error listing albums: retry over")
}

// CreateAlbum makes a `Albums.Create` call to Google Photos API and returns the response.
// It implements retries using exponential back-off.
func (s *GooglePhotosService) CreateAlbum(ctx context.Context, request *photoslibrary.CreateAlbumRequest) (*photoslibrary.Album, error) {
	create := s.service.Albums.Create(request)
	b, cancel := defaultRetryPolicy.Start(ctx)
	defer cancel()

	for backoff.Continue(b) {
		res, err := create.Do()
		switch {
		case err == nil:
			return res, nil
		case IsRetryableError(err):
			s.log.Debugf("Error while creating an album: %s", err)
		case IsRateLimitError(err):
			s.log.Warnf("Rate limit reached.")
			return nil, errors.New("rate limit reached. wait ~30 seconds before trying again")
		default:
			s.log.Errorf("Error while creating album: %s", err)
			return nil, fmt.Errorf("error creating albums: err=%v", err)
		}
	}

	return nil, errors.New("error creating album: retry over")
}

// CreateMediaItems makes a `MediaItems.BatchCreate` call to Google Photos API and returns the response.
// It implements retries using exponential back-off.
func (s *GooglePhotosService) CreateMediaItems(ctx context.Context, request *photoslibrary.BatchCreateMediaItemsRequest) (*photoslibrary.BatchCreateMediaItemsResponse, error) {
	create := s.service.MediaItems.BatchCreate(request)
	b, cancel := defaultRetryPolicy.Start(ctx)
	defer cancel()

	for backoff.Continue(b) {
		res, err := create.Do()
		switch {
		case err == nil:
			return res, nil
		case IsRetryableError(err):
			s.log.Debugf("Error while creating media items: %s", err)
		case IsRateLimitError(err):
			s.log.Warnf("Rate limit reached.")
			return nil, errors.New("rate limit reached. wait ~30 seconds before trying again")
		default:
			s.log.Errorf("Error while creating media items: %s", err)
			return nil, fmt.Errorf("error creating media items: err=%v", err)
		}
	}

	return nil, errors.New("error creating album: retry over")
}
