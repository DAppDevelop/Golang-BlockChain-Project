package BLC

//常量值
const DBName = "blockchain.db"   //数据库的名字
const BlockBucketName = "blocks" //定义bucket
const targetBit = 8              // 挖矿难度(256位Hash里面前面至少要有16个零)
const UTXOSetBucketName = "utxoset"