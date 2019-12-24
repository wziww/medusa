package encrpt

import (
	"testing"
)

var password []byte = []byte("AES256Key-32Characters1234567890")

var aesobj = (&AesGcm{&password, "", nil}).Construct("aes-256-gcm").(*AesGcm)

func TestString(t *testing.T) {
	s := "hellow world!"
	sd := aesobj.Encode([]byte(s))
	s2 := aesobj.Decode(sd)
	if s != string(s2) {
		t.Fatal(s, "!=", string(s2), "fail to encode and decode")
	}
}

func TestBytes(t *testing.T) {
	s := []byte{5, 1, 0, 1, 3, 4, 5, 7, 4, 3, 2, 2, 3, 5, 6, 0, 0, 0, 0, 0, 0, 9}
	sd := aesobj.Encode([]byte(s))
	s2 := aesobj.Decode(sd)
	for i := range s {
		if s[i] != s2[i] {
			t.Fatal(s, "!=", string(s2), "fail to encode and decode")
			return
		}
	}
}

func TestDecodeErrorData(t *testing.T) {
	s := []byte{141, 157, 142, 107, 1, 29, 217, 71, 14, 30, 214, 145, 91, 119,
		207, 69, 127, 232, 75, 185, 1, 172, 169, 27, 212, 174, 150, 72, 192, 10,
		133, 243, 172, 169, 190, 116, 46, 28, 12, 70, 132, 35, 28, 202, 122, 131,
		106, 58, 196, 96, 229, 204, 59, 138, 65, 43, 149, 37, 149, 200, 231, 80,
		144, 141, 27, 245, 233, 81, 227, 108, 51, 85, 185, 91, 71, 174, 123, 27,
		128, 207, 26, 225, 83, 182, 208, 139, 49, 253, 65, 254, 66, 84, 117, 127,
		159, 36, 84, 60, 190, 28, 46, 108, 121, 157, 90, 135, 173, 70, 119, 114,
		103, 201, 140, 125, 235, 8, 51, 237, 132, 28, 56, 129, 57, 33, 76, 43,
		57, 206, 190, 181, 144, 215, 22, 219, 31, 1, 167, 28, 48, 107, 160, 236,
		182, 181, 204, 6, 53, 212, 100, 20, 197, 182, 169, 196, 70, 114, 43, 90,
		152, 25, 18, 159, 50, 12, 196, 150, 86, 142, 197, 190, 66, 190, 111, 80,
		61, 121, 217, 73, 2, 142, 255, 13, 15, 212, 243, 228, 123, 20, 55, 12, 197,
		170, 98, 71, 142, 68, 11, 43, 71, 82, 213, 254, 236, 218, 20, 51, 108, 24,
		161, 111, 21, 52, 123, 198, 169, 154, 140, 52, 192, 39, 86, 83, 189, 49,
		31, 202, 122, 153, 254, 26, 16, 139, 217, 52, 4, 45, 177, 35, 236, 85, 146,
		96, 39, 157, 42, 135, 44, 125, 14, 209, 25, 29, 108, 86, 254, 68, 81, 202,
		53, 24, 197, 153, 253, 130, 35, 134, 180, 246, 137, 21, 130, 87, 231, 102,
		87, 141, 239, 89, 47, 209, 14, 31, 152, 245, 52, 193, 104, 107, 4, 100, 225,
		147, 237, 194, 66, 182, 125, 104, 185, 47, 64, 108, 239, 72, 205, 239, 197,
		244, 35, 34, 206, 178, 225, 49, 130, 27, 205, 43, 78, 194, 11, 136, 16, 202,
		202, 191, 106, 160, 168, 177, 55, 131, 86, 3, 123, 186, 126, 85, 91, 159, 118,
		157, 206, 40, 43, 22, 159, 204, 49, 40, 187, 95, 29, 208, 185, 44, 168, 54, 63,
		98, 191, 162, 142, 35, 181, 85, 109, 15, 30, 196, 100, 3, 237, 43, 40, 230, 59,
		243, 78, 207, 181, 251, 219, 100, 237, 36, 14, 125, 73, 71, 109, 192, 163, 235,
		17, 244, 86, 146, 87, 135, 175, 72, 58, 98, 50, 112, 174, 211, 236, 110, 86, 78,
		94, 37, 142, 153, 165, 215, 13, 201, 91, 138, 175, 120, 6, 234, 11, 91, 92, 231,
		99, 147, 101, 133, 231, 251, 67, 111, 230, 101, 247, 83, 154, 231, 166, 25, 188,
		37, 81, 80, 39, 213, 170, 101, 188, 182, 118, 196, 13, 147, 15, 15, 71, 185, 119,
		163, 74, 97, 184, 93, 155, 220, 121, 50, 164, 72, 90, 169, 186, 204, 222, 248, 123,
		222, 126, 38, 213, 213, 125, 187, 252, 255, 106, 211, 182, 13, 131, 177, 43, 142,
		214, 220, 63, 19, 23, 21, 218, 76, 208, 38, 50, 80, 138, 75, 192, 198, 77, 8, 183,
		35, 100, 57, 78, 239, 26, 174, 128, 255, 251, 93, 27, 7, 165, 246, 153, 194, 174,
		131, 104, 16, 35, 154, 30, 227, 58, 32, 118, 161, 86, 40, 127, 38, 114, 154, 168,
		226, 81, 52, 244, 126, 115, 72, 138, 51, 101, 17, 87, 162, 161, 78, 53, 84, 73,
		107, 237, 42, 197, 1, 191, 77, 2, 123, 184, 57, 71, 70, 255, 42, 57, 149, 188,
		75, 151, 219, 171, 138, 224, 239, 151, 177, 144, 139, 134, 90, 226, 170, 223,
		120, 75, 2, 147, 223, 140, 116, 219, 73, 59, 101, 8, 229, 63, 75, 92, 156, 18,
		54, 46, 204, 52, 72, 189, 153, 141, 232, 92, 15, 46, 235, 199, 141, 130, 171,
		207, 53, 37, 125, 206, 168, 255, 13, 149, 178, 105, 114, 84, 163, 110, 141,
		13, 80, 4, 91, 87, 234, 44, 71, 240, 17, 156, 232, 221, 31, 115, 55, 228, 16,
		188, 163, 102, 213, 11, 208, 27, 142, 128, 103, 90, 169, 10, 0, 198, 252, 95,
		56, 244, 35, 136, 120, 185, 101, 133, 12, 115, 143, 66, 255, 55, 155, 190, 82,
		27, 154, 229, 18, 232, 151, 161, 210, 155, 97, 154, 208, 129, 207, 142, 209, 23,
		131, 52, 243, 241, 196, 13, 152, 247, 80, 129, 33, 59, 112, 140, 43, 183, 211,
		8, 224, 40, 110, 172, 47, 22, 60, 39, 216, 56, 44, 104, 159, 95, 133, 236, 239,
		239, 151, 109, 58, 235, 26, 148, 175, 26, 86, 39, 60, 187, 54, 30, 175, 202,
		166, 176, 240, 3, 88, 25, 248, 169, 77, 105, 231, 243, 181, 246, 190, 60, 124,
		106, 221, 192, 226, 104, 7, 239, 220, 177, 237, 185, 160, 197, 71, 107, 0, 139,
		80, 134, 200, 164, 208, 123, 161, 61, 250, 164, 66, 234, 95, 194, 61, 80, 70,
		94, 235, 203, 44, 148, 220, 91, 40, 69, 117, 94, 189, 225, 99, 184, 155, 254,
		67, 195, 61, 56, 203, 142, 12, 229, 188, 120, 91, 233, 60, 79, 235, 230, 60,
		136, 59, 162, 201, 246, 91, 73, 77, 71, 50, 125, 12, 149, 224, 30, 28, 48, 110,
		181, 2, 77, 213, 10, 176, 175, 121, 48, 175, 169, 14, 98, 193, 151, 27, 141,
		183, 236, 154, 113, 162, 84, 167, 124, 200, 144, 206, 206, 210, 0, 185, 201,
		122, 81, 228, 226, 57, 144, 247, 52, 51, 201, 228, 12, 157, 193, 92, 85, 122,
		60, 186, 191, 105, 33, 82, 75, 248, 10, 109, 39, 199, 228, 100, 39, 168, 171,
		143, 77, 55, 10, 156, 20, 218, 142, 139, 223, 151, 86, 44, 154, 30, 197, 182,
		6, 98, 174, 122, 19, 246, 110, 135, 233, 224, 179, 100, 137, 77, 69, 182, 230,
		175, 69, 80, 139, 209, 200, 59, 172, 2, 4, 121, 47, 24, 102, 142, 236, 116, 210,
		101, 212, 31, 239, 80, 238, 217, 84, 78, 112, 4, 125, 17, 153, 43, 42, 123, 174,
		191, 59, 52, 85, 176}
	s2 := aesobj.Decode(s)
	if len(s2) != 0 {
		t.Fatal("test decode error data fail")
	}
}

func TestDecodeData(t *testing.T) {
	s := []byte{141, 157, 142, 107, 1, 29, 217, 71, 14, 30, 214, 31, 157, 186, 51, 226,
		104, 44, 184, 43, 169, 68, 113, 84, 100, 179, 85, 217, 70, 166, 199, 22, 107, 96,
		206, 26, 237, 137, 27, 26, 150, 214, 13, 202, 88, 214, 9, 49, 0, 94, 67, 21, 216,
		227, 217, 35, 7, 240, 79, 110, 226, 43, 34, 17, 160, 100, 179, 8, 35, 234, 193,
		116, 6, 21, 163, 185, 123, 27, 128, 207, 26, 225, 83, 182, 208, 139, 49, 253, 65,
		47, 49, 255, 52, 39, 9, 174, 135, 101, 255, 162, 50, 91, 111, 204, 209, 30, 62,
		93, 201, 30, 251, 52, 20, 79, 72, 216, 36, 247, 255, 54, 70, 129, 57, 33, 76, 43,
		57, 206, 190, 181, 144, 215, 22, 219, 31, 1, 160, 150, 183, 186, 30, 141, 32, 188,
		107, 88, 164, 73, 126, 252, 166, 149, 3, 96, 148, 253, 231, 212, 34, 254, 19, 55,
		46, 194, 44, 156, 245, 108, 26, 0, 171, 53, 224, 163, 45, 17, 231, 180, 198, 255,
		53, 36, 72, 48, 210, 62, 24, 82, 224, 168, 198, 90, 153, 206, 202, 249, 36, 137,
		146, 24, 98, 169, 62, 6, 76, 126, 223, 151, 39, 95, 62, 35, 129, 142, 141, 136, 74,
		184, 158, 39, 170, 156, 65, 249, 40, 80, 137, 216, 119, 152, 158, 225, 71, 123, 65,
		86, 109, 230, 16, 0, 210, 154, 95, 153, 101, 76, 50, 160, 93, 224, 63, 107, 33, 110,
		181, 248, 80, 196, 2, 33, 147, 214, 79, 89, 168, 81, 66, 142, 196, 62, 107, 0, 165,
		163, 78, 67, 191, 43, 140, 137, 175, 88, 106, 83, 1, 245, 24, 8, 102, 19, 248, 243,
		108, 135, 118, 28, 210, 171, 241, 95, 177, 29, 142, 231, 146, 255, 186, 41, 183, 20,
		131, 6, 50, 168, 96, 156, 182, 138, 216, 153, 28, 46, 28, 147, 219, 139, 104, 185,
		20, 18, 22, 60, 127, 130, 130, 251, 52, 192, 195, 26, 92, 155, 187, 85, 82, 191, 61,
		14, 145, 156, 144, 253, 229, 94, 107, 239, 232, 151, 84, 26, 41, 195, 42, 136, 60,
		209, 57, 22, 186, 167, 148, 147, 176, 33, 197, 0, 192, 71, 117, 243, 14, 87, 23, 91,
		104, 151, 231, 61, 228, 18, 250, 201, 47, 86, 146, 157, 195, 126, 205, 151, 222, 239,
		118, 154, 247, 179, 197, 31, 215, 23, 189, 23, 59, 95, 111, 215, 10, 50, 135, 68, 97,
		69, 200, 103, 86, 28, 189, 101, 212, 23, 160, 127, 3, 145, 243, 253, 223, 115, 23, 37,
		128, 72, 154, 63, 243, 124, 88, 148, 135, 29, 84, 193, 182, 248, 217, 103, 81, 97,
		184, 12, 241, 245, 179, 198, 146, 127, 203, 215, 221, 192, 15, 25, 167, 209, 88, 4,
		224, 100, 122, 94, 78, 0, 105, 240, 173, 54, 243, 162, 23, 198, 231, 0, 159, 132, 151,
		172, 250, 50, 59, 141, 226, 161, 80, 14, 93, 229, 129, 128, 94, 143, 179, 143, 3, 131,
		50, 118, 4, 198, 73, 110, 183, 125, 163, 211, 80, 148, 26, 214, 23, 136, 35, 100, 222,
		235, 112, 175, 168, 119, 214, 41, 147, 43, 146, 43, 179, 146, 151, 97, 52, 100, 48,
		7, 60, 1, 152, 30, 149, 12, 222, 88, 135, 136, 181, 87, 250, 216, 201, 194, 23, 112,
		72, 87, 224, 44, 68, 15, 218, 135, 43, 155, 61, 64, 126, 1, 115, 191, 252, 62, 83, 156,
		101, 11, 163, 211, 234, 104, 132, 16, 33, 254, 42, 194, 158, 236, 223, 53, 102, 204, 74,
		113, 65, 56, 136, 111, 138, 164, 218, 149, 221, 155, 254, 127, 254, 229, 183, 19, 145,
		208, 4, 18, 57, 18, 92, 93, 32, 70, 123, 251, 138, 69, 135, 149, 152, 142, 131, 218, 113,
		127, 18, 45, 236, 78, 254, 217, 27, 112, 10, 130, 84, 49, 66, 174, 148, 16, 132, 42, 151,
		19, 194, 133, 101, 167, 17, 34, 159, 84, 27, 134, 138, 203, 119, 65, 240, 185, 3, 130, 147,
		147, 87, 85, 222, 19, 17, 189, 240, 211, 98, 173, 240, 72, 128, 248, 178, 39, 81, 8, 14,
		235, 189, 205, 42, 61, 88, 236, 103, 231, 247, 218, 87, 46, 77, 48, 104, 38, 118, 227, 210,
		47, 251, 80, 17, 183, 217, 194, 213, 191, 129, 144, 48, 253, 58, 188, 178, 146, 32, 10, 22,
		190, 69, 31, 28, 143, 106, 55, 38, 6, 113, 52, 121, 144, 225, 1, 179, 248, 170, 205, 93,
		188, 221, 25, 62, 103, 15, 106, 70, 51, 244, 142, 89, 171, 19, 13, 197, 119, 206, 69, 221,
		182, 194, 113, 42, 66, 194, 235, 143, 116, 130, 149, 191, 129, 156, 47, 184, 153, 129, 215,
		59, 121, 140, 190, 65, 83, 108, 17, 126, 238, 21, 15, 182, 239, 9, 185, 223, 123, 55, 212,
		106, 142, 147, 115, 64, 87, 235, 164, 150, 195, 197, 98, 190, 51, 23, 143, 117, 54, 35,
		38, 65, 69, 9, 102, 236, 65, 31, 254, 123, 247, 109, 57, 124, 71, 28, 51, 227, 240, 180,
		131, 204, 44, 79, 146, 179, 235, 52, 198, 198, 19, 140, 28, 204, 31, 98, 155, 101, 161,
		245, 98, 60, 230, 72, 204, 208, 252, 50, 31, 38, 55, 158, 99, 26, 248, 43, 69, 144, 14, 255,
		210, 23, 20, 139, 59, 4, 54, 240, 130, 24, 63, 33, 79, 66, 240, 76, 107, 154, 249, 239, 200,
		138, 216, 138, 22, 254, 142, 109, 175, 208, 155, 68, 242, 246, 232, 102, 191, 219, 233, 43,
		30, 172, 147, 13, 115, 79, 12, 201, 15, 66, 119, 100, 227, 98, 27, 210, 168, 163, 13, 241,
		148, 201, 240, 222, 58, 158, 107, 51, 4, 244, 160, 176, 5, 58, 92, 59, 201, 194, 133, 197,
		14, 98, 144, 76, 118, 224, 153, 156, 185, 125, 28, 30, 208, 146, 128, 46, 115, 209, 227,
		31, 142, 131, 173, 97, 203, 163, 242, 89, 85, 225, 12, 152, 210, 230, 170, 82, 64, 9, 98, 53,
		115, 211, 94, 180, 44, 25, 226, 244, 216, 109, 3, 136, 204, 90, 149}
	s2 := aesobj.Decode(s)
	if len(s2) == 0 {
		t.Fatal("test decode data error")
	}
}
