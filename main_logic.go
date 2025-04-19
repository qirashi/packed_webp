// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Qirashi
// Project: packed_webp

package main

import (
	"bytes"
	"encoding/binary"
	"os"
	"strings"
)

const (
	ExtrSignature = "extr"
	RiffHeader    = "RIFF"
)

type errorString string

func (e errorString) Error() string {
	return string(e)
}

func UnpackWebp(packedFile string) error {
	data, err := os.ReadFile(packedFile)
	if err != nil {
		os.Stderr.WriteString("Read error: " + err.Error() + "\n")
		return err
	}

	extrIndex := bytes.LastIndex(data, []byte(ExtrSignature))
	if extrIndex == - 1 {
		return errorString("extr signature not found")
	}

	if len(data) < extrIndex + 8 {
		return errorString("invalid file structure")
	}

	dataLen := int(binary.LittleEndian.Uint32(data[extrIndex + 4: extrIndex + 8]))

	if len(data) < extrIndex + 8 + dataLen {
		return errorString("corrupted extra data")
	}

	extraData := data[extrIndex + 8: extrIndex + 8 + dataLen]

	txtFile := strings.TrimSuffix(packedFile, ".packed.webp") + ".txt"
	if err := os.WriteFile(txtFile, extraData, 0644); err != nil {
		os.Stderr.WriteString("Txt write failed: " + err.Error() + "\n")
		return err
	}

	newWebpData := data[:extrIndex]

	if !bytes.HasPrefix(newWebpData, []byte(RiffHeader)) {
		return errorString("invalid WEBP header")
	}
	riffSize := len(newWebpData) - 8
	binary.LittleEndian.PutUint32(newWebpData[4:8], uint32(riffSize))

	newWebpName := strings.TrimSuffix(packedFile, ".packed.webp") + ".webp"
	if err := os.WriteFile(newWebpName, newWebpData, 0644); err != nil {
		os.Stderr.WriteString("Webp write failed: " + err.Error() + "\n")
		return err
	}

	if err := os.Remove(packedFile); err != nil {
		os.Stderr.WriteString("Cleanup failed: " + err.Error() + "\n")
		return err
	}

	return nil
}

func PackWebp(webpFile string) error {
	txtFile := strings.TrimSuffix(webpFile, ".webp") + ".txt"

	if _, err := os.Stat(txtFile); os.IsNotExist(err) {
		os.Stderr.WriteString("Txt file not found: " + err.Error() + "\n")
		return err
	}

	extraData, err := os.ReadFile(txtFile)
	if err != nil {
		os.Stderr.WriteString("Read error: " + err.Error() + "\n")
		return err
	}

	webpData, err := os.ReadFile(webpFile)
	if err != nil {
		os.Stderr.WriteString("Read error: " + err.Error() + "\n")
		return err
	}

	if !bytes.HasPrefix(webpData, []byte(RiffHeader)) {
		return errorString("invalid RIFF header")
	}

	var buffer bytes.Buffer
	buffer.Write(webpData)

	buffer.WriteString(ExtrSignature)
	lenBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(lenBuf, uint32(len(extraData)))
	buffer.Write(lenBuf)
	buffer.Write(extraData)

	totalSize := buffer.Len() - 8
	if totalSize % 2 != 0 {
		buffer.WriteByte(0x00)
	}

	riffSize := buffer.Len() - 8
	binary.LittleEndian.PutUint32(buffer.Bytes()[4:8], uint32(riffSize))

	packedFile := strings.TrimSuffix(webpFile, ".webp") + ".packed.webp"
	if err := os.WriteFile(packedFile, buffer.Bytes(), 0644); err != nil {
		os.Stderr.WriteString("Write error: " + err.Error() + "\n")
		return err
	}

	if err := os.Remove(webpFile); err != nil {
		os.Stderr.WriteString("Cleanup failed: " + err.Error() + "\n")
		return err
	}

	if err := os.Remove(txtFile); err != nil {
		os.Stderr.WriteString("Cleanup failed: " + err.Error() + "\n")
		return err
	}

	return nil
}