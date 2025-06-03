package utils

import (
	"encoding/base64"
)

func ConvertImageToBase64(src []byte) string{
	str := base64.StdEncoding.EncodeToString(src)
	return str
}
