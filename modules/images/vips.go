package images

import (
	"os/exec"

	"userstyles.world/modules/errors"
	"userstyles.world/modules/log"
)

type ImageType int

const (
	ImageTypeWEBP ImageType = iota
	ImageTypeJPEG
)

// CheckVips will look for Vips binary and exit if it's not found.
func CheckVips() {
	cmd := exec.Command("vips", "--version")
	if err := cmd.Run(); err != nil {
		log.Warn.Fatal(errors.ErrVipsNotFound)
	}
}

func decodeImage(src, out string, imageType ImageType) error {
	var cmd *exec.Cmd

	switch imageType {
	case ImageTypeWEBP:
		cmd = exec.Command("vips", "webpsave", "--strip", "--reduction-effort",
			"4", "-n", "--Q", "80", src, out)
	case ImageTypeJPEG:
		cmd = exec.Command("vips", "jpegsave", "--strip", "--Q", "80",
			"--optimize-coding", "--optimize-scans", "--trellis-quant",
			"--quant-table", "3", src, out)
	}

	if err := cmd.Run(); err != nil {
		return errors.ErrNoImageProcessing
	}

	return nil
}
