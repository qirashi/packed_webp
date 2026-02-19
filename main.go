// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Qirashi
// Project: packed_webp

package main

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		os.Stderr.WriteString("Usage: packed_webp [files & folders...]\n")
		os.Exit(1)
	}

	numWorkers := 2 // Можно использовать runtime.NumCPU() с импортом runtime
	jobs := make(chan string, numWorkers*2)
	wg := sync.WaitGroup{}

	for range numWorkers {
		wg.Go(func() {
			for file := range jobs {
				processFile(file)
			}
		})
	}

	for _, path := range os.Args[1:] {
		info, err := os.Stat(path)
		if err != nil {
			os.Stderr.WriteString("Error accessing path: " + path + "\n")
			continue
		}

		if info.IsDir() {
			filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
				if err != nil {
					os.Stderr.WriteString("Error accessing: " + p + " -> " + err.Error() + "\n")
					return nil
				}
				if !d.IsDir() {
					jobs <- p
				}
				return nil
			})
		} else {
			jobs <- path
		}
	}

	close(jobs)
	wg.Wait()
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

	case ext == ".webp" && !strings.Contains(baseName, ".packed"):
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
