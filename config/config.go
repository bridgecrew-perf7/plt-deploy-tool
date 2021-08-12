package config

import (
	"crypto/ecdsa"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/howeyc/gopass"
	"github.com/palettechain/onRobot/pkg/dao"
	"github.com/palettechain/onRobot/pkg/files"
	"github.com/palettechain/onRobot/pkg/log"
	polysdk "github.com/polynetwork/poly-go-sdk"
)

type pwdSessionType byte

const (
	pwdSessionUnknown pwdSessionType = iota
	pwdSessionETH
	pwdSessionPLT
	pwdSessionPoly
)

var (
	Conf           = new(Config)
	ConfigFilePath string
)

type Config struct {
	PolyRPCUrl     string
	PolyAccountDir string

	EthereumRPCUrl          string
	EthereumCrossChainAdmin string

	PaletteRPCUrl          string
	PaletteCrossChainAdmin string

	// palette side chain
	PaletteSideChainID   uint64
	PaletteSideChainName string
	PaletteECCD          common.Address
	PaletteECCM          common.Address
	PaletteCCMP          common.Address
	PaletteNFTProxy      common.Address
	PaletteWrapper       common.Address

	// ethereum side chain configuration
	EthereumSideChainID   uint64
	EthereumSideChainName string
	EthereumECCD          common.Address
	EthereumECCM          common.Address
	EthereumCCMP          common.Address
	EthereumPLTAsset      common.Address
	EthereumPLTProxy      common.Address
	EthereumNFTProxy      common.Address

	// bind nft asset
	PaletteNFTAsset  common.Address
	EthereumNFTAsset common.Address
}

func (c *Config) DeepCopy() *Config {
	cp := new(Config)
	enc, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(enc, cp); err != nil {
		panic(err)
	}
	return cp
}

func Init(filepath string) {
	ConfigFilePath = filepath
	err := LoadConfig(ConfigFilePath, Conf)
	if err != nil {
		panic(err)
	}

	// init leveldb
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dao.NewDao(dir)
}

func LoadConfig(filepath string, ins interface{}) error {
	data, err := files.ReadFile(filepath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, ins)
	if err != nil {
		return fmt.Errorf("json.Unmarshal TestConfig:%s error:%s", data, err)
	}
	return nil
}

func (c *Config) LoadPLTAdminAccount() (*ecdsa.PrivateKey, error) {
	return getEthAccount(Conf.PaletteCrossChainAdmin, pwdSessionPLT)
}

func (c *Config) LoadETHAdminAccount() (*ecdsa.PrivateKey, error) {
	return getEthAccount(Conf.EthereumCrossChainAdmin, pwdSessionETH)
}

func (c *Config) LoadPolyAccountList() []*polysdk.Account {
	list := make([]*polysdk.Account, 0)

	dir := c.PolyAccountDir
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	fmt.Println("fs length ", len(fs))
	for _, f := range fs {
		fullPath := path.Join(dir, f.Name())
		acc, err := c.LoadPolyAccount(fullPath)
		if err != nil {
			panic(err)
		}
		list = append(list, acc)
	}

	return list
}

func (c *Config) LoadPolyAccount(path string) (*polysdk.Account, error) {
	polySDK := polysdk.NewPolySdk()

	acc, err := getPolyAccountByPassword(polySDK, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get poly account, err: %s", err)
	}
	return acc, nil
}

func (c *Config) StorePaletteECCD(addr common.Address) error {
	c.PaletteECCD = addr
	return SaveConfig(Conf)
}

func (c *Config) StorePaletteECCM(addr common.Address) error {
	c.PaletteECCM = addr
	return SaveConfig(Conf)
}

func (c *Config) StorePaletteCCMP(addr common.Address) error {
	c.PaletteCCMP = addr
	return SaveConfig(Conf)
}

func (c *Config) StorePaletteNFTProxy(addr common.Address) error {
	c.PaletteNFTProxy = addr
	return SaveConfig(Conf)
}

func (c *Config) StorePaletteWrapper(addr common.Address) error {
	c.PaletteWrapper = addr
	return SaveConfig(Conf)
}

func (c *Config) StoreEthereumECCD(addr common.Address) error {
	c.EthereumECCD = addr
	return SaveConfig(Conf)
}

func (c *Config) StoreEthereumECCM(addr common.Address) error {
	c.EthereumECCM = addr
	return SaveConfig(Conf)
}

func (c *Config) StoreEthereumCCMP(addr common.Address) error {
	c.EthereumCCMP = addr
	return SaveConfig(Conf)
}

func (c *Config) StoreEthereumNFTProxy(addr common.Address) error {
	c.EthereumNFTProxy = addr
	return SaveConfig(Conf)
}

func (c *Config) StoreEthereumPLTAsset(addr common.Address) error {
	c.EthereumPLTAsset = addr
	return SaveConfig(Conf)
}

func (c *Config) StoreEthereumPLTProxy(addr common.Address) error {
	c.EthereumPLTProxy = addr
	return SaveConfig(Conf)
}

func getPolyAccountByPassword(sdk *polysdk.PolySdk, path string) (
	*polysdk.Account, error) {
	wallet, err := sdk.OpenWallet(path)
	if err != nil {
		return nil, fmt.Errorf("open wallet error: %v", err)
	}

	return repeatPolyDecrypt(wallet, path)
}

func SaveConfig(c *Config) error {
	enc, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(ConfigFilePath, enc, os.ModePerm)
}

func getEthAccount(path string, typ pwdSessionType) (*ecdsa.PrivateKey, error) {
	enc, err := readWalletFile(path)
	if err != nil {
		return nil, err
	}
	if len(enc) <= 64 {
		bz, err := hex.DecodeString(string(enc))
		if err != nil {
			return nil, err
		}
		return crypto.ToECDSA(bz)
	}

	key, err := repeatEthDecrypt(enc, path, "", typ)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt keyjson: [%v]", err)
	}

	return key.PrivateKey, nil
}

const MaxPwdInputRetry int = 20

func repeatPolyDecrypt(wallet *polysdk.Wallet, path string) (acc *polysdk.Account, err error) {
	var (
		existPwd string
		curPwd   []byte
		typ      = pwdSessionPoly
	)

	if existPwd, err = getPwdSession(path, typ); err == nil {
		acc, err = wallet.GetDefaultAccount([]byte(existPwd))
		return
	}

	log.Infof("please input password for poly account %s", path)

	for i := 0; i < MaxPwdInputRetry; i++ {
		if curPwd, err = gopass.GetPasswd(); err != nil {
			log.Infof("input error, try it again......")
			continue
		}
		if acc, err = wallet.GetDefaultAccount(curPwd); err == nil {
			return
		} else {
			log.Infof("password invalid, err %s, try it again......", err.Error())
		}
	}
	return
}

func repeatEthDecrypt(enc []byte, path string, pwd string, typ pwdSessionType) (key *keystore.Key, err error) {
	var (
		existPwd  string
		curPwdEnc []byte
		curPwd    string
	)
	if existPwd, err = getPwdSession(path, typ); err == nil {
		return keystore.DecryptKey(enc, existPwd)
	}

	if key, err = keystore.DecryptKey(enc, pwd); err == nil {
		_ = setPwdSession(path, pwd, typ)
		return
	}

	log.Infof("please input password for ethereum account %s", path)

	for i := 0; i < MaxPwdInputRetry; i++ {
		if curPwdEnc, err = gopass.GetPasswd(); err != nil {
			log.Infof("input error, try it again......")
			continue
		}
		curPwd = string(curPwdEnc)
		if key, err = keystore.DecryptKey(enc, curPwd); err == nil {
			_ = setPwdSession(path, curPwd, typ)
			return
		} else {
			log.Infof("password invalid, err %s, try it again......", err.Error())
		}
	}
	return
}

func readWalletFile(path string) (enc []byte, err error) {
	return ioutil.ReadFile(path)
}

func setPwdSession(path string, pwd string, typ pwdSessionType) error {
	return dao.SavePwd(byte(typ), pathkey(path), []byte(pwd))
}

func getPwdSession(path string, typ pwdSessionType) (string, error) {
	bz, err := dao.GetPwd(byte(typ), pathkey(path))
	if err != nil {
		return "", err
	}
	return string(bz), nil
}

func pathkey(path string) []byte {
	data := md5.Sum([]byte(path))
	return data[:]
}
