package core

import (
	"io"
	"encoding/binary"
	"errors"
)

type TxFormat uint8
const (
	TxFormatBasic    = TxFormat(0)
	TxFormatExtended = TxFormat(1)
)

type Tx interface {
	//Deserialize(reader io.Reader)
	Format() TxFormat
}

func DeserializeTx(r io.Reader) (Tx, error) {
	// Read transaction format (1 byte)
	var format TxFormat
	err := binary.Read(r, binary.BigEndian, &format)
	if err != nil { return nil, err }

	switch format {
	case TxFormatBasic:
		// Static length: Create and fill buffer
		var txBuf [BasicTxSize]byte
		_, err = io.ReadFull(r, txBuf[:])
		if err != nil { return nil, err }

		// Deserialize basic transaction
		basicTx := new(BasicTx)
		basicTx.Deserialize(&txBuf)

		return basicTx, nil

	case TxFormatExtended:
		// Deserialize extended transaction
		extendedTx := new(ExtendedTx)
		extendedTx.Deserialize(r)

		return extendedTx, nil

	default:
		return nil, errors.New("unsupported tx format")
	}
}
