// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Qirashi
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

func replaceExt(filename, oldExt, newExt string) string {
	return strings.TrimSuffix(filename, oldExt) + newExt
}

func removeFiles(files ...string) error {
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}

func UnpackWebp(packedFile string) error {
	data, err := os.ReadFile(packedFile)
	if err != nil {
		return err
	}

	extrIndex := bytes.LastIndex(data, []byte(ExtrSignature))
	if extrIndex == -1 {
		return os.ErrInvalid
	}

	if len(data) < extrIndex+8 || len(data) < extrIndex+8+int(binary.LittleEndian.Uint32(data[extrIndex+4:extrIndex+8])) {
		return os.ErrInvalid
	}

	extraData := data[extrIndex+8 : extrIndex+8+int(binary.LittleEndian.Uint32(data[extrIndex+4:extrIndex+8]))]

	txtFile := replaceExt(packedFile, ".packed.webp", ".txt")
	if err := os.WriteFile(txtFile, extraData, 0644); err != nil {
		return err
	}

	newWebpData := data[:extrIndex]
	if !bytes.HasPrefix(newWebpData, []byte(RiffHeader)) {
		return os.ErrInvalid
	}

	binary.LittleEndian.PutUint32(newWebpData[4:8], uint32(len(newWebpData)-8))
	if err := os.WriteFile(replaceExt(packedFile, ".packed.webp", ".webp"), newWebpData, 0644); err != nil {
		return err
	}

	return removeFiles(packedFile)
}

func PackWebp(webpFile string) error {
	txtFile := replaceExt(webpFile, ".webp", ".txt")
	extraData, err := os.ReadFile(txtFile)
	if err != nil {
		return err
	}

	webpData, err := os.ReadFile(webpFile)
	if err != nil {
		return err
	}

	if !bytes.HasPrefix(webpData, []byte(RiffHeader)) {
		return os.ErrInvalid
	}

	buf := bytes.NewBuffer(webpData)
	buf.Write([]byte(ExtrSignature))
	binary.Write(buf, binary.LittleEndian, uint32(len(extraData)))
	buf.Write(extraData)

	if (buf.Len()-8)%2 != 0 {
		buf.WriteByte(0)
	}

	binary.LittleEndian.PutUint32(buf.Bytes()[4:8], uint32(buf.Len()-8))

	if err := os.WriteFile(replaceExt(webpFile, ".webp", ".packed.webp"), buf.Bytes(), 0644); err != nil {
		return err
	}

	return removeFiles(webpFile, txtFile)
}
