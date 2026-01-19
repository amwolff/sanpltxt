// Package sanpltxt marshals Santander Bank's transfer import format.
package sanpltxt

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// FormatVersion is the Santander format version.
const (
	FormatVersion = "4120414"
	dateFormat    = "02-01-2006"
)

// Transfer is implemented by all transfer types.
type Transfer interface {
	Marshal() (string, error)
}

// PackageOptions configures encoding behavior.
type PackageOptions struct {
	EncodeUTF8 bool // false (default) = Windows-1250, true = UTF-8
}

// Package is a collection of transfers for export.
type Package struct {
	typ       int
	transfers []Transfer
	options   PackageOptions
}

// NewPackage creates a Package. Type: 1 = regular, 2 = payroll.
func NewPackage(typ int, transfers []Transfer, opts *PackageOptions) *Package {
	p := &Package{
		typ:       typ,
		transfers: transfers,
	}
	if opts != nil {
		p.options = *opts
	}
	return p
}

// Marshal returns the package content as a UTF-8 string.
func (p *Package) Marshal() (string, error) {
	var b strings.Builder
	if err := p.marshal(&b); err != nil {
		return "", err
	}
	return b.String(), nil
}

// MarshalBytes returns the package content with encoding based on options.
func (p *Package) MarshalBytes() ([]byte, error) {
	s, err := p.Marshal()
	if err != nil {
		return nil, err
	}
	if p.options.EncodeUTF8 {
		return []byte(s), nil
	}
	return ToWindows1250(s)
}

func (p *Package) marshal(b *strings.Builder) error {
	if p.typ != 1 && p.typ != 2 {
		return errors.New("package type must be 1 (regular) or 2 (payroll)")
	}

	b.WriteString(FormatVersion)
	b.WriteString("|")
	b.WriteString(strconv.Itoa(p.typ))
	b.WriteString("\n")

	for i, t := range p.transfers {
		// Validate transfer type matches package type
		_, isPayroll := t.(*Payroll)
		if p.typ == 1 && isPayroll {
			return errors.New("payroll transfers (type 5) cannot be in regular packages (type 1)")
		}
		if p.typ == 2 && !isPayroll {
			return errors.New("only payroll transfers (type 5) are allowed in payroll packages (type 2)")
		}

		if m, ok := t.(interface{ marshal(*strings.Builder) error }); ok {
			if err := m.marshal(b); err != nil {
				return fmt.Errorf("transfer %d: %w", i, err)
			}
		} else {
			line, err := t.Marshal()
			if err != nil {
				return fmt.Errorf("transfer %d: %w", i, err)
			}
			b.WriteString(line)
		}
		b.WriteString("\n")
	}

	return nil
}

// TransferMode is the transfer processing mode.
type TransferMode int

// Transfer modes.
const (
	ModeInternal      TransferMode = 0
	ModeElixir        TransferMode = 1
	ModeSORBNET       TransferMode = 6
	ModeExpressElixir TransferMode = 8
)

func (m TransferMode) String() string { return strconv.Itoa(int(m)) }

// Amount is a monetary value in grosze (1/100 PLN).
type Amount int64

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
