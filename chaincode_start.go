package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
    
    err := stub.PutState("init value", []byte(args[0]))
    if err != nil {
        return nil, err
    }

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "winChallenge" {
        return t.winChallenge(stub, args)
    }
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation")
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {											//read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query")
}

// args[]  team, code
func (t *SimpleChaincode) winChallenge(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var key, value string
    var err error
    fmt.Println("running winChallenge()")

    if len(args) != 2 {
        return []byte("ERROR"), errors.New("Incorrect number of arguments. Expecting 2. team validation code")
    }

    key = args[0]                            
    value = "ibmbluemix"

    if(args[1]!="PS12094") {
        err = stub.PutState(key, []byte("Wrong validation code!"))  
        return nil, nil
    }
    err = stub.PutState(key, []byte(value))  
    if err != nil {
        return nil, err
    }
    return []byte("ibmbluemix"), nil   // ibmbluemix
}

func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var key, jsonResp string
    var err error

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    key = args[0]
    valAsbytes, err := stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }

    return valAsbytes, nil
}