package globals

import (
	"bytes"
	"cw/logger"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func SetInit() {
	parsedABI, err := abi.JSON(bytes.NewReader(Erc20JSON))
	if err != nil {
		logger.GlobalLogger.Fatalf("Failed parsing ABI: %v", err)
	}

	Erc20ABI = &parsedABI
	_, success := MaxApproveValue.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
	if !success {
		logger.GlobalLogger.Fatalf("Failed to set MaxRepayBigInt: invalid number")
	}
}

// Interpretation for file paths
var (
	AppConfiguration = "config"
	UserConfig       = "userConfig"
	Withdraw         = "withdraw"
	Bridge           = "bridge"
	Addresses        = "addresses"
	Deposits         = "deposits"
	Proxy            = "proxy"
	PrivateKeys      = "private_keys"
	WallGen          = "wallet_generaor"
	Destinations     = "destination_address"
)

// ####### Default GLOBALS #########
var (
	// Max approve value.
	MaxApproveValue = new(big.Int)

	// Default ABI for erc20 tokens
	Erc20ABI *abi.ABI

	// Console title name.
	ConsoleTitle = "Cheif Machine | cheif.ssq"

	// Maximum permissible percentage impact. It is not recommended to set it higher than this value
	MaxPercent = 3.0

	// convertation user token name to software tokens names
	TokenSymbolToName = map[string]string{
		"BTC":    "bitcoin",
		"ETH":    "ethereum",
		"USDT":   "tether",
		"XRP":    "ripple",
		"BNB":    "bnb",
		"SOL":    "solana",
		"USDC":   "usdcoin",
		"DOGE":   "dogecoin",
		"ADA":    "cardano",
		"TRX":    "tron",
		"stETH":  "lido-staked-ether",
		"WBTC":   "wrapped-bitcoin",
		"TON":    "toncoin",
		"LINK":   "chainlink",
		"LEO":    "leobitfinex",
		"XLM":    "stellar",
		"AVAX":   "avalanche",
		"SUI":    "sui",
		"SHIB":   "shiba-inu",
		"HBAR":   "hedera-hashgraph",
		"LTC":    "litecoin",
		"DOT":    "polkadot",
		"OM":     "mantra-dao",
		"BCH":    "bitcoin-cash",
		"BGB":    "bitget-token",
		"PI":     "pi-network-1",
		"USDe":   "ethena-usde",
		"HYPE":   "hyperliquid",
		"weETH":  "wrapped-eeth",
		"XMR":    "monero",
		"UNI":    "uniswap",
		"DAI":    "multicollateraldai",
		"NEAR":   "near-protocol",
		"APT":    "aptos",
		"PEPE":   "pepe",
		"OKB":    "okb",
		"CRO":    "crypto-com-cro",
		"GT":     "gatechain-token",
		"MNT":    "mantlexyz",
		"ICP":    "internet-computer",
		"ONDO":   "ondo-finance",
		"AAVE":   "aave",
		"ETC":    "ethereum-classic",
		"FDUSD":  "first-digital-usd",
		"TAO":    "bittensor",
		"VET":    "vechainthor",
		"TIA":    "celestia",
		"TRUMP":  "official-trump",
		"ENA":    "ethena",
		"ATOM":   "cosmos",
		"RENDER": "render-token",
		"FIL":    "filecoin",
		"POL":    "polygon-ecosystem-token",
		"KAS":    "kaspa",
		"S":      "sonic",
		"ALGO":   "algorand",
		"FTN":    "fasttoken",
		"ARB":    "arbitrum",
		"FTM":    "fantom",
		"DEXE":   "dexe-network",
		"JUP":    "jupiter-stattion",
		"KCS":    "kucoin-shares",
		"OP":     "optimism",
		"IP":     "story-protocol",
		"FET":    "fetch-ai",
		"MKR":    "maker",
		"MOVE":   "movement-labs",
		"NEXO":   "nexo",
		"XDC":    "xdc-network",
		"IMX":    "immutable-x",
		"QNT":    "quant",
		"WLD":    "world-coin",
		"STX":    "blockstack",
		"BONK":   "bonk-new",
		"SEI":    "sei-network",
		"INJ":    "injective-protocol",
		"GRT":    "the-graph",
		"BERA":   "berachain",
		"FORM":   "binaryx-new",
		"THETA":  "theta-token",
		"EOS":    "eos",
		"FLR":    "flare",
		"LDO":    "lido-dao",
		"PYUSD":  "paypal-usd",
		"XAUT":   "tether-gold",
		"GALA":   "gala",
		"JTO":    "jito-labs",
		"XTZ":    "tezos",
		"SAND":   "sandbox",
		"IOTA":   "iota",
		"CRV":    "curve",
		"AB":     "ab",
		"BTT":    "bittorrent-new",
		"PAXG":   "pax-gold",
		"BSV":    "bitcoin-sv",
		"mSOL":   "marinade",
		"CAKE":   "pancakeswap",
		"KAIA":   "kaia-chain",
		"FLOW":   "flow",
		"FLOKI":  "floki",
	}

	// Bringing to a common standard of names in software due to different network names on different exchanges.
	ChainNameToSymbolCEX = map[string]map[string]string{
		"bybit": {
			"Arbitrum":      "Arbitrum One",
			"Aptos":         "APTOS",
			"Starknet":      "Starknet",
			"BNB":           "BNB Smart Chain",
			"ZkLite":        "zkSync Lite",
			"ZkEra":         "zkSync Era",
			"Tron":          "TRC20",
			"Optimism":      "OP Mainnet",
			"Polygon":       "Polygon PoS",
			"Sol":           "SOL",
			"Ton":           "Ton",
			"KAVAEVM":       "KAVAEVM",
			"Eth":           "Ethereum",
			"Celo":          "CELO",
			"Avax":          "CAVAX",
			"Base":          "Base Mainnet",
			"Mantle":        "Mantle Network",
			"Arbitrum Nova": "Arbitrum Nova",
			"Linea":         "LINEA",
			"BTC":           "BTC",
			"XRP":           "XRP",
			"ADA":           "ADA",
			"Algo":          "ALGO",
			"Atom":          "ATOM",
			"Dot":           "DOT",
			"Bera":          "Berachain",
			"Blast":         "BLAST",
			"Doge":          "Dogecoin",
			"Sui":           "SUI",
			"Story":         "STORY",
			"Ftm":           "FTM",
			"Sonic":         "SONIC",
			"ICP":           "ICP",
		},
		"binance": {
			"Arbitrum": "ARBITRUM",
			"Aptos":    "APT",
			"BNB":      "BSC",
			"ZkEra":    "ZKSYNCERA",
			"Tron":     "TRX",
			"Optimism": "OPTIMISM",
			"Polygon":  "MATIC",
			"Sol":      "SOL",
			"Ton":      "TON",
			"KAVAEVM":  "KAVAEVM",
			"Eth":      "ETH",
			"Celo":     "CELO",
			"Avax":     "AVAXC",
			"Base":     "BASE",
			"BTC":      "BTC",
			"XRP":      "XRP",
			"ADA":      "ADA",
			"Algo":     "ALGO",
			"Atom":     "ATOM",
			"Dot":      "DOT",
			"Bera":     "BERA",
			"Doge":     "DOGE",
			"Sui":      "SUI",
		},
		"mexc": {
			"Arbitrum": "ARB",
			"Aptos":    "APTOS",
			"Starknet": "STARK",
			"BNB":      "BSC",
			"ZkLite":   "ZKSYNC",
			"ZkEra":    "ZKSYNCERA",
			"Tron":     "TRX",
			"Optimism": "OP",
			"Polygon":  "MATIC",
			"Sol":      "SOL",
			"Ton":      "TONCOIN",
			// "KAVAEVM":       "KAVAEVM",
			"Eth":  "ETH",
			"Celo": "CELO",
			"Avax": "AVAX_CCHAIN",
			"Near": "NEAR",
			"Base": "BASE",
			// "Mantle":        "Mantle Network",
			// "Arbitrum Nova": "Arbitrum Nova",
			"Linea": "LINEA",
			"BTC":   "BTC",
			"XRP":   "XRP",
			"ADA":   "ADA",
			// "Algo":          "ALGO",
			"Atom":  "ATOM",
			"Dot":   "DOT",
			"Bera":  "BERACHAIN",
			"Blast": "BLAST",
			"Doge":  "DOGE",
			"Sui":   "SUI",
		},
		"kucoin": {
			"Arbitrum": "ARBITRUM",
			"Aptos":    "APT",
			// "Starknet": "STARK",
			"BNB": "BEP20",
			// "ZkLite":   "ZKSYNC",
			"ZkEra":    "ZKS20",
			"Tron":     "TRC20",
			"Optimism": "OPTIMISM",
			"Polygon":  "Polygon POS",
			"Sol":      "SOL",
			"KAVAEVM":  "KAVA EVM",
			"Eth":      "ERC20",
			"Celo":     "CELO",
			"Avax":     "AVAX C-Chain",
			"Near":     "NEAR",
			"Ton":      "TON",
			// "Base": "BASE",
			// // "Mantle":        "Mantle Network",
			// // "Arbitrum Nova": "Arbitrum Nova",
			// "Linea": "LINEA",
			"BTC":  "BTC",
			"XRP":  "XRP",
			"ADA":  "ADA",
			"Algo": "ALGO",
			"Atom": "ATOM",
			// "Dot":   "DOT",
			"Bera":  "Bera",
			"Blast": "BLAST",
			"Doge":  "DOGE",
			"Sui":   "SUI",
		},
		"okx": {
			"Arbitrum": "ARBONE",
			"Base":     "Base",
			"Starknet": "Starknet",
			"Aptos":    "APT",
			"BNB":      "BEP20",
			// "ZkLite":   "ZKSYNC",
			"ZkEra":    "zkSync Era",
			"Tron":     "TRC20",
			"Optimism": "OPTIMISM",
			"Polygon":  "MATIC",
			"Sol":      "SOL",
			// "KAVAEVM":  "KAVA EVM",
			"Eth":  "ERC20",
			"Celo": "CELO",
			"Avax": "AVAXC",
			"Near": "NEAR",
			"Ton":  "TON",
			// // "Mantle":        "Mantle Network",
			// // "Arbitrum Nova": "Arbitrum Nova",
			"Linea": "Linea",
			"BTC":   "BTC",
			"XRP":   "XRP",
			"ADA":   "ADA",
			"Algo":  "ALGO",
			"Atom":  "ATOM",
			"Dot":   "DOT",
			"Bera":  "Berachain",
			// "Blast": "BLAST",
			"Doge": "DOGE",
			"Sui":  "SUI",
		},
	}

	// ####### TOKEN CONTRACTS GLOBALS #########
	TokenContracts = map[int64]map[string]string{
		1: {
			"ETH":  "0x0000000000000000000000000000000000000000",
			"WETH": "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
			"USDT": "0xdAC17F958D2ee523a2206206994597C13D831ec7",
			"USDC": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		},
		33139: {
			"APE": "0x0000000000000000000000000000000000000000",
		},
		81457: {
			"ETH": "0x0000000000000000000000000000000000000000",
		},
		42161: {
			"ETH":  "0x0000000000000000000000000000000000000000",
			"USDC": "0xaf88d065e77c8cC2239327C5EDb3A432268e5831",
			"USDT": "0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9",
			"WETH": "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1",
		},
		43114: {
			"AVAX": "0x0000000000000000000000000000000000000000",
			"USDC": "0xB97EF9Ef8734C71904D8002F8b6Bc66Dd9c48a6E",
			"USDT": "0x9702230a8ea53601f5cd2dc00fdbc13d4df4a8c7",
		},
		8453: {
			"ETH":  "0x0000000000000000000000000000000000000000",
			"USDC": "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
			"WETH": "0x4200000000000000000000000000000000000006",
		},
		59144: {
			"ETH":  "0x0000000000000000000000000000000000000000",
			"USDC": "0x176211869cA2b568f2A7D4EE941E073a821EE1ff",
			"USDT": "0xA219439258ca9da29E9Cc4cE5596924745e12B93",
		},
		56: {
			"BNB":  "0x0000000000000000000000000000000000000000",
			"USDT": "0x55d398326f99059ff775485246999027b3197955",
			"USDC": "0x8ac76a51cc950d9822d68b83fe1ad97b32cd580d",
		},
		10: {
			"ETH":  "0x0000000000000000000000000000000000000000",
			"USDC": "0x0b2c639c533813f4aa9d7837caf62653d097ff85",
			"USDT": "0x94b008aa00579c1307b0ef2c499ad98a8ce58e58",
		},
		1135: {
			"ETH": "0x0000000000000000000000000000000000000000",
		},
	}

	// Exploers links. Use chain id for Using id to determine the type of exploer
	ExploerLink = map[int64]string{
		1:       "https://etherscan.io",
		42161:   "https://arbiscan.io",
		33139:   "https://apescan.io",
		43114:   "https://snowtrace.io",
		8453:    "https://basescan.org",
		81457:   "https://blastscan.io/",
		56:      "https://bscscan.com",
		59144:   "https://lineascan.build",
		1135:    "https://blockscout.lisk.com",
		5000:    "https://mantlescan.xyz/",
		34443:   "https://modescan.io/",
		10:      "https://optimistic.etherscan.io",
		137:     "https://polygonscan.com",
		1101:    "https://zkevm.polygonscan.com/",
		324:     "https://zkscan.io/",
		7777777: "https://zora.thesuperscan.io/",
	}
)

// ###### Base ERC20 ABI. #######
// # Supplement other ABIs here #
var (
	Erc20JSON = []byte(`[
	{
		
		"constant":true,
		"inputs":[{"name":"account","type":"address"}],
		"name":"balanceOf",
		"outputs":[{"name":"","type":"uint256"}],
		"payable":false,
		"stateMutability":"view",
		"type":"function"
	},
	{
		"constant":true,
		"inputs":[{"name":"spender","type":"address"},{"name":"owner","type":"address"}],
		"name":"allowance",
		"outputs":[{"name":"","type":"uint256"}],
		"payable":false,
		"stateMutability":"view",
		"type":"function"
	},
	{
		"constant": false,
		"inputs": [],
		"name": "deposit",
		"outputs": [],
		"payable": true,
		"stateMutability": "payable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
			"internalType": "uint256",
			"name": "wad",
			"type": "uint256"
			}
		],
		"name": "withdraw",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant":false,
		"inputs":[{"name":"spender","type":"address"},{"name":"amount","type":"uint256"}],
		"name":"approve",
		"outputs":[{"name":"","type":"bool"}],
		"payable":false,
		"stateMutability":"nonpayable",
		"type":"function"
	},
	{
		"constant":false,
		"inputs":[{"name":"recipient","type":"address"},{"name":"amount","type":"uint256"}],
		"name":"transfer",
		"outputs":[{"name":"","type":"bool"}],
		"payable":false,
		"stateMutability":"nonpayable",
		"type":"function"
	},
	{
		"constant":false,
		"inputs":[{"name":"sender","type":"address"},{"name":"recipient","type":"address"},{"name":"amount","type":"uint256"}],
		"name":"transferFrom",
		"outputs":[{"name":"","type":"bool"}],
		"payable":false,
		"stateMutability":"nonpayable",
		"type":"function"
	},
	{
		"constant":true,
		"inputs":[],
		"name":"decimals",
		"outputs":[{"name":"","type":"uint8"}],
		"payable":false,
		"stateMutability":"view",
		"type":"function"
	},
	{
		"constant":true,
		"inputs":[],
		"name":"name",
		"outputs":[{"name":"","type":"string"}],
		"payable":false,
		"stateMutability":"view",
		"type":"function"
	},
	{
		"constant":true,
		"inputs":[],
		"name":"symbol",
		"outputs":[{"name":"","type":"string"}],
		"payable":false,
		"stateMutability":"view",
		"type":"function"
	},
	{
		"constant":true,
		"inputs":[],
		"name":"totalSupply",
		"outputs":[{"name":"","type":"uint256"}],
		"payable":false,
		"stateMutability":"view",
		"type":"function"
	}
]`)
)
