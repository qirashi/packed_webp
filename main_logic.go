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

func replaceExt(filename, oldExt, newExt string) string {
	return strings.TrimSuffix(filename, oldExt) + newExt
}

func logErr(err error, msg string) error {
	if err != nil {
		os.Stderr.WriteString(msg + ": " + err.Error() + "\n")
	}
	return err
}

func removeFiles(files ...string) error {
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			return logErr(err, "Cleanup failed")
		}
	}
	return nil
}

func UnpackWebp(packedFile string) error {
	data, err := os.ReadFile(packedFile)
	if err != nil {
		return logErr(err, "Read error")
	}

	extrIndex := bytes.LastIndex(data, []byte(ExtrSignature))
	if extrIndex == -1 {
		os.Stderr.WriteString("extr signature not found")
		return nil
	}

	if len(data) < extrIndex+8 || len(data) < extrIndex+8+int(binary.LittleEndian.Uint32(data[extrIndex+4:extrIndex+8])) {
		os.Stderr.WriteString("invalid file structure")
		return nil
	}

	extraData := data[extrIndex+8 : extrIndex+8+int(binary.LittleEndian.Uint32(data[extrIndex+4:extrIndex+8]))]
	
	txtFile := replaceExt(packedFile, ".packed.webp", ".txt")
	if err := os.WriteFile(txtFile, extraData, 0644); err != nil {
		return logErr(err, "Txt write failed")
	}

	newWebpData := data[:extrIndex]
	if !bytes.HasPrefix(newWebpData, []byte(RiffHeader)) {
		os.Stderr.WriteString("invalid WEBP header")
		return nil
	}

	binary.LittleEndian.PutUint32(newWebpData[4:8], uint32(len(newWebpData)-8))
	if err := os.WriteFile(replaceExt(packedFile, ".packed.webp", ".webp"), newWebpData, 0644); err != nil {
		return logErr(err, "Webp write failed")
	}

	return removeFiles(packedFile)
}

func PackWebp(webpFile string) error {
	txtFile := replaceExt(webpFile, ".webp", ".txt")
	extraData, err := os.ReadFile(txtFile)
	if err != nil {
		return logErr(err, "Txt read error")
	}

	webpData, err := os.ReadFile(webpFile)
	if err != nil {
		return logErr(err, "Read error")
	}

	if !bytes.HasPrefix(webpData, []byte(RiffHeader)) {
		os.Stderr.WriteString("invalid RIFF header")
		return nil
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
		return logErr(err, "Write error")
	}

	return removeFiles(webpFile, txtFile)
}