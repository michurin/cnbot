package helpers

import (
	"net/http"
)

func ImageType(data []byte) string {
	tp := http.DetectContentType(data)
	switch tp {
	case "image/png", "image/jpeg", "image/gif":
		return tp[6:]
	}
	return ""
}
