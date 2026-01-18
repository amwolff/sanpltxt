package sanpltxt

import "strings"

// Standard represents a Type 1 transfer - external account transfer.
// Used for regular transfers to accounts at other banks.
//
// Example line:
//
//	1|51109010430000000100111111|50102055581111103350100016|Jerzy Kowalski|Warszawa ul. Kaliska 123 00-123|123,12|1|zasielenie konta|01-09-2020|7850000000|
type Standard struct {
	DebitAccount  string       // 26-digit NRB account number (source)
	CreditAccount string       // 26-digit NRB account number (destination)
	RecipientName string       // Recipient name, max 80 chars
	Address       string       // Recipient address, max 60 chars
	Amount        Amount       // Transfer amount in grosze
	Mode          TransferMode // Transfer mode: Internal(0), Elixir(1), SORBNET(6), ExpressElixir(8)
	Title         string       // Transfer description, max 140 chars
	Date          *Date        // Optional execution date
	NIP           string       // Optional recipient NIP for VAT whitelist verification
}

var _ Transfer = (*Standard)(nil)

// Marshal converts the transfer to Santander format string.
func (s *Standard) Marshal() (string, error) {
	if err := s.validate(); err != nil {
		return "", err
	}

	var b strings.Builder
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
		b.WriteString(s.Date.String())
	}
	b.WriteString("|")
	b.WriteString(s.NIP)
	b.WriteString("|")

	return b.String(), nil
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
