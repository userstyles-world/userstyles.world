package images

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"

	"github.com/Kagami/go-avif"
	"github.com/chai2010/webp"
)

func ProcessToJPEG(input string, path string) error {
	imageFile, err := os.Open(input)
	if err != nil {
		return err
	}
	image, _, err := image.Decode(imageFile)
	if err != nil {
		image, err = png.Decode(imageFile)
		if err != nil {
			return err
		}
	}
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, image, &jpeg.Options{
		Quality: 75,
	})
	if err != nil {
		return err
	}

	err = os.WriteFile(path, buf.Bytes(), 0666)
	if err != nil {
		return err
	}
	return nil
}

func ProcessToWebp(imageBytes io.Reader, path string) error {
	image, _, err := image.Decode(imageBytes)
	if err != nil {
		return err
	}
	if err = webp.Save(path, image, &webp.Options{
		Lossless: false,
		Quality:  65,
		Exact:    false,
	}); err != nil {
		return err
	}
	return nil
}

func ProcessToAvif(imageBytes io.Reader, path string) error {
	image, _, err := image.Decode(imageBytes)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err = avif.Encode(&buf, image, &avif.Options{
		Threads: 0,
		Speed:   2,
		Quality: avif.MaxQuality,
	}); err != nil {
		return err
	}
	err = os.WriteFile(path, buf.Bytes(), 0666)
	if err != nil {
		return err
	}
	return nil
}
