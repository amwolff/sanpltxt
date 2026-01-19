package sanpltxt

import (
	"strings"
	"time"
)

// Tax is a Type 3 (TaxOffice=true) or Type 4 (TaxOffice=false) transfer.
type Tax struct {
	TaxOffice      bool // true = type 3, false = type 4
	DebitAccount   string
	CreditAccount  string
	RecipientName  string
	Address        string
	Amount         Amount
	Date           *time.Time
	PayerName      string
	IdentifierType IdentifierType
	Identifier     string
	Year           string
	PeriodType     PeriodType
	PeriodNumber   string
	FormSymbol     string
	ObligationID   string
}

var _ Transfer = (*Tax)(nil)

// Marshal returns the transfer in Santander format.
func (t *Tax) Marshal() (string, error) {
	var b strings.Builder
	if err := t.marshal(&b); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (t *Tax) marshal(b *strings.Builder) error {
	if err := t.validate(); err != nil {
		return err
	}

	if t.TaxOffice {
		b.WriteString("3|")
	} else {
		b.WriteString("4|")
	}
	b.WriteString(t.DebitAccount)
	b.WriteString("|")
	b.WriteString(t.CreditAccount)
	b.WriteString("|")
	b.WriteString(t.RecipientName)
	b.WriteString("|")
	b.WriteString(t.Address)
	b.WriteString("|")
	b.WriteString(t.Amount.String())
	b.WriteString("|")
	if t.Date != nil {
		b.WriteString(t.Date.Format(dateFormat))
	}
	b.WriteString("|")
	b.WriteString(t.PayerName)
	b.WriteString("|")
	b.WriteString(string(t.IdentifierType))
	b.WriteString("|")
	b.WriteString(t.Identifier)
	b.WriteString("|")
	b.WriteString(t.Year)
	b.WriteString("|")
	b.WriteString(string(t.PeriodType))
	b.WriteString("|")
	b.WriteString(t.PeriodNumber)
	b.WriteString("|")
	b.WriteString(t.FormSymbol)
	b.WriteString("|")
	b.WriteString(t.ObligationID)
	b.WriteString("|")

	return nil
}

func (t *Tax) validate() error {
	if err := validateNRB(t.DebitAccount, "debit account"); err != nil {
		return err
	}
	if err := validateNRB(t.CreditAccount, "credit account"); err != nil {
		return err
	}
	if err := validateRecipientName(t.RecipientName); err != nil {
		return err
	}
	if err := validateAddress(t.Address, false); err != nil {
		return err
	}
	if err := validatePayerName(t.PayerName); err != nil {
		return err
	}
	if err := validateIdentifierType(t.IdentifierType); err != nil {
		return err
	}
	if err := validateIdentifier(t.Identifier, t.IdentifierType); err != nil {
		return err
	}
	if err := validateYear(t.Year); err != nil {
		return err
	}
	if err := validatePeriodType(t.PeriodType); err != nil {
		return err
	}
	if err := validatePeriodNumber(t.PeriodNumber, t.PeriodType); err != nil {
		return err
	}
	if err := validateFormSymbol(t.FormSymbol); err != nil {
		return err
	}
	if err := validateObligationID(t.ObligationID); err != nil {
		return err
	}
	return nil
}
