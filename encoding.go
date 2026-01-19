package sanpltxt

import "golang.org/x/text/encoding/charmap"

// ToWindows1250 encodes a UTF-8 string to Windows-1250.
func ToWindows1250(s string) ([]byte, error) {
	return charmap.Windows1250.NewEncoder().Bytes([]byte(s))
}

// FromWindows1250 decodes Windows-1250 bytes to a UTF-8 string.
func FromWindows1250(b []byte) (string, error) {
	decoded, err := charmap.Windows1250.NewDecoder().Bytes(b)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
