package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/common/tools/protolator"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/rwsetutil"
	ledgerUtil "github.com/hyperledger/fabric/core/ledger/util"
	"github.com/spf13/viper"
	//"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	ccpcontext "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	configImpl "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	fabImpl "github.com/hyperledger/fabric-sdk-go/pkg/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric/cmd/common/comm"
	_ "github.com/hyperledger/fabric/cmd/common/comm"
	"github.com/hyperledger/fabric/cmd/common/signer"
	_ "github.com/hyperledger/fabric/cmd/common/signer"
	discovery "github.com/hyperledger/fabric/discovery/client"
	//"github.com/hyperledger/fabric/protos/utils"
	"strconv"
	"strings"

	cb "github.com/hyperledger/fabric/protos/common"
	discoverypb "github.com/hyperledger/fabric/protos/discovery"
	msp "github.com/hyperledger/fabric/protos/msp"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/utils"
	"github.com/lifegoeson/blockchain-explorer/common"
	"github.com/lifegoeson/blockchain-explorer/defaultclient"
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
//func initSDK() *fabsdk.FabricSDK {
//	//// Initialize the SDK with the configuration file
//	configProvider := config.FromFile("config_e2e.yaml")
//	sdk, err := fabsdk.New(configProvider)
//	if err != nil {
//		fmt.Errorf("failed to create sdk: %v", err)
//	}
//
//	return sdk
//}
//
//type ServiceResponse interface {
//	// ForChannel returns a ChannelResponse in the context of a given channel
//	ForChannel(string) discovery.ChannelResponse
//
//	// ForLocal returns a LocalResponse in the context of no channel
//	ForLocal() discovery.LocalResponse
//
//	// Raw returns the raw response from the server
//	Raw() *discoverypb.Response
//}

func initChannels(){
	orgResMgmt := defaultclient.GetInstance().DefaultResmgmt
	sdk := defaultclient.GetInstance().DefaultFabSdk

	configBackend, err := configImpl.FromFile("config/config_e2e.yaml")()
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
		channelGenesisHash := hex.EncodeToString(block.Header.Hash())
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
		constructBlock(block,channelGenesisHash)
		discoveryFunc(sdk,channelName,channelGenesisHash,orgResMgmt)
		syncBlocks(sdk,channelName,channelGenesisHash)
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
	common.CheckErr(err)
	chaincodeInfo,err :=orgResMgmt.QueryInstantiatedChaincodes("mychannel")
	return chaincodeInfo
}

func syncBlocks(sdk *fabsdk.FabricSDK,chlName string,channelGenesisHash string){
	ccp := sdk.ChannelContext(chlName, fabsdk.WithUser("User1"),fabsdk.WithOrg("Org1"))
	ledgerClient, err := ledger.New(ccp)
	chainInfos,err := ledgerClient.QueryInfo()
	var b *cb.Block
	//bb:= new(model.Block)
	var i uint64
	for i = 1 ; i < chainInfos.BCI.Height ; i++{
		b ,err = ledgerClient.QueryBlock(i)
		constructBlock(b,channelGenesisHash)
	}
	common.CheckErr(err)
	//listenBlockEvent(ccp)
}

func constructBlock(b *cb.Block,channelGenesisHash string){
	bb:= new(model.Block)
	bb.BlockNum = int64(b.Header.Number)
	bb.PrevBlockHash = ""
	bb.ChannelGenesisHash = channelGenesisHash
	bb.TxCount = int64(len(b.Data.Data))
	bb.DataHash = hex.EncodeToString(b.Header.DataHash)
	bb.PreHash = hex.EncodeToString(b.Header.PreviousHash)
	//fmt.Println(b.Header.Bytes()
	bb.BlockHash = hex.EncodeToString(b.Header.Hash())
	env, err := utils.GetEnvelopeFromBlock(b.Data.Data[0])
	payload,err := utils.GetPayload(env)
	chdr, err := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	bb.CreateAt = chdr.Timestamp.AsTime()
	saveBlock(bb)
	txsfltr := getBlockMetaData(b)
	for j := 0 ; j < len(b.Data.Data) ; j++{
		e, _ := utils.GetEnvelopeFromBlock(b.Data.Data[j])
		syncTx(e,txsfltr,bb.BlockNum,channelGenesisHash,j,chdr.Extension,chdr.Type)
	}
	common.CheckErr(err)
}

func syncTx(env *cb.Envelope,txsfltr ledgerUtil.TxValidationFlags,blockId int64,channelGenesisHash string,txIndex int,payload_extension []byte,header_type int32){
	tx := new(model.Transaction)
	payload,_ := utils.GetPayload(env)
	chdr, _ := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	tx.BlockId = blockId
	//fmt.Println(blockId)
	tx.TxHash = chdr.TxId
	tx.CreateAt = chdr.Timestamp.AsTime()
	t, err := utils.GetTransaction(payload.Data)
	if err != nil {
		log.Fatal(err)
	}
	sdr, err := utils.UnmarshalSignatureHeader(payload.Header.SignatureHeader)
	var sId2 msp.SerializedIdentity
	_ = proto.Unmarshal(sdr.Creator, &sId2)
	tx.EnvelopeSignature = hex.EncodeToString(env.Signature)
	tx.PayloadExtension = hex.EncodeToString(payload_extension)
	tx.CreatorIdBytes = string(sId2.IdBytes)
	tx.CreatorNonce = hex.EncodeToString(sdr.Nonce)
	var sId msp.SerializedIdentity
	_ = proto.Unmarshal(sdr.Creator, &sId)
	tx.CreatorMspId = sId.Mspid
	//fmt.Println(header_type)
	tx.Type = cb.HeaderType(header_type).String()
	if header_type == int32(3) {
		ccActionPayload, err := utils.GetChaincodeActionPayload(t.Actions[0].Payload)
		if err != nil {
			log.Fatal(err)
		}
		//ccActionHeader,err := utils.UnmarshalSignatureHeader(t.Actions[0].Header)
		prp, err := utils.GetProposalResponsePayload(ccActionPayload.Action.ProposalResponsePayload)
		if err != nil {
			log.Fatal(err)
		}
		caPayload, err := utils.GetChaincodeAction(prp.Extension)
		tx.ChaincodeName = caPayload.ChaincodeId.Name
		tx.Status = caPayload.Response.Status
		tx.EndorserMspId = getMSPIdFromEndorsement(ccActionPayload.Action.Endorsements)
		tx.ChaincodeId = ""
		//tx.Type = string(cb.HeaderType(chdr.Type))
		//var results
		t2 := &rwsetutil.TxRwSet{}
		_ = t2.FromProtoBytes(caPayload.Results)
		//var results lb.TxReadWriteSet
		//_ = proto.Unmarshal(caPayload.Results, &t2)
		tx.ReadSet,tx.WriteSet = getRwSet(t2.NsRwSets)
		tx.ChannelGenesisHash = channelGenesisHash
		tx.ValidationCode = txsfltr.Flag(txIndex).String()
		tx.PayloadProposalHash = hex.EncodeToString(prp.ProposalHash)
		endorser := ccActionPayload.Action.Endorsements[0].Endorser
		//var sId msp.SerializedIdentity
		_ = proto.Unmarshal(endorser, &sId)
		tx.EndorserIdBytes = string(sId.IdBytes)
		tx.EndorserSignature = hex.EncodeToString(ccActionPayload.Action.Endorsements[0].Signature)
		input := ""
		cpp := &pb.ChaincodeProposalPayload{}
		_ = proto.Unmarshal(ccActionPayload.ChaincodeProposalPayload, cpp)
		cis := &pb.ChaincodeInvocationSpec{}
		err = proto.Unmarshal(cpp.Input, cis)
		args := cis.ChaincodeSpec.Input.Args
		for k := 0; k < len(args);k++ {
			if k == len(args) - 1  {
				input += string(args[k])
			}else {
				input += string(args[k]) + ","
			}
		}
		tx.ChaincodeProposalInput = input
		tx.TxResponse = ""
	}else {
		rset,err := json.Marshal([]byte(""))
		if err != nil {
			log.Fatal(err)
		}
		wset,err := json.Marshal([]byte(""))

		if err != nil {
			log.Fatal(err)
		}
		tx.ReadSet = string(rset)
		tx.WriteSet = string(wset)
		tx.ChannelGenesisHash = channelGenesisHash
	}
	saveTransaction(tx)

}

func getBlockMetaData(b *cb.Block) ledgerUtil.TxValidationFlags{
	//md := &cb.Metadata{}
	//_ = proto.Unmarshal(blockMetadata[0], md)
	var txsfltr ledgerUtil.TxValidationFlags
	txsfltr = b.Metadata.Metadata[cb.BlockMetadataIndex_TRANSACTIONS_FILTER]
	return txsfltr
}

func getRwSet(nss []*rwsetutil.NsRwSet)(string,string){
	var reads []map[string]interface{}
	var writes []map[string]interface{}

	for i := 0 ; i < len(nss) ;i++{
		//var nsRwSet lb.NsReadWriteSet
		//_ = proto.Unmarshal(nss[i].KvRwSet.,&nsRwSet)
		rm := make(map[string]interface{})
		wm := make(map[string]interface{})
		//var rwset kvb.KVRWSet/
		//_ = proto.Unmarshal(nss[i].KvRwSet.,&rwset)
		rm["chaincode"] = nss[i].NameSpace
		rm["set"] = nss[i].KvRwSet.Reads
		wm["chaincode"] = nss[i].NameSpace
		wm["set"] = nss[i].KvRwSet.Writes
		reads = append(reads,rm)
		writes = append(writes,wm)
	}

	rss,_ := json.Marshal(reads)
	wss,_ := json.Marshal(writes)

	return string(rss),string(wss)
}


func getMSPIdFromEndorsement(endorsements []*pb.Endorsement) string {
	mspid := "{"
	for i := 0 ; i < len(endorsements) ; i++{
		var sId msp.SerializedIdentity
		_ = proto.Unmarshal(endorsements[i].Endorser, &sId)
		if i == len(endorsements) - 1 {
			mspid += "\""+ sId.GetMspid() + "\""
		}else {
			mspid += "\""+ sId.GetMspid() + "\","
		}
	}
	mspid += "}"
	return mspid
}

func blockFromProto2Struct(b *cb.Block) *viper.Viper{
	buf := new (bytes.Buffer)
	err := protolator.DeepMarshalJSON(buf, b)
	//err = protolator.DeepMarshalJSON(os.Stdout, b)
	v := viper.New()
	v.SetConfigType("json")
	_ = v.MergeConfig(buf)
	common.CheckErr(err)
	return v
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

func discoveryFunc(sdk *fabsdk.FabricSDK,channelName string,channelGenesisHash string,orgResMgmt *resmgmt.Client){
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

func listenBlockEvent(ccp ccpcontext.ChannelProvider){
	ec,err := event.New(ccp,event.WithBlockEvents())

	if err !=nil {
		fmt.Errorf("init event client error %s",err)
	}

	reg, notifier, err :=ec.RegisterBlockEvent()

	if err != nil {
		fmt.Printf("Failed to register block event: %s", err)
		return
	}
	defer ec.Unregister(reg)

	var bEvent *fab.BlockEvent

	for  {
		select {
		case bEvent = <-notifier:
			b := bEvent.Block
			constructBlock(b,"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
		}
	}
}

