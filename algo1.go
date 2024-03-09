package main

import (
	//"fmt"
	"image"
	"image/color"
)

type Counter struct {
	num uint8
	counter int
}

// Function that gets the average RGB values form a square in a image
func GetAverageColorSquare(background_img image.RGBA, posx, posy, width, height int) color.RGBA {

	r, g, b, _ := background_img.At(posx, posy).RGBA()
	av_color := color.RGBA{uint8(r), uint8(g), uint8(b), 255}

	// Initialize three arrays for holding the RGB values
	red_counter := []Counter{ {av_color.R, 1} }
	green_counter := []Counter{ {av_color.G, 1} }
	blue_counter := []Counter{ {av_color.B, 1} }

	// Loop over all the pixels in the square
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			if x == 0 && y == 0 {
				break;
			}

			// extract the color from the pixel
			r, g, b, _ := background_img.At(posx + x, posy + y).RGBA()
			checked_color := color.RGBA{uint8(r), uint8(g), uint8(b), 255}

			// extract the red value and add it to the array if it is not yet seen
			for index, color_info := range red_counter {

				if (checked_color.R == color_info.num) {
					color_info.counter += 1
					break

				} else if index == len(red_counter) - 1 {
					red_counter = append(red_counter, Counter{checked_color.R, 1})
				}
					
			}

			// extract the green value and add it to the array if it is not yet seen
			for index, color_info := range green_counter {

				if (checked_color.G == color_info.num) {
					color_info.counter += 1
					break

				} else if index == len(green_counter) - 1 {
					green_counter = append(green_counter, Counter{checked_color.G, 1})
				}

			}

			// extract the blue value and add it to the array if it is not yet seen
			for index, color_info := range blue_counter {

				if (checked_color.B == color_info.num) {
					color_info.counter += 1
					break

				} else if index == len(blue_counter) - 1 {
					blue_counter = append(blue_counter, Counter{checked_color.B, 1})
				}

			}

			
		}
	}

	// Get the most used value from each of the arrays
	var max_num uint32
	new_color := color.RGBA{0, 0, 0, 0}

	for _, color_info := range red_counter {
		if uint32(color_info.counter) > max_num {
			new_color.R = uint8(color_info.num)
		}
	}

	max_num = 0
	for _, color_info := range green_counter {
		if uint32(color_info.counter) > max_num {
			new_color.G = uint8(color_info.num)
		}
	}

	max_num = 0
	for _, color_info := range blue_counter {
		if uint32(color_info.counter) > max_num {
			new_color.B = uint8(color_info.num)
		}
	}

	return new_color
}

