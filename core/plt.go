package core

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/palettechain/onRobot/config"
	"github.com/palettechain/onRobot/pkg/log"
	"github.com/palettechain/onRobot/pkg/poly"
)

// 在palette合约部署成功后由三本合约:
// eccd: 管理epoch
// eccm: 管理跨链转账
// ccmp: 记录eccm地址及升级等
// 加入跨链事件从poly回到palette，事件流如下:
// relayer:
// 1. 执行palette eccm合约的verifyProofAndExecuteTx，这个方法会进入到palette native PLT合约的unlock方法
// 2. palette native PLT unlock 取出ccmp地址，并进入该合约查询eccm地址，比较从relayer过来的eccm地址与该地址是否匹配
// 3. 进入unlock资金逻辑

func PLTDeployECCD() (succeed bool) {
	cli, err := getPaletteCli()
	if err != nil {
		log.Errorf("get palette cross chain admin client failed")
		return
	}

	eccd, err := cli.DeployECCD()
	if err != nil {
		log.Errorf("deploy eccd on palette failed, err: %s", err.Error())
		return
	}

	log.Infof("deploy eccd %s on palette success!", eccd.Hex())

	if err := config.Conf.StorePaletteECCD(eccd); err != nil {
		log.Error("store palette eccd failed")
		return
	}

	return true
}

func PLTDeployECCM() (succeed bool) {
	cli, err := getPaletteCli()
	if err != nil {
		log.Errorf("get palette cross chain admin client failed")
		return
	}

	eccd := config.Conf.PaletteECCD
	sideChainID := config.Conf.PaletteSideChainID
	eccm, err := cli.DeployECCM(eccd, sideChainID)
	if err != nil {
		log.Errorf("deploy eccm on palette failed, err: %s", err.Error())
		return
	}

	log.Infof("deploy eccm %s on palette success!", eccm.Hex())

	if err := config.Conf.StorePaletteECCM(eccm); err != nil {
		log.Error("store palette eccm failed")
		return
	}

	return true
}

func PLTDeployCCMP() (succeed bool) {
	cli, err := getPaletteCli()
	if err != nil {
		log.Errorf("get palette cross chain admin client failed")
		return
	}

	eccm := config.Conf.PaletteECCM
	ccmp, err := cli.DeployCCMP(eccm)
	if err != nil {
		log.Errorf("deploy ccmp on palette failed, err: %s", err.Error())
		return
	}

	log.Infof("deploy ccmp %s on palette success!", ccmp.Hex())

	if err := config.Conf.StorePaletteCCMP(ccmp); err != nil {
		log.Error("store palette ccmp failed")
		return
	}

	return true
}

func PLTTransferECCDOwnerShip() (succeed bool) {
	cli, err := getPaletteCli()
	if err != nil {
		log.Errorf("get palette cross chain admin client failed")
		return
	}

	eccd := config.Conf.PaletteECCD
	eccm := config.Conf.PaletteECCM

	cur, _ := cli.ECCDOwnership(eccd)
	if bytes.Equal(eccm.Bytes(), cur.Bytes()) {
		log.Infof("eccd %s owner is %s already", eccd.Hex(), eccm.Hex())
		return true
	}

	hash, err := cli.ECCDTransferOwnerShip(eccd, eccm)
	if err != nil {
		log.Error(err)
		return
	}
	actual, err := cli.ECCDOwnership(eccd)
	if err != nil {
		log.Error(err)
		return
	}
	if !bytes.Equal(eccm.Bytes(), actual.Bytes()) {
		log.Error("new owner %s != acutal %s", eccm.Hex(), actual.Hex())
		return
	}
	log.Infof("transfer eccd %s to eccm %s success! hash %s", eccd.Hex(), eccm.Hex(), hash.Hex())

	return true
}

func PLTTransferECCMOwnerShip() (succeed bool) {
	cli, err := getPaletteCli()
	if err != nil {
		log.Errorf("get palette cross chain admin client failed")
		return
	}

	eccm := config.Conf.PaletteECCM
	ccmp := config.Conf.PaletteCCMP

	cur, _ := cli.ECCMOwnership(eccm)
	if bytes.Equal(ccmp.Bytes(), cur.Bytes()) {
		log.Infof("eccm %s owner is %s already", eccm.Hex(), ccmp.Hex())
		return true
	}

	hash, err := cli.ECCMTransferOwnerShip(eccm, ccmp)
	if err != nil {
		log.Error(err)
		return
	}
	actual, err := cli.ECCMOwnership(eccm)
	if err != nil {
		log.Error(err)
		return
	}
	if !bytes.Equal(ccmp.Bytes(), actual.Bytes()) {
		log.Error("new owner %s != acutal %s", ccmp.Hex(), actual.Hex())
		return
	}
	log.Infof("transfer eccm %s to ccmp %s success! hash %s", eccm.Hex(), ccmp.Hex(), hash.Hex())

	return true
}

func PLTSetCCMP() (succeed bool) {
	cli, err := getPaletteCli()
	if err != nil {
		log.Errorf("get palette cross chain admin client failed")
		return
	}

	ccmp := config.Conf.PaletteCCMP

	cur, _ := cli.GetPLTCCMP("latest")
	if cur == ccmp {
		log.Infof("PLT proxy already managed by %s", ccmp.Hex())
		return true
	}

	hash, err := cli.SetPLTCCMP(ccmp)
	if err != nil {
		log.Error(err)
		return
	}

	actual, err := cli.GetPLTCCMP("latest")
	if err != nil {
		log.Error(err)
		return
	}
	if actual != ccmp {
		log.Errorf("set proxy manager failed, expect %s != actual %s", ccmp.Hex(), actual.Hex())
		return
	}

	log.Infof("set PLT ccmp success! hash %s", hash.Hex())
	return true
}

// 在palette native合约上记录以太坊localProxy地址,
// 这里我们将实现palette->poly->palette的循环，不走ethereum，那么proxy就直接是plt地址，
// asset的地址也是palette plt地址
func PLTBindPLTProxy() (succeed bool) {
	cli, err := getPaletteCli()
	if err != nil {
		log.Errorf("get palette cross chain admin client failed")
		return
	}

	proxy := config.Conf.EthereumPLTProxy
	sideChainID := config.Conf.EthereumSideChainID

	cur, _ := cli.GetBindPLTProxy(sideChainID, "latest")
	if cur == proxy {
		log.Infof("PLT proxy already bound to by %s", proxy.Hex())
		return true
	}

	hash, err := cli.BindPLTProxy(sideChainID, proxy)
	if err != nil {
		log.Error(err)
		return
	}

	actual, err := cli.GetBindPLTProxy(sideChainID, "latest")
	if err != nil {
		log.Error(err)
		return
	}
	if actual != proxy {
		log.Errorf("bind PLT proxy failed, expect  %s != actual %s", proxy.Hex(), actual.Hex())
		return
	}

	log.Infof("bind PLT proxy to %s on palette success! hash %s", proxy.Hex(), hash.Hex())
	return true
}

// 在palette native合约上记录以太坊erc20资产地址
func PLTBindPLTAsset() (succeed bool) {
	cli, err := getPaletteCli()
	if err != nil {
		log.Errorf("get palette cross chain admin client failed")
		return
	}

	asset := config.Conf.EthereumPLTAsset
	sideChainID := config.Conf.EthereumSideChainID

	cur, _ := cli.GetBindPLTAsset(sideChainID, "latest")
	if cur == asset {
		log.Infof("PLT asset already bound to by %s", asset.Hex())
		return true
	}

	hash, err := cli.BindPLTAsset(sideChainID, asset)
	if err != nil {
		log.Error(err)
		return
	}

	actual, err := cli.GetBindPLTAsset(sideChainID, "latest")
	if err != nil {
		log.Error(err)
		return
	}
	if actual != asset {
		log.Errorf("bind PLT asset err, expect %s != actual %s", asset.Hex(), actual.Hex())
		return
	}

	log.Infof("bind PLT asset to %s on palette success! hash %s", asset.Hex(), hash.Hex())
	return true
}

func PLTDeployNFTProxy() (succeed bool) {
	cli, err := getPaletteCli()
	if err != nil {
		log.Errorf("get palette cross chain admin client failed")
		return
	}

	proxy, err := cli.DeployNFTProxy()
	if err != nil {
		log.Errorf("deploy NFT proxy on palette failed, err: %s", err.Error())
		return
	}

	log.Infof("deploy NFT proxy %s on palette success!", proxy.Hex())

	if err := config.Conf.StorePaletteNFTProxy(proxy); err != nil {
		log.Error("store palette nft proxy failed")
		return
	}

	return true
}

func PLTBindNFTProxy() (succeed bool) {
	cli, err := getPaletteCli()
	if err != nil {
		log.Errorf("get palette cross chain admin client failed")
		return
	}

	localLockproxy := config.Conf.PaletteNFTProxy
	targetLockProxy := config.Conf.EthereumNFTProxy
	targetSideChainID := config.Conf.EthereumSideChainID

	cur, _ := cli.GetBoundNFTProxy(localLockproxy, targetSideChainID)
	if cur == targetLockProxy {
		log.Infof("NFT proxy %s already bound to by %s", localLockproxy.Hex(), targetLockProxy.Hex())
		return true
	}

	hash, err := cli.BindNFTProxy(localLockproxy, targetLockProxy, targetSideChainID)
	if err != nil {
		log.Errorf("bind NFT proxy on palette failed, err: %s", err.Error())
		return
	}

	actual, err := cli.GetBoundNFTProxy(localLockproxy, targetSideChainID)
	if err != nil {
		log.Error(err)
		return
	}
	if !bytes.Equal(targetLockProxy.Bytes(), actual.Bytes()) {
		log.Errorf("asset err, expect %s != actual %s", targetLockProxy.Hex(), actual.Hex())
		return
	}

	log.Infof("bind NFT proxy %s to %s on palette success! hash %s", localLockproxy.Hex(), targetLockProxy.Hex(), hash.Hex())
	return true
}

func PLTSetNFTCCMP() (succeed bool) {
	cli, err := getPaletteCli()
	if err != nil {
		log.Errorf("get palette cross chain admin client failed")
		return
	}

	proxy := config.Conf.PaletteNFTProxy
	ccmp := config.Conf.PaletteCCMP

	cur, _ := cli.GetNFTCCMP(proxy)
	if bytes.Equal(ccmp.Bytes(), cur.Bytes()) {
		log.Infof("NFT proxy %s already managed by %s", proxy.Hex(), ccmp.Hex())
		return true
	}

	hash, err := cli.SetNFTCCMP(proxy, ccmp)
	if err != nil {
		log.Errorf("set ccmp on palette failed, err: %s", err.Error())
		return
	}

	actual, err := cli.GetNFTCCMP(proxy)
	if err != nil {
		log.Error(err)
		return
	}
	if !bytes.Equal(ccmp.Bytes(), actual.Bytes()) {
		log.Errorf("asset err, expect %s, actual %s", ccmp.Hex(), actual.Hex())
		return
	}
	log.Infof("set NFT proxy manager %s for nft proxy %s on palette success! hash %s", actual.Hex(), proxy.Hex(), hash.Hex())
	return true
}

// 同步palette区块头到poly链上
// 1. 环境准备，palette cli: 使用任意palette签名者对应的cli, poly cli: 必须是poly验证节点的validators作为多签地址
// 2. 获取palette当前块高的区块头, 并使用json序列化为bytes
// 3. 使用poly cli同步第二步的bytes以及palette network id到poly native管理合约,
//	  这笔交易发出后等待poly当前块高超过交易块高, 作为落账的判断条件
// 4. 获取poly当前块高作为写入palette管理合约的genesis块高，获取对应的block，将block header及block book keeper
//    序列化，提交到palette管理合约
func PLTSyncPLTGenesis() (succeed bool) {
	// 1. prepare
	polyRPC := config.Conf.PolyRPCUrl
	polyValidators := config.Conf.LoadPolyAccountList()
	polyCli, err := poly.NewPolyClient(polyRPC, polyValidators)
	if err != nil {
		log.Errorf("failed to generate poly client, err: %s", err)
		return
	} else {
		log.Infof("generate poly client success!")
	}

	// 2. get palette current block header
	logsplit()
	cli, err := getPaletteCli()
	if err != nil {
		log.Errorf("get palette cross chain admin client failed")
		return
	}
	curr, hdr, err := cli.GetCurrentBlockHeader()
	if err != nil {
		log.Errorf("failed to get block header, err: %s", err)
		return
	}
	pltHeaderEnc, err := hdr.MarshalJSON()
	if err != nil {
		log.Errorf("marshal header failed, err: %s", err)
		return
	}
	log.Infof("get palette block header with current height %d, header %s", curr, hexutil.Encode(pltHeaderEnc))

	logsplit()
	crossChainID := config.Conf.PaletteSideChainID
	if err := polyCli.SyncGenesisBlock(crossChainID, pltHeaderEnc); err != nil {
		log.Errorf("SyncEthGenesisHeader failed: %v", err)
		return
	}
	log.Infof("sync palette genesis header to poly success, txhash %s, block number %d",
		hdr.Hash().Hex(), hdr.Number.Uint64())

	return true
}

// 同步poly区块头到palette
func PLTSyncPolyGenesis() (succeed bool) {
	polyRPC := config.Conf.PolyRPCUrl
	polyCli, err := poly.NewPolyClient(polyRPC, nil)
	if err != nil {
		log.Errorf("failed to generate poly client, err: %s", err)
		return
	} else {
		log.Infof("generate poly client success!")
	}

	// `epoch` related with the poly validators changing,
	// we can set it as 0 if poly validators never changed on develop environment.
	var hasValidatorsBlockNumber uint32 = 0
	gB, err := polyCli.GetBlockByHeight(hasValidatorsBlockNumber)
	if err != nil {
		log.Errorf("failed to get block, err: %s", err)
		return
	}
	bookeepers, err := poly.GetBookeeper(gB)
	if err != nil {
		log.Errorf("failed to get bookeepers, err: %s", err)
		return
	}
	bookeepersEnc := poly.AssembleNoCompressBookeeper(bookeepers)
	headerEnc := gB.Header.ToArray()

	cli, err := getPaletteCli()
	if err != nil {
		log.Errorf("get palette cross chain admin client failed")
		return
	}
	eccm := config.Conf.PaletteECCM
	txhash, err := cli.InitGenesisBlock(eccm, headerEnc, bookeepersEnc)
	if err != nil {
		log.Errorf("failed to initGenesisBlock, err: %s", err)
		return
	}

	log.Infof("sync poly genesis header to palette success, txhash %s, block number %d",
		txhash.Hex(), gB.Header.Height)

	return true
}

func PLTBindNFTAsset() (succeed bool) {
	cli, err := getPaletteCli()
	if err != nil {
		log.Errorf("get palette cross chain admin client failed")
		return
	}

	proxy := config.Conf.PaletteNFTProxy
	fromAsset := config.Conf.PaletteNFTAsset
	toAsset := config.Conf.EthereumNFTAsset
	targetSideChainID := config.Conf.EthereumSideChainID

	curAddr, _ := cli.GetBoundNFTAsset(proxy, fromAsset, targetSideChainID)
	if curAddr != utils.EmptyAddress {
		if curAddr == toAsset {
			log.Infof("ethereum NFT asset %s bound already", toAsset.Hex())
			return
		} else {
			log.Infof("ethereum NFT asset %s bound != asset %s", curAddr.Hex(), toAsset.Hex())
		}
	}

	hash, err := cli.BindNFTAsset(
		proxy,
		fromAsset,
		toAsset,
		targetSideChainID,
	)
	if err != nil {
		log.Errorf("bind NFT proxy on palette failed, err: %s", err.Error())
		return
	}

	actual, err := cli.GetBoundNFTAsset(proxy, fromAsset, targetSideChainID)
	if err != nil {
		log.Error(err)
		return
	}
	if actual != toAsset {
		log.Errorf("asset err, expect %s, actual %s", toAsset.Hex(), actual.Hex())
		return
	}

	log.Infof("bind NFT asset %s to %s on palette success, hash %s", fromAsset.Hex(), toAsset.Hex(), hash.Hex())
	return true
}

func PLTDeployWrap() (succeed bool) {
	cli, err := getPaletteCli()
	if err != nil {
		log.Errorf("get palette cross chain admin client failed")
		return
	}

	feeToken := common.HexToAddress(native.PLTContractAddress)
	chainId := new(big.Int).SetUint64(config.Conf.PaletteSideChainID)

	contractAddr, err := cli.DeployPaletteWrapper(cli.Address(), feeToken, chainId)
	if err != nil {
		log.Errorf("deploy wrap on palette failed, err: %s", err.Error())
		return
	}

	if err := config.Conf.StorePaletteWrapper(contractAddr); err != nil {
		log.Error("store palette wrapper failed")
		return
	}

	log.Infof("deploy wrap %s on palette success!", contractAddr.Hex())
	return true
}

func PLTWrapperSetLockProxy() (succeed bool) {
	cli, err := getPaletteCli()
	if err != nil {
		log.Errorf("get palette cross chain admin client failed")
		return
	}

	wrapAddr := config.Conf.PaletteWrapper
	targetLockProxy := common.HexToAddress(native.PLTContractAddress)

	cur, _ := cli.GetPaletteWrapLockProxy(wrapAddr)
	if bytes.Equal(cur.Bytes(), targetLockProxy.Bytes()) {
		log.Infof("wrapper proxy %s already settled", targetLockProxy.Hex())
		return true
	}

	if _, err := cli.PaletteWrapSetLockProxy(wrapAddr, targetLockProxy); err != nil {
		log.Errorf("wrapper set lock proxy failed, err: %v", err)
		return false
	}

	got, _ := cli.GetPaletteWrapLockProxy(wrapAddr)
	if bytes.Equal(cur.Bytes(), targetLockProxy.Bytes()) {
		log.Infof("wrapper proxy set failed, expect %s, got %s", targetLockProxy.Hex(), got.Hex())
		return true
	}

	log.Infof("wrap set lock proxy %s on palette success!", targetLockProxy.Hex())
	return true
}
