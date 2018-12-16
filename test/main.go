package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/ageapps/Peerster/pkg/file"

	"github.com/ageapps/Peerster/pkg/chain"

	"github.com/ageapps/Peerster/pkg/data"
	"github.com/ageapps/Peerster/pkg/logger"
	"github.com/ageapps/Peerster/pkg/utils"
)

type Person struct {
	Name string
	age  int
}

type House struct {
	Person *Person
	street string
}

type Family struct {
	Person   *Person
	children int
}

func NewFamily(person *Person, child int) *Family {
	return &Family{
		Person:   person,
		children: child,
	}
}

func testHashValue() {
	var pepe utils.HashValue
	pepe.Set("0c515910c21c81b00d899705c2da2afc70db2d0c5b29d4293f5e698fd5afa5c0")
	fmt.Println(pepe.String())
}
func testFiles() {
	logger.CreateLogger("file", "0.0.0.0", true)
	// f, err := file.NewFileFromLocalSync("test.png")
	// err = f.Reconstruct()
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

func testPointers() {
	house := &House{
		Person: &Person{Name: "adri", age: 20},
		street: "calle",
	}

	family := NewFamily(house.Person, 2)

	fmt.Println(*house.Person)

	fmt.Println(*family.Person)

	house.Person.Name = "Pepe"
	fmt.Println(*house.Person)
	fmt.Println(*family.Person)

}

func createFile(name string) *file.File {
	b, err := file.NewBlobFromLocalSync(name)
	if err != nil {
		fmt.Println(err)
	}
	f := file.NewFile(b.GetName(), b.GetBlobSize(), b.GetMetaHash())
	return f
}
func createTx(name string) *data.TxPublish {
	f := createFile(name)
	return data.NewTXPublish(*f, uint32(10))
}

func addTx(bc *chain.BlockChain, tx *data.TxPublish) {
	if !bc.IsTransactionSaved(tx) {
		bc.TxChannel <- tx
	} else {
		logger.Logf("Transaction for %v already indexed", tx.File.Name)
	}
}
func addBlock(bc *chain.BlockChain, bl *data.Block) {
	if bc.CanAddBlock(bl) {
		bc.BlockChannel <- bl
	} else {
		logger.Logf("Block for %v already indexed", bl.PrintPrev())
	}
}
func createBlock(fake string) ([32]byte, *data.Block) {
	hashT, _ := hex.DecodeString(fake)
	hash := [32]byte{}
	for index := 0; index < len(hash); index++ {
		hash[index] = hashT[index]
	}
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	randIndex := r1.Int() % 32
	nonce := hash
	nonce[randIndex] = 1
	return nonce,
		&data.Block{
			PrevHash: hash,
			Nonce:    nonce,
		}
}

func main() {
	logger.CreateLogger("file", "0.0.0.0", true)
	bc := chain.NewBlockChain()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		bc.Start(func() {
			fmt.Println("Stopped")
		})
	}()
	addTx(bc, createTx("test1.png"))
	addTx(bc, createTx("test2.png"))

	time.Sleep(time.Second * 5)

	_, b := createBlock("00000d87b29b25e2c9dd794f8884454c6d27696fb5421430ce9de0566bbe418d")
	addBlock(bc, b)
	_, b2 := createBlock("00000d87b29b25e2c9dd794f8884454c6d27696fb5421430ce9de0566bbe418d")
	addBlock(bc, b2)
	_, b3 := createBlock(b2.String())
	addBlock(bc, b3)
	addTx(bc, createTx("test3.png"))
	// addTx(bc, createTx("test4.png"))
	// addTx(bc, createTx("test5.png"))
	// addTx(bc, createTx("test6.png"))

	time.Sleep(time.Second * 5)
	bc.Stop()
	wg.Wait()
	fmt.Println("Stop")

}
