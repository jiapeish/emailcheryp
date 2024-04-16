package emailcheryp

import (
	"fmt"
	"os"
	"testing"
)

var (
	samples = []struct {
		mail    string
		format  bool
		account bool // local part + domain
	}{
		// 1. Syntax validation: OK, Account: exist
		{mail: "jiapeish@gmail.com", format: true, account: true},
		{mail: "florian@carrere.cc", format: true, account: true},
		{mail: " florian@carrere.cc", format: true, account: true},
		{mail: "florian@carrere.cc ", format: true, account: true},

		// 2. Syntax validation: OK, Account: not exist
		{mail: "support@g2mail.com", format: true, account: false},
		{mail: "test@912-wrong-domain902.com", format: true, account: false},
		{mail: "admin@notarealdomain12345.com", format: true, account: false},
		{mail: "a@gmail.xyz", format: true, account: false},

		// this email address is syntax validate, but not ISP-Specific syntax validate
		// https://verifalia.com/validate-email
		// we need add rules later
		{mail: "0932910-qsdcqozuioqkdmqpeidj8793@gmail.com", format: true, account: false},
		{mail: " test@gmail.com", format: true, account: false},

		// 3. Syntax validation: not OK, Account: not exist
		{mail: "@gmail.com", format: false, account: false},
		{mail: "test@gmail@gmail.com", format: false, account: false},
		{mail: "test test@gmail.com", format: false, account: false},
		{mail: "test@wrong domain.com", format: false, account: false},
		{mail: "", format: false, account: false},
		{mail: "not-a-valid-email", format: false, account: false},

		// some validate tools consider these addresses as Syntax validation OK.
		// https://verifalia.com/validate-email
		{mail: "é&ààà@gmail.com", format: false, account: false},
	}
)

func TestValidateFormat(t *testing.T) {
	for _, s := range samples {
		err := ValidateFormat(s.mail)
		if err != nil && s.format == true {
			t.Errorf(`"%s" => unexpected error: "%v"`, s.mail, err)
		}
		if err == nil && s.format == false {
			t.Errorf(`"%s" => expected error`, s.mail)
		}
	}
}

func TestValidateDomain(t *testing.T) {
	for _, s := range samples {
		if !s.format {
			continue
		}

		err := ValidateDomain(s.mail)
		if err != nil && s.account == true {
			t.Errorf(`"%s" => unexpected error: "%v"`, s.mail, err)
		}
		if err == nil && s.account == false {
			t.Errorf(`"%s" => expected error`, s.mail)
		}
	}
}

func TestValidateMX(t *testing.T) {
	for _, s := range samples {
		if !s.format {
			continue
		}

		err := ValidateMX(s.mail)
		if err != nil && s.account == true {
			t.Errorf(`"%s" => unexpected error: "%v"`, s.mail, err)
		}
		if err == nil && s.account == false {
			t.Errorf(`"%s" => expected error`, s.mail)
		}
	}
}

func TestValidateLocalAndDomain(t *testing.T) {
	var (
		serverHostName    = getenv(t, "self_hostname")
		serverMailAddress = getenv(t, "self_mail")
	)
	for _, s := range samples {
		if !s.format {
			continue
		}

		err := ValidateLocalAndDomain(serverHostName, serverMailAddress, s.mail)
		if err != nil && s.account == true {
			t.Errorf(`"%s" => unexpected error: "%v"`, s.mail, err)
		}
		if err == nil && s.account == false {
			t.Errorf(`"%s" => expected error`, s.mail)
		}
	}
}

func getenv(t *testing.T, name string) (value string) {
	name = "test_checkmail_" + name
	if value = os.Getenv(name); value == "" {
		panic(fmt.Errorf("enviroment variable %q is not defined", name))
	}
	return
}
