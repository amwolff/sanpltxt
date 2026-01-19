package sanpltxt

import (
	"strings"
	"time"
)

// ZUS is a Type 2 transfer (social insurance payment). Mode is always Elixir.
type ZUS struct {
	DebitAccount  string
	CreditAccount string
	RecipientName string
	Address       string
	Amount        Amount
	Title         string
	Date          *time.Time
}

var _ Transfer = (*ZUS)(nil)

// Marshal returns the transfer in Santander format.
func (z *ZUS) Marshal() (string, error) {
	var b strings.Builder
	if err := z.marshal(&b); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (z *ZUS) marshal(b *strings.Builder) error {
	if err := z.validate(); err != nil {
		return err
	}

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
		b.WriteString(z.Date.Format(dateFormat))
	}
	b.WriteString("|")

	return nil
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
