package colorthief

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/Achno/gowall/internal/backends/colorthief/mediancut"
)

var DefaultMaxCubes = 6

// returns the base color from the image file
func GetColorFromFile(imgPath string) (color.Color, error) {
	colors, err := GetPaletteFromFile(imgPath, DefaultMaxCubes)
	if err != nil {
		return color.RGBA{}, nil
	}
	return colors[0], nil
}

// returns cluster similar colors from the image file
func GetPaletteFromFile(imgPath string, maxCubes int) ([]color.Color, error) {
	f, err := os.Open(imgPath)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return GetPalette(img, maxCubes)
}

// returns cluster similar colors by the median cut algorithm
func GetPalette(img image.Image, maxCubes int) ([]color.Color, error) {
	return mediancut.GetPalette(img, maxCubes)
}

// returns the base color from the image.Image
func GetColor(img image.Image) (color.Color, error) {
	colors, err := GetPalette(img, DefaultMaxCubes)
	if err != nil {
		return color.RGBA{}, nil
	}
	return colors[0], nil
}

func PrintColor(colors []color.Color, filename string) error {
	imgWidth := 100 * len(colors)
	imgHeight := 200
	if imgWidth == 0 {
		return errors.New("colors empty")
	}

	paletted := image.NewPaletted(image.Rect(0, 0, imgWidth, imgHeight), colors)

	for x := 0; x < imgWidth; x++ {
		idx := x / 100
		for y := 0; y < imgHeight; y++ {
			paletted.SetColorIndex(x, y, uint8(idx))
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, paletted)
}