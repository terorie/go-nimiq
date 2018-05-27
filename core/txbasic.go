package core

import (
	bu "github.com/terorie/go-nimiq/bufferutils"
)

type BasicTx struct {
	SenderPublicKey PublicKey
	Recipient Address
	Value Satoshi
	Fee Satoshi
	ValidityStartHeight uint32
	NetworkId uint8
	Signature Signature
}

const BasicTxSize =
	32 + // PublicKey
	20 + // Address
	8 +  // Value
	8 +  // Fee
	4 +  // ValidityStartHeight
	1 +  // NetworkId
	64   // Signature

func (_ *BasicTx) Format() TxFormat {
	return TxFormatBasic
}

func (t *BasicTx) Deserialize(buf *[BasicTxSize]byte) {
	sl := bu.ReadBytes(buf[:])

	sl.CopyNext(t.SenderPublicKey[:])
	sl.CopyNext(t.Recipient[:])
	t.Value = Satoshi(sl.Uint64())
	t.Fee = Satoshi(sl.Uint64())
	t.ValidityStartHeight = sl.Uint32()
	t.NetworkId = sl.Uint8()
	sl.CopyNext(t.Signature[:])
}

func (t *BasicTx) Serialize(buf [BasicTxSize]byte) {
	sl := bu.WriteBytes(buf[:])

	sl.Uint8(uint8(TxFormatBasic))
	sl.WriteNext(t.SenderPublicKey[:])
	sl.WriteNext(t.Recipient[:])
	sl.Uint64(uint64(t.Value))
	sl.Uint64(uint64(t.Fee))
	sl.Uint32(t.ValidityStartHeight)
	sl.Uint8(t.NetworkId)
}
