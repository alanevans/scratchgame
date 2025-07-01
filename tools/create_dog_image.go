package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

func main() {
	// Create a simple 48x48 black and white dog-like image
	img := image.NewRGBA(image.Rect(0, 0, 48, 48))
	
	// Colors - Black and White only
	white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 255}
	pink := color.RGBA{255, 192, 203, 255}
	transparent := color.RGBA{0, 0, 0, 0}
	
	// Fill with transparent background
	for y := 0; y < 48; y++ {
		for x := 0; x < 48; x++ {
			img.Set(x, y, transparent)
		}
	}
	
	// Draw a black and white dog shape (horizontally longer)
	// Body (white)
	for y := 16; y < 40; y++ {
		for x := 8; x < 40; x++ {
			img.Set(x, y, white)
		}
	}
	
	// Black patches on body for spotted pattern
	for y := 18; y < 30; y++ {
		for x := 12; x < 24; x++ {
			img.Set(x, y, black)
		}
	}
	for y := 24; y < 32; y++ {
		for x := 28; x < 36; x++ {
			img.Set(x, y, black)
		}
	}
	
	// Head (white)
	for y := 4; y < 18; y++ {
		for x := 20; x < 36; x++ {
			img.Set(x, y, white)
		}
	}
	
	// Snout (white)
	for y := 10; y < 16; y++ {
		for x := 32; x < 44; x++ {
			img.Set(x, y, white)
		}
	}
	
	// Black patch on head
	for y := 6; y < 16; y++ {
		for x := 22; x < 30; x++ {
			img.Set(x, y, black)
		}
	}
	
	// Ears (black, floppy)
	for y := 6; y < 16; y++ {
		for x := 16; x < 21; x++ {
			img.Set(x, y, black)
		}
	}
	for y := 4; y < 12; y++ {
		for x := 30; x < 34; x++ {
			img.Set(x, y, black)
		}
	}
	
	// Eye (black)
	img.Set(26, 8, black)
	img.Set(27, 8, black)
	img.Set(26, 9, black)
	img.Set(27, 9, black)
	
	// Nose (pink)
	img.Set(42, 12, pink)
	img.Set(43, 12, pink)
	
	// Legs (white)
	for y := 40; y < 46; y++ {
		for x := 16; x < 20; x++ {
			img.Set(x, y, white)
		}
		for x := 28; x < 32; x++ {
			img.Set(x, y, white)
		}
	}
	
	// Paws (black)
	for x := 16; x < 20; x++ {
		img.Set(x, 46, black)
		img.Set(x, 47, black)
	}
	for x := 28; x < 32; x++ {
		img.Set(x, 46, black)
		img.Set(x, 47, black)
	}
	
	// Tail (black and white striped)
	for y := 20; y < 28; y++ {
		for x := 4; x < 8; x++ {
			img.Set(x, y, black)
		}
	}
	for y := 16; y < 20; y++ {
		for x := 4; x < 8; x++ {
			img.Set(x, y, white)
		}
	}
	
	// Create the file
	file, err := os.Create("../internal/game/assets/dog.png")
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
