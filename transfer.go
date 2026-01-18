// Package sanpltxt implements marshaling for Santander Bank's transfer import format.
// The format is documented in the bank's "Formaty danych w pliku do importu przelewów" PDF.
//
// The package supports all transfer types: Standard (1), ZUS (2), Tax (3/4), Payroll (5),
// and SplitPayment (6). Output is UTF-8 encoded; use ToWindows1250 for file export.
package sanpltxt

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// FormatVersion is the current version of the Santander import format.
const FormatVersion = "4120414"

// Transfer is the interface implemented by all transfer types.
type Transfer interface {
	Marshal() (string, error)
}

// Package represents a collection of transfers to be exported.
// Type must be 1 (regular) or 2 (payroll).
// Regular packages (Type=1) can contain Standard, ZUS, Tax, and SplitPayment transfers.
// Payroll packages (Type=2) can only contain Payroll transfers.
type Package struct {
	Type      int
	Transfers []Transfer
}

// Marshal produces the complete file content for the package.
// The output is UTF-8 encoded; use ToWindows1250 to convert for file export.
func (p *Package) Marshal() (string, error) {
	if p.Type != 1 && p.Type != 2 {
		return "", errors.New("package type must be 1 (regular) or 2 (payroll)")
	}

	var b strings.Builder
	b.WriteString(FormatVersion)
	b.WriteString("|")
	b.WriteString(strconv.Itoa(p.Type))
	b.WriteString("\n")

	for i, t := range p.Transfers {
		// Validate transfer type matches package type
		_, isPayroll := t.(*Payroll)
		if p.Type == 1 && isPayroll {
			return "", errors.New("payroll transfers (type 5) cannot be in regular packages (type 1)")
		}
		if p.Type == 2 && !isPayroll {
			return "", errors.New("only payroll transfers (type 5) are allowed in payroll packages (type 2)")
		}

		line, err := t.Marshal()
		if err != nil {
			return "", fmt.Errorf("transfer %d: %w", i, err)
		}
		b.WriteString(line)
		b.WriteString("\n")
	}

	return b.String(), nil
}

// TransferMode represents the transfer processing mode.
type TransferMode int

const (
	// ModeInternal is for internal bank transfers (0).
	ModeInternal TransferMode = 0
	// ModeElixir is for standard interbank transfers via Elixir (1).
	ModeElixir TransferMode = 1
	// ModeSORBNET is for high-value RTGS transfers (6).
	ModeSORBNET TransferMode = 6
	// ModeExpressElixir is for instant transfers via Express Elixir (8).
	ModeExpressElixir TransferMode = 8
)

// String returns the numeric string representation of the mode.
func (m TransferMode) String() string {
	return strconv.Itoa(int(m))
}

// Amount represents a monetary value in grosze (1/100 PLN).
// Use grosze to avoid floating-point precision issues.
type Amount int64

// String formats the amount as "złoty,grosze" (e.g., "123,45" or "100").
func (a Amount) String() string {
	if a < 0 {
		a = -a
	}
	zloty := a / 100
	grosze := a % 100
	if grosze == 0 {
		return strconv.FormatInt(int64(zloty), 10)
	}
	var b strings.Builder
	b.WriteString(strconv.FormatInt(int64(zloty), 10))
	b.WriteString(",")
	if grosze < 10 {
		b.WriteString("0")
	}
	b.WriteString(strconv.FormatInt(int64(grosze), 10))
	return b.String()
}

// Date represents a date for transfer execution.
type Date struct {
	time.Time
}

// String formats the date as DD-MM-RRRR (e.g., "01-09-2020").
func (d Date) String() string {
	day := d.Day()
	month := int(d.Month())
	year := d.Year()

	var b strings.Builder
	if day < 10 {
		b.WriteString("0")
	}
	b.WriteString(strconv.Itoa(day))
	b.WriteString("-")
	if month < 10 {
		b.WriteString("0")
	}
	b.WriteString(strconv.Itoa(month))
	b.WriteString("-")
	b.WriteString(strconv.Itoa(year))
	return b.String()
}
