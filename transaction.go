package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const subsidy = 10

// Note:
// 1. There are outputs that are not linked to inputs.
// 2. In one transaction, inputs can reference outputs from
//    multiple transactions.
// 3. An input must reference an output.

// TXOutput is an output for scripts.
type TXOutput struct {
	Value        int
	ScriptPubKey string
}

// TXInput is an input for scripts.
type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

// Transaction is a struct representing a transaction.
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}
