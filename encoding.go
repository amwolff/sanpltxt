package sanpltxt

import "golang.org/x/text/encoding/charmap"

// ToWindows1250 converts a UTF-8 string to Windows-1250 encoding.
// The Santander import format requires Windows-1250 encoding for Polish characters.
func ToWindows1250(s string) ([]byte, error) {
	return charmap.Windows1250.NewEncoder().Bytes([]byte(s))
}

// FromWindows1250 converts Windows-1250 encoded bytes to a UTF-8 string.
func FromWindows1250(b []byte) (string, error) {
	decoded, err := charmap.Windows1250.NewDecoder().Bytes(b)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
