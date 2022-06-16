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

func (x imageKind) String() string {
	return [...]string{"WebP", "JPEG", "WebP thumb", "JPEG thumb"}[x]
}

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
	var args []string
	switch imageType {
	case imageFullWebP:
		args = []string{"vips", "webpsave", "--strip", "--reduction-effort",
			"4", "-n", "--Q", "80", src, out}

	case imageFullJPEG:
		args = []string{"vips", "jpegsave", "--strip", "--Q", "80",
			"--optimize-coding", "--optimize-scans", "--trellis-quant",
			"--quant-table", "3", src, out}

	case imageThumbWebP:
		args = []string{"vipsthumbnail", "--size", "300", "-o", "%st.webp", out}

	case imageThumbJPEG:
		args = []string{"vipsthumbnail", "--size", "300", "-o", "%st.jpeg", out}
	}

	return exec.Command(args[0], args[1:]...).Run()
}
