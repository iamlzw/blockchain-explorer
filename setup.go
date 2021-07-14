package main

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	configImpl "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	fabImpl "github.com/hyperledger/fabric-sdk-go/pkg/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric/cmd/common/comm"
	_ "github.com/hyperledger/fabric/cmd/common/comm"
	"github.com/hyperledger/fabric/cmd/common/signer"
	_ "github.com/hyperledger/fabric/cmd/common/signer"
	discovery "github.com/hyperledger/fabric/discovery/client"
	"github.com/hyperledger/fabric/protos/utils"
	"strconv"
	"strings"

	discoverypb "github.com/hyperledger/fabric/protos/discovery"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/lifegoeson/blockchain-explorer/common"
	"github.com/lifegoeson/blockchain-explorer/model"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	//"github.com/hyperledger/fabric/protos/utils"
)

const (
	orgName = "Org1"
	orgUser = "User1"
	orgAdmin = "Admin"
	peerUrl = "grpcs://192.168.126.128:7051"
	serverName = "peer0.org1.example.com"
	channelName = "mychannel"
	tlscapath = "E:\\workspace\\go\\src\\github.com\\lifegoeson\\blockchain-explorer\\crypto-config\\peerOrganizations\\org1.example.com\\tlsca\\tlsca.org1.example.com-cert.pem"
	)

	//init the sdk
func initSDK() *fabsdk.FabricSDK {
	//// Initialize the SDK with the configuration file
	configProvider := config.FromFile("config_e2e.yaml")
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		fmt.Errorf("failed to create sdk: %v", err)
	}

	return sdk
}

type ServiceResponse interface {
	// ForChannel returns a ChannelResponse in the context of a given channel
	ForChannel(string) discovery.ChannelResponse

	// ForLocal returns a LocalResponse in the context of no channel
	ForLocal() discovery.LocalResponse

	// Raw returns the raw response from the server
	Raw() *discoverypb.Response
}

func initChannels(sdk *fabsdk.FabricSDK){
	adminContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(orgName))
	
	orgResMgmt, err := resmgmt.New(adminContext)
	if err != nil {
		fmt.Println("Failed to create new resource management client: %s", err)
	}
	configBackend, err := configImpl.FromFile("config_e2e.yaml")()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := fabImpl.ConfigFromBackend(configBackend...)
	if err != nil {
		log.Fatal(err)
	}
	p,err  := peer.New(cfg,peer.WithURL(peerUrl),peer.WithTLSCert(loadCertificate(tlscapath)),peer.WithServerName(serverName))
	chlInfos,err := orgResMgmt.QueryChannels(resmgmt.WithTargets(p))
	common.CheckErr(err)
	chls := chlInfos.Channels
	var i int
	for i = 0 ; i < len(chls) ; i++ {
		chlName := chls[i].ChannelId
		ccp := sdk.ChannelContext(chlName, fabsdk.WithUser(orgUser),fabsdk.WithOrg(orgName))
		ledgerClient, err := ledger.New(ccp)
		common.CheckErr(err)
		block ,err := ledgerClient.QueryBlock(0)
		if err != nil{
			fmt.Println(err)
		}

		output,err := asn1.Marshal(bHeader{Number: int64(block.Header.Number),PreviousHash: string(block.Header.GetPreviousHash()),DataHash: string(block.Header.DataHash)})
		common.CheckErr(err)

		channelGenesisHash := hex.EncodeToString(sha256.New().Sum(output))
		chl := model.Channel{Name:chlName,
			Blocks: 0,
			Trans: 1,
			CreateAt: time.Now(),
			ChannelGenesisHash: channelGenesisHash,
			ChannelHash: "",
			ChannelConfig: nil,
			ChannelBlock: nil,
			ChannelTx: nil,
			ChannelVersion: nil,
		}
		saveChannel(chl)
		discoveryTest(sdk,channelName,channelGenesisHash,orgResMgmt)
	}
}

func queryChaincodeInfo(sdk *fabsdk.FabricSDK,channelName string) *pb.ChaincodeQueryResponse{
	//ledgerClient,err := ledger.New(ccp)

	//prepare context
	adminContext := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("Org1"))

	// Org resource management client
	orgResMgmt, err := resmgmt.New(adminContext)
	if err != nil {
		fmt.Println("Failed to create new resource management client: %s", err)
	}

	//configBackend, err := configImpl.FromFile("config_e2e.yaml")()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//cfg, err := fabImpl.ConfigFromBackend(configBackend...)
	//if err != nil {
	//	log.Fatal(err)
	//}
	common.CheckErr(err)
	//p,err  := peer.New(cfg,peer.WithURL(peerUrl),peer.WithTLSCert(loadCertificate(tlscapath)),peer.WithServerName(serverName))
	//chlInfos,err := orgResMgmt.QueryChannels(resmgmt.WithTargets(p))
	//common.CheckErr(err)
	//chls := chlInfos.Channels
	//var i int
	//for i = 0 ; i < len(chls) ; i++ {
	//	fmt.Println(chls[i].ChannelId)
	//}
	chaincodeInfo,err :=orgResMgmt.QueryInstantiatedChaincodes("mychannel")
	return chaincodeInfo
}

type bHeader struct {
	Number int64
	PreviousHash string
	DataHash string
}

func queryGenesisBlock(sdk *fabsdk.FabricSDK) string{
	ccp := sdk.ChannelContext(channelName, fabsdk.WithUser("User1"),fabsdk.WithOrg("Org1"))
	ledgerClient, err := ledger.New(ccp)
	common.CheckErr(err)
	block ,err := ledgerClient.QueryBlock(0)
	if err != nil{
		fmt.Println(err)
	}

	output,err := asn1.Marshal(bHeader{Number: int64(block.Header.Number),PreviousHash: string(block.Header.GetPreviousHash()),DataHash: string(block.Header.DataHash)})
	common.CheckErr(err)

	//"2dfaf3fa74316ef1b0b476d5535de673ab2516cab93664237bdf3e441558cf6d"
	//"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	return hex.EncodeToString(sha256.New().Sum(output))
}

func loadCertificate(path string) *x509.Certificate{
	cf, e := ioutil.ReadFile(path)
	if e != nil {
		fmt.Println("cfload:", e.Error())
		os.Exit(1)
	}
	cpb, _ := pem.Decode(cf)
	crt, e := x509.ParseCertificate(cpb.Bytes)

	if e != nil {
		fmt.Println("parsex509:", e.Error())
		os.Exit(1)
	}

	return crt
}

const defaultTimeout = time.Second * 5
type ConfigResponseParser struct {
	io.Writer
}

func discoveryTest(sdk *fabsdk.FabricSDK,channelName string,channelGenesisHash string,orgResMgmt *resmgmt.Client){
	const (
		server             = "peer0.org1.example.com:7051"
		discoveryConfigPath = "config/discovery_config.yaml"
	)
	conf,err := ConfigFromFile(discoveryConfigPath)

	client, err := comm.NewClient(conf.TLSConfig)

	siger, err := signer.NewSigner(conf.SignerConfig)
	timeout, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()


	disc := discovery.NewClient(client.NewDialer(server), siger.Sign, 0)
	req := discovery.NewRequest()
	req = req.OfChannel(channelName)
	req = req.AddPeersQuery()
	req = req.AddConfigQuery()
	resp, err := disc.Send(timeout, req, &discoverypb.AuthInfo{
		ClientIdentity:    siger.Creator,
		ClientTlsCertHash: client.TLSCertHash,
	})

	var chlConfig *discoverypb.ConfigResult
	chlConfig, err = resp.ForChannel("mychannel").Config()
	orderers := chlConfig.Orderers
	var p model.Peer
	var prc model.PeerRefChannel
	for k,v := range orderers{
		p.ChannelGenesisHash = channelGenesisHash
		p.CreateAt = time.Now()
		p.ServerHostName = v.GetEndpoint()[0].Host
		p.PeerType = "Orderer"
		p.MspId = k
		p.Requests = "grpcs://"+v.GetEndpoint()[0].Host+":"+ strconv.Itoa(int(v.GetEndpoint()[0].Port))
		p.Org = 0
		prc.CreateAt = time.Now()
		prc.PeerType = "Orderer"
		prc.ChannelId = channelName
		prc.PeerId = v.GetEndpoint()[0].Host
		savePeer(p)
		savePeerChannelRef(prc)
	}
	//fmt.Println(chlConfig.Orderers)
	//jsonBytes, _ := json.MarshalIndent(chlConfig, "", "\t")
	//fmt.Fprintln(os.Stdout, string(jsonBytes))

	//fmt.Println(resp)
	var peers []*discovery.Peer
	peers,err  = resp.ForChannel("mychannel").Peers()
	cqi := queryChaincodeInfo(sdk,channelName)
	var cc model.Chaincode
	fmt.Println(len(peers))
	for i := 0 ; i < len(peers) ; i++ {
		p.MspId = peers[i].MSPID
		p.Org = 0
		p.Requests = "grpcs://"+peers[i].AliveMessage.GetAliveMsg().Membership.Endpoint
		p.PeerType = "PEER"
		p.ServerHostName = strings.Split(peers[i].AliveMessage.GetAliveMsg().Membership.Endpoint,":")[0]
		p.CreateAt = time.Now()
		p.ChannelGenesisHash = channelGenesisHash
		p.Events = ""
		prc.PeerId = strings.Split(peers[i].AliveMessage.GetAliveMsg().Membership.Endpoint,":")[0]
		prc.ChannelId = channelName
		prc.PeerType = "PEER"
		prc.CreateAt = time.Now()
		savePeer(p)
		savePeerChannelRef(prc)
		var j int
		var prcc model.PeerRefChaincode
		for j = 0 ; j < len(cqi.Chaincodes);j++{
			cc.ChannelGenesisHash = channelGenesisHash
			cc.CreateAt = time.Now()
			cc.Name = cqi.Chaincodes[j].Name
			cc.Version = cqi.Chaincodes[j].Version
			cc.Path = cqi.Chaincodes[j].Path
			cc.TxCount = 0
			prcc.CreateAt = time.Now()
			prcc.ChannelId = channelName
			prcc.CCVersion = cqi.Chaincodes[j].Version
			prcc.ChaincodeId = cqi.Chaincodes[j].Name
			prcc.PeerId = strings.Split(peers[i].AliveMessage.GetAliveMsg().Membership.Endpoint,":")[0]
			saveChaincode(cc)
			saveChaincodPeerRef(prcc)
		}

	}
	common.CheckErr(err)
}

type response struct {
	raw *discoverypb.Response
	discovery.Response
}

func (r response) Raw() *discoverypb.Response {
		return r.raw
}

func discoveryRaw(){
	server := "peer0.org1.example.com:7051"
	conf,err := ConfigFromFile("config/discovery_config.yaml")

	client, err := comm.NewClient(conf.TLSConfig)

	siger, err := signer.NewSigner(conf.SignerConfig)
	timeout, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	req := discovery.NewRequest()
	req = req.OfChannel("mychannel")
	req = req.AddPeersQuery()
	req.Authentication = &discoverypb.AuthInfo{
		ClientIdentity:    siger.Creator,
		ClientTlsCertHash: client.TLSCertHash,
	}
	payload := utils.MarshalOrPanic(req.Request)
	sig, err := siger.Sign(payload)

	cc, err := client.NewDialer(server)()

	timeout, cancel = context.WithTimeout(context.Background(), defaultTimeout)

	resp,err := discoverypb.NewDiscoveryClient(cc).Discover(timeout,&discoverypb.SignedRequest{Payload: payload,Signature: sig})

	fmt.Println(resp)
	common.CheckErr(err)
}

