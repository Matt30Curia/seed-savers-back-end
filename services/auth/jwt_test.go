package auth

import "testing"

func TestJwtCreation(t *testing.T){
	secret := []byte("super secret")
	token, err := CreateJWT(secret, 1)

	if err != nil{
		t.Errorf("error on creation of jwt %v", err)
	}

	if token == ""{
		t.Errorf("token must not be empty")
	}
}