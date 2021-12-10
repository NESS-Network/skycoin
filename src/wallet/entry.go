package wallet

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ness-network/privateness/src/cipher"
)

// Entry represents the wallet entry
type Entry struct {
	Address     cipher.Addresser
	Public      cipher.PubKey
	Secret      cipher.SecKey
	ChildNumber uint32 // For bip32/bip44
	Change      uint32 // For bip44
}

// SkycoinAddress returns the Skycoin address of an entry. Panics if Address is not a Skycoin address
func (we Entry) SkycoinAddress() cipher.Address {
	return we.Address.(cipher.Address)
}

// BitcoinAddress returns the Skycoin address of an entry. Panics if Address is not a Bitcoin address
func (we Entry) BitcoinAddress() cipher.BitcoinAddress {
	return we.Address.(cipher.BitcoinAddress)
}

// Verify checks that the public key is derivable from the secret key,
// and that the public key is associated with the address
func (we *Entry) Verify() error {
	pk, err := cipher.PubKeyFromSecKey(we.Secret)
	if err != nil {
		return err
	}

	if pk != we.Public {
		return errors.New("invalid public key for secret key")
	}

	return we.VerifyPublic()
}

// VerifyPublic checks that the public key is associated with the address
func (we *Entry) VerifyPublic() error {
	if err := we.Public.Verify(); err != nil {
		return err
	}
	return we.Address.Verify(we.Public)
}

// Entries are an array of wallet entries
type Entries []Entry

// Clone make an copy of the entire entries.
func (entries Entries) Clone() Entries {
	if len(entries) == 0 {
		return nil
	}
	return append(Entries{}, entries...)
}

// Has checks if entries contains the entry with specified address
func (entries Entries) Has(a cipher.Addresser) bool {
	// This doesn't use getEntry() to avoid copying an Entry in the return value,
	// which may contain a secret key
	for _, e := range entries {
		if e.Address == a {
			return true
		}
	}
	return false
}

// Get returns the entry of specific address
func (entries Entries) Get(a cipher.Addresser) (Entry, bool) {
	for _, e := range entries {
		if e.Address == a {
			return e, true
		}
	}
	return Entry{}, false
}

// GetAddresses returns all addresses
func (entries Entries) GetAddresses() []cipher.Addresser {
	addrs := make([]cipher.Addresser, len(entries))
	for i, e := range entries {
		addrs[i] = e.Address
	}
	return addrs
}

// Erase wipes private keys in entries
func (entries Entries) Erase() {
	for i := range entries {
		for j := range entries[i].Secret {
			entries[i].Secret[j] = 0
		}
		entries[i].Secret = cipher.SecKey{}
	}
}

// UnpackSecretKeys for each entry, look for the secret key in the Secrets dict, keyed by address
func (entries Entries) UnpackSecretKeys(ss Secrets) error {
	for i, e := range entries {
		sstr, ok := ss.Get(e.Address.String())
		if !ok {
			return fmt.Errorf("secret of address %s doesn't exist in secrets", e.Address)
		}

		s, err := hex.DecodeString(sstr)
		if err != nil {
			return fmt.Errorf("decode secret hex string failed: %v", err)
		}

		copy(entries[i].Secret[:], s[:])
	}

	return nil
}
