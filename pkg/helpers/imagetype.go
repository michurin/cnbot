package helpers

import (
	"bytes"
)

var ( /* const */
	// Details: https://en.wikipedia.org/wiki/List_of_file_signatures
	fpGif  = []byte{0x47, 0x49, 0x46, 0x38}
	fpPng  = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	fpJpgA = []byte{0xFF, 0xD8, 0xFF, 0xDB}
	fpJpgB = []byte{0xFF, 0xD8, 0xFF, 0xE0}
	fpJpgC = []byte{0xFF, 0xD8, 0xFF, 0xE1}
)

func ImageType(data []byte) (imgExt string) {
	if bytes.HasPrefix(data, fpGif) {
		imgExt = "gif"
		return
	}
	if bytes.HasPrefix(data, fpPng) {
		imgExt = "png"
		return
	}
	if bytes.HasPrefix(data, fpJpgA) ||
		bytes.HasPrefix(data, fpJpgB) ||
		bytes.HasPrefix(data, fpJpgC) {
		imgExt = "jpeg"
		return
	}
	return
}
