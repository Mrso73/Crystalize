package generator

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"
	"time"

	"github.com/Mrso73/Crystalize/internal/config"
)

type Generator struct {
	cfg      *config.ArtConfig
	original image.Image
	current  *image.RGBA
	rng      *rand.Rand
	bounds   image.Rectangle
}

func NewGenerator(cfg *config.ArtConfig, img image.Image) *Generator {
	bounds := img.Bounds()
	current := image.NewRGBA(bounds)

	// Initialize with average color of original image instead of white
	avgColor := calculateImageAverage(img)
	draw.Draw(current, bounds, &image.Uniform{avgColor}, image.Point{}, draw.Src)

	return &Generator{
		cfg:      cfg,
		original: img,
		current:  current,
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
		bounds:   bounds,
	}
}

// Calculate average color of entire image for better initial state
func calculateImageAverage(img image.Image) color.RGBA {
	bounds := img.Bounds()
	var rSum, gSum, bSum int64
	count := 0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rSum += int64(r)
			gSum += int64(g)
			bSum += int64(b)
			count++
		}
	}

	return color.RGBA{
		R: uint8((rSum / int64(count)) >> 8),
		G: uint8((gSum / int64(count)) >> 8),
		B: uint8((bSum / int64(count)) >> 8),
		A: 255,
	}
}

func (g *Generator) Generate(progress chan<- float32) {
	defer close(progress)

	for i := 0; i < g.cfg.Iterations; i++ {
		rect, rectColor := g.createRectangle()
		g.drawRectangle(rect, rectColor)
		progress <- float32(i+1) / float32(g.cfg.Iterations)
	}
}

func (g *Generator) createRectangle() (image.Rectangle, color.RGBA) {
	// Use smaller rectangles for better detail
	size := g.cfg.MinRectSize +
		g.rng.Intn(g.cfg.MaxRectSize-g.cfg.MinRectSize+1)

	// Get position using color difference sampling
	x, y := g.findBestPosition(size)
	rect := image.Rect(x, y, x+size, y+size)

	// Get color with strict adherence to original
	avgColor := g.sampleStrictColor(rect)

	return rect, avgColor
}

func (g *Generator) findBestPosition(size int) (int, int) {
	bounds := g.original.Bounds()
	maxX := bounds.Max.X - size
	maxY := bounds.Max.Y - size

	// Try a few positions and pick the one with highest color difference
	bestDiff := -1.0
	bestX, bestY := 0, 0

	for i := 0; i < 10; i++ {
		x := g.rng.Intn(maxX)
		y := g.rng.Intn(maxY)

		diff := g.calculateColorDifference(image.Rect(x, y, x+size, y+size))
		if diff > bestDiff {
			bestDiff = diff
			bestX = x
			bestY = y
		}
	}

	return bestX, bestY
}

func (g *Generator) calculateColorDifference(rect image.Rectangle) float64 {
	origColor := g.sampleStrictColor(rect)
	currColor := g.sampleCurrentColor(rect)

	rDiff := float64(origColor.R) - float64(currColor.R)
	gDiff := float64(origColor.G) - float64(currColor.G)
	bDiff := float64(origColor.B) - float64(currColor.B)

	return math.Abs(rDiff) + math.Abs(gDiff) + math.Abs(bDiff)
}

func (g *Generator) sampleCurrentColor(rect image.Rectangle) color.RGBA {
	var rSum, gSum, bSum int64
	count := 0

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			c := g.current.RGBAAt(x, y)
			rSum += int64(c.R)
			gSum += int64(c.G)
			bSum += int64(c.B)
			count++
		}
	}

	return color.RGBA{
		R: uint8(rSum / int64(count)),
		G: uint8(gSum / int64(count)),
		B: uint8(bSum / int64(count)),
		A: 255,
	}
}

func (g *Generator) sampleStrictColor(rect image.Rectangle) color.RGBA {
	var rSum, gSum, bSum uint64
	count := 0

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			r, g, b, _ := g.original.At(x, y).RGBA()
			rSum += uint64(r)
			gSum += uint64(g)
			bSum += uint64(b)
			count++
		}
	}

	alpha := uint8(float64(g.cfg.AlphaMin) * float64(g.cfg.BlendingFactor))

	return color.RGBA{
		R: uint8((rSum / uint64(count)) >> 8),
		G: uint8((gSum / uint64(count)) >> 8),
		B: uint8((bSum / uint64(count)) >> 8),
		A: alpha,
	}
}

func (g *Generator) drawRectangle(rect image.Rectangle, c color.RGBA) {
	draw.Draw(
		g.current,
		rect,
		&image.Uniform{c},
		image.Point{},
		draw.Over,
	)
}

func (g *Generator) Result() image.Image {
	return g.current
}
