package fiber

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher/bip44"
)

// TODO(therealssj): write better tests
func TestNewConfig(t *testing.T) {
	coinConfig, err := NewConfig("test.fiber.toml", "./testdata")
	require.NoError(t, err)
	require.Equal(t, Config{
		Node: NodeConfig{
			GenesisSignatureStr: "0b0661652a064c48f5ec565b596cf3be1a438e9e1bd1de551f16f76172ae0a02628a5cecdd366aaba070786c2040c32113da871ca3a80d26902eb7566a319d6f00",
			GenesisAddressStr:   "2LuKycJ7SQwbSxoX65Bu8BwZ66RegDaWGK",
			BlockchainPubkeyStr: "0278d26405ec24c8bf8998cf767b65c29f0dfcdff3542cdb1de44ed0539e9c9d9b",
			BlockchainSeckeyStr: "",
			GenesisTimestamp:    1426562704,
			GenesisCoinVolume:   100e12,
			DefaultConnections: []string{
				"118.178.135.93:6000",
				"47.88.33.156:6000",
				"104.237.142.206:6000",
				"176.58.126.224:6000",
				"172.104.85.6:6000",
				"139.162.7.132:6000",
			},
			Port:                           6000,
			PeerListURL:                    "https://downloads.skycoin.com/blockchain/peers.txt",
			WebInterfacePort:               6420,
			UnconfirmedBurnFactor:          10,
			UnconfirmedMaxTransactionSize:  777,
			UnconfirmedMaxDropletPrecision: 3,
			CreateBlockBurnFactor:          9,
			CreateBlockMaxTransactionSize:  1234,
			CreateBlockMaxDropletPrecision: 4,
			MaxBlockTransactionsSize:       1111,
			DisplayName:                    "Testcoin",
			Ticker:                         "TST",
			CoinHoursName:                  "Testcoin Hours",
			CoinHoursNameSingular:          "Testcoin Hour",
			CoinHoursTicker:                "TCH",
			QrURIPrefix:                    "skycoin",
			ExplorerURL:                    "https://explorer.testcoin.com",
			VersionURL:                     "https://version.testcoin.com/testcoin/version.txt",
			Bip44Coin:                      bip44.CoinTypeSkycoin,
		},
		Params: ParamsConfig{
			MaxCoinSupply:           1e8,
			UnlockAddressRate:       5,
			InitialUnlockedCount:    25,
			UnlockTimeInterval:      60 * 60 * 24 * 365,
			UserBurnFactor:          3,
			UserMaxTransactionSize:  999,
			UserMaxDropletPrecision: 2,
		},
	}, coinConfig)
}
