package core

import (
	"io"
	"encoding/binary"
)

type ExtendedTx struct {
	Data []byte
	Sender Address
	SenderType AccountType
	Recipient Address
	RecipientType AccountType
	Value Satoshi
	Fee Satoshi
	ValidityStartHeight uint32
	NetworkId uint8
	Flags uint8
	Proof []byte
}

func (_ *ExtendedTx) Format() TxFormat {
	return TxFormatExtended
}

func (e *ExtendedTx) Deserialize(r io.Reader) error {
	bo := binary.BigEndian
	var err error

	// Variable length part

	// Read length of data segment
	var dataLen uint16
	err = binary.Read(r, bo, &dataLen)
	if err != nil { return err }

	// Read data segment
	e.Data = make([]byte, dataLen)
	_, err = io.ReadFull(r, e.Data)
	if err != nil { return err }

	// Static length part

	// Read all static fields
	staticFields := []interface{}{
		e.Sender[:],
		&e.SenderType,
		e.Recipient[:],
		&e.RecipientType,
		&e.Value,
		&e.Fee,
		&e.ValidityStartHeight,
		&e.NetworkId,
		&e.Flags,
	}

	for _, field := range staticFields {
		err = binary.Read(r, bo, field)
		if err != nil { return err }
	}

	// Variable length part

	// Read length of proof data segment
	var proofLen uint16
	err = binary.Read(r, bo, &proofLen)
	if err != nil { return err }

	// Read data segment
	e.Proof = make([]byte, proofLen)
	_, err = io.ReadFull(r, e.Proof)
	if err != nil { return err }

	// No error
	return nil
}
