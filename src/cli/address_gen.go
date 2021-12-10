package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ness-network/privateness/src/cipher"
	"github.com/ness-network/privateness/src/cipher/bip39"
	"github.com/ness-network/privateness/src/cipher/crypto"
	"github.com/ness-network/privateness/src/wallet"
)

func addressGenCmd() *cobra.Command {
	addressGenCmd := &cobra.Command{
		Short: "Generate skycoin or bitcoin addresses",
		Use:   "addressGen",
		Long: `Use caution when using the "-p" command. If you have command history enabled
    your wallet encryption password can be recovered from the history log. If you
    do not include the "-p" option you will be prompted to enter your password
    after you enter your command.`,
		SilenceUsage: true,
		RunE: func(c *cobra.Command, _ []string) error {
			numAddresses, err := c.Flags().GetInt("num")
			if err != nil {
				return err
			}
			if numAddresses <= 0 {
				return errors.New("num must be > 0")
			}

			coinName, err := c.Flags().GetString("coin")
			if err != nil {
				return err
			}

			coinType, err := wallet.ResolveCoinType(coinName)
			if err != nil {
				return err
			}

			encrypt, err := c.Flags().GetBool("encrypt")
			if err != nil {
				return err
			}

			mode, err := c.Flags().GetString("mode")
			if err != nil {
				return nil
			}

			// default label
			label, err := c.Flags().GetString("label")
			if err != nil {
				return nil
			}
			if label == "" {
				label = "default"
			}

			hideSecrets, err := c.Flags().GetBool("hide-secrets")
			if err != nil {
				return nil
			}

			seed, err := resolveSeed(c)
			if err != nil {
				return err
			}

			var password []byte
			if encrypt {
				switch strings.ToLower(mode) {
				case "json", "wallet":
				default:
					return errors.New("Encrypt flag requires -mode to be json")
				}

				var err error
				password, err = PasswordFromTerm{}.Password()
				if err != nil {
					return err
				}
			}

			w, err := wallet.NewWallet(wallet.NewWalletFilename(), label, seed, wallet.Options{
				Coin:       coinType,
				Encrypt:    encrypt,
				Password:   password,
				CryptoType: crypto.DefaultCryptoType,
				GenerateN:  uint64(numAddresses),
				Type:       wallet.WalletTypeDeterministic,
			})
			if err != nil {
				return err
			}

			if hideSecrets {
				w.Erase()
			}

			//rw := w.ToReadable()

			switch strings.ToLower(mode) {
			case "json", "wallet":
				output, err := w.Serialize()
				if err != nil {
					return err
				}

				fmt.Println(string(output))
			case "addrs", "addresses":
				es, err := w.GetEntries()
				if err != nil {
					return err
				}
				for _, e := range es {
					fmt.Println(e.Address)
				}
			case "secrets":
				if hideSecrets {
					return errors.New("secrets mode selected but hideSecrets enabled")
				}
				es, err := w.GetEntries()
				if err != nil {
					return err
				}

				for _, e := range es {
					switch coinType {
					case wallet.CoinTypeSkycoin:
						fmt.Println(e.Secret.Hex())
					case wallet.CoinTypeBitcoin:
						fmt.Println(cipher.BitcoinWalletImportFormatFromSeckey(e.Secret))
					}
				}
			default:
				return errors.New("invalid mode")
			}

			return nil
		},
	}

	addressGenCmd.Flags().IntP("num", "n", 1, "Number of addresses to generate")
	addressGenCmd.Flags().StringP("coin", "c", "skycoin", "Coin type. Must be skycoin or bitcoin. If bitcoin, secret keys are in Wallet Import Format instead of hex.")
	addressGenCmd.Flags().StringP("label", "l", "", "Wallet label to use when printing or writing a wallet file")
	addressGenCmd.Flags().Bool("hex", false, "Use hex(sha256sum(rand(1024))) (CSPRNG-generated) as the seed if not seed is not provided")
	addressGenCmd.Flags().StringP("seed", "s", "", "Seed for deterministic key generation. Will use bip39 as the seed if not provided.")
	addressGenCmd.Flags().BoolP("strict-seed", `t`, false, "Seed should be a valid bip39 mnemonic seed.")
	addressGenCmd.Flags().IntP("entropy", "e", 128, "Entropy of the autogenerated bip39 seed, when the seed is not provided. Can be 128 or 256")
	addressGenCmd.Flags().BoolP("hide-secrets", "i", false, "Hide the secret key and seed from the output when printing a JSON wallet file")
	addressGenCmd.Flags().StringP("mode", "m", "wallet", "Output mode. Options are wallet (prints a full JSON wallet), addresses (prints addresses in plain text), secrets (prints secret keys in plain text)")
	addressGenCmd.Flags().BoolP("encrypt", "x", false, "Encrypt the wallet when printing a JSON wallet")

	return addressGenCmd
}

func resolveSeed(c *cobra.Command) (string, error) {
	entropy, err := c.Flags().GetInt("entropy")
	if err != nil {
		return "", nil
	}

	switch entropy {
	case 128, 256:
	default:
		return "", errors.New("entropy must be 128 or 256")
	}

	seed, err := c.Flags().GetString("seed")
	if err != nil {
		return "", nil
	}

	strictSeed, err := c.Flags().GetBool("strict-seed")
	if err != nil {
		return "", nil
	}

	if seed != "" {
		if strictSeed {
			if err := bip39.ValidateMnemonic(seed); err != nil {
				return "", fmt.Errorf("seed is not a valid bip39 seed: %v", err)
			}
		}

		return seed, nil
	}

	useHex, err := c.Flags().GetBool("hex")
	if err != nil {
		return "", nil
	}

	if useHex {
		seed = cipher.SumSHA256(cipher.RandByte(1024)).Hex()
	} else {
		e, err := bip39.NewEntropy(entropy)
		if err != nil {
			return "", err
		}

		seed, err = bip39.NewMnemonic(e)
		if err != nil {
			return "", err
		}
	}

	return seed, nil
}

func fiberAddressGenCmd() *cobra.Command {
	fiberAddressGenCmd := &cobra.Command{
		Use:   "fiberAddressGen",
		Short: "Generate addresses and seeds for a new fiber coin",
		Long: `Addresses are written in a format that can be copied into fiber.toml
    for configuring distribution addresses. Addresses along with their seeds are written to a csv file,
    these seeds can be imported into the wallet to access distribution coins.`,
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) (err error) {
			if len(args) > 0 {
				return errors.New("this command does not take any positional arguments")
			}

			fiberNumAddresses, err := c.Flags().GetInt("num")
			if err != nil {
				return nil
			}

			entropy, err := c.Flags().GetInt("entropy")
			if err != nil {
				return nil
			}

			addrsFilename, err := c.Flags().GetString("addrs-file")
			if err != nil {
				return nil
			}

			seedsFilename, err := c.Flags().GetString("seeds-file")
			if err != nil {
				return nil
			}

			overwrite, err := c.Flags().GetBool("overwrite")
			if err != nil {
				return nil
			}

			if fiberNumAddresses < 1 {
				return errors.New("num must be > 0")
			}

			switch entropy {
			case 128, 256:
			default:
				return errors.New("entropy must be 128 or 256")
			}

			addrs := make([]cipher.Address, fiberNumAddresses)
			seeds := make([]string, fiberNumAddresses)

			for i := 0; i < fiberNumAddresses; i++ {
				e, err := bip39.NewEntropy(entropy)
				if err != nil {
					return err
				}

				seed, err := bip39.NewMnemonic(e)
				if err != nil {
					return err
				}

				_, seckey, err := cipher.GenerateDeterministicKeyPair([]byte(seed))
				if err != nil {
					return err
				}
				addr := cipher.MustAddressFromSecKey(seckey)

				seeds[i] = seed
				addrs[i] = addr
			}

			_, err = os.Stat(addrsFilename)
			if err != nil {
				if !os.IsNotExist(err) {
					return err
				}
			} else if !overwrite {
				return fmt.Errorf("-addrs-file %q already exists. Use -overwrite to force writing", addrsFilename)
			}

			_, err = os.Stat(seedsFilename)
			if err != nil {
				if !os.IsNotExist(err) {
					return err
				}
			} else if !overwrite {
				return fmt.Errorf("-seeds-file %q already exists. Use -overwrite to force writing", seedsFilename)
			}

			addrsF, err := os.Create(addrsFilename)
			if err != nil {
				return err
			}
			defer addrsF.Close()

			seedsF, err := os.Create(seedsFilename)
			if err != nil {
				return err
			}
			defer seedsF.Close()

			for i, a := range addrs {
				if _, err := fmt.Fprintf(addrsF, "\"%s\",\n", a); err != nil {
					return err
				}
				if _, err := fmt.Fprintf(seedsF, "\"%s\",\"%s\"\n", a, seeds[i]); err != nil {
					return err
				}
			}

			if err := addrsF.Sync(); err != nil {
				return err
			}

			return seedsF.Sync()
		},
	}

	fiberAddressGenCmd.Flags().IntP("num", "n", 100, "Number of addresses to generate")
	fiberAddressGenCmd.Flags().IntP("entropy", "e", 128, "Entropy of the autogenerated bip39 seeds. Can be 128 or 256")
	fiberAddressGenCmd.Flags().StringP("addrs-file", "a", "addresses.txt", "Output file for the generated addresses in fiber.toml format")
	fiberAddressGenCmd.Flags().StringP("seeds-file", "s", "seeds.csv", "Output file for the generated addresses and seeds in a csv")
	fiberAddressGenCmd.Flags().BoolP("overwrite", "o", false, "Allow overwriting any existing addrs-file or seeds-file")

	return fiberAddressGenCmd
}
