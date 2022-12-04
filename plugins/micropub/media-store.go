package micropub

import (
	"context"
	"log"
)

type mediaStore interface {
	// Save the files in the micropub raw post and insert the urls into the micropub post
	SaveMediaFiles(ctx context.Context, file MicropubFile) (string, error)
}

type nopMediaStore struct{}

func (n nopMediaStore) SaveMediaFiles(ctx context.Context, file MicropubFile) (string, error) {
	defer file.Reader.Close()
	log.Println("Skipping file: ", file.Name)
	return "", nil
}
