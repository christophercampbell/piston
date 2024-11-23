# piston
Utility to drive Ethereum's engine api without consensus

#### build

```shell
git clone git@github.com:christophercampbell/piston.git
cd piston
make build
```

#### help

```shell
build/piston run --help
NAME:
   piston run - Run the program

USAGE:
   piston run [command options] [arguments...]

OPTIONS:
   --jwt value        JWT file for engine API auth
   --ethUrl value     URL to engine's public API (default: http://localhost:8545)
   --engineUrl value  URL to engine's secure API (default: http://localhost:8551)
   --help, -h         show help
```

----

#### zeth

Configure genesis for zeth such that it supports engine api v3. For example

```json
{
  "config": {
    "chainId": 17000,
    "homesteadBlock": 0,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "istanbulBlock": 0,
    "berlinBlock": 0,
    "londonBlock": 0,
    "mergeNetsplitBlock": 0,
    "terminalTotalDifficulty": 0,
    "terminalTotalDifficultyPassed": true,
    "shanghaiTime": 1696000704,
    "cancunTime":1707305664
  },
  "alloc": {
  },
  "coinbase": "0x0000000000000000000000000000000000000000",
  "difficulty": "0x01",
  "extraData": "",
  "gasLimit": "0x17D7840",
  "nonce": "0x1234",
  "mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "timestamp": "1695902100"
}
```

#### run zeth 
```shell
zeth node --chain ./genesis.json --http --disable-discovery --datadir ./zeth-data --zeth.db-path ./zeth-data/polygon-zero.db
```

#### run piston

```shell
build/piston run --jwt [PATH_TO_ZETH]/zeth-data/jwt.hex
```
