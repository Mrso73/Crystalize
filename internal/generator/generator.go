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
}

func NewGenerator(cfg *config.ArtConfig, img image.Image) *Generator {
    bounds := img.Bounds()
    current := image.NewRGBA(bounds)
    draw.Draw(current, bounds, image.NewUniform(color.White), image.Point{}, draw.Src)
    return &Generator{
        cfg:      cfg,
        original: img,
        current:  current,
        rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
    }
}

func (g *Generator) Generate(progress chan<- float32) {
    defer close(progress)
    for i := 0; i < g.cfg.Iterations; i++ {
        rect, rectColor := g.createRectangle()
        g.drawRectangle(rect, rectColor)
        progress <- float32(i+1)/float32(g.cfg.Iterations)
    }
}

func (g *Generator) createRectangle() (image.Rectangle, color.RGBA) {
    // Random rectangle size with config constraints
    size := g.rng.Intn(g.cfg.MaxRectSize - g.cfg.MinRectSize + 1) + g.cfg.MinRectSize
    
    // Ensure rectangle stays within image bounds
    maxX := g.original.Bounds().Max.X - size
    maxY := g.original.Bounds().Max.Y - size
    x := g.rng.Intn(maxX)
    y := g.rng.Intn(maxY)
    
    rect := image.Rect(x, y, x+size, y+size)
    avgColor := g.sampleAverageColor(rect)
    
    // Apply color variation
    avgColor.R = g.applyVariation(avgColor.R)
    avgColor.G = g.applyVariation(avgColor.G)
    avgColor.B = g.applyVariation(avgColor.B)
    
    // Calculate alpha with blending factor
    alphaBase := uint8(g.rng.Intn(int(g.cfg.AlphaMax-g.cfg.AlphaMin+1))) + g.cfg.AlphaMin
    alpha := uint8(float32(alphaBase) * g.cfg.BlendingFactor)
    alpha = uint8(math.Max(float64(g.cfg.AlphaMin), math.Min(float64(g.cfg.AlphaMax), float64(alpha))))
    
    return rect, color.RGBA{
        R: avgColor.R,
        G: avgColor.G,
        B: avgColor.B,
        A: alpha,
    }
}

func (g *Generator) applyVariation(base uint8) uint8 {
    variation := int(g.rng.Intn(int(g.cfg.ColorVariation)*2+1)) - int(g.cfg.ColorVariation)
    newVal := int(base) + variation
    return uint8(math.Max(0, math.Min(255, float64(newVal))))
}

func (g *Generator) sampleAverageColor(rect image.Rectangle) color.RGBA {
    var rSum, gSum, bSum uint32 // Changed variable names
    count := 0
    for y := rect.Min.Y; y < rect.Max.Y; y++ {
        for x := rect.Min.X; x < rect.Max.X; x++ {
            pr, pg, pb, _ := g.original.At(x, y).RGBA() // Now correctly references Generator
            rSum += pr >> 8
            gSum += pg >> 8  // Now using gSum instead of g
            bSum += pb >> 8
            count++
        }
    }
    return color.RGBA{
        R: uint8(rSum / uint32(count)),
        G: uint8(gSum / uint32(count)),  // Fixed reference
        B: uint8(bSum / uint32(count)),
        A: 255,
    }
}

func (g *Generator) drawRectangle(rect image.Rectangle, c color.RGBA) {
    // Draw semi-transparent rectangle
    draw.Draw(
        g.current,
        rect,
        &image.Uniform{c},
        image.Point{},
        draw.Over,
    )
    
    // Draw border if enabled
    if g.cfg.BorderEnabled {
        borderColor := color.RGBA{
            R: g.cfg.BorderColor[0],
            G: g.cfg.BorderColor[1],
            B: g.cfg.BorderColor[2],
            A: 255,
        }
        drawBorder(g.current, rect, borderColor, g.cfg.BorderThickness)
    }
}

func drawBorder(img *image.RGBA, rect image.Rectangle, c color.RGBA, thickness int) {
    if thickness <= 0 {
        return
    }
    
    // Top border
    top := image.Rect(
        rect.Min.X,
        rect.Min.Y,
        rect.Max.X,
        rect.Min.Y+thickness,
    )
    draw.Draw(img, top, &image.Uniform{c}, image.Point{}, draw.Over)
    
    // Bottom border
    bottom := image.Rect(
        rect.Min.X,
        rect.Max.Y-thickness,
        rect.Max.X,
        rect.Max.Y,
    )
    draw.Draw(img, bottom, &image.Uniform{c}, image.Point{}, draw.Over)
    
    // Left border
    left := image.Rect(
        rect.Min.X,
        rect.Min.Y,
        rect.Min.X+thickness,
        rect.Max.Y,
    )
    draw.Draw(img, left, &image.Uniform{c}, image.Point{}, draw.Over)
    
    // Right border
    right := image.Rect(
        rect.Max.X-thickness,
        rect.Min.Y,
        rect.Max.X,
        rect.Max.Y,
    )
    draw.Draw(img, right, &image.Uniform{c}, image.Point{}, draw.Over)
}

func (g *Generator) Result() image.Image {
    return g.current
}