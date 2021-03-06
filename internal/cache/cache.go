package cache

import (
	"context"
	"errors"
	"time"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
)

var (
	ErrCacheMiss = errors.New("item could not be found in the cache")
)

// Cache is used to store and retrieve previously obtained objects.
type Cache interface {
	albumsCache
}

// albumsCache is used to store and retrieve previously obtained Albums.
type albumsCache interface {
	// GetAlbum returns Album data from the cache corresponding to the specified title.
	// It will return ErrCacheMiss if there is no cached Album.
	GetAlbum(ctx context.Context, title string) (photoslibrary.Album, error)

	// PutAlbum stores the Album data in the cache using the title as key.
	// Underlying implementations may use any data storage format,
	// as long as the reverse operation, GetAlbum, results in the original data.
	PutAlbum(ctx context.Context, album photoslibrary.Album, ttl time.Duration) error

	// DeleteAlbum removes the Album data from the cache corresponding to the specified title.
	// If there's no such Album in the cache, it will return nil.
	InvalidateAlbum(ctx context.Context, title string) error
}
