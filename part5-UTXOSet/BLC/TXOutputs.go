package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

type TxOutputs struct {
	UTXOs []*UTXO
}

func (outs *TxOutputs) Serialize() []byte {
	var buff bytes.Buffer

	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(outs)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func DeserializeTxOutputs(data []byte) *TxOutputs  {
	outs := TxOutputs{}

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&outs)

	if err != nil {
		log.Panic(err)
	}

	return &outs
}