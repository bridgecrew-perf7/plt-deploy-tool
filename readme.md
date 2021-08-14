# plt-deploy-tool

this tool used to deploy and bind contracts for polynetwork.

1. deploy nft-proxy, prepare for eccm white list
```bash
make tool m=plt-deploy-nft-proxy
```

2. deploy eccd, eccm, ccmp contracts
```bash
make tool m=plt-deploy-eccd
make tool m=plt-deploy-eccm
make tool m=plt-recover-eccm  // failed
make tool m=plt-deploy-ccmp
```

3. transfer eccd ownership to eccm, transfer eccm ownership to ccmp.
```bash
make tool m=plt-eccd-ownership
make tool m=plt-eccm-ownership
```

4. set proxy upgrade manager contract
```bash
make tool m=plt-plt-ccmp
make tool m=plt-nft-ccmp
```

## register palette chain and deploy contracts
1. register side chain id to poly chain and approve it with 4 poly validators' wallet file.
```bash
make tool m=plt-register-sidechain
make tool m=plt-approve-sidechain
```

2. sync palette header to palette chain and store poly book keepers in the palette chain
```bash
make tool m=plt-sync-plt-genesis
make tool m=plt-sync-poly-genesis
```

## bind proxies and PLT asset on ethereum chain
```bash
make tool m=eth-bind-plt-proxy
make tool m=eth-bind-nft-proxy
make tool m=eth-bind-plt-asset
```

## bind proxies and PLT asset on palette chain
```bash
make tool m=plt-bind-plt-proxy
make tool m=plt-bind-nft-proxy
make tool m=plt-bind-plt-asset
```

## bind nft asset
```bash
make tool m=plt-bind-nft-asset
make tool m=eth-bind-nft-asset
```

## deploy wrap and set lock proxy
```bash
make tool m=plt-deploy-plt-wrap  
// proxy is plt asset 0x000000000000000000000000000103

make tool m=plt-deploy-nft-wrap
make tool m=plt-set-nft-wrap-proxy
```