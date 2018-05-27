package core

type Account struct {
	Type AccountType
	Balance Satoshi
}

const AccountSize =
	1 + // AccountType uint8
	8   // Balance uint64

// Account type enum
type AccountType uint8
const (
	BasicAccount = AccountType(0)
	VestingAccount = AccountType(1)
	HtlcAccount = AccountType(2)
)

func (a AccountType) String() string {
	switch a {
	case BasicAccount: return "Basic account"
	case VestingAccount: return "Vesting contract"
	case HtlcAccount: return "Hash time-locked contract"
	default: return "Invalid account type"
	}
}
