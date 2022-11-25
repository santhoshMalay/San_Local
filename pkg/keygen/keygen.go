// Package keygen wraps a key derivation function with the goal of conveniently generating cryptographically strong
// encryption keys based on secret strings
// The current implementation is based on HKDF (RFC 5869)
package keygen

import (
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/hkdf"
	"io"
)

var (
	initError error
	salt      []byte
)

// saltString is an arbitrary seed value used for all generated keys. It is non-secret. Changing this constant will
// change all key data sequences generated from secret strings
// Can be generated using the following code (bytesCount is sha512.New384().BlockSize() == 128):
//
//	func GetRandomString(bytesCount int) (string, error) {
//		bytes := make([]byte, bytesCount)
//		_, err := rand.Read(bytes)  // "crypto/rand"
//		if err != nil {
//			return "", err
//		}
//		return base64.RawURLEncoding.EncodeToString(bytes), nil
//	}
//
//goland:noinspection SpellCheckingInspection
const saltString = "liL1g3zrBg6UwNIR_1R1JC1txDfIzjGQPzVzkrhQtDgPkNe0iEOuaKbAY1zR7HRpz8KvA0aypSM8UDg5-rleJB5yBmehYxTSkc1kyge8I-fG-lyVZQX4KKGKvsSErRgvxOBQ9puxbFHnnonEOsmE6pimqu6vRx7PxXFwG-NSxYo"

func init() {
	bytes, err := base64.RawURLEncoding.DecodeString(saltString)
	if err != nil {
		initError = fmt.Errorf("keygen init failure: %w", err)
		return
	}
	salt = bytes
}

// KeyGen represents a source of deterministic byte stream which can be used as cryptographically strong encryption
// keys. The byte stream is determined by the secret string and optional contextInfo passed to keygen.New
// Implements io.Reader interface
type KeyGen struct {
	hkdf io.Reader
}

// New creates an instance of KeyGen parametrized with the provided secret string and optional context information
//
// secret must be a sufficiently long text string which is normally read from a configuration file or an env variable,
//
// contextInfo is either an empty string or some arbitrary text. Using different contextInfo values with the same secret
// allows to generate different sets of keys from the same secret value. contextInfo is not a sensitive info and can be
// hardcoded
func New(secret, contextInfo string) (*KeyGen, error) {
	if initError != nil {
		return nil, initError
	}
	var info []byte
	if len(contextInfo) > 0 {
		info = []byte(contextInfo)
	}
	gen := hkdf.New(sha512.New384, []byte(secret), salt, info)
	kg := &KeyGen{
		hkdf: gen,
	}
	return kg, nil
}

func (kg *KeyGen) Read(p []byte) (int, error) {
	return kg.hkdf.Read(p)
}

// Generate generates a single key of specified size (bytesCount) using the provided secret string and optional
// context information. This is a helper wrapper for the more generic KeyGen
//
// secret must be a sufficiently long text string which is normally read from a configuration file, an env variable,
// a key vault
//
// contextInfo is either an empty string or some arbitrary text. Using different contextInfo values with the same secret
// allows to generate different sets of keys from the same secret value. contextInfo is not a sensitive info and can be
// hardcoded
//
// bytesCount is the size of key in bytes
func Generate(secret, contextInfo string, bytesCount int) ([]byte, error) {
	if bytesCount <= 0 {
		return nil, errors.New("bad argument: bytesCount")
	}
	kg, err := New(secret, contextInfo)
	if err != nil {
		return nil, err
	}
	bytes := make([]byte, bytesCount)
	_, err = kg.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
