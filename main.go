package main

import (
	"log"
	"os"

	"image"
	"image/png"
) 	

func main() {

	// Check arguments given
	if len(os.Args[1:]) != 1 {
		log.Fatal("Error: Program should only contain one argument, the path to the input image")
	}

	// Open the image file 
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal("Error: File could not be opened")
	}
	defer file.Close()

	input_img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal("Error: Image could not be decoded")
	}

	// ----------------------

	input_pixels := GetImageInfo(input_img)

	// Create an empty RGBA image
	upLeft := image.Point{0, 0}
	lowRight := image.Point{input_img.Bounds().Dx(), input_img.Bounds().Dy()}
	output_img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for y := 0; y < input_img.Bounds().Dy(); y++ {
		for x := 0; x < input_img.Bounds().Dx(); x++ {
			output_img.Set(x, y, input_pixels[y][x])

		}
	}

	// ----------------------

	for y := 0; y < input_img.Bounds().Dy(); y += 15 {
		for x := 0; x < input_img.Bounds().Dx(); x += 15 {
			avr_color := GetAverageColorSquare(*output_img, x, y, 15, 15)
			avr_color.A = 255
			DrawLayeredSquare(output_img, x, y, 15, 15, avr_color)
		}
	}

	// ----------------------

	// Save the image to a file
	outputFile, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, output_img)
	if err != nil {
		panic(err)
	}
}

