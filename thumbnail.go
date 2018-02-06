package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

func imageResizeName(imageFileName string, width, height int) string {
	infile := imageFileName
	ext := filepath.Ext(infile)
	size := fmt.Sprintf(".%dx%d", width, height)
	return strings.TrimSuffix(infile, ext) + size + ".thumb" + ext
}

func resizeImage(width, height int, infile string) {
	var e Err
	var dn DirNav

	fmt.Printf("Part 1")

	outfile := imageResizeName(infile, width, height)

	dn.setImgResizedName(outfile)

	fmt.Printf("Part 2")

	file, err := os.Open(dn.getOrigImgPath())
	if err != nil {
		e.setInternalServerError(err)
	}

	fmt.Printf("Part 3")

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		e.setInternalServerError(err)
	}

	file.Close()

	fmt.Printf("Part 4")

	// Documentaion on interpolation opions:
	// https://www.cambridgeincolour.com/tutorials/image-resize-for-web.htm
	m := resize.Thumbnail(uint(width), uint(height), img, resize.Bicubic)

	mb := m.Bounds()
	origWidth := mb.Dx()
	origHeight := mb.Dy()
	paddingTop := (height / 2) - (origHeight / 2)
	paddingLeft := (width / 2) - (origWidth / 2)

	myimageOffset := image.Pt(0, 0)
	myimageRect := image.Rect(0, 0, width, height)
	myimage := image.NewRGBA(myimageRect)
	black := color.RGBA{0, 0, 0, 0}
	draw.Draw(myimage, myimage.Bounds().Add(myimageOffset), &image.Uniform{black}, image.ZP, draw.Src)

	fmt.Printf("Part 5")

	b := myimage.Bounds().Add(image.Pt(-paddingLeft, -paddingTop))
	image3 := image.NewRGBA(b)
	draw.Draw(image3, m.Bounds(), m, image.ZP, draw.Over)

	fmt.Printf("Part 6")

	out, err := os.Create(dn.getResizedImgPath())
	if err != nil {
		e.setInternalServerError(err)
	}

	defer out.Close()

	// write new image to file
	jpeg.Encode(out, image3, &jpeg.Options{jpeg.DefaultQuality})
}
