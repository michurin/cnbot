package datatype

import "bytes"

var (
	fpGif   = []byte{0x47, 0x49, 0x46, 0x38}
	fpPng   = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	fpJpegA = []byte{0xFF, 0xD8, 0xFF, 0xDB}
	fpJpegB = []byte{0xFF, 0xD8, 0xFF, 0xE0}
	fpJpegC = []byte{0xFF, 0xD8, 0xFF, 0xE1}
)

func ImageType(data []byte) string {
	if bytes.HasPrefix(data, fpGif) {
		return "gif"
	} else if bytes.HasPrefix(data, fpPng) {
		return "png"
	} else if bytes.HasPrefix(data, fpJpegA) ||
		bytes.HasPrefix(data, fpJpegB) ||
		bytes.HasPrefix(data, fpJpegC) {
		return "jpeg"
	}
	return ""
}
