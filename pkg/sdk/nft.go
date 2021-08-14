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

package sdk

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native/nft"
	"github.com/ethereum/go-ethereum/contracts/native/nftmanager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

func (c *Client) NFTDeploy(name string, symbol string) (common.Hash, common.Address, error) {
	payload, err := c.packNFTManager(nftmanager.MethodDeploy, name, symbol, c.Address())
	if err != nil {
		return utils.EmptyHash, utils.EmptyAddress, err
	}

	hash, err := c.sendNFTManager(payload)
	if err != nil {
		return utils.EmptyHash, utils.EmptyAddress, err
	}

	receipts, err := c.GetReceipt(hash)
	if err != nil {
		return utils.EmptyHash, utils.EmptyAddress, fmt.Errorf("nft depoly - get receipt %s err: %s", hash.Hex(), err)
	}
	if len(receipts.Logs) == 0 {
		return utils.EmptyHash, utils.EmptyAddress, fmt.Errorf("invalid tx %s, no receipts events", hash.Hex())
	}

	for _, event := range receipts.Logs {
		if event.Topics[0] == NFTABI.Events[nft.EventDeploy].ID() {
			return hash, event.Address, nil
		}
	}

	return utils.EmptyHash, utils.EmptyAddress, fmt.Errorf("no valid nft address")
}

func (c *Client) NFTName(asset common.Address, blockNum string) (string, error) {
	payload, err := c.packNFT(nft.MethodName)
	if err != nil {
		return "", err
	}
	data, err := c.callNFT(asset, payload, blockNum)
	if err != nil {
		return "", err
	}
	result := &nft.NameResult{}
	err = c.unpackNFT(nft.MethodName, result, data)
	if err != nil {
		return "", err
	}

	return result.Name, nil
}

func (c *Client) NFTSymbol(asset common.Address, blockNum string) (string, error) {
	payload, err := c.packNFT(nft.MethodSymbol)
	if err != nil {
		return "", err
	}
	data, err := c.callNFT(asset, payload, blockNum)
	if err != nil {
		return "", err
	}
	result := &nft.SymbolResult{}
	err = c.unpackNFT(nft.MethodSymbol, result, data)
	if err != nil {
		return "", err
	}

	return result.Symbol, nil
}

func (c *Client) NFTAssetOwner(asset common.Address, blockNum string) (common.Address, error) {
	payload, err := c.packNFT(nft.MethodOwner)
	if err != nil {
		return utils.EmptyAddress, err
	}
	data, err := c.callNFT(asset, payload, blockNum)
	if err != nil {
		return utils.EmptyAddress, err
	}
	result := &nft.OwnerResult{}
	if err = c.unpackNFT(nft.MethodOwner, result, data); err != nil {
		return utils.EmptyAddress, err
	}

	return result.Owner, nil
}

// NFT
func (c *Client) packNFT(method string, args ...interface{}) ([]byte, error) {
	return utils.PackMethod(NFTABI, method, args...)
}
func (c *Client) unpackNFT(method string, output interface{}, enc []byte) error {
	return utils.UnpackOutputs(NFTABI, method, output, enc)
}
func (c *Client) sendNFT(nftAddr common.Address, payload []byte) (common.Hash, error) {
	hash, err := c.SendTransaction(nftAddr, payload)
	if err != nil {
		return utils.EmptyHash, err
	}
	if err := c.WaitTransaction(hash); err != nil {
		return utils.EmptyHash, err
	}
	return hash, nil
}
func (c *Client) callNFT(nftAddr common.Address, payload []byte, blockNum string) ([]byte, error) {
	return c.CallContract(c.Address(), nftAddr, payload, blockNum)
}

// nft manager
func (c *Client) packNFTManager(method string, args ...interface{}) ([]byte, error) {
	return utils.PackMethod(NFTManagerABI, method, args...)
}
func (c *Client) unpackNFTManager(method string, output interface{}, enc []byte) error {
	return utils.UnpackOutputs(NFTManagerABI, method, output, enc)
}
func (c *Client) sendNFTManager(payload []byte) (common.Hash, error) {
	hash, err := c.SendTransaction(NFTMangerAddress, payload)
	if err != nil {
		return utils.EmptyHash, err
	}
	if err := c.WaitTransaction(hash); err != nil {
		return utils.EmptyHash, err
	}
	return hash, nil
}
func (c *Client) callNFTManager(payload []byte, blockNum string) ([]byte, error) {
	return c.CallContract(c.Address(), NFTMangerAddress, payload, blockNum)
}
