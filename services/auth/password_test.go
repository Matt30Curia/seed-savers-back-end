package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("password")

	if err != nil {
		t.Errorf("error hashing password")
	}
	if hash == "password" {
		t.Errorf("hash is expected instead plain password")
	}
	if hash == "" {
		t.Errorf("the expected hash must not be empty")
	}
}

func TestComparePasswords(t *testing.T) {
	hash, _ := HashPassword("password")

	if !ComparePasswords(hash, []byte("password")) {
		t.Errorf("error compare correct password")
	}
	if ComparePasswords(hash, []byte("not password")) {
		t.Errorf("expects an error because it compares different passwords")
	}

}
