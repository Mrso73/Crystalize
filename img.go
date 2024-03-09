package main

import (
	"image"
	"image/color"
)

// Function to get the color of every pixel in the image
func GetImageInfo(img image.Image) [][]color.RGBA {
	
	// Get all the pixels in the image
	bounds := img.Bounds().Max
    width, height := bounds.X, bounds.Y

    var pixels[][]color.RGBA
    for y := 0; y < height; y++ {
		
        var row[]color.RGBA
        for x := 0; x < width; x++ {
            r, g, b, a := img.At(x, y).RGBA()

            row = append(row, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
        }
        pixels = append(pixels, row)
    }
	
	return pixels
}


// This functions draws a square to a image using the blend function to give a layering effect
func DrawLayeredSquare(img *image.RGBA, posx, posy, width, height int, square_color color.RGBA) {

	// loop over every pixel in the square
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			var new_color color.RGBA

			// If this isn't the border of the square blend the color else draw a border
			if !(x == 0 || x == width - 1 || y == 0 || y == height - 1) {

				r, g, b, a := img.At(posx + x, posy + y).RGBA()
				new_color = BlendColors(square_color, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
			} else {
				// draw a border around the square
				new_color = color.RGBA{10, 10, 10, 255}
			}

			img.Set(posx + x, posy + y, new_color)
		}
	}
}


// This function blends two colors to give the illusion over layering colors using transparency
func BlendColors(foreground, background color.RGBA) color.RGBA {

	// putting each var in a variable so i won't have to cast them to a different type later
	f_r, f_g, f_b, f_a := float32(foreground.R), float32(foreground.G), float32(foreground.B), float32(foreground.A)
	b_r, b_g, b_b, b_a := float32(background.R), float32(background.G), float32(background.B), float32(background.A)

	f_a_normalized := f_a / 255
	b_a_normalized := b_a / 255

	new_r := f_r * f_a_normalized + b_r * b_a_normalized * (1 - f_a_normalized)
	new_g := f_g * f_a_normalized + b_g * b_a_normalized * (1 - f_a_normalized)
	new_b := f_b * f_a_normalized + b_b * b_a_normalized * (1 - f_a_normalized)

	blended_color := color.RGBA{ uint8( new_r ), uint8( new_g ), uint8( new_b ), 255, }

	return blended_color
}