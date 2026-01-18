package sanpltxt

import "strings"

// Payroll represents a Type 5 transfer - salary payment.
// Can only be used in Package.Type=2 (payroll packages).
//
// Example line:
//
//	5|51109010430000000100111111|50102055581111103350100016|Jan Nowak|Poznań ul. Swojska 17 06-123|1000,12|1|Wynagrodzenie za miesiąc|01-09-2020|
type Payroll struct {
	DebitAccount  string       // 26-digit NRB account number (source)
	CreditAccount string       // 26-digit NRB account number (destination)
	RecipientName string       // Employee name, max 80 chars
	Address       string       // Employee address, max 60 chars
	Amount        Amount       // Transfer amount in grosze
	Mode          TransferMode // Transfer mode: Internal(0), Elixir(1), SORBNET(6), ExpressElixir(8)
	Title         string       // Transfer description (e.g., "Wynagrodzenie za miesiąc"), max 140 chars
	Date          *Date        // Optional execution date
}

var _ Transfer = (*Payroll)(nil)

// Marshal converts the transfer to Santander format string.
func (p *Payroll) Marshal() (string, error) {
	if err := p.validate(); err != nil {
		return "", err
	}

	var b strings.Builder
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
		b.WriteString(p.Date.String())
	}
	b.WriteString("|")

	return b.String(), nil
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
