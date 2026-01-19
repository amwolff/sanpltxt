package sanpltxt

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const polishChars = "ąćęłńóśźżĄĆĘŁŃÓŚŹŻ"

var (
	charsRecipientName = buildCharSet("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`!@#$%^&*()_+-=[]{}; :.?/" + polishChars)
	charsAddress       = buildCharSet("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-.,:;/ " + polishChars)
	charsTitle         = buildCharSet("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`!@#$%^&*()_+-=[]{}; :,.?/" + polishChars)
	charsPayerName     = buildCharSet("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-;:.,/ " + polishChars)
	charsFormSymbol    = buildCharSet("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ-")
	charsObligationID  = buildCharSet("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-.,:; " + polishChars)
	charsInvoice       = buildCharSet("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-.,:;/ " + polishChars)
	charsFreeText      = charsInvoice
)

func buildCharSet(s string) map[rune]struct{} {
	m := make(map[rune]struct{})
	for _, r := range s {
		m[r] = struct{}{}
	}
	return m
}

func containsOnly(s string, allowed map[rune]struct{}) bool {
	for _, r := range s {
		if _, ok := allowed[r]; !ok {
			return false
		}
	}
	return true
}

func isDigitsOnly(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func validateNRB(account, fieldName string) error {
	if len(account) != 26 {
		return fmt.Errorf("%s must be exactly 26 digits, got %d", fieldName, len(account))
	}
	if !isDigitsOnly(account) {
		return fmt.Errorf("%s must contain only digits", fieldName)
	}
	return nil
}

func validateNIP(nip string) error {
	if len(nip) != 10 {
		return fmt.Errorf("NIP must be exactly 10 digits, got %d", len(nip))
	}
	if !isDigitsOnly(nip) {
		return errors.New("NIP must contain only digits")
	}
	return nil
}

func validateRecipientName(name string) error {
	if name == "" {
		return errors.New("recipient name is required")
	}
	if len(name) > 80 {
		return fmt.Errorf("recipient name must be at most 80 characters, got %d", len(name))
	}
	if !containsOnly(name, charsRecipientName) {
		return errors.New("recipient name contains invalid characters")
	}
	return nil
}

func validateAddress(address string, required bool) error {
	if address == "" {
		if required {
			return errors.New("address is required")
		}
		return nil
	}
	if len(address) > 60 {
		return fmt.Errorf("address must be at most 60 characters, got %d", len(address))
	}
	if !containsOnly(address, charsAddress) {
		return errors.New("address contains invalid characters")
	}
	return nil
}

func validateTitle(title string) error {
	if title == "" {
		return errors.New("title is required")
	}
	if len(title) > 140 {
		return fmt.Errorf("title must be at most 140 characters, got %d", len(title))
	}
	if !containsOnly(title, charsTitle) {
		return errors.New("title contains invalid characters")
	}
	return nil
}

func validatePayerName(name string) error {
	if name == "" {
		return errors.New("payer name is required")
	}
	if len(name) > 50 {
		return fmt.Errorf("payer name must be at most 50 characters, got %d", len(name))
	}
	if !containsOnly(name, charsPayerName) {
		return errors.New("payer name contains invalid characters")
	}
	return nil
}

func validateFormSymbol(symbol string) error {
	if symbol == "" {
		return errors.New("form symbol is required")
	}
	if len(symbol) > 6 {
		return fmt.Errorf("form symbol must be at most 6 characters, got %d", len(symbol))
	}
	if !containsOnly(symbol, charsFormSymbol) {
		return errors.New("form symbol contains invalid characters (allowed: 0-9, A-Z, -)")
	}
	return nil
}

func validateObligationID(id string) error {
	if id == "" {
		return nil // optional field
	}
	if len(id) > 20 {
		return fmt.Errorf("obligation ID must be at most 20 characters, got %d", len(id))
	}
	if !containsOnly(id, charsObligationID) {
		return errors.New("obligation ID contains invalid characters")
	}
	return nil
}

func validateInvoiceNumber(num string) error {
	if num == "" {
		return errors.New("invoice number is required")
	}
	if len(num) > 35 {
		return fmt.Errorf("invoice number must be at most 35 characters, got %d", len(num))
	}
	if !containsOnly(num, charsInvoice) {
		return errors.New("invoice number contains invalid characters")
	}
	return nil
}

func validateFreeText(text string) error {
	if text == "" {
		return nil // optional field
	}
	if len(text) > 33 {
		return fmt.Errorf("free text must be at most 33 characters, got %d", len(text))
	}
	if !containsOnly(text, charsFreeText) {
		return errors.New("free text contains invalid characters")
	}
	return nil
}

func validateTransferMode(mode TransferMode, allowedModes ...TransferMode) error {
	for _, m := range allowedModes {
		if mode == m {
			return nil
		}
	}
	var modes []string
	for _, m := range allowedModes {
		modes = append(modes, m.String())
	}
	return errors.New("transfer mode must be one of: " + strings.Join(modes, ", "))
}

// IdentifierType is the type of tax identifier.
type IdentifierType string

// Identifier types.
const (
	IdentifierNIP      IdentifierType = "N"
	IdentifierREGON    IdentifierType = "R"
	IdentifierPESEL    IdentifierType = "P"
	IdentifierID       IdentifierType = "1"
	IdentifierPassport IdentifierType = "2"
	IdentifierOther    IdentifierType = "3"
)

func validateIdentifierType(t IdentifierType) error {
	switch t {
	case IdentifierNIP, IdentifierREGON, IdentifierPESEL, IdentifierID, IdentifierPassport, IdentifierOther:
		return nil
	}
	return errors.New("identifier type must be one of: N (NIP), R (REGON), P (PESEL), 1 (ID), 2 (Passport), 3 (Other)")
}

func validateIdentifier(id string, idType IdentifierType) error {
	if id == "" {
		return errors.New("identifier is required")
	}

	switch idType {
	case IdentifierNIP:
		if len(id) != 10 || !isDigitsOnly(id) {
			return errors.New("NIP must be exactly 10 digits")
		}
	case IdentifierREGON:
		if (len(id) != 9 && len(id) != 14) || !isDigitsOnly(id) {
			return errors.New("REGON must be 9 or 14 digits")
		}
	case IdentifierPESEL:
		if len(id) != 11 || !isDigitsOnly(id) {
			return errors.New("PESEL must be exactly 11 digits")
		}
	case IdentifierID:
		if len(id) < 8 || len(id) > 9 {
			return errors.New("ID card number must be 8 or 9 characters")
		}
	case IdentifierPassport, IdentifierOther:
		if len(id) > 14 {
			return errors.New("identifier must be at most 14 characters")
		}
	}
	return nil
}

// PeriodType is the type of tax period.
type PeriodType string

// Period types.
const (
	PeriodYear    PeriodType = "R"
	PeriodHalf    PeriodType = "P"
	PeriodQuarter PeriodType = "K"
	PeriodMonth   PeriodType = "M"
	PeriodDecade  PeriodType = "D"
	PeriodDay     PeriodType = "J"
)

func validatePeriodType(t PeriodType) error {
	if t == "" {
		return nil // optional
	}
	switch t {
	case PeriodYear, PeriodHalf, PeriodQuarter, PeriodMonth, PeriodDecade, PeriodDay:
		return nil
	}
	return errors.New("period type must be one of: R (year), P (half), K (quarter), M (month), D (decade), J (day)")
}

func validatePeriodNumber(num string, periodType PeriodType) error {
	if periodType == "" || periodType == PeriodYear {
		if num != "" {
			return errors.New("period number should be empty for yearly periods")
		}
		return nil
	}

	if num == "" {
		return errors.New("period number is required for non-yearly periods")
	}

	if len(num) > 4 || !isDigitsOnly(num) {
		return errors.New("period number must be 1-4 digits")
	}

	n, _ := strconv.Atoi(num)

	switch periodType {
	case PeriodHalf:
		if n < 1 || n > 2 {
			return errors.New("half-year number must be 01 or 02")
		}
	case PeriodQuarter:
		if n < 1 || n > 4 {
			return errors.New("quarter number must be 01, 02, 03, or 04")
		}
	case PeriodMonth:
		if n < 1 || n > 12 {
			return errors.New("month number must be 01-12")
		}
	case PeriodDecade:
		if n < 1 || n > 3 {
			return errors.New("decade number must be 01, 02, or 03")
		}
	case PeriodDay:
		// Format DDMM, validated separately
		if len(num) != 4 {
			return errors.New("day format must be DDMM (4 digits)")
		}
	}

	return nil
}

func validateYear(year string) error {
	if year == "" {
		return nil // optional
	}
	if len(year) != 2 || !isDigitsOnly(year) {
		return errors.New("year must be 2 digits (e.g., 05 for 2005)")
	}
	return nil
}
