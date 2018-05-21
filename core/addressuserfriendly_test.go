package core

import (
	"testing"
	"bytes"
)

func TestAddress_FromUserFriendly(t *testing.T) {
	var addr Address

	// Correct decoded
	correctAddr := Address{0x21, 0xa9, 0x34, 0xfe, 0x3d, 0x6a, 0x68, 0xbd, 0xb6, 0x44, 0x47, 0xc5, 0x71, 0xc8, 0x8c, 0x19, 0xe3, 0x9f, 0xb6, 0x85}

	// 1. Test if conversion with works
	err1 := addr.FromUserFriendly("NQ1946LK9YHVD9LBTDJ48Y2P3J4C37HRYDL5")
	if err1 != nil {
		t.Errorf("Decoding errored: %s\n", err1)
	} else if bytes.Compare(correctAddr[:], addr[:]) != 0 {
		t.Errorf("Invalid result decoded: %p\n", addr)
	}

	// 2. Test if conversion with spaces works
	err2 := addr.FromUserFriendly("NQ19 46LK 9YHV D9LB TDJ4 8Y2P 3J4C 37HR YDL5")
	if err2 != nil {
		t.Errorf("Decoding with spaces errored: %s\n", err2)
	} else if bytes.Compare(correctAddr[:], addr[:]) != 0 {
		t.Errorf("Invalid result decoded with spaces: %p\n", addr)
	}

	// 3. Test if detecting errors in checksum works
	err3 := addr.FromUserFriendly("NQ20 46LK 9YHV D9LB TDJ4 8Y2P 3J4C 37HR YDL5")
	if err3 != ErrUserFriendlyAddress_Checksum {
		t.Error("Wrong checksum ignored!")
	}

	// 4. Test if other country codes are accepted
	err4 := addr.FromUserFriendly("LD20 46LK 9YHV D9LB TDJ4 8Y2P 3J4C 37HR YDL5")
	if err4 != ErrUserFriendlyAddress_InvalidCountryCode {
		t.Error("Other country code than \"NQ\" accepted")
	}

	// 5. Test if invalid Base32 is rejected
	err5 := addr.FromUserFriendly("NQ20 46LK 9YHV D9LB IIII 8Y2P 3J4C 37HR YDL5")
	if err5 == nil {
		t.Error("Invalid Base32 accepted!")
	}

	// 6. Test if detecting errors in length works (too short)
	err6 := addr.FromUserFriendly("NQ20 46LK 9YHV D9LB IIII 8Y2P 3J4C 37HR YDL")
	if err6 != ErrUserFriendlyAddress_InvalidLength {
		t.Error("Length not checked")
	}

	// 7. Test if detecting errors in length works (too long)
	err7 := addr.FromUserFriendly("NQ20 46LK 9YHV D9LB IIII 8Y2P 3J4C 37HR YDL5 A")
	if err7 != ErrUserFriendlyAddress_InvalidLength {
		t.Error("Length not checked")
	}

	// 8. Check reaction on empty input
	err8 := addr.FromUserFriendly("                                            ")
	if err8 != ErrUserFriendlyAddress_InvalidLength {
		t.Error("Length not checked")
	}

	// 9. Test if Unicode is rejected / check for buffer overflow
	//    (А is a cyrillic character here and needs >1 byte, google it)
	err9 := addr.FromUserFriendly("NQ20 ААB 9YHV D9LB IIII 8Y2P 3J4C 37HR YDL")
	if err9 != ErrUserFriendlyAddress_InvalidCharacter && err9 != ErrUserFriendlyAddress_InvalidLength {
		t.Errorf("Non-ASCII character accepted: %p\n", addr)
	}
}

func TestAddress_ToUserFriendly(t *testing.T) {
	addrInput := Address{0x21, 0xa9, 0x34, 0xfe, 0x3d, 0x6a, 0x68, 0xbd, 0xb6, 0x44, 0x47, 0xc5, 0x71, 0xc8, 0x8c, 0x19, 0xe3, 0x9f, 0xb6, 0x85}
	correctOutput := "NQ19 46LK 9YHV D9LB TDJ4 8Y2P 3J4C 37HR YDL5"

	// 1. Test if conversion works
	str1, err1 := addrInput.ToUserFriendly()
	if err1 != nil {
		t.Errorf("Encoding errored: %s\n", err1)
	} else if str1 != correctOutput {
		t.Errorf("Invalid user friendly address encoded: %s\n", str1)
	}
}
