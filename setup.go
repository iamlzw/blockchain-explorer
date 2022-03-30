package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/platforms"
	"github.com/hyperledger/fabric/core/chaincode/platforms/car"
	"github.com/hyperledger/fabric/core/chaincode/platforms/golang"
	"github.com/hyperledger/fabric/core/chaincode/platforms/java"
	"github.com/hyperledger/fabric/core/chaincode/platforms/node"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/rwsetutil"
	ledgerUtil "github.com/hyperledger/fabric/core/ledger/util"
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
	cmdcommon "github.com/hyperledger/fabric/cmd/common"
	"github.com/hyperledger/fabric/cmd/common/comm"
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
	"github.com/lifegoeson/blockchain-explorer/service"
	"log"
	"time"
)
//type Client struct {
//	ctx       ccpcontext.Channel
//	filter    fab.TargetFilter
//	ledger    *channel.Ledger
//	verifier  channel.ResponseVerifier
//	discovery fab.DiscoveryService
//}
func initChannels(){
	orgResMgmt := defaultclient.GetInstance().DefaultResmgmt
	//sdk := defaultclient.GetInstance().DefaultFabSdk

	config := configImpl.FromFile("config/config_e2e.yaml")
	configBackend , err := config()

	if err != nil {
		log.Println(err)
	}

	cfg, err := fabImpl.ConfigFromBackend(configBackend...)
	if err != nil {
		log.Println(err)
	}
	p,err  := peer.New(cfg,peer.WithURL(cfg.NetworkPeers()[0].URL),peer.WithTLSCert(cfg.NetworkPeers()[0].TLSCACert))
	chlInfos,err := orgResMgmt.QueryChannels(resmgmt.WithTargets(p))
	//p.ProcessTransactionProposal()

	if err != nil {
		log.Println(err)
	}
	chls := chlInfos.Channels

	sdk, err := fabsdk.New(config)
	if err != nil {
		log.Println(err)
	}
	var i int
	for i = 0 ; i < len(chls) ; i++ {
		chlName := chls[i].ChannelId
		ccp := sdk.ChannelContext(chlName, fabsdk.WithUser(defaultclient.GetInstance().DefaultOrgUser),fabsdk.WithOrg(defaultclient.GetInstance().DefaultOrg))
		ledgerClient, err := ledger.New(ccp)
		if err != nil {
			log.Println(err)
		}

		//ledgerclient := (*Client)(unsafe.Pointer(ledgerClient))
		//peers, _ := ledgerclient.discovery.GetPeers()
		//fmt.Println(peers[0].URL())

		block ,err := ledgerClient.QueryBlock(0)
		//var client interface{}

		if err != nil{
			fmt.Println(err)
		}
		channelGenesisHash := hex.EncodeToString(block.Header.Hash())
		chl := model.Channel{Name:chlName,
			Blocks: 0,
			Trans: 0,
			CreateAt: time.Now(),
			ChannelGenesisHash: channelGenesisHash,
			ChannelHash: "",
			ChannelConfig: nil,
			ChannelBlock: nil,
			ChannelTx: nil,
			ChannelVersion: nil,
		}
		service.SaveChannel(chl)
		constructBlock(block,channelGenesisHash)
		discoveryNetwork(sdk,chlName,channelGenesisHash)
		bc := make(chan *cb.Block,100)
		go processBlockChannel(bc,channelGenesisHash)
		go syncBlocks(sdk,chlName,channelGenesisHash,bc)
		go listenBlockEvent(ccp,channelGenesisHash)
		//go catchupBlocks(sdk,channelName,channelGenesisHash)
	}
}

func queryChaincodeInfo(sdk *fabsdk.FabricSDK,channelName string) *pb.ChaincodeQueryResponse{
	//ledgerClient,err := ledger.New(ccp)

	//prepare context
	adminContext := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("Org1"))

	// Org resource management client
	orgResMgmt, err := resmgmt.New(adminContext)
	if err != nil {
		log.Println(err)
	}
	chaincodeInfo,err :=orgResMgmt.QueryInstantiatedChaincodes(channelName)
	return chaincodeInfo
}
//同步区块,仅在初始化channel时同步一次
func syncBlocks(sdk *fabsdk.FabricSDK,chlName string,channelGenesisHash string,bc chan *cb.Block){
	ccp := sdk.ChannelContext(chlName, fabsdk.WithUser("User1"),fabsdk.WithOrg("Org1"))
	ledgerClient, err := ledger.New(ccp)
	chainInfos,err := ledgerClient.QueryInfo()
	if err != nil {
		log.Println(err)
	}
	var b *cb.Block
	//bb:= new(model.Block)
	var i uint64
	for i = 1 ; i < chainInfos.BCI.Height ; i++{
		b ,err = ledgerClient.QueryBlock(i)
		for i = 1 ; i < chainInfos.BCI.Height ; i++{
			b ,err = ledgerClient.QueryBlock(i)
			if err != nil {
				log.Println(err)
			}
			bc <- b
		}
	}

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
	if err != nil {
		log.Println(err)
	}
	payload,err := utils.GetPayload(env)
	if err != nil {
		log.Println(err)
	}
	chdr, err := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		log.Println(err)
	}
	x :=  chdr.Timestamp
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	//fmt.Println("SH : ", time.Now().In(cstSh).Format("2006-01-02 15:04:05"))
	bb.CreateAt = time.Unix(int64(x.GetSeconds()), int64(x.GetNanos())).In(cstSh)
	txsfltr := getBlockMetaData(b)
	for j := 0 ; j < len(b.Data.Data) ; j++{
		e, _ := utils.GetEnvelopeFromBlock(b.Data.Data[j])
		syncTx(e,txsfltr,bb.BlockNum,channelGenesisHash,j,chdr.Extension,chdr.Type)
	}
	service.SaveBlock(bb)
}

func syncTx(env *cb.Envelope,txsfltr ledgerUtil.TxValidationFlags,blockId int64,channelGenesisHash string,txIndex int,payload_extension []byte,header_type int32){
	tx := new(model.Transaction)
	payload,_ := utils.GetPayload(env)
	chdr, _ := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	tx.BlockId = blockId
	//fmt.Println(blockId)
	if chdr.TxId == ""{
		tx.TxHash = "nil"
	}else {
		tx.TxHash = chdr.TxId
	}
	x :=  chdr.Timestamp
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	//fmt.Println("SH : ", time.Now().In(cstSh).Format("2006-01-02 15:04:05"))
	tx.CreateAt = time.Unix(int64(x.GetSeconds()), int64(x.GetNanos())).In(cstSh)
	//tx.CreateAt = chdr.Timestamp.AsTime()
	t, err := utils.GetTransaction(payload.Data)
	if err != nil {
		log.Println(err)
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
			log.Println(err)
		}
		//ccActionHeader,err := utils.UnmarshalSignatureHeader(t.Actions[0].Header)
		prp, err := utils.GetProposalResponsePayload(ccActionPayload.Action.ProposalResponsePayload)

		if err != nil {
			log.Println(err)
		}
		caPayload, err := utils.GetChaincodeAction(prp.Extension)
		if err != nil {
			log.Println(err)
		}
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
		lsccFunc := string(args[0])
		lsccArgs := args[1:]
		if lsccFunc == "upgrade" || lsccFunc == "deploy" {

			cds := &pb.ChaincodeDeploymentSpec{}
			_ = proto.Unmarshal(lsccArgs[1], cds)

			cdsArgs, _ := utils.GetChaincodeDeploymentSpec(lsccArgs[1], platforms.NewRegistry(
				// XXX We should definitely _not_ have this external dependency in VSCC
				// as adding a platform could cause non-determinism.  This is yet another
				// reason why all of this custom LSCC validation at commit time has no
				// long term hope of staying deterministic and needs to be removed.
				&golang.Platform{},
				&node.Platform{},
				&java.Platform{},
				&car.Platform{},
			))

			inputArgs := cdsArgs.ChaincodeSpec.Input.Args
			input += lsccFunc + "," +cds.ChaincodeSpec.ChaincodeId.Name + ","

			for k := 0; k < len(inputArgs);k++ {
				if k == len(inputArgs) - 1  {
					//fmt.Println(string(inputArgs[k]))
					input += string(inputArgs[k])
				}else {
					//fmt.Println(string(inputArgs[k]))
					input += string(inputArgs[k]) + ","
				}
			}

			fmt.Println(input)
		}else {
			//fmt.Println(lsccArgs)
			args := cis.ChaincodeSpec.Input.Args
			for k := 0; k < len(args);k++ {
				if k == len(args) - 1  {
					//fmt.Println(string(args[k]))
					input += string(args[k])
				}else {
					//fmt.Println(string(args[k]))
					input += string(args[k]) + ","
				}
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
	service.SaveTransaction(tx)
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

//func blockFromProto2Struct(b *cb.Block) *viper.Viper{
//	buf := new (bytes.Buffer)
//	err := protolator.DeepMarshalJSON(buf, b)
//	//err = protolator.DeepMarshalJSON(os.Stdout, b)
//	v := viper.New()
//	v.SetConfigType("json")
//	_ = v.MergeConfig(buf)
//	common.CheckErr(err)
//	return v
//}

//type bHeader struct {
//	Number int64
//	PreviousHash string
//	DataHash string
//}


func listenBlockEvent(ccp ccpcontext.ChannelProvider,channelGenesisHash string){
	ec,err := event.New(ccp,event.WithBlockEvents())

	if err !=nil {
		log.Printf("init event client error %s\n", err)
	}

	reg, notifier, err :=ec.RegisterBlockEvent()

	if err != nil {
		log.Printf("Failed to register block event: %s", err)
	}
	defer ec.Unregister(reg)

	var bEvent *fab.BlockEvent

	for  {
		select {
		case bEvent = <-notifier:
			b := bEvent.Block
			constructBlock(b,channelGenesisHash)
		}
	}
}

func processBlockChannel(bc chan  *cb.Block,channelGenesisHash string){
	var b *cb.Block
	for  {
		select {
		case b = <-bc:
			constructBlock(b,channelGenesisHash)
		}
	}
}

func discoveryNetwork(sdk *fabsdk.FabricSDK,channelName string,channelGenesisHash string){
	const (
		server             = "peer0.org1.example.com:7051"
		discoveryConfigPath = "config/discovery_config.yaml"
	)
	conf,err := cmdcommon.ConfigFromFile(discoveryConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	client, err := comm.NewClient(conf.TLSConfig)

	siger, err := signer.NewSigner(conf.SignerConfig)
	timeout, cancel := context.WithTimeout(context.Background(), time.Second*5)
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
	chlConfig, err = resp.ForChannel(channelName).Config()
	orderers := chlConfig.Orderers
	var p model.Peer
	var prc model.PeerRefChannel
	for k,v := range orderers{
		p.ChannelGenesisHash = channelGenesisHash
		p.CreateAt = time.Now()
		for i := 0 ; i < len(v.GetEndpoint()) ; i++{
			p.ServerHostName = v.GetEndpoint()[i].Host
			p.PeerType = "Orderer"
			p.MspId = k
			p.Requests = "grpcs://"+v.GetEndpoint()[i].Host+":"+ strconv.Itoa(int(v.GetEndpoint()[i].Port))
			p.Org = 0
			prc.CreateAt = time.Now()
			prc.PeerType = "Orderer"
			prc.ChannelId = channelGenesisHash
			prc.PeerId = v.GetEndpoint()[i].Host
			service.SavePeer(p)
			service.SavePeerChannelRef(prc)
		}
	}
	//fmt.Println(chlConfig.Orderers)
	//jsonBytes, _ := json.MarshalIndent(chlConfig, "", "\t")
	//fmt.Fprintln(os.Stdout, string(jsonBytes))

	//fmt.Println(resp)
	var peers []*discovery.Peer
	peers,err  = resp.ForChannel(channelName).Peers()
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
		prc.ChannelId = channelGenesisHash
		prc.PeerType = "PEER"
		prc.CreateAt = time.Now()
		service.SavePeer(p)
		service.SavePeerChannelRef(prc)
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
			prcc.ChannelId = channelGenesisHash
			prcc.CCVersion = cqi.Chaincodes[j].Version
			prcc.ChaincodeId = cqi.Chaincodes[j].Name
			prcc.PeerId = strings.Split(peers[i].AliveMessage.GetAliveMsg().Membership.Endpoint,":")[0]
			service.SaveChaincode(cc)
			service.SaveChaincodPeerRef(prcc)
		}

	}
	common.CheckErr(err)
}

//func queryGenesisBlock(sdk *fabsdk.FabricSDK) string{
//	ccp := sdk.ChannelContext(channelName, fabsdk.WithUser("User1"),fabsdk.WithOrg("Org1"))
//	ledgerClient, err := ledger.New(ccp)
//	common.CheckErr(err)
//	block ,err := ledgerClient.QueryBlock(0)
//	if err != nil{
//		fmt.Println(err)
//	}
//
//	output,err := asn1.Marshal(bHeader{Number: int64(block.Header.Number),PreviousHash: string(block.Header.GetPreviousHash()),DataHash: string(block.Header.DataHash)})
//	common.CheckErr(err)
//
//	return hex.EncodeToString(sha256.New().Sum(output))
//}

//func loadCertificate(path string) *x509.Certificate{
//	cf, e := ioutil.ReadFile(path)
//	if e != nil {
//		fmt.Println("cfload:", e.Error())
//		os.Exit(1)
//	}
//	cpb, _ := pem.Decode(cf)
//	crt, e := x509.ParseCertificate(cpb.Bytes)
//
//	if e != nil {
//		fmt.Println("parsex509:", e.Error())
//		os.Exit(1)
//	}
//
//	return crt
//}

//const defaultTimeout = time.Second * 5
//type ConfigResponseParser struct {
//	io.Writer
//}

//
//type response struct {
//	raw *discoverypb.Response
//	discovery.Response
//}
//
//func (r response) Raw() *discoverypb.Response {
//		return r.raw
//}

//func discoveryRaw(){
//	server := "peer0.org1.example.com:7051"
//	conf,err := cmdcommon.ConfigFromFile("config/discovery_config.yaml")
//
//	client, err := comm.NewClient(conf.TLSConfig)
//
//	siger, err := signer.NewSigner(conf.SignerConfig)
//	timeout, cancel := context.WithTimeout(context.Background(), time.Second*5)
//	defer cancel()
//
//	req := discovery.NewRequest()
//	req = req.OfChannel("channel2")
//	req = req.AddPeersQuery()
//	req.Authentication = &discoverypb.AuthInfo{
//		ClientIdentity:    siger.Creator,
//		ClientTlsCertHash: client.TLSCertHash,
//	}
//	payload := utils.MarshalOrPanic(req.Request)
//	sig, err := siger.Sign(payload)
//
//	cc, err := client.NewDialer(server)()
//
//	timeout, cancel = context.WithTimeout(context.Background(), time.Second*5)
//
//	resp,err := discoverypb.NewDiscoveryClient(cc).Discover(timeout,&discoverypb.SignedRequest{Payload: payload,Signature: sig})
//
//	fmt.Println(resp)
//	common.CheckErr(err)
//}


