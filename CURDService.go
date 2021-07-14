package main

import (
	"github.com/lifegoeson/blockchain-explorer/common"
	"github.com/lifegoeson/blockchain-explorer/model"
	"log"
	"time"
)

//根据创世hash以及区块号获取交易数量
func getTxCountByBlockNum(channel_genesis_hash string, blockNum int64) model.GetTxCountByBlockNumResultModel{
	//更新数据
	stmt := `select blocknum ,txcount from blocks where channel_genesis_hash=$1 and blocknum=$2;`
	row ,err:= db.Query(stmt, channel_genesis_hash,blockNum)
	common.CheckErr(err)
	var (
		blocknum  int64
		txcount int64
	)
	for row.Next() {
		if err := row.Scan(&blocknum, &txcount); err != nil {
			log.Fatal(err)
		}
	}
	return model.GetTxCountByBlockNumResultModel{BlockNum: blocknum, TxCount: txcount}
}

//根据交易ID获取交易信息
func getTransactionByID(txhash string) model.GetTransactionByIDResultModel{
	stmt := ` select t.txhash,t.validation_code,t.payload_proposal_hash,t.creator_msp_id,t.endorser_msp_id,t.chaincodename,t.type,t.createdt,t.read_set,
        t.write_set,channel.name as channelName from transactions as t inner join channel on t.channel_genesis_hash=channel.channel_genesis_hash where t.txhash = $1 ;`
	row ,err:= db.Query(stmt,txhash)
	common.CheckErr(err)
	var tx model.Transaction
	var channelName string
	for row.Next() {
		if err := row.Scan(&tx.TxHash,&tx.ValidationCode,&tx.PayloadProposalHash,&tx.CreatorMspId,&tx.EndorserMspId,&tx.ChaincodeName,&tx.Type,&tx.CreateAt,&tx.ReadSet,&tx.WriteSet,&channelName); err != nil {
			log.Fatal(err)
		}
	}
	return model.GetTransactionByIDResultModel{Tx: tx, ChannelName: channelName}
}

//获取channel最近的3个block
func getBlockActivityList(channelGenesisHash string) []model.GetBlockActivityListResultModel{
	stmt := `select blocks.blocknum,blocks.txcount ,blocks.datahash ,blocks.blockhash ,blocks.prehash,blocks.createdt,(
      SELECT  array_agg(txhash) as txhash FROM transactions where blockid = blocks.blocknum and
       channel_genesis_hash = $1 group by transactions.blockid ),
      channel.name as channelname  from blocks inner join channel on blocks.channel_genesis_hash = channel.channel_genesis_hash  where
       blocks.channel_genesis_hash = $1 and blocknum >= 0
       order by blocks.blocknum desc limit 3;`
	row ,err:= db.Query(stmt, channelGenesisHash)
	common.CheckErr(err)
	var blks []model.GetBlockActivityListResultModel
	var blk model.GetBlockActivityListResultModel
	for row.Next() {
		if err := row.Scan(&blk.BlockNum,&blk.TxCount,&blk.DataHash,&blk.BlockHash,&blk.PreHash,&blk.CreateAt,&blk.TxHash,&blk.ChannelName); err != nil {
			log.Fatal(err)
		}
		blks = append(blks, blk)
	}
	return blks
}
//Returns the list of transactions by channel, organization, date range and greater than a block and transaction id.
func getTxList(channelGenesisHash string , blockNum int64, txId string , from time.Time, to time.Time, organizations string) []model.GetTxListResultModel{
	txListSql := ""
	if len(organizations) != 0 {
		txListSql = "and t.creator_msp_id in ("+"'"+organizations+"')"
	}
	queryText := `select t.creator_msp_id,t.txhash,t.type,t.chaincodename,t.createdt,channel.name as channelName from transactions as t
       inner join channel on t.channel_genesis_hash=channel.channel_genesis_hash where t.blockid >= $2 and t.id >= $3 `+txListSql+`and
       t.channel_genesis_hash = $1  and t.createdt between $4 and $5  order by  t.id desc;`
	rows ,err:= db.Query(queryText, channelGenesisHash,blockNum,txId,from,to,txListSql)
	common.CheckErr(err)
	var txListResultModels []model.GetTxListResultModel
	for rows.Next() {
		var tx model.GetTxListResultModel
		if err := rows.Scan(&tx.CreatorMspId,&tx.TxHash,&tx.Type,&tx.ChaincodeName,&tx.CreateAt,&tx.ChannelName); err != nil {
			log.Fatal(err)
		}
		txListResultModels = append(txListResultModels,tx)
	}
	return txListResultModels
}

func getBlockAndTxList(channelGenesisHash string, from time.Time, to time.Time, organizations string) []model.GetBlockAndTxListResultModel{
	blockAndTxList := ""
	if len(organizations) != 0 {
		blockAndTxList = "and t.creator_msp_id in ("+"'"+organizations+"')"
	}
	queryText := `select a.* from  (
      select (select c.name from channel c where c.channel_genesis_hash =
         $1 ) as channelname, blocks.blocknum,blocks.txcount ,blocks.datahash ,blocks.blockhash ,blocks.prehash,blocks.createdt,(
        SELECT  array_agg(txhash) as txhash FROM transactions where blockid = blocks.blocknum `+blockAndTxList+` and
         channel_genesis_hash = $1 and createdt between $2 and $3) from blocks where
         blocks.channel_genesis_hash =$1 and blocknum >= 0 and blocks.createdt between $2 and $3
         order by blocks.blocknum desc)  a where  a.txhash IS NOT NULL;`
	rows,err := db.Query(queryText, channelGenesisHash,from,to)
	common.CheckErr(err)
	var blockAndTxs []model.GetBlockAndTxListResultModel
	for rows.Next() {
		var blockAndTx model.GetBlockAndTxListResultModel
		if err := rows.Scan(&blockAndTx.ChannelName,&blockAndTx.BlockNum,&blockAndTx.TxHash,&blockAndTx.DataHash,&blockAndTx.BlockHash,&blockAndTx.PreHash,&blockAndTx.CreateAt,&blockAndTx.TxHash); err != nil {
			log.Fatal(err)
		}
		blockAndTxs = append(blockAndTxs,blockAndTx)
	}
	return blockAndTxs
}

func getChannelConfig(channelGenesisHash string) model.Channel{
	queryText := ` select * from channel where channel_genesis_hash = $1 `

	row := db.QueryRow(queryText, channelGenesisHash)
	var channel model.Channel
	if err := row.Scan(&channel.Id,&channel.Name,&channel.Blocks,&channel.Trans,&channel.CreateAt,&channel.ChannelGenesisHash,&channel.ChannelHash,&channel.ChannelConfig,&channel.ChannelBlock,&channel.ChannelTx,&channel.ChannelVersion); err != nil {
		log.Fatal(err)
	}
	return channel
}

func getChannel(channelName string, channelGenesisHash string) model.Channel{
	queryText := ` select * from channel where name= $1 and channel_genesis_hash=$2`
	row := db.QueryRow(queryText,channelName, channelGenesisHash)
	var channel model.Channel
	if err := row.Scan(&channel.Id,&channel.Name,&channel.Blocks,&channel.Trans,&channel.CreateAt,&channel.ChannelGenesisHash,&channel.ChannelHash,&channel.ChannelConfig,&channel.ChannelBlock,&channel.ChannelTx,&channel.ChannelVersion); err != nil {
		log.Fatal(err)
	}
	return channel
}

func existChannel(channelName string) bool {
	queryText := ` select count(1) from channel where name= $1`
	row := db.QueryRow(queryText,channelName)
	var count int
	if err := row.Scan(&count); err != nil {
		log.Fatal(err)
	}
	return count > 0
}

func saveBlock(block model.Block) bool {
	//判断区块是否存在
	queryText := `select count(1) as c from blocks where blocknum= $1 and txcount= $2 and channel_genesis_hash= $3 and prehash=$4 and datahash= $5`
	row := db.QueryRow(queryText,block.BlockNum,block.TxCount,block.ChannelGenesisHash,block.PrevBlockHash,block.DataHash)
	var count int
	if err := row.Scan(&count); err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		return false
	}
	//插入区块
	insertText := `insert into blocks(blocknum, datahash, prehash, txcount, createdt, prev_blockhash, blockhash, channel_genesis_hash)  values ($1,$2,$3,$4,$5,$6,$7,$8)`
	result, err := db.Exec(insertText,block.BlockNum,block.DataHash,block.PreHash,block.TxCount,block.CreateAt,block.PrevBlockHash,block.BlockHash,block.ChannelGenesisHash)
	common.CheckErr(err)
	if _,err = result.RowsAffected() ; err != nil {
		log.Fatal(err)
	}

	//更新channel的区块数量
	insertText2 := `update channel set blocks =blocks+1 where channel_genesis_hash=$1`
	result2, err := db.Exec(insertText2,block.ChannelGenesisHash)
	common.CheckErr(err)
	if _,err = result2.RowsAffected() ; err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func saveTransaction(tx model.Transaction) bool {
	//判断交易是否存在
	queryText := `select count(1) as c from transactions where blockid= $1 and txhash= $2 and channel_genesis_hash= $3`
	row := db.QueryRow(queryText,tx.BlockId,tx.TxHash,tx.ChannelGenesisHash)
	var count int
	if err := row.Scan(&count); err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		return false
	}
	//插入交易
	insertText := `insert into transactions(blockid, txhash, createdt, chaincodename, status, creator_msp_id, endorser_msp_id, chaincode_id, type, read_set, write_set, channel_genesis_hash, validation_code, envelope_signature, payload_extension, creator_id_bytes, creator_nonce, chaincode_proposal_input, tx_response, payload_proposal_hash, endorser_id_bytes, endorser_signature) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)`
	result, err := db.Exec(insertText,tx.BlockId,tx.TxHash,tx.CreateAt,tx.ChaincodeName,tx.Status,tx.CreatorMspId,tx.EndorserMspId,tx.ChaincodeId,tx.Type,tx.ReadSet,tx.WriteSet,tx.ChannelGenesisHash,tx.ValidationCode,tx.EnvelopeSignature,tx.PayloadExtension,tx.CreatorIdBytes,tx.CreatorNonce,tx.ChaincodeProposalInput,tx.TxResponse,tx.PayloadProposalHash,tx.EndorserIdBytes,tx.EndorserSignature)
	common.CheckErr(err)
	if _,err = result.RowsAffected() ; err != nil {
		log.Fatal(err)
	}

	//更新chaincode的交易数量
	insertText2 := `update chaincodes set txcount = txcount+1 where channel_genesis_hash=$1`
	result2, err := db.Exec(insertText2,tx.ChannelGenesisHash)
	common.CheckErr(err)
	if _,err = result2.RowsAffected() ; err != nil {
		log.Fatal(err)
		return false
	}
	//更新channel中的交易数量
	insertText3 := `update channel set trans = trans+1 where channel_genesis_hash=$1`
	result3, err := db.Exec(insertText3,tx.ChannelGenesisHash)
	common.CheckErr(err)
	if _,err = result3.RowsAffected() ; err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func getCurBlockNum(channelGenesisHash string) int64{
	queryText := `select max(blocknum) as blocknum from blocks  where channel_genesis_hash=$1`
	row := db.QueryRow(queryText, channelGenesisHash)
	var max int64
	var curBlockNum int64
	if row.Err() != nil {
		log.Fatal(row.Err())
	}
	if err := row.Scan(&max); err != nil {
		curBlockNum = -1;
	}
	curBlockNum = max
	return curBlockNum
}

func saveChaincode(chaincode model.Chaincode) bool {
	queryText := `select count(1) as c from chaincodes where name= $1 and channel_genesis_hash= $2 and version= $3 and path=$4`
	row := db.QueryRow(queryText,chaincode.Name,chaincode.ChannelGenesisHash,chaincode.Version,chaincode.Path)
	var count int64
	if err := row.Scan(&count); err != nil {
		log.Fatal(err)
	}
	if count > 0{
		return false
	}
	insertText := `insert into chaincodes(name, version, path, channel_genesis_hash, txcount, createdt) VALUES ($1,$2,$3,$4,$5,$6)`
	result, err := db.Exec(insertText,chaincode.Name,chaincode.Version,chaincode.Path,chaincode.ChannelGenesisHash,chaincode.TxCount,chaincode.CreateAt)
	common.CheckErr(err)
	if _,err = result.RowsAffected() ; err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func getChannelByGenesisBlockHash(channel_genesis_hash string) string{
	queryText := `select name from channel where channel_genesis_hash=$1`
	row := db.QueryRow(queryText,channel_genesis_hash)
	var name string
	if err := row.Scan(&name); err != nil {
		log.Fatal(err)
	}
	return name
}

func saveChaincodPeerRef(prc model.PeerRefChaincode) bool{
	queryText := `select count(1) as c from peer_ref_chaincode prc where prc.peerid= $1 and prc.chaincodeid=$2 and cc_version= $3 and channelid=$4`
	row := db.QueryRow(queryText,prc.PeerId,prc.ChaincodeId,prc.CCVersion,prc.ChannelId)
	var count int64
	if err := row.Scan(&count); err != nil {
		log.Fatal(err)
	}
	if count >0 {
		return false
	}

	insertText := `insert into peer_ref_chaincode(peerid, chaincodeid, cc_version, channelid, createdt) VALUES ($1,$2,$3,$4,$5)`
	result, err := db.Exec(insertText,prc.PeerId,prc.ChaincodeId,prc.CCVersion,prc.ChannelId,prc.CreateAt)
	common.CheckErr(err)
	if _,err = result.RowsAffected() ; err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func saveChannel(channel model.Channel) bool {
	queryText := `select count(1) as c from channel where name= $1 and channel_genesis_hash=$2`

	row := db.QueryRow(queryText,channel.Name,channel.ChannelGenesisHash)
	var count int64
	if err := row.Scan(&count); err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		insertText := `insert into channel(name, blocks, trans, createdt, channel_genesis_hash, channel_hash, channel_config, channel_block, channel_tx, channel_version) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
		result, err := db.Exec(insertText,channel.Name,channel.Blocks,channel.Trans,channel.CreateAt,channel.ChannelGenesisHash,channel.ChannelHash,channel.ChannelConfig,channel.ChannelBlock,channel.ChannelTx,channel.ChannelVersion)
		common.CheckErr(err)
		if _,err = result.RowsAffected() ; err != nil {
			log.Fatal(err)
			return false
		}
		return true
	}
	insertText := `update channel set blocks = $1 ,trans = $2,channel_hash=$3 where name=$4 and channel_genesis_hash=$5`
	result, err := db.Exec(insertText,channel.Blocks,channel.Trans,channel.ChannelHash,channel.Name,channel.ChannelGenesisHash)
	common.CheckErr(err)
	if _,err = result.RowsAffected() ; err != nil {
		log.Fatal(err)
		return false
	}
	return true

}

func savePeer(peer model.Peer) bool{
	queryText := `select count(1) as c from peer where channel_genesis_hash=$1 and server_hostname=$2 `
	row := db.QueryRow(queryText,peer.ChannelGenesisHash,peer.ServerHostName)
	var count int64
	if err := row.Scan(&count); err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		return false
	}
	insertText := `insert into peer(org, channel_genesis_hash, mspid, requests, events, server_hostname, createdt, peer_type) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	result, err := db.Exec(insertText,peer.Org,peer.ChannelGenesisHash,peer.MspId,peer.Requests,peer.Events,peer.ServerHostName,peer.CreateAt,peer.PeerType)
	common.CheckErr(err)
	if _,err = result.RowsAffected() ; err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func savePeerChannelRef(prc model.PeerRefChannel) bool {
	queryText := `select count(1) as c from peer_ref_channel prc where prc.peerid = $1 and prc.channelid= $2 `
	row := db.QueryRow(queryText,prc.PeerId,prc.ChannelId)
	var count int64
	if err := row.Scan(&count); err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		return false
	}
	insertText := `insert into peer_ref_channel(createdt, peerid, channelid, peer_type) VALUES  ($1,$2,$3,$4)`
	result, err := db.Exec(insertText,prc.CreateAt,prc.PeerId,prc.ChannelId,prc.PeerType)
	common.CheckErr(err)
	if _,err = result.RowsAffected() ; err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func getChannelsInfo(peerid string) []model.Channel{
	queryText := ` select c.id as id,c.name as channelName,c.blocks as blocks ,c.channel_genesis_hash as channel_genesis_hash,c.trans as transactions,c.createdt as createdat,c.channel_hash as channel_hash from channel c,
        peer_ref_channel pc where c.channel_genesis_hash = pc.channelid and pc.peerid= $1 group by c.id ,c.name ,c.blocks  ,c.trans ,c.createdt ,c.channel_hash,c.channel_genesis_hash order by c.name `
	rows,err := db.Query(queryText,peerid)
	common.CheckErr(err)
	var chls []model.Channel
	var chl model.Channel
	for rows.Next() {
		if err := rows.Scan(&chl); err != nil {
			log.Fatal(err)
		}
		chls = append(chls,chl)
	}
	return chls
}

func saveOrderer(orderer model.Orderer) bool {
	queryText := `select count(1) as c from orderer where requests= $1`
	row := db.QueryRow(queryText,orderer.Requests)
	var count int64
	if err := row.Scan(&count); err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		return false
	}

	insertText := `insert into orderer(requests, server_hostname, createdt) VALUES  ($1,$2,$3)`
	result, err := db.Exec(insertText,orderer.Requests,orderer.ServerHostName,orderer.CreateAt)
	common.CheckErr(err)
	if _,err = result.RowsAffected() ; err != nil {
		log.Fatal(err)
		return false
	}
	return true
}










