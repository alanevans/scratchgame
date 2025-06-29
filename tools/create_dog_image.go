package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

func main() {
	// Create a simple 32x48 dog-like image
	img := image.NewRGBA(image.Rect(0, 0, 32, 48))
	
	// Colors
	black := color.RGBA{0, 0, 0, 255}
	brown := color.RGBA{139, 69, 19, 255}
	pink := color.RGBA{255, 192, 203, 255}
	transparent := color.RGBA{0, 0, 0, 0}
	
	// Fill with transparent background
	for y := 0; y < 48; y++ {
		for x := 0; x < 32; x++ {
			img.Set(x, y, transparent)
		}
	}
	
	// Draw a simple dog shape
	// Body (brown)
	for y := 16; y < 40; y++ {
		for x := 6; x < 26; x++ {
			img.Set(x, y, brown)
		}
	}
	
	// Head (brown)
	for y := 4; y < 18; y++ {
		for x := 12; x < 28; x++ {
			img.Set(x, y, brown)
		}
	}
	
	// Snout (brown)
	for y := 10; y < 16; y++ {
		for x := 24; x < 30; x++ {
			img.Set(x, y, brown)
		}
	}
	
	// Ears (black)
	for y := 6; y < 16; y++ {
		for x := 10; x < 15; x++ {
			img.Set(x, y, black)
		}
	}
	for y := 4; y < 12; y++ {
		for x := 22; x < 26; x++ {
			img.Set(x, y, black)
		}
	}
	
	// Eye (black)
	img.Set(20, 8, black)
	img.Set(21, 8, black)
	img.Set(20, 9, black)
	img.Set(21, 9, black)
	
	// Nose (pink)
	img.Set(28, 12, pink)
	img.Set(29, 12, pink)
	
	// Legs (brown)
	for y := 40; y < 46; y++ {
		for x := 10; x < 14; x++ {
			img.Set(x, y, brown)
		}
		for x := 20; x < 24; x++ {
			img.Set(x, y, brown)
		}
	}
	
	// Paws (black)
	for x := 10; x < 14; x++ {
		img.Set(x, 46, black)
		img.Set(x, 47, black)
	}
	for x := 20; x < 24; x++ {
		img.Set(x, 46, black)
		img.Set(x, 47, black)
	}
	
	// Tail (brown with black tip)
	for y := 20; y < 28; y++ {
		for x := 2; x < 6; x++ {
			img.Set(x, y, brown)
		}
	}
	for y := 16; y < 20; y++ {
		for x := 2; x < 6; x++ {
			img.Set(x, y, black)
		}
	}
	
	// Create the file
	file, err := os.Create("internal/game/assets/dog.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	
	// Encode as PNG
	err = png.Encode(file, img)
	if err != nil {
		log.Fatal(err)
	}
	
	log.Println("Created sample dog.png in internal/game/assets/")
	log.Println("You can now rebuild the game to use this image!")
}
