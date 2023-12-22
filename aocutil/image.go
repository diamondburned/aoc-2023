package aocutil

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// SaveImage saves the given image to a PNG image.
func SaveImage(img image.Image, dst string) {
	if ext := filepath.Ext(dst); ext != ".png" {
		panic(fmt.Errorf("invalid extension %q, only .png is supported", ext))
	}
	f := E2(os.Create(dst))
	defer f.Close()
	E1(png.Encode(f, img))
}

// OpenImage saves the given image to a temporary PNG image and opens it.
func OpenImage(img image.Image) {
	f := E2(os.CreateTemp("", "aocutil-*.png"))
	defer f.Close()
	E1(png.Encode(f, img))
	log.Printf("opening PNG at %q", f.Name())
	cmd := exec.Command("xdg-open", f.Name())
	cmd.Stderr = os.Stderr
	E1(cmd.Start())
}
