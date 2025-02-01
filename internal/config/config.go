package config

import (
	"encoding/json"
	"os"
)

type ArtConfig struct {
	InputPath  string `json:"-"`
	OutputPath string `json:"-"`

	// Rectangle parameters
	MinRectSize int `json:"minRectSize"`
	MaxRectSize int `json:"maxRectSize"`
	Iterations  int `json:"iterations"`

	// Transparency controls
	AlphaMin       uint8   `json:"alphaMin"`
	AlphaMax       uint8   `json:"alphaMax"`
	BlendingFactor float32 `json:"blendingFactor"`

	// Color variation
	ColorVariation uint8 `json:"colorVariation"`

	// Border styling
	BorderEnabled   bool     `json:"borderEnabled"`
	BorderColor     [3]uint8 `json:"borderColor"`
	BorderThickness int      `json:"borderThickness"`
}

func NewDefaultConfig() *ArtConfig {
	return &ArtConfig{
		MinRectSize:     10,
		MaxRectSize:     30,
		Iterations:      5000,
		AlphaMin:        20,
		AlphaMax:        60,
		BlendingFactor:  0.3,
		ColorVariation:  2,
		BorderEnabled:   false,
		BorderColor:     [3]uint8{20, 20, 20},
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
