package game

import (
	"bytes"
	"embed"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*
var assets embed.FS

// AssetManager manages game assets like images and sounds
type AssetManager struct {
	images map[string]*ebiten.Image
}

// NewAssetManager creates a new asset manager
func NewAssetManager() *AssetManager {
	return &AssetManager{
		images: make(map[string]*ebiten.Image),
	}
}

// LoadAssets loads all game assets
func (am *AssetManager) LoadAssets() {
	// Create placeholder colored rectangles for now
	// In a real game, you would load actual image files
	
	// Player sprite (blue rectangle)
	am.images["player"] = am.createColoredRect(32, 48, color.RGBA{0, 100, 255, 255})
	
	// Platform sprite (green rectangle)
	am.images["platform"] = am.createColoredRect(128, 32, color.RGBA{0, 200, 0, 255})
	
	// Background (light blue)
	am.images["background"] = am.createColoredRect(800, 600, color.RGBA{135, 206, 235, 255})
	
	log.Println("Assets loaded successfully")
}

// GetImage returns an image by key
func (am *AssetManager) GetImage(key string) *ebiten.Image {
	return am.images[key]
}

// LoadImageFromFile loads an image from the embedded assets
func (am *AssetManager) LoadImageFromFile(filename, key string) error {
	data, err := assets.ReadFile("assets/" + filename)
	if err != nil {
		return err
	}
	
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return err
	}
	
	am.images[key] = ebiten.NewImageFromImage(img)
	return nil
}

// createColoredRect creates a colored rectangle image (for prototyping)
func (am *AssetManager) createColoredRect(width, height int, c color.RGBA) *ebiten.Image {
	img := ebiten.NewImage(width, height)
	img.Fill(c)
	return img
}

// AddImage adds a pre-created image to the asset manager
func (am *AssetManager) AddImage(key string, image *ebiten.Image) {
	am.images[key] = image
}

// HasImage checks if an image exists
func (am *AssetManager) HasImage(key string) bool {
	_, exists := am.images[key]
	return exists
}

// GetImageSize returns the size of an image
func (am *AssetManager) GetImageSize(key string) (int, int) {
	if img, exists := am.images[key]; exists {
		return img.Size()
	}
	return 0, 0
}
