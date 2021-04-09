package images

import (
	"os"

	"github.com/davidbyttow/govips/v2/vips"
)

func DecodeImage(input, output string, imageType vips.ImageType) error {
	buf, err := os.ReadFile(input)
	if err != nil {
		return err
	}
	buffer, err := vips.NewImageFromBuffer(buf)
	if err != nil {
		return err
	}

	newImage, _, err := buffer.Export(&vips.ExportParams{
		Format:        imageType,
		Quality:       60,
		Compression:   6,
		Effort:        4,
		StripMetadata: true,
	})

	if err != nil {
		return err
	}
	err = os.WriteFile(output, newImage, 0644)
	if err != nil {
		return err
	}

	return nil
}
