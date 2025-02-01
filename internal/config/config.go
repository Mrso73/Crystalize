package config

import (
    "encoding/json"
    "os"
)

type ArtConfig struct {
    InputPath  string         `json:"-"`
    OutputPath string         `json:"-"`

    // Rectangle parameters
    MinRectSize    int        `json:"minRectSize"`
    MaxRectSize    int        `json:"maxRectSize"`
    Iterations     int        `json:"iterations"`
     
    // Transparency controls
    AlphaMin       uint8      `json:"alphaMin"`
    AlphaMax       uint8      `json:"alphaMax"`
    BlendingFactor float32    `json:"blendingFactor"`
     
    // Color variation   
    ColorVariation uint8      `json:"colorVariation"`
    
    // Border styling
    BorderEnabled   bool      `json:"borderEnabled"`
    BorderColor     [3]uint8  `json:"borderColor"`
    BorderThickness int       `json:"borderThickness"`
}

func NewDefaultConfig() *ArtConfig {
    return &ArtConfig{
        MinRectSize:    10,     // Minimum rectangle width/height in pixels
        MaxRectSize:    30,     // Maximum rectangle width/height
        Iterations:     5000,   // Number of rectangles to draw
        AlphaMin:       30,     // Minimum transparency (0 = fully transparent)
        AlphaMax:       150,    // Maximum transparency (255 = fully opaque)
        BlendingFactor: 0.5,    // Opacity accumulation multiplier (0.1-1.0)
        ColorVariation: 10,     // Max color channel variation (±15)
        BorderEnabled:  true,
        BorderColor:    [3]uint8{20, 20, 20},
        BorderThickness: 1,
    }
}

func (c *ArtConfig) SaveToFile(path string) error {
    data, err := json.MarshalIndent(c, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(path, data, 0644)
}

func LoadFromFile(path string) (*ArtConfig, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    config := &ArtConfig{}
    if err := json.Unmarshal(data, config); err != nil {
        return nil, err
    }
    return config, nil
}