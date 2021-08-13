# plt-deploy-tool

this tool used to deploy and bind contracts for polynetwork.

1. deploy eccd, eccm, ccmp contracts
```bash
make tool m=plt-deploy-eccd
make tool m=plt-deploy-eccm
make tool m=plt-recover-eccm
make tool m=plt-deploy-ccmp
```

2. transfer eccd ownership to eccm, transfer eccm ownership to ccmp.
```bash
make tool m=plt-eccd-ownership
make tool m=plt-eccm-ownership
```

3. set proxy upgrade manager contract
```bash
make tool m=plt-plt-ccmp
make tool m=plt-nft-ccmp
```

4. deploy nft-proxy
```bash
make tool m=plt-deploy-nft-proxy
```

## register palette chain and deploy contracts
1. register side chain id to poly chain and approve it with 4 poly validators' wallet file.
```bash
make tool m=plt-registerSideChain
make tool m=plt-approveRegisterSideChain
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
