# 使用说明


* 分别设置3个节点的id为  3000（主节点） 、3001（普通节点）、 3002（挖矿节点）

`export NODE_ID=3000``export NODE_ID=3001``export NODE_ID=3002`

* 主节点钱包：
	1. 1XtLrwjcCnaBfE3Hwuypzchnsz7PKQLxnDyfba67cBkmXG1XYa
	2. 12Mqbd1ALZQ9k4kqDvuYkCJgcSgAJNdRLVz7eWx43Ddq5jScG3p
	3. 1fQVvtoGtUeBNj2VvfpVL4ouB2jQgW5LWbzLh5LCYwpX4GXgDJ

* 普通节点钱包：
	1. 1RAbXZJVTYvPfdrXRd274vjrBE6XWxeyMLSxNYZixsuT7Uetrc
	2. 121whDSbqnDBCRUo6M47QnPz4aJLEAYpJMdp7KrenXiGKYBdh5x

* 创世区块已经给了121whDSbqnDBCRUo6M47QnPz4aJLEAYpJMdp7KrenXiGKYBdh5x 10个币

* 所有区块已经设置默认的挖矿奖励地址1XtLrwjcCnaBfE3Hwuypzchnsz7PKQLxnDyfba67cBkmXG1XYa（参考coinbase_300*文件| coinbase命令）

## 操作步骤：
1、启动主节点和挖矿节点 `./bc startnode`

2、普通本地直接挖矿交易（这个可以正常交易）
`send -from '["121whDSbqnDBCRUo6M47QnPz4aJLEAYpJMdp7KrenXiGKYBdh5x"]' -to '["1RAbXZJVTYvPfdrXRd274vjrBE6XWxeyMLSxNYZixsuT7Uetrc"]' -amount '["1"]' `

3、 挖矿节点交易（普通交易添加 -mine f）
`send -from '["121whDSbqnDBCRUo6M47QnPz4aJLEAYpJMdp7KrenXiGKYBdh5x"]' -to '["1RAbXZJVTYvPfdrXRd274vjrBE6XWxeyMLSxNYZixsuT7Uetrc"]' -amount '["1"]' -mine f`

问题：本地直接挖矿交易时，验证签名可以正常通过。但是进行挖矿节点挖矿，挖矿节点运行`func (tx *Transaction) NewTxID() []byte`时，hash出来的TxID跟普通节点签名时的TxID不一致，个人分析应该是这个问题导致验证签名失败的原因。打印hash之前的交易副本内容，签名处和验证签名处看上去是完全一致的。找不出哪里出了问题。。










	