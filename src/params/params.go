package params

/*
CODE GENERATED AUTOMATICALLY WITH FIBER COIN CREATOR
AVOID EDITING THIS MANUALLY
*/

var (
	// MainNetDistribution Skycoin mainnet coin distribution parameters
	MainNetDistribution = Distribution{
		MaxCoinSupply:        200000000,
		InitialUnlockedCount: 10,
		UnlockAddressRate:    5,
		UnlockTimeInterval:   31536000,
		Addresses: []string{
    "GxQmHiuSeD4qZNPFZaoqKsvnLZ48oepr7y",
    "23EHCCywNmG5sWQwHbu6MxEfsfBMSCGnk2R",
    "22481dKneqmHzs1NbEi8Xoazc58WwRF4usL",
    "28NS66nUjzzWGKf7S2bHFxdfssNzPBLVav4",
    "2P8nPWCZLkCE7BrjuHDNJtiUELUtVGWntry",
    "6Cr1bHEsPYUgFJ5hAhNVDNjush1UouHwqr",
    "29Cbm1jTn97Uzr3mfNxgxdZEsizufp15jxb",
    "Sku4R4xLR3dACamSuNUriYv5VzHY7fyEKE",
    "2G3HdJrJoRBZ2Jc86D7MYoLXQ1uv6bVfKed",
    "oD53aAJYXo3wD9RMZBoseEXh57222oUEzP",
		},
	}

	// UserVerifyTxn transaction verification parameters for user-created transactions
	UserVerifyTxn = VerifyTxn{
		// BurnFactor can be overriden with `USER_BURN_FACTOR` env var
		BurnFactor: 2,
		// MaxTransactionSize can be overriden with `USER_MAX_TXN_SIZE` env var
		MaxTransactionSize: 32768, // in bytes
		// MaxDropletPrecision can be overriden with `USER_MAX_DECIMALS` env var
		MaxDropletPrecision: 3,
	}
)
