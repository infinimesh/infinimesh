package credentials

import (
	"context"
	"testing"

	"github.com/infinimesh/proto/node/accounts"
	"go.uber.org/zap"
)

func TestLDAPAuth(t *testing.T) {
	// TODO: Fix this test
	t.SkipNow()

	ctrl := NewCredentialsController(context.TODO(), zap.NewExample(), nil)

	t.Log("LDAP Configured", LDAP_CONFIGURED)
	if !LDAP_CONFIGURED {
		t.Fail()
	}

	t.Log("Test Make Credentials with no provider key")
	_, err := ctrl.MakeCredentials(&accounts.Credentials{
		Type: "ldap", Data: []string{"user"},
	})
	if err == nil {
		t.Fatalf("Expected error but could create credentials with no Provider Key")
	}

	t.Log("Test Make Credentials with wrong provider key")
	_, err = ctrl.MakeCredentials(&accounts.Credentials{
		Type: "ldap", Data: []string{"user", "unexistent"},
	})
	if err == nil {
		t.Fatalf("Expected error but could create credentials with wrong Provider Key")
	}

	t.Log("Test Make Credentials")
	cred, err := ctrl.MakeCredentials(&accounts.Credentials{
		Type: "ldap", Data: []string{"user", "local"},
	})
	if err != nil {
		t.Fatalf("Couldn't create credentials: %v", err)
	}

	t.Log("Test Authorize with Wrong Password")
	ok := cred.Authorize("user", "wrongpassword")
	if ok {
		t.Fatalf("Could authorize with wrong credentials")
	}

	t.Log("Test Authorize with Correct Password")
	ok = cred.Authorize("user", "password")
	if !ok {
		t.Fatalf("Couldn't authorize with correct credentials")
	}
}
