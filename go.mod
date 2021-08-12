module github.com/palettechain/onRobot

go 1.14

require (
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/btcsuite/btcutil v1.0.2
	github.com/btcsuite/goleveldb v1.0.0
	github.com/ethereum/go-ethereum v1.9.15
	github.com/howeyc/gopass v0.0.0-20190910152052-7cb4b85ec19c
	github.com/ontio/ontology-crypto v1.0.9
	github.com/palettechain/palette_token v0.0.0-20210120103528-1db803afdd45
	github.com/polynetwork/eth-contracts v0.0.1
	github.com/polynetwork/poly v1.7.2-0.20210802025248-aaa66443deb5
	github.com/polynetwork/poly-go-sdk v0.0.0-20200817120957-365691ad3493
	github.com/stretchr/testify v1.7.0
	github.com/tyler-smith/go-bip39 v1.0.2
)

replace (
	github.com/ethereum/go-ethereum v1.9.15 => /Users/dylen/workspace/gohome/src/github.com/palettechain/palette
	github.com/polynetwork/eth-contracts v0.0.1 => github.com/zouxyan/eth-contracts v0.0.0-20210115072359-e4cac6edc20c
)
