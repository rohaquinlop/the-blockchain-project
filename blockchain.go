package main

import (
	"log"

	"go.mills.io/bitcask/v2"
)

const dbFile = "blockchain.db"

type Blockchain struct {
	tip []byte
	db  bitcask.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          bitcask.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	data, err := i.db.Get(i.currentHash)
	if err != nil {
		log.Fatalf("error getting block: %v", err)
	}

	block = DeserializeBlock(data)
	i.currentHash = block.PrevBlockHash

	return block
}

func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	lastHash, err := bc.db.Get([]byte("l"))
	if err != nil {
		log.Fatalf("error getting last hash: %v", err)
	}

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Put(newBlock.Hash, newBlock.Serialize())
	if err != nil {
		log.Fatalf("error putting new block: %v", err)
	}

	err = bc.db.Put([]byte("l"), newBlock.Hash)
	if err != nil {
		log.Fatalf("error putting last hash: %v", err)
	}

	bc.tip = newBlock.Hash
}

func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bitcask.Open(dbFile)

	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}

	_, err = db.Get([]byte("l"))

	if err != nil {
		genesis := NewGenesisBlock()
		err = db.Put(genesis.Hash, genesis.Serialize())

		if err != nil {
			log.Fatalf("error putting genesis block: %v", err)
		}

		err = db.Put([]byte("l"), genesis.Hash)

		if err != nil {
			log.Fatalf("error putting last hash: %v", err)
		}

		tip = genesis.Hash
	} else {
		tip, _ = db.Get([]byte("l"))
	}

	bc := Blockchain{tip, db}

	return &bc
}
