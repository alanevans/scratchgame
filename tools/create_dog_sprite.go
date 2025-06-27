package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

// This is a helper program to create dog sprite images
func main() {
	createDogSprite(false, "dog_sprite_right.png") // Right-facing (default)
	createDogSprite(true, "dog_sprite_left.png")   // Left-facing (flipped)
	log.Println("Dog sprites created: dog_sprite_right.png and dog_sprite_left.png")
}

func createDogSprite(flipHorizontal bool, filename string) {
	// Create a 32x32 image
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))

	// Fill with transparent background
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 0}) // Transparent
		}
	}

	// Helper function to draw a rectangle (with optional horizontal flip)
	drawRect := func(x, y, w, h int, c color.RGBA) {
		for dy := 0; dy < h; dy++ {
			for dx := 0; dx < w; dx++ {
				finalX := x + dx
				if flipHorizontal {
					finalX = 31 - (x + dx) // Flip horizontally
				}
				if finalX < 32 && y+dy < 32 && finalX >= 0 && y+dy >= 0 {
					img.Set(finalX, y+dy, c)
				}
			}
		}
	}

	// Colors
	white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 255}

	// Dog body (white)
	drawRect(8, 16, 16, 12, white)

	// Dog head (white)
	drawRect(4, 8, 12, 12, white)

	// Dog ears (black)
	drawRect(2, 6, 4, 6, black)  // Left ear
	drawRect(14, 6, 4, 6, black) // Right ear

	// Dog nose (black)
	drawRect(8, 12, 2, 2, black)

	// Dog eyes (black)
	drawRect(6, 10, 1, 1, black)  // Left eye
	drawRect(13, 10, 1, 1, black) // Right eye

	// Dog legs (white)
	drawRect(10, 28, 3, 4, white) // Front left leg
	drawRect(15, 28, 3, 4, white) // Front right leg
	drawRect(19, 28, 3, 4, white) // Back left leg
	drawRect(22, 28, 3, 4, white) // Back right leg

	// Dog paws (black)
	drawRect(10, 30, 3, 2, black) // Front left paw
	drawRect(15, 30, 3, 2, black) // Front right paw
	drawRect(19, 30, 3, 2, black) // Back left paw
	drawRect(22, 30, 3, 2, black) // Back right paw

	// Dog tail (black)
	drawRect(24, 18, 4, 2, black)

	// Dog spots (black patches for pattern)
	drawRect(12, 18, 4, 4, black) // Body spot
	drawRect(16, 10, 2, 3, black) // Head spot

	// Save the image
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		log.Fatal(err)
	}
}
