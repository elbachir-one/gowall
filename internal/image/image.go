package image

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/Achno/gowall/utils"

	"github.com/chai2010/webp"
)

// Available formats to Encode an image in
var encoders = map[string]func(file *os.File, img image.Image) error{
    "png": func(file *os.File, img image.Image) error {
        return png.Encode(file, img)
    },
    "jpg": func(file *os.File, img image.Image) error {
        return jpeg.Encode(file, img, nil)
    },
    "jpeg": func(file *os.File, img image.Image) error {
        return jpeg.Encode(file, img, nil)
    },
    "webp": func(file *os.File, img image.Image) error {
        return webp.Encode(file, img, nil)
    },
}

func LoadImage(filePath string) (image.Image , error){
	
	file,err := os.Open(filePath)

	if err != nil {
		return nil,err
	}

	defer file.Close()

	img,_,err := image.Decode(file)

	return img,err
}

func SaveImage(img image.Image, filePath string, format string) error{

	file,err := os.Create(filePath)

	if err != nil {
		return err
	}

	defer file.Close()

	
    encoder, ok := encoders[strings.ToLower(format)]

    if !ok {
        return fmt.Errorf("unsupported format: %s", format)
    }

    return encoder(file, img)

}

func ProcessImg(imgPath string, theme string )  {

	img, err := LoadImage(imgPath)

	if err != nil {
		fmt.Println("Error loading image unsupported format:", err)
		return
	}

	var selectedTheme Theme

	switch theme {
	case "catpuccin":
		selectedTheme = Catpuccin
	
	default: 
		fmt.Println("Unknown theme:", theme)
		return

	}

	newImg, err := convertImage(img,selectedTheme)

	if err != nil {
		fmt.Println("Error Converting image:", err)
		return
	}

	//Extract file extension from imgPath
	extension := strings.ToLower(filepath.Ext(imgPath))

	if extension == "" {
		fmt.Println("Error: Could not determine file extension.")
		return
	}

	// remove '.' from the extension
	extension = extension[1:]

	dirPath, err := utils.CreateDirectory()
	nameOfFile := filepath.Base(imgPath)
	
    outputFilePath := filepath.Join(dirPath, nameOfFile)


	if err != nil{
		fmt.Println("Error creating Directory or getting path")
		return
	}

	err = SaveImage(newImg, outputFilePath, extension)

	if err != nil {
		fmt.Println("Error saving image:", err, outputFilePath)
		return
	}

	fmt.Printf("Image processed and saved as %s\n", outputFilePath)

}