# Go Faucet
![Crowd](https://user-images.githubusercontent.com/104348282/182021454-d3ffebb5-a27b-4d81-930c-486557b5569c.png)
Faucet for Crowdcontrol written in go.

## Usage
```shell
go build
./go-faucet
```

## Configuration
```shell
SECRET_KEY=key  # The key used to communicate with the hCaptcha server
RPC_NODE=tcp://127.0.0.1:26657  # Blockchains rpc node
BLOCKCHAIN_USER=alice  # The users name that send messages to the chain, keys must be present in local keyring
```
