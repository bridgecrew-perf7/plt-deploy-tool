package core

import (
	"github.com/palettechain/deploy-tool/pkg/frame"
)

func Endpoint() {
	// palette side chain register and init
	frame.Tool.RegMethod("plt-register-sidechain", PLTRegisterSideChain)
	frame.Tool.RegMethod("plt-approve-sidechain", PLTApproveRegisterSideChain)
	frame.Tool.RegMethod("plt-sync-plt-genesis", PLTSyncPLTGenesis)
	frame.Tool.RegMethod("plt-sync-poly-genesis", PLTSyncPolyGenesis)

	// palette contract binding relationship
	frame.Tool.RegMethod("plt-deploy-eccd", PLTDeployECCD)
	frame.Tool.RegMethod("plt-deploy-eccm", PLTDeployECCM)
	frame.Tool.RegMethod("plt-recover-eccm", PLTRecoverBookeeper)
	frame.Tool.RegMethod("plt-deploy-ccmp", PLTDeployCCMP)
	frame.Tool.RegMethod("plt-deploy-wrap", PLTDeployWrap)
	frame.Tool.RegMethod("plt-eccd-ownership", PLTTransferECCDOwnerShip)
	frame.Tool.RegMethod("plt-eccm-ownership", PLTTransferECCMOwnerShip)
	frame.Tool.RegMethod("plt-plt-ccmp", PLTSetCCMP)
	frame.Tool.RegMethod("plt-bind-plt-proxy", PLTBindPLTProxy)
	frame.Tool.RegMethod("plt-bind-plt-asset", PLTBindPLTAsset)
	frame.Tool.RegMethod("plt-wrap-set-proxy", PLTWrapperSetLockProxy)
	frame.Tool.RegMethod("plt-deploy-nft-proxy", PLTDeployNFTProxy)
	frame.Tool.RegMethod("plt-bind-nft-proxy", PLTBindNFTProxy)
	frame.Tool.RegMethod("plt-bind-nft-asset", PLTBindNFTAsset)
	frame.Tool.RegMethod("plt-nft-ccmp", PLTSetNFTCCMP)

	// ethereum bind proxy and asset
	frame.Tool.RegMethod("eth-bind-plt-proxy", ETHBindPLTProxy)
	frame.Tool.RegMethod("eth-bind-plt-asset", ETHBindPLTAsset)
	frame.Tool.RegMethod("eth-bind-nft-proxy", ETHBindNFTProxy)
	frame.Tool.RegMethod("eth-bind-nft-asset", ETHBindNFTAsset)
}
