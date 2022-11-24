package micropub

import (
	"context"
	"log"
)

type MediaStore interface {
	// Save the files in the micropub raw post and insert the urls into the micropub post
	SaveMediaFiles(ctx context.Context, mpr MicropubPostRaw, mp *MicropubPost) error
}

type NopMediaStore struct{}

func (n NopMediaStore) SaveMediaFiles(ctx context.Context, mpr MicropubPostRaw, mp *MicropubPost) error {
	for _, file := range mpr.Files {
		log.Println("Skipping file:", file)
		file.Reader.Close()
	}
	return nil
}
