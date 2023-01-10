package micropub

import (
	"context"
	"log"
)

type mediaStore interface {
	// Save the files in the micropub raw post and insert the urls into the micropub post
	SaveMediaFiles(ctx context.Context, file MicropubFile) (string, error)
}

type nopMediaStore struct {
	logger *log.Logger
}

func (n nopMediaStore) SaveMediaFiles(ctx context.Context, file MicropubFile) (string, error) {
	defer file.Reader.Close()
	n.logger.Println("Skipping file: ", file.Name)
	return "", nil
}
