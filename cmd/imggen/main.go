package main

import (
    "flag"
    "fmt"
    "image"
    "image/png"
    "log"
    "os"
    
    "github.com/Mrso73/Crystalize/internal/config"
    "github.com/Mrso73/Crystalize/internal/generator"
)

func main() {
    // Parse command line flags
    inputPath := flag.String("input", "", "Path to input image")
    outputPath := flag.String("output", "output.png", "Path to output image")
    configPath := flag.String("config", "", "Path to configuration file (optional)")
    saveConfig := flag.String("save-config", "", "Save default configuration to file")
    flag.Parse()

    // Handle configuration
    var cfg *config.ArtConfig
    var err error
    
    if *saveConfig != "" {
        cfg = config.NewDefaultConfig()
        if err := cfg.SaveToFile(*saveConfig); err != nil {
            log.Fatalf("Failed to save config: %v", err)
        }
        fmt.Println("Configuration saved to:", *saveConfig)
        return
    }
    
    if *configPath != "" {
        cfg, err = config.LoadFromFile(*configPath)
        if err != nil {
            log.Fatalf("Failed to load config: %v", err)
        }
    } else {
        cfg = config.NewDefaultConfig()
    }
    
    // Validate input
    if *inputPath == "" {
        log.Fatal("Input path is required")
    }
    
    // Load input image
    file, err := os.Open(*inputPath)
    if err != nil {
        log.Fatalf("Failed to open input file: %v", err)
    }
    defer file.Close()
    
    img, _, err := image.Decode(file)
    if err != nil {
        log.Fatalf("Failed to decode image: %v", err)
    }
    
    // Create and run generator
    gen := generator.NewGenerator(cfg, img)
    
    // Create progress channel and start progress reporting
    progress := make(chan float32)
    go func() {
        for p := range progress {
            fmt.Printf("\rProgress: %.1f%%", p*100)
        }
        fmt.Println()
    }()
    
    gen.Generate(progress)
    
    // Save result
    output, err := os.Create(*outputPath)
    if err != nil {
        log.Fatalf("Failed to create output file: %v", err)
    }
    defer output.Close()
    
    if err := png.Encode(output, gen.Result()); err != nil {
        log.Fatalf("Failed to encode output image: %v", err)
    }
    
    fmt.Println("Generation complete! Output saved to:", *outputPath)
}