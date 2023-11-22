package models

type WalletType int

const (
	Personal WalletType = iota
	Family
)

const (
	personalStr = "personal"
	familyStr   = "family"
)

var WalletTypesStr = map[string]WalletType{
	personalStr: Personal,
	familyStr:   Family,
}

func (wt WalletType) String() string {
	switch wt {
	case 0:
		return personalStr
	case 1:
		return familyStr
	default:
		return "undefined"
	}
}

func (wt WalletType) FromString(walletType string) WalletType {
	return WalletTypesStr[walletType]
}

type Wallet struct {
	ID      string  `db:"id"`
	UserID  string  `db:"user_id"`
	Card    string  `db:"card_number"`
	Balance float64 `db:"balance"`
}

type FamilyWallet struct {
	ID           string  `db:"id"`
	WalletID     string  `db:"wallet_id"`
	Balance      float64 `db:"balance"`
	FixedBalance float64 `db:"fixed_balance"`
}
