package auth

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestTokenRoundTrip(t *testing.T) {
	// Create a token, read it back and compare the IDs
	wantUserID := uuid.New() 
	secret := "super-secret-test-key-32-bytes!!"
	validFor := 5 * time.Minute
	t.Run("valid round trip", func(t *testing.T) {
		tokenStr, err := MakeJWT(wantUserID, secret,validFor)
		if err != nil {
			t.Fatalf("create token: %v", err)
		}
		gotUserID, err := ValidateJWT(tokenStr, secret)
		if err != nil {
			t.Fatalf("failed, token does not match, %v", err)
		}
		if gotUserID != wantUserID {
			t.Errorf("User ID got %v, want %v", gotUserID, wantUserID)
		}
	})
}


func TestExpiredToken(t *testing.T) {
	// Create a token, read it back and compare the IDs
	wantUserID := uuid.New() 
	secret := "super-secret-test-key-32-bytes!!"
	validFor := -time.Minute
	t.Run("valid round trip", func(t *testing.T) {
		tokenStr, err := MakeJWT(wantUserID, secret,validFor)
		if err != nil {
			t.Fatalf("create token: %v", err)
		}
		_, err = ValidateJWT(tokenStr, secret)
		if err == nil {
			t.Fatal("expired token was accepted")
		}
	})
}


func TestDifferentKeys(t *testing.T) {
	// Create a token, read it back and compare the IDs
	wantUserID := uuid.New() 
	secret := "super-secret-test-key-32-bytes!!"
	scannerKey := "different-key"
	validFor := 5 * time.Minute
	t.Run("valid round trip", func(t *testing.T) {
		tokenStr, err := MakeJWT(wantUserID, secret,validFor)
		if err != nil {
			t.Fatalf("create token: %v", err)
		}
		_, err = ValidateJWT(tokenStr, scannerKey)
		if err == nil {
			t.Fatal("failed, token was validated with different keys")
		}
	})
}
