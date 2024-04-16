package emailcheryp

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

var (
	samples = []struct {
		mail             string
		syntaxValidated  bool
		domainValidated  bool
		mxValidated      bool
		accountValidated bool // local part + domain
	}{
		// 1. Syntax validation: OK, Account: exist
		{
			mail:             "jiapeish@gmail.com",
			syntaxValidated:  true,
			domainValidated:  true,
			mxValidated:      true,
			accountValidated: true,
		},
		{
			mail:             "florian@carrere.cc",
			syntaxValidated:  true,
			domainValidated:  true,
			mxValidated:      true,
			accountValidated: true,
		},
		{
			mail:             " florian@carrere.cc",
			syntaxValidated:  true,
			domainValidated:  true,
			mxValidated:      true,
			accountValidated: true,
		},
		{
			mail:             "florian@carrere.cc ",
			syntaxValidated:  true,
			domainValidated:  true,
			mxValidated:      true,
			accountValidated: true,
		},

		// 2. Syntax validation: OK, Account: not exist
		{
			mail:             "support@g2mail.com",
			syntaxValidated:  true,
			domainValidated:  false,
			mxValidated:      false,
			accountValidated: false,
		},
		{
			mail:             "test@912-wrong-domain902.com",
			syntaxValidated:  true,
			domainValidated:  false,
			mxValidated:      false,
			accountValidated: false,
		},
		{
			mail:             "admin@notarealdomain12345.com",
			syntaxValidated:  true,
			domainValidated:  false,
			mxValidated:      false,
			accountValidated: false,
		},
		{
			mail:             "a@gmail.xyz",
			syntaxValidated:  true,
			domainValidated:  false,
			mxValidated:      false,
			accountValidated: false,
		},

		// this email address is syntax validate, but not ISP-Specific syntax validate
		// https://verifalia.com/validate-email
		// we need add rules later
		{
			mail:             "0932910-qsdcqozuioqkdmqpeidj8793@gmail.com",
			syntaxValidated:  true,
			domainValidated:  true,
			mxValidated:      true,
			accountValidated: false,
		},
		{
			mail:             " test@gmail.com",
			syntaxValidated:  true,
			domainValidated:  true,
			mxValidated:      true,
			accountValidated: false,
		},

		// 3. Syntax validation: not OK, Account: not exist
		{
			mail:             "@gmail.com",
			syntaxValidated:  false,
			domainValidated:  true,
			mxValidated:      true,
			accountValidated: false,
		},
		{
			mail:             "test@gmail@gmail.com",
			syntaxValidated:  false,
			domainValidated:  false,
			mxValidated:      false,
			accountValidated: false,
		},
		{
			mail:             "test test@gmail.com",
			syntaxValidated:  false,
			domainValidated:  true,
			mxValidated:      true,
			accountValidated: false,
		},
		{
			mail:             "test@wrong domain.com",
			syntaxValidated:  false,
			domainValidated:  false,
			mxValidated:      false,
			accountValidated: false,
		},
		{
			mail:             "",
			syntaxValidated:  false,
			domainValidated:  false,
			mxValidated:      false,
			accountValidated: false,
		},
		{
			mail:             "not-a-valid-email",
			syntaxValidated:  false,
			domainValidated:  false,
			mxValidated:      false,
			accountValidated: false,
		},

		// some validate tools consider these addresses as Syntax validation OK.
		// https://verifalia.com/validate-email
		{
			mail:             "é&ààà@gmail.com",
			syntaxValidated:  false,
			domainValidated:  true,
			mxValidated:      true,
			accountValidated: false,
		},
	}
)

func TestValidateFormat(t *testing.T) {
	for _, s := range samples {
		err := ValidateFormat(s.mail)
		if err != nil && s.syntaxValidated == true {
			t.Errorf(`"%s" => unexpected error: "%v"`, s.mail, err)
		}
		if err == nil && s.syntaxValidated == false {
			t.Errorf(`"%s" => expected error`, s.mail)
		}
	}
}

func TestValidateDomain(t *testing.T) {
	for _, s := range samples {
		if !s.syntaxValidated {
			continue
		}

		err := ValidateDomain(s.mail)
		if err != nil && s.domainValidated == true {
			t.Errorf(`"%s" => unexpected error: "%v"`, s.mail, err)
		}
		if err == nil && s.domainValidated == false {
			t.Errorf(`"%s" => expected error`, s.mail)
		}
	}
}

func TestValidateMX(t *testing.T) {
	for _, s := range samples {
		if !s.syntaxValidated {
			continue
		}

		err := ValidateMX(s.mail)
		if err != nil && s.mxValidated == true {
			t.Errorf(`"%s" => unexpected error: "%v"`, s.mail, err)
		}
		if err == nil && s.mxValidated == false {
			t.Errorf(`"%s" => expected error`, s.mail)
		}
	}
}

func TestValidateLocalAndDomain(t *testing.T) {
	var (
		serverHostName    = "gmail.com"
		serverMailAddress = "jiapeish2@gmail.com"
	)
	for _, s := range samples {
		mail := strings.TrimSpace(s.mail)
		if !s.syntaxValidated {
			continue
		}

		err := ValidateLocalAndDomain(serverHostName, serverMailAddress, mail)
		if err != nil && s.accountValidated == true {
			t.Errorf(`"%s" => unexpected error: "%v"`, mail, err)
		}
		if err == nil && s.accountValidated == false {
			t.Errorf(`"%s" => expected error`, mail)
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
