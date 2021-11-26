package config

import (
	"bytes"
	"crypto/ecdsa"
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
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/palettechain/deploy-tool/pkg/dao"
	"github.com/palettechain/deploy-tool/pkg/files"
	"github.com/palettechain/deploy-tool/pkg/log"
	"github.com/palettechain/deploy-tool/pkg/poly"
	"github.com/palettechain/deploy-tool/pkg/sdk"
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
	LevelDB string

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
	PalettePLTWrapper    common.Address
	PaletteNFTWrapper    common.Address
	PaletteNFTQuery  	  common.Address

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

	sdk.Init()

	// init leveldb
	dao.NewDao(Conf.LevelDB)
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

func SaveConfig(c *Config) error {
	type XConfig struct {
		LevelDB string

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
		PalettePLTWrapper    common.Address
		PaletteNFTWrapper    common.Address
		PaletteNFTQuery  	  common.Address

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

	x := new(XConfig)
	x.LevelDB = c.LevelDB
	x.PolyRPCUrl = c.PolyRPCUrl
	x.PolyAccountDir = c.PolyAccountDir

	x.EthereumRPCUrl = c.EthereumRPCUrl
	x.EthereumCrossChainAdmin = c.EthereumCrossChainAdmin

	x.PaletteRPCUrl = c.PaletteRPCUrl
	x.PaletteCrossChainAdmin = c.PaletteCrossChainAdmin

	x.PaletteSideChainID = c.PaletteSideChainID
	x.PaletteSideChainName = c.PaletteSideChainName
	x.PaletteECCD = c.PaletteECCD
	x.PaletteECCM = c.PaletteECCM
	x.PaletteCCMP = c.PaletteCCMP
	x.PaletteNFTProxy = c.PaletteNFTProxy
	x.PalettePLTWrapper = c.PalettePLTWrapper
	x.PaletteNFTWrapper = c.PaletteNFTWrapper
	x.PaletteNFTQuery = c.PaletteNFTQuery

	x.EthereumSideChainID = c.EthereumSideChainID
	x.EthereumSideChainName = c.EthereumSideChainName
	x.EthereumECCD = c.EthereumECCD
	x.EthereumECCM = c.EthereumECCM
	x.EthereumCCMP = c.EthereumCCMP
	x.EthereumPLTAsset = c.EthereumPLTAsset
	x.EthereumPLTProxy = c.EthereumPLTProxy
	x.EthereumNFTProxy = c.EthereumNFTProxy

	x.PaletteNFTAsset = c.PaletteNFTAsset
	x.EthereumNFTAsset = c.PaletteNFTAsset

	enc, err := json.Marshal(x)
	if err != nil {
		return err
	}
	var out bytes.Buffer
	json.Indent(&out, enc, "", "\t")
	return ioutil.WriteFile(ConfigFilePath, out.Bytes(), os.ModePerm)
}

func (c *Config) LoadPLTAdminAccount() (*ecdsa.PrivateKey, error) {
	return getEthAccount(Conf.PaletteCrossChainAdmin, pwdSessionPLT)
}

func (c *Config) LoadETHAdminAccount() (*ecdsa.PrivateKey, error) {
	return getEthAccount(Conf.EthereumCrossChainAdmin, pwdSessionETH)
}

func (c *Config) LoadPolyAccountList() []*polysdk.Account {
	list := make([]*polysdk.Account, 0)

	//dir := c.PolyAccountDir
	//fs, err := ioutil.ReadDir(dir)
	//if err != nil {
	//	panic(err)
	//}
	acc, err := c.LoadPolyAccount(c.PolyAccountDir)
	if err != nil {
		panic(err)
	}
	list = append(list, acc)

	return list
}

func (c *Config) LoadPolyCurBookeeperBytes() []byte {
	accs := c.LoadPolyAccountList()
	keepers := []keypair.PublicKey{}
	for _, v := range accs {
		keepers = append(keepers, v.PublicKey)
	}
	sink, _ := poly.AssemblePubKeyList(keepers)
	return sink.Bytes()
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

func (c *Config) StorePalettePLTWrapper(addr common.Address) error {
	c.PalettePLTWrapper = addr
	return SaveConfig(Conf)
}

func (c *Config) StorePaletteNFTWrapper(addr common.Address) error {
	c.PaletteNFTWrapper = addr
	return SaveConfig(Conf)
}

func (c *Config) StorePaletteNFTQuery(addr common.Address) error {
	c.PaletteNFTQuery = addr
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

func repeatPolyDecrypt(wallet *polysdk.Wallet, filepath string) (acc *polysdk.Account, err error) {
	var (
		existPwd string
		curPwd   []byte
		typ      = pwdSessionPoly
	)

	_, fn := path.Split(filepath)
	if existPwd, err = getPwdSession(fn, typ); err == nil {
		acc, err = wallet.GetDefaultAccount([]byte(existPwd))
		return
	}

	log.Infof("please input password for poly account %s", fn)

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

func repeatEthDecrypt(enc []byte, filepath string, pwd string, typ pwdSessionType) (key *keystore.Key, err error) {
	var (
		existPwd  string
		curPwdEnc []byte
		curPwd    string
	)
	_, fn := path.Split(filepath)
	if existPwd, err = getPwdSession(fn, typ); err == nil {
		return keystore.DecryptKey(enc, existPwd)
	}

	if key, err = keystore.DecryptKey(enc, pwd); err == nil {
		_ = setPwdSession(fn, pwd, typ)
		return
	}

	log.Infof("please input password for ethereum account %s", fn)

	for i := 0; i < MaxPwdInputRetry; i++ {
		if curPwdEnc, err = gopass.GetPasswd(); err != nil {
			log.Infof("input error, try it again......")
			continue
		}
		curPwd = string(curPwdEnc)
		if key, err = keystore.DecryptKey(enc, curPwd); err == nil {
			_ = setPwdSession(fn, curPwd, typ)
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

func setPwdSession(fn string, pwd string, typ pwdSessionType) error {
	return dao.SavePwd(byte(typ), []byte(fn), []byte(pwd))
}

func getPwdSession(fn string, typ pwdSessionType) (string, error) {
	bz, err := dao.GetPwd(byte(typ), []byte(fn))
	if err != nil {
		return "", err
	}
	return string(bz), nil
}
