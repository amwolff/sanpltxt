package sanpltxt_test

import (
	"strings"
	"testing"
	"time"

	"github.com/zeebo/assert"

	"github.com/amwolff/sanpltxt"
)

func date(year, month, day int) *sanpltxt.Date {
	return &sanpltxt.Date{Time: time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)}
}

func TestAmount_String(t *testing.T) {
	tests := []struct {
		amount sanpltxt.Amount
		want   string
	}{
		{12312, "123,12"},
		{31994, "319,94"},
		{100000, "1000"},
		{12350, "123,50"},
		{100012, "1000,12"},
		{2309, "23,09"},
		{100, "1"},
		{1, "0,01"},
		{10, "0,10"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.amount.String()
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestDate_String(t *testing.T) {
	tests := []struct {
		date *sanpltxt.Date
		want string
	}{
		{date(2020, 9, 1), "01-09-2020"},
		{date(2020, 9, 30), "30-09-2020"},
		{date(2020, 12, 31), "31-12-2020"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.date.String()
			assert.Equal(t, got, tt.want)
		})
	}
}

// PDF example: 1|51109010430000000100111111|50102055581111103350100016|Jerzy Kowalski|Warszawa ul. Kaliska 123 00-123|123,12|1|zasielenie konta|01-09-2020|7850000000|
func TestStandard_Marshal(t *testing.T) {
	s := &sanpltxt.Standard{
		DebitAccount:  "51109010430000000100111111",
		CreditAccount: "50102055581111103350100016",
		RecipientName: "Jerzy Kowalski",
		Address:       "Warszawa ul. Kaliska 123 00-123",
		Amount:        12312, // 123,12 PLN
		Mode:          sanpltxt.ModeElixir,
		Title:         "zasielenie konta",
		Date:          date(2020, 9, 1),
		NIP:           "7850000000",
	}

	got, err := s.Marshal()
	assert.NoError(t, err)

	want := "1|51109010430000000100111111|50102055581111103350100016|Jerzy Kowalski|Warszawa ul. Kaliska 123 00-123|123,12|1|zasielenie konta|01-09-2020|7850000000|"
	assert.Equal(t, got, want)
}

// PDF example: 2|51109010430000000100111111|82600000020260111122223333|ZUS|Warszawa ul. Szamocka 3,5 01748|319,94|1|Skladka ZUS|01-09-2020|
func TestZUS_Marshal(t *testing.T) {
	z := &sanpltxt.ZUS{
		DebitAccount:  "51109010430000000100111111",
		CreditAccount: "82600000020260111122223333",
		RecipientName: "ZUS",
		Address:       "Warszawa ul. Szamocka 3,5 01748",
		Amount:        31994, // 319,94 PLN
		Title:         "Skladka ZUS",
		Date:          date(2020, 9, 1),
	}

	got, err := z.Marshal()
	assert.NoError(t, err)

	want := "2|51109010430000000100111111|82600000020260111122223333|ZUS|Warszawa ul. Szamocka 3,5 01748|319,94|1|Skladka ZUS|01-09-2020|"
	assert.Equal(t, got, want)
}

// PDF example: 3|51109010430000000100111111|06101014690039392223000000|Urzad Skarbowy Poznan Winogrady|Poznan Wojciechowskiego 3/5 60-685|1000|01-09-2020|Jan Kowalski|N|9721230101|05|M|07|PIT5|id.zobowiazania|
func TestTax_Marshal(t *testing.T) {
	tax := &sanpltxt.Tax{
		TaxOffice:      true,
		DebitAccount:   "51109010430000000100111111",
		CreditAccount:  "06101014690039392223000000",
		RecipientName:  "Urzad Skarbowy Poznan Winogrady",
		Address:        "Poznan Wojciechowskiego 3/5 60-685",
		Amount:         100000, // 1000 PLN
		Date:           date(2020, 9, 1),
		PayerName:      "Jan Kowalski",
		IdentifierType: sanpltxt.IdentifierNIP,
		Identifier:     "9721230101",
		Year:           "05",
		PeriodType:     sanpltxt.PeriodMonth,
		PeriodNumber:   "07",
		FormSymbol:     "PIT5",
		ObligationID:   "id.zobowiazania",
	}

	got, err := tax.Marshal()
	assert.NoError(t, err)

	want := "3|51109010430000000100111111|06101014690039392223000000|Urzad Skarbowy Poznan Winogrady|Poznan Wojciechowskiego 3/5 60-685|1000|01-09-2020|Jan Kowalski|N|9721230101|05|M|07|PIT5|id.zobowiazania|"
	assert.Equal(t, got, want)
}

// PDF example: 5|51109010430000000100111111|50102055581111103350100016|Jan Nowak|Poznań ul. Swojska 17 06-123|1000,12|1|Wynagrodzenie za miesiąc|01-09-2020|
func TestPayroll_Marshal(t *testing.T) {
	p := &sanpltxt.Payroll{
		DebitAccount:  "51109010430000000100111111",
		CreditAccount: "50102055581111103350100016",
		RecipientName: "Jan Nowak",
		Address:       "Poznań ul. Swojska 17 06-123",
		Amount:        100012, // 1000,12 PLN
		Mode:          sanpltxt.ModeElixir,
		Title:         "Wynagrodzenie za miesiąc",
		Date:          date(2020, 9, 1),
	}

	got, err := p.Marshal()
	assert.NoError(t, err)

	want := "5|51109010430000000100111111|50102055581111103350100016|Jan Nowak|Poznań ul. Swojska 17 06-123|1000,12|1|Wynagrodzenie za miesiąc|01-09-2020|"
	assert.Equal(t, got, want)
}

// PDF example: 6|51109010430000000100111111|50102055581111103350100011|Jan Nowak|Warszawa ul. Mickiewicza 11 02-222|123,5|1|/VAT/23,09/IDC/8960005670/INV/5/2018/TXT/Faktura 5/2018|30-09-2020|
func TestSplitPayment_Marshal(t *testing.T) {
	sp := &sanpltxt.SplitPayment{
		DebitAccount:  "51109010430000000100111111",
		CreditAccount: "50102055581111103350100011",
		RecipientName: "Jan Nowak",
		Address:       "Warszawa ul. Mickiewicza 11 02-222",
		GrossAmount:   12350, // 123,50 PLN
		Mode:          sanpltxt.ModeElixir,
		VATAmount:     2309, // 23,09 PLN
		RecipientNIP:  "8960005670",
		InvoiceNumber: "5/2018",
		FreeText:      "Faktura 5/2018",
		Date:          date(2020, 9, 30),
	}

	got, err := sp.Marshal()
	assert.NoError(t, err)

	want := "6|51109010430000000100111111|50102055581111103350100011|Jan Nowak|Warszawa ul. Mickiewicza 11 02-222|123,50|1|/VAT/23,09/IDC/8960005670/INV/5/2018/TXT/Faktura 5/2018|30-09-2020|"
	assert.Equal(t, got, want)
}

func TestPackage_Marshal_Regular(t *testing.T) {
	pkg := &sanpltxt.Package{
		Type: 1,
		Transfers: []sanpltxt.Transfer{
			&sanpltxt.Standard{
				DebitAccount:  "51109010430000000100111111",
				CreditAccount: "50102055581111103350100016",
				RecipientName: "Jerzy Kowalski",
				Address:       "Warszawa ul. Kaliska 123 00-123",
				Amount:        12312,
				Mode:          sanpltxt.ModeElixir,
				Title:         "zasielenie konta",
				Date:          date(2020, 9, 1),
				NIP:           "7850000000",
			},
			&sanpltxt.ZUS{
				DebitAccount:  "51109010430000000100111111",
				CreditAccount: "82600000020260111122223333",
				RecipientName: "ZUS",
				Address:       "Warszawa ul. Szamocka 3,5 01748",
				Amount:        31994,
				Title:         "Skladka ZUS",
				Date:          date(2020, 9, 1),
			},
		},
	}

	got, err := pkg.Marshal()
	assert.NoError(t, err)

	lines := strings.Split(strings.TrimSuffix(got, "\n"), "\n")
	assert.Equal(t, len(lines), 3)
	assert.Equal(t, lines[0], "4120414|1")
	assert.True(t, strings.HasPrefix(lines[1], "1|"))
	assert.True(t, strings.HasPrefix(lines[2], "2|"))
}

func TestPackage_Marshal_Payroll(t *testing.T) {
	pkg := &sanpltxt.Package{
		Type: 2,
		Transfers: []sanpltxt.Transfer{
			&sanpltxt.Payroll{
				DebitAccount:  "51109010430000000100111111",
				CreditAccount: "50102055581111103350100016",
				RecipientName: "Jan Nowak",
				Address:       "Poznań ul. Swojska 17 06-123",
				Amount:        100012,
				Mode:          sanpltxt.ModeElixir,
				Title:         "Wynagrodzenie za miesiąc",
				Date:          date(2020, 9, 1),
			},
		},
	}

	got, err := pkg.Marshal()
	assert.NoError(t, err)

	lines := strings.Split(strings.TrimSuffix(got, "\n"), "\n")
	assert.Equal(t, len(lines), 2)
	assert.Equal(t, lines[0], "4120414|2")
	assert.True(t, strings.HasPrefix(lines[1], "5|"))
}

func TestPackage_Marshal_RejectsPayrollInRegular(t *testing.T) {
	pkg := &sanpltxt.Package{
		Type: 1, // Regular package
		Transfers: []sanpltxt.Transfer{
			&sanpltxt.Payroll{ // Payroll transfer - should be rejected
				DebitAccount:  "51109010430000000100111111",
				CreditAccount: "50102055581111103350100016",
				RecipientName: "Jan Nowak",
				Address:       "Poznań ul. Swojska 17 06-123",
				Amount:        100012,
				Mode:          sanpltxt.ModeElixir,
				Title:         "Wynagrodzenie za miesiąc",
			},
		},
	}

	_, err := pkg.Marshal()
	assert.Error(t, err)
}

func TestPackage_Marshal_RejectsStandardInPayroll(t *testing.T) {
	pkg := &sanpltxt.Package{
		Type: 2, // Payroll package
		Transfers: []sanpltxt.Transfer{
			&sanpltxt.Standard{ // Standard transfer - should be rejected
				DebitAccount:  "51109010430000000100111111",
				CreditAccount: "50102055581111103350100016",
				RecipientName: "Jerzy Kowalski",
				Address:       "Warszawa ul. Kaliska 123 00-123",
				Amount:        12312,
				Mode:          sanpltxt.ModeElixir,
				Title:         "zasielenie konta",
			},
		},
	}

	_, err := pkg.Marshal()
	assert.Error(t, err)
}

func TestValidation_InvalidNRB(t *testing.T) {
	s := &sanpltxt.Standard{
		DebitAccount:  "12345", // Too short
		CreditAccount: "50102055581111103350100016",
		RecipientName: "Test",
		Address:       "Test",
		Amount:        100,
		Mode:          sanpltxt.ModeElixir,
		Title:         "Test",
	}

	_, err := s.Marshal()
	assert.Error(t, err)
}

func TestValidation_InvalidRecipientName(t *testing.T) {
	s := &sanpltxt.Standard{
		DebitAccount:  "51109010430000000100111111",
		CreditAccount: "50102055581111103350100016",
		RecipientName: strings.Repeat("a", 81), // Too long
		Address:       "Test",
		Amount:        100,
		Mode:          sanpltxt.ModeElixir,
		Title:         "Test",
	}

	_, err := s.Marshal()
	assert.Error(t, err)
}

func TestValidation_InvalidNIP(t *testing.T) {
	s := &sanpltxt.Standard{
		DebitAccount:  "51109010430000000100111111",
		CreditAccount: "50102055581111103350100016",
		RecipientName: "Test",
		Address:       "Test",
		Amount:        100,
		Mode:          sanpltxt.ModeElixir,
		Title:         "Test",
		NIP:           "123", // Too short
	}

	_, err := s.Marshal()
	assert.Error(t, err)
}

func TestEncoding_Windows1250(t *testing.T) {
	// Test Polish characters encoding
	input := "Poznań ąćęłńóśźż ĄĆĘŁŃÓŚŹŻ"

	encoded, err := sanpltxt.ToWindows1250(input)
	assert.NoError(t, err)

	decoded, err := sanpltxt.FromWindows1250(encoded)
	assert.NoError(t, err)
	assert.Equal(t, decoded, input)
}
