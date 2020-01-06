package encrypt

import "github/wziww/medusa/encrypt/padding"

// HandlePadding ...
func HandlePadding(paddingMode string) func([]byte, int) []byte {
	switch paddingMode {
	case "ANSI X9.23":
		return padding.ANSIX923Padding
	case "ISO 10126":
		return padding.ISO10126Padding
	case "PKCS#7":
		return padding.PKCS7Padding
	case "ISO/IEC 7816-4":
		return padding.ISOIEC7816_4Padding
	case "Zero":
		return padding.ZeroPadding
	default:
		return func(buf []byte, size int) []byte {
			return buf
		}
	}
}

// HandleUnPadding ...
func HandleUnPadding(paddingMode string) func([]byte, int) ([]byte, error) {
	switch paddingMode {
	case "ANSI X9.23":
		return padding.ANSIX923UnPadding
	case "ISO 10126":
		return padding.ISO10126UnPadding
	case "PKCS#7":
		return padding.PKCS7UnPadding
	case "ISO/IEC 7816-4":
		return padding.ISOIEC7816_4UnPadding
	case "Zero":
		return padding.ZeroUnPadding
	default:
		return func(buf []byte, size int) ([]byte, error) {
			return buf, nil
		}
	}
}
