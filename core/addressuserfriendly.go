package core

import (
	"strconv"
	"bytes"
	bu "github.com/terorie/go-nimiq/bufferutils"
)

// The Nimiq user friendly address format is a 36 character
// ASCII string separated by spaces.
// Spaces (0x20) are not counted as characters and will be ignored.
// Characters 0-2 (country code) must be "NQ".
// Characters 2-4 are a IBAN-like checksum of the following chars.
// Characters 4-36 are Base32 (bufferutils/base32nimiq.go) and
//     encode the 20 byte Address type.
// Example: NQ19 46LK 9YHV D9LB TDJ4 8Y2P 3J4C 37HR YDL5


// Decodes the user friendly address (ufs) into *a.
// ufs without spaces must have a length of 36.
func (a *Address) FromUserFriendly(ufs string) error {
	// Encoded in ASCII (as opposed to the Unicode-encoded string type)
	var uf [36]byte

	// Remove all spaces from ufs and write to uf
	err := truncateWhitespaceFromUserFriendly(&uf, ufs)
	if err != nil { return err }

	// Verify string: Country code
	if uf[0] != 'N' && uf[1] != 'Q' {
		return ErrUserFriendlyAddress_InvalidCountryCode
	}

	// Verify string: Checksum and characters
	check, err := calculateUserFriendlyAddressCheck(&uf)
	if err != nil { return err }

	if check != 1 {
		return ErrUserFriendlyAddress_Checksum
	}

	// Decode everything past the first 4 chars (data region)
	_, err = bu.Base32NimiqEncoder.Decode(a[:], uf[4:])
	if err != nil { return err }

	return nil
}

// Encodes Address *a into a user friendly string (ASCII)
func (a *Address) ToUserFriendly() (string, error) {
	// User friendly string without spaces
	var noSpaces [36]byte

	// Use "00" as temporary checksum
	copy(noSpaces[0:4], "NQ00")
	// Encode the rest of the address
	bu.Base32NimiqEncoder.Encode(noSpaces[4:], a[:])

	// Checksum is (98 - "current check result")
	check, err := calculateUserFriendlyAddressCheck(&noSpaces)
	check = 98 - check
	if err != nil { return "", err }

	// Allocate a new buffer with spaces
	var withSpaces [44]byte
	withSpacesWtr := bu.WriteBytes(withSpaces[:])
	withSpacesWtr.WriteNext([]byte("NQ"))

	// Write checksum digits to buffer
	// (9 + 0x30 == '9') turns digit to ASCII character
	withSpacesWtr.WriteNext([]byte{
		0x30 + (uint8(check % 100) / 10), // Digit X of XY
		0x30 + uint8(check % 10),         // Digit Y of XY
	})

	// Copy noSpaces to withSpaces
	for i := 4; i < 36; i += 4 {
		// Insert a space every 4 characters
		withSpacesWtr.Uint8(' ')
		withSpacesWtr.WriteNext(noSpaces[i : i+4])
	}

	return string(withSpaces[:]), nil
}

// Removes spaces from ufs and writes it into uf
// Fails if length of ufs (as ASCII) != 36
func truncateWhitespaceFromUserFriendly(uf *[36]byte, ufs string) error {
	// Treat ufs as an ASCII string
	ufsb := []byte(ufs)

	// Moving slices pointer arithmetic
	suf := uf[:]
	sufsb := ufsb

	for {
		if len(sufsb) == 0 { // End of ufs reached
			if len(suf) == 0 { // End of uf reached
				// Read successful
				return nil
			} else {
				return ErrUserFriendlyAddress_InvalidLength
			}
		} else if len(suf) == 0 {
			return ErrUserFriendlyAddress_InvalidLength
		}

		char := sufsb[0]
		if char != 0x20 {
			// If not space write to uf
			suf[0] = char
			suf = suf[1:]
		}

		sufsb = sufsb[1:]
	}
}

// Calculates a checksum over the UFA
// Returns values from 0-98
func calculateUserFriendlyAddressCheck(userFriendly *[36]byte) (uint8, error) {
	var sumBuffer bytes.Buffer

	// Writes a bunch of chars to sumBuffer for check loop later
	// e.g. "6789ABcd" will append "678910111213" to the sumBuffer
	nextChars := func(slice []byte) error {
		for _, char := range slice {
			switch {

			// Lower case character
			case char > 0x60 && char <= 0x7A:
				char -= 0x20 // convert to upper
				fallthrough

			// Upper case character
			case char > 0x40 && char <= 0x5A:
				// Subtract 0x37
				// Rationale: Numerical values are 0-9
				// Subtracting 0x37 will make letters to number 10…
				num := char - 0x37
				// Represent num as a decimal number
				numStr := strconv.FormatUint(uint64(num), 10)
				// Append to sum buffer
				sumBuffer.WriteString(numStr)
				break

			// Numerical character
			case char >= 0x30 && char <= 0x39:
				// Append numerical character (not code) to buffer
				sumBuffer.WriteByte(char)
				break

			// Unknown character
			default:
				return ErrUserFriendlyAddress_InvalidCharacter
				break

			}
		}

		return nil
	}

	// Iterate over the chars in addrBuf, the first 4 chars at last
	if err := nextChars(userFriendly[4:]); err != nil { return 0, err }
	if err := nextChars(userFriendly[0:4]); err != nil { return 0, err }

	// Extract string (here as bytes) from sumBuffer
	sum := sumBuffer.Bytes()

	// Create new buffer for tmp variable
	var tmpBuffer bytes.Buffer

	// Rounding up division
	blockCount := (len(sum) + 5) / 6

	// Iterate over sum in blocks of 6
	for i := 0; true; i++ {
		offset := i * 6
		var stop int
		if len(sum) <= offset + 6 {
			stop = len(sum)
		} else {
			stop = offset + 6
		}

		block := sum[offset:stop]

		// Append the block to the buffer
		tmpBuffer.Write(block)

		// Read string out of buffer
		tmp := tmpBuffer.String()

		// Convert to integer, ignore errors
		tmpNum, _ := strconv.ParseUint(tmp, 10, 64)
		tmpNum %= 97 // magic :P

		if (i+1) < blockCount {
			// Put string back into buffer
			tmpBuffer.Reset()
			tmpBuffer.WriteString(strconv.FormatUint(tmpNum, 10))
		} else {
			// No more blocks, return final value
			return uint8(tmpNum), nil
		}
	}

	// Unreachable
	panic(0)
	return 0, nil
}


// Error codes
type UserFriendlyAddressError uint8

const (
	_ = UserFriendlyAddressError(iota)
	ErrUserFriendlyAddress_InvalidCountryCode
	ErrUserFriendlyAddress_InvalidLength
	ErrUserFriendlyAddress_Checksum
	ErrUserFriendlyAddress_InvalidCharacter
)

func (u UserFriendlyAddressError) Error() string {
	switch u {
	case ErrUserFriendlyAddress_InvalidCountryCode:
		return "invalid address: should start with \"NQ\""
	case ErrUserFriendlyAddress_InvalidLength:
		return "invalid address: length without spaces not 36"
	case ErrUserFriendlyAddress_Checksum:
		return "invalid address: invalid checksum"
	case ErrUserFriendlyAddress_InvalidCharacter:
		return "invalid address: unexpected character"
	default:
		return ""
	}
}
