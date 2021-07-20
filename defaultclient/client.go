package defaultclient

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/spf13/viper"
	"log"
	"sync"
)

type defaultclient struct {
	DefaultPeer       string
	DefaultOrg        string
	DefaultOrgUser    string
	DefaultOrgAdmin   string
	DefaultServerName string
	DefaultChanel     string
	DefaultFabSdk     *fabsdk.FabricSDK
	DefaultResmgmt    *resmgmt.Client
	DefaultCCP        context.ChannelProvider
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
		instance.DefaultPeer = v.GetString("default_peer_url")
		instance.DefaultOrgUser = v.GetString("default_org_user")
		instance.DefaultOrgAdmin = v.GetString("default_org_admin")
		cfp := config.FromFile("config/config_e2e.yaml")
		sdk, err := fabsdk.New(cfp)
		if err != nil {
			fmt.Errorf("failed to create sdk: %v", err)
		}
		adminContext := sdk.Context(fabsdk.WithUser(instance.DefaultOrgAdmin), fabsdk.WithOrg(instance.DefaultOrg))

		orgResMgmt, err := resmgmt.New(adminContext)

		ccp := sdk.ChannelContext(instance.DefaultChanel, fabsdk.WithUser("User1"),fabsdk.WithOrg("Org1"))
		instance.DefaultCCP = ccp
		instance.DefaultFabSdk = sdk
		instance.DefaultResmgmt = orgResMgmt
	})
	return instance
}
