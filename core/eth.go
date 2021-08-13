/*
 * Copyright (C) 2021 The Zion Authors
 * This file is part of The Zion library.
 *
 * The Zion is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The Zion is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The Zion.  If not, see <http://www.gnu.org/licenses/>.
 */

package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	"github.com/palettechain/deploy-tool/config"
	"github.com/palettechain/deploy-tool/pkg/log"
)

func ETHBindPLTProxy() (succeed bool) {
	cli, err := getEthereumCli()
	if err != nil {
		log.Errorf("get eth cross chain admin failed, err: %v", err)
		return
	}
	localLockProxy := config.Conf.EthereumPLTProxy
	targetLockProxy := common.HexToAddress(native.PLTContractAddress)
	targetSideChainID := config.Conf.PaletteSideChainID

	cur, _ := cli.GetBoundPLTProxy(localLockProxy, targetSideChainID)
	if cur == targetLockProxy {
		log.Infof("PLT proxy %s already bound to %s", localLockProxy.Hex(), targetLockProxy.Hex())
		return true
	}

	hash, err := cli.BindPLTProxy(localLockProxy, targetLockProxy, targetSideChainID)
	if err != nil {
		log.Errorf("bind PLT proxy on ethereum failed, err: %s", err.Error())
		return
	}

	actual, err := cli.GetBoundPLTProxy(localLockProxy, targetSideChainID)
	if err != nil {
		log.Error(err)
		return
	}
	if actual != targetLockProxy {
		log.Errorf("proxy bind failed, expect %s, got %s", targetLockProxy.Hex(), actual.Hex())
		return
	}

	log.Infof("bind PLT proxy %s to %s on ethereum success, hash %s", localLockProxy.Hex(), targetLockProxy.Hex(), hash.Hex())
	return true
}

func ETHBindPLTAsset() (succeed bool) {
	cli, err := getEthereumCli()
	if err != nil {
		log.Errorf("get eth cross chain admin failed, err: %v", err)
		return
	}

	localLockProxy := config.Conf.EthereumPLTProxy
	fromAsset := config.Conf.EthereumPLTAsset
	toAsset := common.HexToAddress(native.PLTContractAddress)
	toChainId := config.Conf.PaletteSideChainID

	cur, _ := cli.GetBoundPLTAsset(localLockProxy, fromAsset, toChainId)
	if cur == toAsset {
		log.Infof("PLT asset %s already bound to %s", fromAsset.Hex(), toAsset.Hex())
		return true
	}

	hash, err := cli.BindPLTAsset(localLockProxy, fromAsset, toAsset, toChainId)
	if err != nil {
		log.Errorf("bind PLT asset on ethereum failed, err: %s", err.Error())
		return
	}

	actual, err := cli.GetBoundPLTAsset(localLockProxy, fromAsset, toChainId)
	if err != nil {
		log.Error(err)
		return
	}
	if actual != toAsset {
		log.Errorf("bind plt asset on ethereum failed, expect %s, got %s", toAsset.Hex(), actual.Hex())
		return
	}

	log.Infof("bind PLT asset %s to %s on ethereum success, hash %s", fromAsset.Hex(), toAsset.Hex(), hash.Hex())
	return true
}

func ETHBindNFTProxy() (succeed bool) {
	cli, err := getEthereumCli()
	if err != nil {
		log.Errorf("get eth cross chain admin failed, err: %v", err)
		return
	}

	localLockProxy := config.Conf.EthereumNFTProxy
	targetLockProxy := config.Conf.PaletteNFTProxy
	targetSideChainID := config.Conf.PaletteSideChainID

	cur, _ := cli.GetBoundNFTProxy(localLockProxy, targetSideChainID)
	if cur == targetLockProxy {
		log.Infof("NFT proxy %s already bound to %s", localLockProxy.Hex(), targetLockProxy.Hex())
		return true
	}

	hash, err := cli.BindNFTProxy(localLockProxy, targetLockProxy, targetSideChainID)
	if err != nil {
		log.Errorf("bind NFT proxy on ethereum failed, err: %s", err.Error())
		return
	}

	actual, err := cli.GetBoundNFTProxy(localLockProxy, targetSideChainID)
	if err != nil {
		log.Error(err)
		return
	}
	if actual != targetLockProxy {
		log.Errorf("bind NFT proxy to ccmp failed, expect %s, got %s", targetLockProxy.Hex(), actual.Hex())
		return
	}

	log.Infof("bind NFT proxy %s to %s on ethereum success, tx %s", localLockProxy.Hex(), targetLockProxy.Hex(), hash.Hex())
	return true
}

func ETHBindNFTAsset() (succeed bool) {
	cli, err := getEthereumCli()
	if err != nil {
		log.Errorf("get eth cross chain admin failed, err: %v", err)
		return
	}

	proxy := config.Conf.EthereumNFTProxy
	fromAsset := config.Conf.EthereumNFTAsset
	toAsset := config.Conf.PaletteNFTAsset
	chainID := config.Conf.PaletteSideChainID

	cur, _ := cli.GetBoundNFTAsset(proxy, fromAsset, chainID)
	if cur == toAsset {
		log.Infof("NFT asset %s already bound to %s", fromAsset.Hex(), toAsset.Hex())
		return true
	}

	hash, err := cli.BindNFTAsset(
		proxy,
		fromAsset,
		toAsset,
		chainID,
	)
	if err != nil {
		log.Errorf("bind NFT asset on ethereum failed, err: %s", err.Error())
		return
	}

	actual, err := cli.GetBoundNFTAsset(proxy, fromAsset, chainID)
	if err != nil {
		log.Error(err)
		return
	}
	if actual != toAsset {
		log.Errorf("bind NFT asset failed, expect %s, got %s", toAsset.Hex(), actual.Hex())
		return
	}

	log.Infof("bind NFT asset %s to %s on ethereum success, hash %s", fromAsset.Hex(), toAsset.Hex(), hash.Hex())
	return true
}
