package sanpltxt

import "strings"

// SplitPayment represents a Type 6 transfer - split VAT payment.
// The title field is auto-generated from VATAmount, RecipientNIP, InvoiceNumber, and FreeText.
//
// Example line:
//
//	6|51109010430000000100111111|50102055581111103350100011|Jan Nowak|Warszawa ul. Mickiewicza 11 02-222|123,5|1|/VAT/23,09/IDC/8960005670/INV/5/2018/TXT/Faktura 5/2018|30-09-2020|
type SplitPayment struct {
	DebitAccount  string       // 26-digit NRB account number (source)
	CreditAccount string       // 26-digit NRB account number (destination)
	RecipientName string       // Recipient name, max 80 chars
	Address       string       // Optional recipient address, max 60 chars
	GrossAmount   Amount       // Gross transfer amount in grosze
	Mode          TransferMode // Transfer mode: Internal(0), Elixir(1), SORBNET(6), ExpressElixir(8)
	VATAmount     Amount       // VAT amount in grosze
	RecipientNIP  string       // Recipient NIP (10 digits) for VAT whitelist verification
	InvoiceNumber string       // Invoice number, max 35 chars
	FreeText      string       // Optional free text description, max 33 chars
	Date          *Date        // Optional execution date
}

var _ Transfer = (*SplitPayment)(nil)

// Marshal converts the transfer to Santander format string.
func (s *SplitPayment) Marshal() (string, error) {
	if err := s.validate(); err != nil {
		return "", err
	}

	var b strings.Builder
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
	b.WriteString(s.formatTitle())
	b.WriteString("|")
	if s.Date != nil {
		b.WriteString(s.Date.String())
	}
	b.WriteString("|")

	return b.String(), nil
}

// formatTitle generates the split payment title in format:
// /VAT/{kwota}/IDC/{nip}/INV/{nr faktury}/TXT/{tytuł}
// or /VAT/{kwota}/IDC/{nip}/INV/{nr faktury} if FreeText is empty.
func (s *SplitPayment) formatTitle() string {
	var b strings.Builder
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
	return b.String()
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
