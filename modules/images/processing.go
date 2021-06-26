package images

import (
	"os/exec"

	"userstyles.world/modules/errors"
)

type ImageType int

const (
	ImageTypeAVIF ImageType = iota
	ImageTypeWEBP
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
	err := vipsLocation.Run()
	if err != nil {
		vipsStatus = notInstalled
	} else {
		vipsStatus = installed
	}

	return vipsStatus
}

func DecodeImage(original, newPath string, imageType ImageType) error {
	// Ensure vips installed
	if isVipsInstalled() == notInstalled {
		return errors.ErrVipsNotFound
	}

	var vipsCommand *exec.Cmd

	switch imageType {
	case ImageTypeWEBP:
		vipsCommand = exec.Command("vips", "webpsave", "--strip", "--reduction-effort", "4", "--Q", "50", original, newPath)
	case ImageTypeAVIF:
		vipsCommand = exec.Command("vips", "heifsave", "--strip", "--compression", "av1", "--Q", "50", original, newPath)
	case ImageTypeJPEG:
		vipsCommand = exec.Command("vips", "jpegsave", "--strip", "--Q", "50", "--optimize-coding", "--optimize-scans", "--trellis-quant", "--quant-table", "3", original, newPath)
	}

	err := vipsCommand.Run()
	if err != nil || vipsCommand.ProcessState.ExitCode() == 1 {
		return errors.ErrNoImageProcessing
	}

	return nil
}
