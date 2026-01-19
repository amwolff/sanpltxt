package sanpltxt

import (
	"strings"
	"time"
)

// Payroll is a Type 5 transfer (salary payment). Only allowed in payroll packages.
type Payroll struct {
	DebitAccount  string
	CreditAccount string
	RecipientName string
	Address       string
	Amount        Amount
	Mode          TransferMode
	Title         string
	Date          *time.Time
}

var _ Transfer = (*Payroll)(nil)

// Marshal returns the transfer in Santander format.
func (p *Payroll) Marshal() (string, error) {
	var b strings.Builder
	if err := p.marshal(&b); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (p *Payroll) marshal(b *strings.Builder) error {
	if err := p.validate(); err != nil {
		return err
	}

	b.WriteString("5|")
	b.WriteString(p.DebitAccount)
	b.WriteString("|")
	b.WriteString(p.CreditAccount)
	b.WriteString("|")
	b.WriteString(p.RecipientName)
	b.WriteString("|")
	b.WriteString(p.Address)
	b.WriteString("|")
	b.WriteString(p.Amount.String())
	b.WriteString("|")
	b.WriteString(p.Mode.String())
	b.WriteString("|")
	b.WriteString(p.Title)
	b.WriteString("|")
	if p.Date != nil {
		b.WriteString(p.Date.Format(dateFormat))
	}
	b.WriteString("|")

	return nil
}

func (p *Payroll) validate() error {
	if err := validateNRB(p.DebitAccount, "debit account"); err != nil {
		return err
	}
	if err := validateNRB(p.CreditAccount, "credit account"); err != nil {
		return err
	}
	if err := validateRecipientName(p.RecipientName); err != nil {
		return err
	}
	if err := validateAddress(p.Address, true); err != nil {
		return err
	}
	if err := validateTransferMode(p.Mode, ModeInternal, ModeElixir, ModeSORBNET, ModeExpressElixir); err != nil {
		return err
	}
	if err := validateTitle(p.Title); err != nil {
		return err
	}
	return nil
}
