// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Qirashi
// Project: packed_webp

package main

import (
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		os.Stderr.WriteString("Usage: packed_webp <files>\n")
		os.Exit(1)
	}

	for _, path := range os.Args[1:] {
		_, err := os.Stat(path)
		if err != nil {
			os.Stderr.WriteString("Error accessing path: " + path + "\n")
			continue
		}
		processFile(path)
	}
}

func processFile(filePath string) {
    baseName := filepath.Base(filePath)
    ext := filepath.Ext(filePath)

    switch {
    case strings.HasSuffix(baseName, ".packed.webp"):
        if err := UnpackWebp(filePath); err != nil {
            os.Stderr.WriteString("Unpack error [" + filePath + "]: " + err.Error() + "\n")
        } else {
            newName := strings.TrimSuffix(filePath, ".packed.webp") + ".webp"
            os.Stderr.WriteString("Successfully unpacked: " + filePath + " -> " + newName + "\n")
        }

    case (ext == ".txt" || ext == ".webp") && !strings.Contains(baseName, ".packed"):
        if ext == ".txt" {
            filePath = strings.TrimSuffix(filePath, ".txt") + ".webp"
        }

        if err := PackWebp(filePath); err != nil {
            os.Stderr.WriteString("Pack error [" + filePath + "]: " + err.Error() + "\n")
        } else {
            packedName := strings.TrimSuffix(filePath, ".webp") + ".packed.webp"
            os.Stderr.WriteString("Successfully packed: " + filePath + " -> " + packedName + "\n")
        }

    default:
        os.Stderr.WriteString("Skipping unsupported file: " + filePath + "\n")
    }
}
