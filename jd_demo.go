/*
	author:liulinlin
	emial:liulinlin@daixiaomi.com
	MIT License
*/

package main

import (
	"errors"
	"fmt"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	//"time"
	//"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/hyperledger/fabric/gossip/proto"
)

type SimpleChaincode struct {
}

//var BackGroundNo int = 0
//var RecordNo int = 0

type Assert struct{
	Id string
	Status string
	IssueTime string
	Owner string
	AddInfo string
}


//func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) ([]byte, error) {
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	//function, args := stub.GetFunctionAndParameters()
	//if function == "createSchool"{
	//	return t.createSchool(stub,args)
	//}else if function == "createStudent"{
	//	return t.createStudent(stub,args)
	//}
	fmt.Printf("deploy code success and do nothing")
	return nil, nil
}

//func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) ([]byte, error) {
	//function, args := stub.GetFunctionAndParameters()
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "createAssert"{
		if len(args)!= 2{
			return nil, errors.New("Incorrect number of arguments. Expecting 2")
		}
		return t.createAssert(stub,args)
	}else if function == "updateAssertStatus"{
		if len(args)!= 2{
			return nil, errors.New("Incorrect number of arguments. Expecting 2")
		}
		return t.updateAssertStatus(stub,args)
	}
	return nil,nil
}

//func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface) ([]byte, error) {
//	function, args := stub.GetFunctionAndParameters()
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "getAssertByAddress"{
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting 1")
		}
		_,stuBytes, err := getAssertByAddress(stub,args[0])
		if err != nil {
			fmt.Println("Error get centerBank")
			return nil, err
		}
		return stuBytes, nil
	}
	return nil,nil
}



//生成Address
func GetAddress() (string,string,string) {
	var address,priKey,pubKey string
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "","",""
	}

	h := md5.New()
	h.Write([]byte(base64.URLEncoding.EncodeToString(b)))

	address = hex.EncodeToString(h.Sum(nil))
	priKey = address+"1"
	pubKey = address+"2"
	return address,priKey,pubKey
}


func (t *SimpleChaincode) createAssert(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2{
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	var assert Assert
	var schoolBytes []byte
	var address string
	address,_,_ = GetAddress()

	assert = Assert {Id:address,Owner:args[0],IssueTime:args[1],AddInfo:""}
	err := writeAssert(stub,assert)
	if err != nil{
		return nil, errors.New("write Error" + err.Error())
	}

	schoolBytes ,err = json.Marshal(&assert)
	fmt.Printf("Assert ID=%s\n", assert.Id)
	if err!= nil{
		return nil,errors.New("Error retrieving assertBytes")
	}

	return schoolBytes,nil
}


func (t *SimpleChaincode) updateAssertStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var recordBytes []byte
	assert,_,err:=getAssertByAddress(stub,args[0])
	if err != nil{
		return nil,errors.New("Error get data")
	}
	status := args[1]
	assert.Status = status
	err = writeAssert(stub,assert)
	if err != nil{
		return nil,errors.New("Error write data")
	}

	return recordBytes,nil
}


func getAssertByAddress(stub shim.ChaincodeStubInterface,address string)(Assert,[]byte,error){
	var assert Assert
	schBytes,err := stub.GetState(address)
	if err != nil{
		fmt.Println("Error retrieving data")
	}

	err = json.Unmarshal(schBytes,&assert)
	if err != nil{
		fmt.Println("Error unmarshalling data")
	}
	return assert,schBytes,nil
}


func writeAssert(stub shim.ChaincodeStubInterface,assert Assert)(error){
	schBytes ,err := json.Marshal(&assert)
	if err != nil{
		return err
	}

	err = stub.PutState(assert.Id,schBytes)
	if err !=nil{
		return errors.New("PutState Error" + err.Error())
	}
	fmt.Printf("Assert has been written into world stats:%s\n", schBytes)
	return nil
}


func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
