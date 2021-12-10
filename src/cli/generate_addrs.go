package cli

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ness-network/privateness/src/cipher"
	"github.com/ness-network/privateness/src/wallet"
)

func walletAddAddressesCmd() *cobra.Command {
	walletAddAddressesCmd := &cobra.Command{
		Args:  cobra.ExactArgs(1),
		Use:   "walletAddAddresses [wallet]",
		Short: "Generate additional addresses for a deterministic, bip44 or xpub wallet",
		Long: `Generate additional addresses for a deterministic, bip44 or xpub wallet.
    Addresses are generated according to the wallet type's generation mechanism.

    Warning: if you generate long (over 20) sequences of empty addresses and use
    a later address this can cause the wallet history scanner to miss your addresses,
    if you load the wallet from seed elsewhere. In that case, you'll have to manually
    generate addresses to cover the gap of unused addresses in the sequence.

    BIP44 wallets generate their addresses on the external (0'/0) chain.

    Use caution when using the "-p" command. If you have command
    history enabled your wallet encryption password can be recovered from the
    history log. If you do not include the "-p" option you will be prompted to
    enter your password after you enter your command.`,
		RunE: generateAddrs,
	}

	walletAddAddressesCmd.Flags().Uint64P("num", "n", 1, "Number of addresses to generate")
	walletAddAddressesCmd.Flags().StringP("password", "p", "", "wallet password")
	walletAddAddressesCmd.Flags().BoolP("json", "j", false, "Returns the results in JSON format")
	walletAddAddressesCmd.Flags().StringP("private-keys", "", "", "wallet private keys for collection wallet")

	return walletAddAddressesCmd
}

func generateAddrs(c *cobra.Command, args []string) error {
	// get number of address that are need to be generated.
	var opts []wallet.Option
	jsonFmt, err := c.Flags().GetBool("json")
	if err != nil {
		return err
	}

	wltID := args[0]

	// get the wallet to check if it is encrypted
	wlt, err := apiClient.Wallet(wltID)
	if err != nil {
		return err
	}

	switch wlt.Meta.Type {
	case wallet.WalletTypeCollection:
		s, err := c.Flags().GetString("private-keys")
		if err != nil {
			return err
		}
		privateKeys, err := wallet.ParsePrivateKeys(s)
		if err != nil {
			return err
		}
		opts = append(opts, wallet.OptionCollectionPrivateKeys(privateKeys))
	default:
		num, err := c.Flags().GetUint64("num")
		if err != nil {
			return err
		}

		if num == 0 {
			return errors.New("-n must > 0")
		}

		opts = append(opts, wallet.OptionGenerateN(num))
	}

	var pwd []byte
	pr := NewPasswordReader([]byte(c.Flag("password").Value.String()))
	if wlt.Meta.Encrypted && wlt.Meta.Type != wallet.WalletTypeBip44 {
		pwd, err = pr.Password()
		if err != nil {
			return err
		}
	}

	addrs, err := apiClient.NewWalletAddress(wltID, string(pwd), opts...)
	if err != nil {
		return err
	}

	if jsonFmt {
		s, err := FormatAddressesAsJSON(addrs)
		if err != nil {
			return err
		}
		fmt.Println(s)
	} else {
		fmt.Println(FormatAddressesAsJoinedArray(addrs))
	}

	return nil
}

// GenerateAddressesInFile generates addresses in given wallet file
func GenerateAddressesInFile(walletFile string, num uint64, pr PasswordReader) ([]cipher.Addresser, error) {
	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return nil, WalletLoadError{err}
	}

	switch pr.(type) {
	case nil:
		if wlt.IsEncrypted() {
			return nil, wallet.ErrWalletEncrypted
		}
	case PasswordFromBytes:
		p, err := pr.Password()
		if err != nil {
			return nil, err
		}

		if !wlt.IsEncrypted() && len(p) != 0 {
			return nil, wallet.ErrWalletNotEncrypted
		}
	}

	genAddrsInWallet := func(w wallet.Wallet, n uint64) ([]cipher.Addresser, error) {
		return w.GenerateAddresses(wallet.OptionGenerateN(n))
	}

	if wlt.IsEncrypted() {
		genAddrsInWallet = func(w wallet.Wallet, n uint64) ([]cipher.Addresser, error) {
			password, err := pr.Password()
			if err != nil {
				return nil, err
			}

			var addrs []cipher.Addresser
			if err := wallet.GuardUpdate(w, password, func(wlt wallet.Wallet) error {
				var err error
				addrs, err = wlt.GenerateAddresses(wallet.OptionGenerateN(n))
				return err
			}); err != nil {
				return nil, err
			}

			return addrs, nil
		}
	}

	addrs, err := genAddrsInWallet(wlt, num)
	if err != nil {
		return nil, err
	}

	dir, err := filepath.Abs(filepath.Dir(walletFile))
	if err != nil {
		return nil, err
	}

	if err := wallet.Save(wlt, dir); err != nil {
		return nil, WalletSaveError{err}
	}

	return addrs, nil
}

// FormatAddressesAsJSON converts []cipher.Address to strings and formats the array into a standard JSON object wrapper
func FormatAddressesAsJSON(addrs []string) (string, error) {
	d, err := formatJSON(struct {
		Addresses []string `json:"addresses"`
	}{
		Addresses: addrs,
	})

	if err != nil {
		return "", err
	}

	return string(d), nil
}

// FormatAddressesAsJoinedArray converts []cipher.Address to strings and concatenates them with a comma
func FormatAddressesAsJoinedArray(addrs []string) string {
	return strings.Join(addrs, ",")
}

// AddressesToStrings converts []cipher.Address to []string
func AddressesToStrings(addrs []cipher.Addresser) []string {
	if addrs == nil {
		return nil
	}

	addrsStr := make([]string, len(addrs))
	for i, a := range addrs {
		addrsStr[i] = a.String()
	}

	return addrsStr
}
