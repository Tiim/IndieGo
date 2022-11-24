package micropub

import (
	"context"
	"fmt"
	"io"
	"tiim/go-comment-api/mfobjects"

	"github.com/google/uuid"
	"storj.io/uplink"
	"storj.io/uplink/edge"
)

type storjMediaStore struct {
	accessGrant string
	bucketName  string
	prefix      string
}

func NewStorjMediaStore(accessGrant, bucketName, prefix string) storjMediaStore {
	return storjMediaStore{
		accessGrant: accessGrant,
		bucketName:  bucketName,
		prefix:      prefix,
	}
}

func (s storjMediaStore) SaveMediaFiles(ctx context.Context, mpr MicropubPostRaw, mp *MicropubPost) error {
	access, err := uplink.ParseAccess(s.accessGrant)
	if err != nil {
		return fmt.Errorf("could not request access grant: %v", err)
	}

	// Open up the Project we will be working with.
	project, err := uplink.OpenProject(ctx, access)
	if err != nil {
		return fmt.Errorf("could not open project: %v", err)
	}
	defer project.Close()

	// Ensure the desired Bucket within the Project is created.
	_, err = project.EnsureBucket(ctx, s.bucketName)
	if err != nil {
		return fmt.Errorf("could not ensure bucket: %v", err)
	}

	for _, file := range mpr.Files {
		// Intitiate the upload of our Object to the specified bucket and key.
		key := s.uploadKey(file.ContentType)
		upload, err := project.UploadObject(ctx, s.bucketName, key, nil)
		if err != nil {
			return fmt.Errorf("could not initiate upload: %v", err)
		}

		// Copy the data to the upload.
		_, err = io.Copy(upload, file.Reader)
		if err != nil {
			_ = upload.Abort()
			return fmt.Errorf("could not upload data: %v", err)
		}
		// Commit the uploaded object.
		err = upload.Commit()
		if err != nil {
			return fmt.Errorf("could not commit uploaded object: %v", err)
		}

		urlAccess, err := access.Share(uplink.Permission{AllowDownload: true}, uplink.SharePrefix{Bucket: s.bucketName, Prefix: s.prefix})
		if err != nil {
			return fmt.Errorf("could not create url access: %v", err)
		}
		accessConfig := edge.Config{
			//https://forum.storj.io/t/s3-compatability-copy/17426/5
			AuthServiceAddress: "auth.eu1.storjshare.io:7777",
		}
		cred, err := accessConfig.RegisterAccess(ctx, urlAccess, &edge.RegisterAccessOptions{Public: true})
		if err != nil {
			return fmt.Errorf("could not register access: %v", err)
		}
		url, err := edge.JoinShareURL("https://link.storjshare.io", cred.AccessKeyID, s.bucketName, key, &edge.ShareURLOptions{Raw: true})
		if err != nil {
			return fmt.Errorf("could not join share url: %v", err)
		}
		fmt.Printf("Uploaded %s to %s\n", file.Name, url)

		addFile(mp, url, file.Name, file.ContentType)
	}
	return nil
}

func (s storjMediaStore) uploadKey(mimeType string) string {
	key := s.prefix
	if key != "" && key[len(key)-1] != '/' {
		key += "/"
	}
	extension := "bin"
	switch mimeType {
	case "image/jpeg":
		extension = "jpg"
	case "image/png":
		extension = "png"
	case "image/gif":
		extension = "gif"
	}
	id := uuid.New()
	key += id.String() + "." + extension
	return key
}

func addFile(mp *MicropubPost, url, name, contentType string) {
	switch contentType {
	case "image/jpeg", "image/png", "image/gif":
		mp.Entry.Photos = append(mp.Entry.Photos, mfobjects.MF2Photo{
			Url: url,
		})
	}
}
