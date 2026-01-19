package sanpltxt

import (
	"strings"
	"time"
)

// Standard is a Type 1 transfer (external account transfer).
type Standard struct {
	DebitAccount  string
	CreditAccount string
	RecipientName string
	Address       string
	Amount        Amount
	Mode          TransferMode
	Title         string
	Date          *time.Time
	NIP           string
}

var _ Transfer = (*Standard)(nil)

// Marshal returns the transfer in Santander format.
func (s *Standard) Marshal() (string, error) {
	var b strings.Builder
	if err := s.marshal(&b); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (s *Standard) marshal(b *strings.Builder) error {
	if err := s.validate(); err != nil {
		return err
	}

	b.WriteString("1|")
	b.WriteString(s.DebitAccount)
	b.WriteString("|")
	b.WriteString(s.CreditAccount)
	b.WriteString("|")
	b.WriteString(s.RecipientName)
	b.WriteString("|")
	b.WriteString(s.Address)
	b.WriteString("|")
	b.WriteString(s.Amount.String())
	b.WriteString("|")
	b.WriteString(s.Mode.String())
	b.WriteString("|")
	b.WriteString(s.Title)
	b.WriteString("|")
	if s.Date != nil {
		b.WriteString(s.Date.Format(dateFormat))
	}
	b.WriteString("|")
	b.WriteString(s.NIP)
	b.WriteString("|")

	return nil
}

func (s *Standard) validate() error {
	if err := validateNRB(s.DebitAccount, "debit account"); err != nil {
		return err
	}
	if err := validateNRB(s.CreditAccount, "credit account"); err != nil {
		return err
	}
	if err := validateRecipientName(s.RecipientName); err != nil {
		return err
	}
	if err := validateAddress(s.Address, true); err != nil {
		return err
	}
	if err := validateTransferMode(s.Mode, ModeInternal, ModeElixir, ModeSORBNET, ModeExpressElixir); err != nil {
		return err
	}
	if err := validateTitle(s.Title); err != nil {
		return err
	}
	if s.NIP != "" {
		if err := validateNIP(s.NIP); err != nil {
			return err
		}
	}
	return nil
}
