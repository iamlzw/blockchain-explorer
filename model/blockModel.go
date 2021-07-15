package model

import "time"

//data
type Blk struct {
	Data BlockData `json:"data"`
	Header Header `json:"header"`
	Metadata Metadata `json:"metadata"`
}

type Header struct {
	DataHash string `json:"data_hash"`
	Number string `json:"number"`
	PreviousHash string `json:"previous_hash"`
}

type Metadata struct {
	Metadata []byte `json:"metadata"`
}
//data---
//----data
type BlockData struct {
	Data []Envelope `json:"data"`
}
//data---
//----data
//--------
type Envelope struct {
	Payload Payload `json:"payload"`
	Signature []byte `json:"signature"`
}

type Payload struct {
	PayloadHeader PayloadHeader `json:"header"`
	PayloadData PayloadData `json:"data"`
}

type PayloadHeader struct {
	ChannelHeader ChannelHeader `json:"channel_header"`
	SignatureHeader SignatureHeader `json:"signature_header"`
}

type ChannelHeader struct {
	ChannelId string `json:"channel_id"`
	Epoch string `json:"epoch"`
	Extension []byte `json:"extension"`
	Timestamp time.Time `json:"timestamp"`
	TlsCertHash string `json:"tls_cert_hash"`
	TxId string `json:"tx_id"`
	Type int64 `json:"type"`
	Version int64 `json:"version"`
}

type CommonHeader struct {
	HeaderCreator SignatureHeaderCreator `json:"creator"`
	HeaderNonce string `json:"nonce"`
}

type SignatureHeader struct {
	SignatureHeaderCreator SignatureHeaderCreator `json:"creator"`
	SignatureHeaderNonce string `json:"nonce"`
}
type SignatureHeaderCreator struct {
	IdBytes string `json:"id_bytes"`
	MspId string `json:"mspid"`
}

type PayloadData struct {
	Actions []Action `json:"actions"`
}

type Action struct {
	ActionHeader CommonHeader `json:"header"`
	ActionPayload ActionPayload `json:"payload"`
}

type ActionPayload struct {
	Action PayloadAction `json:"action"`
	ChaincodeProposalPayload ChaincodeProposalPayload `json:"chaincode_proposal_payload"`
}

type PayloadAction struct {
	Endorsements []Endorsement `json:"endorsements"`
	ProposalResponsePayload ProposalResponsePayload `json:"proposal_response_payload"`
}

type Endorsement struct {
	Endorser []byte `json:"endorser"`
	Signature string `json:"signature"`
}

type ProposalResponsePayload struct {
	Extension Extension `json:"extension"`
	ProposalHash string `json:"proposal_hash"`
}

type Extension struct {
	ChaincodeId ChaincodeId `json:"chaincode_id"`
	Events interface{} `json:"events"`
	Response Response `json:"response"`
	Results Results `json:"results"`
	TokenExpectation interface{} `json:"token_expectation"`
}
type ChaincodeId struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Version string `json:"version"`
}
type Response struct {
	Message string `json:"message"`
	Payload interface{} `json:"payload"`
	Status int64 `json:"status"`
}

type Results struct {
	DataModel string `json:"data_model"`
	NsRwSet []NsRwSet `json:"ns_rwset"`
}

type NsRwSet struct {
	CollectionHashedRwSet interface{} `json:"collection_hashed_rwset"`
	Namespace string `json:"namespace"`
	RwSet RwSet `json:"rwset"`
}


type RwSet struct {
	MetadataWrites []interface{} `json:"metadata_writes"`
	RangeQueriesInfo []interface{} `json:"range_queries_info"`
	Reads []ReadSet `json:"reads"`
	Writes []WriteSet `json:"writes"`
}

type ReadSet struct {
	Key string `json:"key"`
	Version ReadSetVersion `json:"version"`
}

type ReadSetVersion struct {
	BlockNum string `json:"block_num"`
	TxNum string `json:"tx_num"`
}

type WriteSet struct {
	IsDelete bool `json:"is_delete"`
	Key string `json:"key"`
	Value string `json:"value"`
}
type ChaincodeProposalPayload struct {
	TransientMap interface{} `json:"TransientMap"`
	Input Input `json:"input"`
}

type Input struct {
	ChaincodeSpec ChaincodeSpec `json:"chaincode_spec"`
}

type ChaincodeSpec struct {
	ChaincodeId ChaincodeId `json:"chaincode_id"`
	ChaincodeInput ChaincodeInput `json:"input"`
	Timeout int64 `json:"timeout"`
	Type string `json:"type"`
}
type ChaincodeInput struct {
	Args []string `json:"args"`
	Decorations interface{} `json:"decorations"`
}













