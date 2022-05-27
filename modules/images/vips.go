package images

import (
	"os/exec"

	"userstyles.world/modules/log"
)

type imageKind int

const (
	imageFullWebP imageKind = iota
	imageFullJPEG
	imageThumbWebP
	imageThumbJPEG
)

// CheckVips will look for Vips binaries and exit if they're not found.
func CheckVips() {
	for _, name := range []string{"vips", "vipsthumbnail"} {
		cmd := exec.Command(name, "--version")
		if err := cmd.Run(); err != nil {
			log.Warn.Fatalf("%q binary not found on the $PATH.\n", name)
		}
	}
}

func decodeImage(src, out string, imageType imageKind) error {
	var cmd *exec.Cmd

	switch imageType {
	case imageFullWebP:
		cmd = exec.Command("vips", "webpsave", "--strip", "--reduction-effort",
			"4", "-n", "--Q", "80", src, out)

	case imageFullJPEG:
		cmd = exec.Command("vips", "jpegsave", "--strip", "--Q", "80",
			"--optimize-coding", "--optimize-scans", "--trellis-quant",
			"--quant-table", "3", src, out)

	case imageThumbWebP:
		cmd = exec.Command("vipsthumbnail", "--size", "300", "--export-profile",
			"srgb", "-o", "%st.webp[profile=none]", out)

	case imageThumbJPEG:
		cmd = exec.Command("vipsthumbnail", "--size", "300", "--export-profile",
			"srgb", "-o", "%st.jpeg[profile=none]", out)
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
