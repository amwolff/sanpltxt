package sanpltxt

import "strings"

// ZUS represents a Type 2 transfer - ZUS/KRUS social insurance payment.
// Transfer mode is always Elixir (1) for ZUS transfers.
//
// Example line:
//
//	2|51109010430000000100111111|82600000020260111122223333|ZUS|Warszawa ul. Szamocka 3,5 01748|319,94|1|Skladka ZUS|01-09-2020|
type ZUS struct {
	DebitAccount  string // 26-digit NRB account number (source)
	CreditAccount string // 26-digit NRB account number (destination, ZUS account)
	RecipientName string // Recipient name (e.g., "ZUS"), max 80 chars
	Address       string // Recipient address, max 60 chars
	Amount        Amount // Transfer amount in grosze
	Title         string // Transfer description, max 140 chars
	Date          *Date  // Optional execution date
}

var _ Transfer = (*ZUS)(nil)

// Marshal converts the transfer to Santander format string.
func (z *ZUS) Marshal() (string, error) {
	if err := z.validate(); err != nil {
		return "", err
	}

	var b strings.Builder
	b.WriteString("2|")
	b.WriteString(z.DebitAccount)
	b.WriteString("|")
	b.WriteString(z.CreditAccount)
	b.WriteString("|")
	b.WriteString(z.RecipientName)
	b.WriteString("|")
	b.WriteString(z.Address)
	b.WriteString("|")
	b.WriteString(z.Amount.String())
	b.WriteString("|")
	b.WriteString(ModeElixir.String()) // ZUS transfers always use Elixir
	b.WriteString("|")
	b.WriteString(z.Title)
	b.WriteString("|")
	if z.Date != nil {
		b.WriteString(z.Date.String())
	}
	b.WriteString("|")

	return b.String(), nil
}

func (z *ZUS) validate() error {
	if err := validateNRB(z.DebitAccount, "debit account"); err != nil {
		return err
	}
	if err := validateNRB(z.CreditAccount, "credit account"); err != nil {
		return err
	}
	if err := validateRecipientName(z.RecipientName); err != nil {
		return err
	}
	if err := validateAddress(z.Address, true); err != nil {
		return err
	}
	if err := validateTitle(z.Title); err != nil {
		return err
	}
	return nil
}
