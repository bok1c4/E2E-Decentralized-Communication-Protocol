package util

import (
	"bytes"
	"errors"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
)

func IsValidPublicKey(key string) error {
	block, err := armor.Decode(bytes.NewReader([]byte(key)))
	if err != nil {
		return errors.New("invalid PGP key format")
	}
	if block.Type != openpgp.PublicKeyType {
		return errors.New("not a public PGP key")
	}

	_, err = openpgp.ReadKeyRing(block.Body)
	if err != nil {
		return errors.New("failed to parse PGP key ring")
	}
	return nil
}
