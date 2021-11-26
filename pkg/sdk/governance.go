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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native/governance"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"math/big"
)

func (c *Client) GetDelegateFactor(validator common.Address, blockNum string) (*big.Int, error) {
	payload, err := c.packGovernance(governance.MethodGetDelegateRewardFactor, validator)
	if err != nil {
		return nil, err
	}
	data, err := c.CallGovernance(payload, blockNum)
	if err != nil {
		return nil, err
	}
	output := new(governance.MethodGetDelegateRewardFactorOutput)
	if err := utils.UnpackOutputs(GovernanceABI, governance.MethodGetDelegateRewardFactor, output, data); err != nil {
		return nil, err
	}
	return output.Factor, nil
}

func (c *Client) packGovernance(method string, args ...interface{}) ([]byte, error) {
	return utils.PackMethod(GovernanceABI, method, args...)
}
func (c *Client) unpackGovernance(method string, output interface{}, enc []byte) error {
	return utils.UnpackOutputs(GovernanceABI, method, output, enc)
}
func (c *Client) SendGovernanceTx(payload []byte) (common.Hash, error) {
	hash, err := c.SendTransaction(GovernanceAddress, payload)
	if err != nil {
		return utils.EmptyHash, err
	}
	if err := c.WaitTransaction(hash); err != nil {
		return utils.EmptyHash, err
	}
	return hash, nil
}
func (c *Client) CallGovernance(payload []byte, blockNum string) ([]byte, error) {
	return c.CallContract(c.Address(), GovernanceAddress, payload, blockNum)
}
