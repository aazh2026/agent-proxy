package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	encryptor, err := NewEncryptor(key)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	plaintext := []byte("test token data")

	encrypted, err := encryptor.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	decrypted, err := encryptor.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Decrypted data doesn't match original")
	}
}

func TestEncryptDecryptString(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	encryptor, err := NewEncryptor(key)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	plaintext := []byte("test token data")

	encrypted, err := encryptor.EncryptToString(plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	decrypted, err := encryptor.DecryptFromString(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Decrypted data doesn't match original")
	}
}

func TestInvalidKeyLength(t *testing.T) {
	_, err := NewEncryptor([]byte("short"))
	if err != ErrInvalidKeyLength {
		t.Errorf("Expected ErrInvalidKeyLength, got %v", err)
	}
}

func TestMaskToken(t *testing.T) {
	tests := []struct {
		token    string
		expected string
	}{
		{"sk-1234567890abcdef", "sk-****cdef"},
		{"short", "****"},
		{"12345678901234567890", "1234****34567890"},
	}

	for _, test := range tests {
		result := MaskToken(test.token)
		if result != test.expected {
			t.Errorf("MaskToken(%s) = %s, expected %s", test.token, result, test.expected)
		}
	}
}

func TestGenerateKey(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}
	if len(key) != 32 {
		t.Errorf("Expected key length 32, got %d", len(key))
	}
}

func TestKeyFromBase64(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	encoded, err := GenerateKeyBase64()
	if err != nil {
		t.Fatalf("Failed to generate base64 key: %v", err)
	}

	decoded, err := KeyFromBase64(encoded)
	if err != nil {
		t.Fatalf("Failed to decode base64 key: %v", err)
	}

	if len(decoded) != 32 {
		t.Errorf("Expected key length 32, got %d", len(decoded))
	}

	_ = key
}
