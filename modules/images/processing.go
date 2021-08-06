package images

import (
	"os/exec"

	"userstyles.world/modules/errors"
)

type ImageType int

const (
	ImageTypeWEBP ImageType = iota
	ImageTypeJPEG
)

type VipsStatus int

const (
	notKnown VipsStatus = iota
	notInstalled
	installed
)

var vipsStatus = notKnown

func isVipsInstalled() VipsStatus {
	if vipsStatus != notKnown {
		return vipsStatus
	}
	vipsLocation := exec.Command("which", "vips")
	if vipsLocation.Run() != nil {
		vipsStatus = notInstalled
	} else {
		vipsStatus = installed
	}

	return vipsStatus
}

func decodeImage(original, newPath string, imageType ImageType) error {
	// Ensure vips installed
	if isVipsInstalled() == notInstalled {
		return errors.ErrVipsNotFound
	}

	var vipsCommand *exec.Cmd

	switch imageType {
	case ImageTypeWEBP:
		vipsCommand = exec.Command("vips", "webpsave", "--strip",
			"--reduction-effort", "4", "--Q", "80", original, newPath)
	case ImageTypeJPEG:
		vipsCommand = exec.Command("vips", "jpegsave", "--strip",
			"--Q", "80", "--optimize-coding", "--optimize-scans",
			"--trellis-quant", "--quant-table", "3", original, newPath)
	}

	err := vipsCommand.Run()
	if err != nil || vipsCommand.ProcessState.ExitCode() == 1 {
		return errors.ErrNoImageProcessing
	}

	return nil
}
