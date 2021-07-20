package defaultclient

import (
	"encoding/hex"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/spf13/viper"
	"log"
	"sync"
)

type defaultclient struct {
	DefaultClientName string
	DefaultPeer       string
	DefaultOrg        string
	DefaultOrgUser    string
	DefaultOrgAdmin   string
	DefaultServerName string
	DefaultChanel     string
	DefaultFabSdk     *fabsdk.FabricSDK
	DefaultResmgmt    *resmgmt.Client
	DefaultCCP        context.ChannelProvider
	DefaultChannelGenHash string
	DefautlMSPId string
}

var instance *defaultclient
var once sync.Once

func GetInstance() *defaultclient {
	once.Do(func() {
		instance = &defaultclient{}
		v := viper.New()
		v.SetConfigType("yaml")
		v.SetConfigFile("config/default_config.yaml")
		err := v.ReadInConfig()
		if err != nil {
			log.Fatal(err)
		}
		instance.DefaultOrg = v.GetString("default_org_name")
		instance.DefaultServerName = v.GetString("default_server_name")
		fmt.Println(instance.DefaultServerName)
		instance.DefaultPeer = v.GetString("default_peer_url")
		instance.DefaultOrgUser = v.GetString("default_org_user")
		instance.DefaultOrgAdmin = v.GetString("default_org_admin")
		instance.DefaultChanel = v.GetString("default_channel_name")
		cfp := config.FromFile("config/config_e2e.yaml")
		sdk, err := fabsdk.New(cfp)
		if err != nil {
			fmt.Errorf("failed to create sdk: %v", err)
		}
		adminContext := sdk.Context(fabsdk.WithUser(instance.DefaultOrgAdmin), fabsdk.WithOrg(instance.DefaultOrg))

		orgResMgmt, err := resmgmt.New(adminContext)

		ccp := sdk.ChannelContext(instance.DefaultChanel, fabsdk.WithUser(instance.DefaultOrgUser),fabsdk.WithOrg(instance.DefaultOrg))
		instance.DefaultCCP = ccp
		instance.DefaultFabSdk = sdk
		instance.DefaultResmgmt = orgResMgmt
		ledgerClient, err := ledger.New(ccp)
		if err != nil {
			log.Fatal(err)
		}
		block ,err := ledgerClient.QueryBlock(0)
		if err != nil{
			log.Fatal(err)
		}
		channelGenesisHash := hex.EncodeToString(block.Header.Hash())
		instance.DefaultChannelGenHash = channelGenesisHash

	})
	return instance
}
