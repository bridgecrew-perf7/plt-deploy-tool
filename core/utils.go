package core

import (
	"github.com/palettechain/onRobot/config"
	"github.com/palettechain/onRobot/pkg/eth"
	"github.com/palettechain/onRobot/pkg/log"
	"github.com/palettechain/onRobot/pkg/sdk"
)

func getPaletteCli() (*sdk.Client, error) {
	url := config.Conf.PaletteRPCUrl
	privateKey, err := config.Conf.LoadPLTAdminAccount()
	if err != nil {
		return nil, err
	}
	return sdk.NewSender(url, privateKey), nil
}

func getEthereumCli() (*eth.EthInvoker, error) {
	url := config.Conf.EthereumRPCUrl
	privateKey, err := config.Conf.LoadETHAdminAccount()
	if err != nil {
		return nil, err
	}

	return eth.NewEInvoker(url, privateKey), nil
}

func logsplit() {
	log.Info("------------------------------------------------------------------")
}
