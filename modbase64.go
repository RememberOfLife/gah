package gah

import (
	"encoding/base64"
)

var (
	// MB64E provides a modified encoding of base64 that should work even in urls and filenames
	MB64E = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+-")
)
