package micropub

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/google/uuid"
	"storj.io/uplink"
)

type FormatUrl func(name, contentType, prefix, bucket string) string

type storjMediaStore struct {
	// the access grant from the storj.io dashboard
	accessGrant string
	// the name of the storage bucket
	bucketName string
	// the prefix to use for the uploaded files
	prefix string

	formatUrl FormatUrl
	logger    *log.Logger
}

func newStorjMediaStore(accessGrant, bucketName, prefix string, formatUrl FormatUrl, logger *log.Logger) storjMediaStore {
	if prefix != "" && prefix[len(prefix)-1] != '/' {
		prefix += "/"
	}

	return storjMediaStore{
		accessGrant: accessGrant,
		bucketName:  bucketName,
		prefix:      prefix,
		formatUrl:   formatUrl,
		logger:      logger,
	}
}

func (s storjMediaStore) SaveMediaFiles(ctx context.Context, file MicropubFile) (string, error) {
	access, err := uplink.ParseAccess(s.accessGrant)
	if err != nil {
		return "", fmt.Errorf("could not request access grant: %v", err)
	}

	// Open up the Project we will be working with.
	project, err := uplink.OpenProject(ctx, access)
	if err != nil {
		return "", fmt.Errorf("could not open project: %v", err)
	}
	defer project.Close()

	// Ensure the desired Bucket within the Project is created.
	_, err = project.EnsureBucket(ctx, s.bucketName)
	if err != nil {
		return "", fmt.Errorf("could not ensure bucket: %v", err)
	}

	// Intitiate the upload of our Object to the specified bucket and key.
	key, name := s.uploadKey(file.ContentType)
	upload, err := project.UploadObject(ctx, s.bucketName, key, nil)
	if err != nil {
		return "", fmt.Errorf("could not initiate upload: %v", err)
	}

	// Copy the data to the upload.
	_, err = io.Copy(upload, file.Reader)
	if err != nil {
		_ = upload.Abort()
		return "", fmt.Errorf("could not upload data: %v", err)
	}
	// Commit the uploaded object.
	err = upload.Commit()
	if err != nil {
		return "", fmt.Errorf("could not commit uploaded object: %v", err)
	}
	url := s.formatUrl(name, file.ContentType, s.prefix, s.bucketName)
	log.Printf("Uploaded %s to %s\n", file.Name, url)

	return url, nil
}

func (s storjMediaStore) uploadKey(mimeType string) (string, string) {
	key := s.prefix
	extension := "bin"
	switch mimeType {
	case "image/jpeg":
		extension = "jpg"
	case "image/png":
		extension = "png"
	case "image/gif":
		extension = "gif"
	case "image/webp":
		extension = "webp"
	case "image/avif":
		extension = "avif"
	default:
		s.logger.Println("Unknown mime type for micropub upload: ", mimeType)
	}
	id := uuid.New()
	name := id.String() + "." + extension
	key += name

	return key, name
}
