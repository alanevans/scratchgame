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

	// Player sprite (black and white fluffy dog)
	am.images["player"] = am.createFluffyDogSprite(32, 48)

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

// createFluffyDogSprite creates a black and white fluffy dog sprite facing right
func (am *AssetManager) createFluffyDogSprite(width, height int) *ebiten.Image {
	img := ebiten.NewImage(width, height)

	// Colors for the dog
	white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 255}
	pink := color.RGBA{255, 192, 203, 255} // For nose

	// Clear to transparent first
	img.Fill(color.RGBA{0, 0, 0, 0})

	// Create a right-facing dog with clear features

	// Main body (white base) - made horizontally longer
	bodyImg := ebiten.NewImage(26, 24)
	bodyImg.Fill(white)
	bodyOpts := &ebiten.DrawImageOptions{}
	bodyOpts.GeoM.Translate(3, 16)
	img.DrawImage(bodyImg, bodyOpts)

	// Random black patches on body for fluffy pattern
	// Large patch on back
	backPatch := ebiten.NewImage(12, 12)
	backPatch.Fill(black)
	backPatchOpts := &ebiten.DrawImageOptions{}
	backPatchOpts.GeoM.Translate(5, 18)
	img.DrawImage(backPatch, backPatchOpts)

	// Small patch on side
	sidePatch := ebiten.NewImage(8, 8)
	sidePatch.Fill(black)
	sidePatchOpts := &ebiten.DrawImageOptions{}
	sidePatchOpts.GeoM.Translate(20, 24)
	img.DrawImage(sidePatch, sidePatchOpts)

	// Head (white base) - positioned for right-facing
	headImg := ebiten.NewImage(16, 14)
	headImg.Fill(white)
	headOpts := &ebiten.DrawImageOptions{}
	headOpts.GeoM.Translate(12, 4)
	img.DrawImage(headImg, headOpts)

	// Snout (longer, extending further to the right)
	snoutImg := ebiten.NewImage(12, 6)
	snoutImg.Fill(white)
	snoutOpts := &ebiten.DrawImageOptions{}
	snoutOpts.GeoM.Translate(24, 10)
	img.DrawImage(snoutImg, snoutOpts)

	// Black patch on head (random pattern)
	headPatch := ebiten.NewImage(8, 10)
	headPatch.Fill(black)
	headPatchOpts := &ebiten.DrawImageOptions{}
	headPatchOpts.GeoM.Translate(14, 6)
	img.DrawImage(headPatch, headPatchOpts)

	// Black eye patch (covering one eye area)
	eyePatch := ebiten.NewImage(6, 6)
	eyePatch.Fill(black)
	eyePatchOpts := &ebiten.DrawImageOptions{}
	eyePatchOpts.GeoM.Translate(18, 7)
	img.DrawImage(eyePatch, eyePatchOpts)

	// Left ear (floppy black, hanging down on left side)
	leftEar := ebiten.NewImage(5, 10)
	leftEar.Fill(black)
	leftEarOpts := &ebiten.DrawImageOptions{}
	leftEarOpts.GeoM.Translate(10, 6)
	img.DrawImage(leftEar, leftEarOpts)

	// Right ear (floppy black, partially visible since facing right)
	rightEar := ebiten.NewImage(4, 8)
	rightEar.Fill(black)
	rightEarOpts := &ebiten.DrawImageOptions{}
	rightEarOpts.GeoM.Translate(22, 4)
	img.DrawImage(rightEar, rightEarOpts)

	// Eye (visible eye since facing right)
	eye := ebiten.NewImage(2, 2)
	eye.Fill(black)
	eyeOpts := &ebiten.DrawImageOptions{}
	eyeOpts.GeoM.Translate(16, 8)
	img.DrawImage(eye, eyeOpts)

	// Nose on snout (pink) - positioned at end of longer snout
	nose := ebiten.NewImage(2, 2)
	nose.Fill(pink)
	noseOpts := &ebiten.DrawImageOptions{}
	noseOpts.GeoM.Translate(33, 12)
	img.DrawImage(nose, noseOpts)

	// Black dot at tip of snout
	snoutTip := ebiten.NewImage(1, 1)
	snoutTip.Fill(black)
	snoutTipOpts := &ebiten.DrawImageOptions{}
	snoutTipOpts.GeoM.Translate(35, 13)
	img.DrawImage(snoutTip, snoutTipOpts)

	// Mouth line (small black line under nose)
	mouth := ebiten.NewImage(3, 1)
	mouth.Fill(black)
	mouthOpts := &ebiten.DrawImageOptions{}
	mouthOpts.GeoM.Translate(32, 14)
	img.DrawImage(mouth, mouthOpts)

	// Front legs (white with black paws)
	// Front leg (closer to viewer)
	frontLeg := ebiten.NewImage(4, 8)
	frontLeg.Fill(white)
	frontLegOpts := &ebiten.DrawImageOptions{}
	frontLegOpts.GeoM.Translate(20, 40)
	img.DrawImage(frontLeg, frontLegOpts)

	// Front paw (black)
	frontPaw := ebiten.NewImage(4, 2)
	frontPaw.Fill(black)
	frontPawOpts := &ebiten.DrawImageOptions{}
	frontPawOpts.GeoM.Translate(20, 46)
	img.DrawImage(frontPaw, frontPawOpts)

	// Back leg (further from viewer)
	backLeg := ebiten.NewImage(4, 8)
	backLeg.Fill(white)
	backLegOpts := &ebiten.DrawImageOptions{}
	backLegOpts.GeoM.Translate(10, 40)
	img.DrawImage(backLeg, backLegOpts)

	// Back paw (black)
	backPaw := ebiten.NewImage(4, 2)
	backPaw.Fill(black)
	backPawOpts := &ebiten.DrawImageOptions{}
	backPawOpts.GeoM.Translate(10, 46)
	img.DrawImage(backPaw, backPawOpts)

	// Tail (black and white striped, pointing up and to the right)
	// Black base of tail
	tailBase := ebiten.NewImage(3, 8)
	tailBase.Fill(black)
	tailBaseOpts := &ebiten.DrawImageOptions{}
	tailBaseOpts.GeoM.Translate(4, 20)
	img.DrawImage(tailBase, tailBaseOpts)

	// White stripe in middle of tail
	tailMiddle := ebiten.NewImage(3, 3)
	tailMiddle.Fill(white)
	tailMiddleOpts := &ebiten.DrawImageOptions{}
	tailMiddleOpts.GeoM.Translate(4, 22)
	img.DrawImage(tailMiddle, tailMiddleOpts)

	// Black tip curving right
	tailTip := ebiten.NewImage(4, 3)
	tailTip.Fill(black)
	tailTipOpts := &ebiten.DrawImageOptions{}
	tailTipOpts.GeoM.Translate(2, 16)
	img.DrawImage(tailTip, tailTipOpts)

	// White spot on tail tip
	tailTipSpot := ebiten.NewImage(2, 2)
	tailTipSpot.Fill(white)
	tailTipSpotOpts := &ebiten.DrawImageOptions{}
	tailTipSpotOpts.GeoM.Translate(3, 17)
	img.DrawImage(tailTipSpot, tailTipSpotOpts)

	return img
}
