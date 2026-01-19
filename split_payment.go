package sanpltxt

import (
	"strings"
	"time"
)

// SplitPayment is a Type 6 transfer (split VAT payment).
type SplitPayment struct {
	DebitAccount  string
	CreditAccount string
	RecipientName string
	Address       string
	GrossAmount   Amount
	Mode          TransferMode
	VATAmount     Amount
	RecipientNIP  string
	InvoiceNumber string
	FreeText      string
	Date          *time.Time
}

var _ Transfer = (*SplitPayment)(nil)

// Marshal returns the transfer in Santander format.
func (s *SplitPayment) Marshal() (string, error) {
	var b strings.Builder
	if err := s.marshal(&b); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (s *SplitPayment) marshal(b *strings.Builder) error {
	if err := s.validate(); err != nil {
		return err
	}

	b.WriteString("6|")
	b.WriteString(s.DebitAccount)
	b.WriteString("|")
	b.WriteString(s.CreditAccount)
	b.WriteString("|")
	b.WriteString(s.RecipientName)
	b.WriteString("|")
	b.WriteString(s.Address)
	b.WriteString("|")
	b.WriteString(s.GrossAmount.String())
	b.WriteString("|")
	b.WriteString(s.Mode.String())
	b.WriteString("|")
	s.formatTitle(b)
	b.WriteString("|")
	if s.Date != nil {
		b.WriteString(s.Date.Format(dateFormat))
	}
	b.WriteString("|")

	return nil
}

func (s *SplitPayment) formatTitle(b *strings.Builder) {
	b.WriteString("/VAT/")
	b.WriteString(s.VATAmount.String())
	b.WriteString("/IDC/")
	b.WriteString(s.RecipientNIP)
	b.WriteString("/INV/")
	b.WriteString(s.InvoiceNumber)
	if s.FreeText != "" {
		b.WriteString("/TXT/")
		b.WriteString(s.FreeText)
	}
}

func (s *SplitPayment) validate() error {
	if err := validateNRB(s.DebitAccount, "debit account"); err != nil {
		return err
	}
	if err := validateNRB(s.CreditAccount, "credit account"); err != nil {
		return err
	}
	if err := validateRecipientName(s.RecipientName); err != nil {
		return err
	}
	if err := validateAddress(s.Address, false); err != nil {
		return err
	}
	if err := validateTransferMode(s.Mode, ModeInternal, ModeElixir, ModeSORBNET, ModeExpressElixir); err != nil {
		return err
	}
	if err := validateNIP(s.RecipientNIP); err != nil {
		return err
	}
	if err := validateInvoiceNumber(s.InvoiceNumber); err != nil {
		return err
	}
	if err := validateFreeText(s.FreeText); err != nil {
		return err
	}
	return nil
}
