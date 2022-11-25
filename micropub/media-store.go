package micropub

import (
	"context"
	"log"
)

type MediaStore interface {
	// Save the files in the micropub raw post and insert the urls into the micropub post
	SaveMediaFiles(ctx context.Context, file MicropubFile) (string, error)
}

type NopMediaStore struct{}

func (n NopMediaStore) SaveMediaFiles(ctx context.Context, file MicropubFile) (string, error) {
	defer file.Reader.Close()
	log.Println("Skipping file: ", file.Name)
	return "", nil
}
