package sanpltxt

import "strings"

// Tax represents a Type 3 (tax office) or Type 4 (other tax authority) transfer.
// The TaxOffice field determines the type: true=3 (urząd skarbowy), false=4 (other organ).
//
// Example line (type 3):
//
//	3|51109010430000000100111111|06101014690039392223000000|Urzad Skarbowy Poznan Winogrady|Poznan Wojciechowskiego 3/5 60-685|1000|01-09-2020|Jan Kowalski|N|9721230101|05|M|07|PIT5|id.zobowiazania|
type Tax struct {
	TaxOffice      bool           // true = type 3 (tax office), false = type 4 (other tax authority)
	DebitAccount   string         // 26-digit NRB account number (source)
	CreditAccount  string         // 26-digit NRB account number (destination)
	RecipientName  string         // Tax authority name, max 80 chars
	Address        string         // Optional address, max 60 chars
	Amount         Amount         // Transfer amount in grosze
	Date           *Date          // Optional execution date
	PayerName      string         // Payer name, max 50 chars
	IdentifierType IdentifierType // Type of identifier (N=NIP, R=REGON, P=PESEL, 1=ID, 2=Passport, 3=Other)
	Identifier     string         // Identifier value (NIP=10, REGON=9/14, PESEL=11, ID=8-9, Passport≤14)
	Year           string         // Optional 2-digit year (e.g., "05" for 2005)
	PeriodType     PeriodType     // Optional period type (R=Year, P=Half, K=Quarter, M=Month, D=Decade, J=Day)
	PeriodNumber   string         // Optional period number (depends on PeriodType)
	FormSymbol     string         // Tax form symbol (e.g., "PIT5"), max 6 chars
	ObligationID   string         // Optional obligation identifier, max 20 chars
}

var _ Transfer = (*Tax)(nil)

// Marshal converts the transfer to Santander format string.
func (t *Tax) Marshal() (string, error) {
	if err := t.validate(); err != nil {
		return "", err
	}

	var b strings.Builder
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
		b.WriteString(t.Date.String())
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

	return b.String(), nil
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
