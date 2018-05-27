package core

import (
	"testing"
	"bytes"
)

func TestBasicTx_Deserialize(t *testing.T) {
	buf := []byte {
		// Tx Format (Basic)
		0x00,
		// SenderPubKey
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F,
		// Recipient
		0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39,
		0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49,
		// Value
		0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57,
		// Fee
		0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67,
		// ValidityStartHeight
		0x00, 0x00, 0x00, 0x00,
		// NetworkId
		42,
		// Signature
		0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7A, 0x7B, 0x7C, 0x7D, 0x7E, 0x7F,
		0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8A, 0x8B, 0x8C, 0x8D, 0x8E, 0x8F,
		0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9A, 0x9B, 0x9C, 0x9D, 0x9E, 0x9F,
		0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0xA6, 0xA7, 0xA8, 0xA9, 0xAA, 0xAB, 0xAC, 0xAD, 0xAE, 0xAF,
	}

	r := bytes.NewReader(buf)
	txi, err := DeserializeTx(r)

	if err != nil {
		t.Fatal(err)
	}

	if txi.Format() != TxFormatBasic {
		t.Fatalf("Deserialized tx format not basic: 0x%x", txi.Format())
	}

	tx := txi.(*BasicTx)

	senderPubKey := PublicKey{
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
		0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F,
	}

	if !bytes.Equal(tx.SenderPublicKey[:], senderPubKey[:]) {
		t.Fatalf("Failed deserializing SenderPublicKey.\n" +
			"Expected: hex(%x)\n" +
			"Actual: hex(%x)",
			&senderPubKey, &tx.SenderPublicKey,
		)
	}

	recipient := Address{
		0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39,
		0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49,
	}

	if !bytes.Equal(tx.Recipient[:], recipient[:]) {
		t.Fatalf("Failed deserializing Recipient.\n" +
			"Expected: hex(%x)\n" +
			"Actual: hex(%x)",
			&recipient, &tx.Recipient,
		)
	}

	if tx.Value != 0x5051525354555657 {
		t.Fatalf("Failed deserializing TxValue.\n" +
			"Expected: 0x%x\n" +
			"Actual: 0x%x",
			0x5051525354555657, tx.Value,
		)
	}

	if tx.Fee != 0x6061626364656667 {
		t.Fatalf("Failed deserializing TxFee.\n" +
			"Expected: 0x%x\n" +
			"Actual: 0x%x",
			0x6061626364656667, tx.Fee,
		)
	}

	if tx.ValidityStartHeight != 0 {
		t.Fatalf("Failed deserializing TxFee.\n" +
			"Expected: 0\n" +
			"Actual: %d",
			tx.ValidityStartHeight,
		)
	}

	if tx.NetworkId != 42 {
		t.Fatalf("Failed deserializing NetworkId.\n" +
			"Expected: 42\n" +
			"Actual: %d",
			tx.NetworkId,
		)
	}

	signature := Signature{
		0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7A, 0x7B, 0x7C, 0x7D, 0x7E, 0x7F,
		0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8A, 0x8B, 0x8C, 0x8D, 0x8E, 0x8F,
		0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9A, 0x9B, 0x9C, 0x9D, 0x9E, 0x9F,
		0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0xA6, 0xA7, 0xA8, 0xA9, 0xAA, 0xAB, 0xAC, 0xAD, 0xAE, 0xAF,
	}

	if !bytes.Equal(tx.Signature[:], signature[:]) {
		t.Fatal("Failed deserializing Signature.")
	}
}