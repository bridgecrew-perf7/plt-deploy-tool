package core

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/palettechain/deploy-tool/config"
	"github.com/palettechain/deploy-tool/pkg/log"
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

	if err := config.Conf.StorePaletteECCD(eccd); err != nil {
		log.Errorf("store palette eccd err: %v", err)
		return
	}

	log.Infof("deploy eccd %s on palette success!", eccd.Hex())

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
	whiteList := []common.Address{
		common.HexToAddress(native.PLTContractAddress),
		config.Conf.PaletteNFTProxy,
	}
	keepers := config.Conf.LoadPolyCurBookeeperBytes()
	eccm, err := cli.DeployECCM(eccd, sideChainID, whiteList, keepers)
	if err != nil {
		log.Errorf("deploy eccm on palette failed, err: %s", err.Error())
		return
	}

	if err := config.Conf.StorePaletteECCM(eccm); err != nil {
		log.Errorf("store palette eccm err: %v", err)
		return
	}

	log.Infof("deploy eccm %s on palette success!", eccm.Hex())

	return true
}

func PLTRecoverBookeeper() (succeed bool) {
	cli, err := getPaletteCli()
	if err != nil {
		log.Errorf("get palette cross chain admin client failed")
		return
	}

	eccm := config.Conf.PaletteECCM
	keepers := config.Conf.LoadPolyCurBookeeperBytes()
	if _, err := cli.RecoverECCM(eccm, keepers); err != nil {
		log.Errorf("recover eccm on palette failed, err: %s", err.Error())
		return
	}

	log.Info("recover eccm bookeepers success")
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

	if err := config.Conf.StorePaletteCCMP(ccmp); err != nil {
		log.Errorf("store palette ccmp err: %v", err)
		return
	}

	log.Infof("deploy ccmp %s on palette success!", ccmp.Hex())

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
		log.Errorf("PLT set proxy ccmp failed, err: %v", err)
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

	if err := config.Conf.StorePaletteNFTProxy(proxy); err != nil {
		log.Errorf("store palette nft proxy err: %v", err)
		return
	}

	log.Infof("deploy NFT proxy %s on palette success!", proxy.Hex())

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

func PLTDeployPLTWrap() (succeed bool) {
	cli, err := getPaletteCli()
	if err != nil {
		log.Errorf("get palette cross chain admin client failed")
		return
	}

	proxy := common.HexToAddress(native.PLTContractAddress)
	chainId := new(big.Int).SetUint64(config.Conf.PaletteSideChainID)

	contractAddr, err := cli.DeployPalettePLTWrapper(cli.Address(), proxy, chainId)
	if err != nil {
		log.Errorf("deploy plt wrap on palette failed, err: %s", err.Error())
		return
	}

	if err := config.Conf.StorePalettePLTWrapper(contractAddr); err != nil {
		log.Errorf("store plt wrap failed, err: %v", err)
		return
	}

	log.Infof("deploy plt wrap %s on palette success!", contractAddr.Hex())
	return true
}

func PLTDeployNFTWrap() (succeed bool) {
	cli, err := getPaletteCli()
	if err != nil {
		log.Errorf("get palette cross chain admin client failed")
		return
	}

	chainId := new(big.Int).SetUint64(config.Conf.PaletteSideChainID)
	feeToken := common.HexToAddress(native.PLTContractAddress)
	contractAddr, err := cli.DeployPaletteNFTWrapper(cli.Address(), feeToken, chainId)
	if err != nil {
		log.Errorf("deploy nft wrap on palette failed, err: %s", err.Error())
		return
	}

	if err := config.Conf.StorePaletteNFTWrapper(contractAddr); err != nil {
		log.Errorf("store nft wrap failed, err: %v", err)
		return
	}

	log.Infof("deploy nft wrap %s on palette success!", contractAddr.Hex())
	return true
}

func PLTDeployNFTQuery() (succeed bool) {
	cli, err := getPaletteCli()
	if err != nil {
		log.Errorf("get palette cross chain admin client failed")
		return
	}

	var limit uint64 = 36
	contractAddr, err := cli.DeployPaletteNFTQuery(cli.Address(), limit)
	if err != nil {
		log.Errorf("deploy nft query on palette failed, err: %s", err.Error())
		return
	}

	if err := config.Conf.StorePaletteNFTQuery(contractAddr); err != nil {
		log.Errorf("store nft query failed, err: %v", err)
		return
	}

	log.Infof("deploy nft query %s on palette success!", contractAddr.Hex())
	return true
}

func PLTNFTWrapperSetLockProxy() (succeed bool) {
	cli, err := getPaletteCli()
	if err != nil {
		log.Errorf("get palette cross chain admin client failed")
		return
	}

	wrapAddr := config.Conf.PaletteNFTWrapper
	targetLockProxy := config.Conf.PaletteNFTProxy

	cur, _ := cli.GetPaletteNFTWrapLockProxy(wrapAddr)
	if cur == targetLockProxy {
		log.Infof("nft wrapper proxy %s already settled", targetLockProxy.Hex())
		return true
	}

	if _, err := cli.PaletteNFTWrapSetLockProxy(wrapAddr, targetLockProxy); err != nil {
		log.Errorf("nft wrapper set lock proxy failed, err: %v", err)
		return false
	}

	got, _ := cli.GetPaletteNFTWrapLockProxy(wrapAddr)
	if got != targetLockProxy {
		log.Infof("nft wrapper proxy set failed, expect %s, got %s", targetLockProxy.Hex(), got.Hex())
		return true
	}

	log.Infof("nft wrap set lock proxy %s on palette success!", targetLockProxy.Hex())
	return true
}
