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

func decodeImage(original, newPath string, imageType ImageType) error {
	var vipsCommand *exec.Cmd

	switch imageType {
	case ImageTypeWEBP:
		vipsCommand = exec.Command("vips", "webpsave", "--strip",
			"--reduction-effort", "4", "-n", "--Q", "80", original, newPath)
	case ImageTypeJPEG:
		vipsCommand = exec.Command("vips", "jpegsave", "--strip",
			"--Q", "80", "--optimize-coding", "--optimize-scans",
			"--trellis-quant", "--quant-table", "3", original, newPath)
	}

	err := vipsCommand.Run()
	if err != nil || vipsCommand.ProcessState.ExitCode() == 1 {
		log.Warn.Printf("Failed to run vips: %v\n", err)
		return errors.ErrNoImageProcessing
	}

	return nil
}
