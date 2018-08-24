package BLC

//常量值
const DBName = "blockchain_%s.db" //数据库的名字
const BlockBucketName = "blocks" //定义bucket
const targetBit = 8              // 挖矿难度(256位Hash里面前面至少要有16个零)
const UTXOSetBucketName = "utxoset"
const walletsFile = "Wallets_%s.dat"//存储钱包数据的本地文件名
const txPollFile = "TxsPool_%s.dat"//本地交易池

//网络相关
const NODE_VERSION = 1    //版本
const COMMAND_LENGTH = 12 //命令长度[]byte
const BLOCK_TYPE = "BLOCK_TYPE"
const TX_TYPE = "TX_TYPE"

//具体的命令
const COMMAND_VERSION = "version"
const COMMAND_GETBLOCKS = "getblocks"
const COMMAND_INV = "inv"
const COMMAND_GETDATA = "getdata"
const COMMAND_BLOCKDATA = "blockdata"
const COMMAND_TXS = "transactions"
const COMMAND_REQUIREMINE = "requiremine"
const COMMAND_VERIFYBLOCK  = "verifyblock"

//钱包
const version = byte(0x00)
const addressCheckSumLen = 4


