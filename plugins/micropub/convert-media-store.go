package micropub

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"

	"github.com/chai2010/webp"
)

func jpegEncode(w io.Writer, m image.Image) error {
	return jpeg.Encode(w, m, &jpeg.Options{Quality: 75})
}
func gifEncode(w io.Writer, m image.Image) error {
	return gif.Encode(w, m, &gif.Options{NumColors: 256})
}

func webpEncode(w io.Writer, m image.Image) error {
	return webp.Encode(w, m, &webp.Options{Lossless: false, Quality: 75})
}

var encoders = map[string]func(io.Writer, image.Image) error{
	"image/jpeg": jpegEncode,
	"image/png":  png.Encode,
	"image/gif":  gifEncode,
	"image/webp": webpEncode,
}

type convertMediaStore struct {
	childMediaStore mediaStore
	convertMap      map[string]string
	logger          *log.Logger
}

func (s convertMediaStore) SaveMediaFiles(ctx context.Context, file MicropubFile) (string, error) {
	destType := "-"
	if ct, ok := s.convertMap[file.ContentType]; ok {
		destType = ct
	} else if ct, ok := s.convertMap["*"]; ok {
		destType = ct
	}

	if destType == "-" {
		return s.childMediaStore.SaveMediaFiles(ctx, file)
	}

	img, _, err := image.Decode(file.Reader)
	if err != nil {
		return "", fmt.Errorf("could not decode image: %v", err)
	}
	if encoder, ok := encoders[destType]; ok {
		buffer := &bytes.Buffer{}
		err = encoder(buffer, img)
		if err != nil {
			return "", fmt.Errorf("could not encode image: %v", err)
		}
		file.Reader = io.NopCloser(buffer)
		file.ContentType = destType
		return s.childMediaStore.SaveMediaFiles(ctx, file)
	}

	return "", fmt.Errorf("could not find encoder for %s, available encoders: %v", destType, encoders)
}
