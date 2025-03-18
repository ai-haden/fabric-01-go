package main

import (
    "fmt"
    "github.com/hyperledger/fabric-chaincode-go/shim"
    pb "github.com/hyperledger/fabric-protos-go/peer"
)

type RobotChaincode struct {
}

// Init initializes the chaincode
func (t *RobotChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
    return shim.Success(nil)
}

// Invoke handles transactions
func (t *RobotChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
    function, args := stub.GetFunctionAndParameters()

    if function == "logTask" {
        return t.logTask(stub, args)
    }
    return shim.Error("Invalid function name")
}

// logTask logs a robot's task completion
func (t *RobotChaincode) logTask(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) != 3 {
        return shim.Error("Expected 3 args: robotID, taskID, timestamp")
    }
    robotID := args[0]
    taskID := args[1]
    timestamp := args[2]

    // Store the task in the ledger
    taskData := fmt.Sprintf("Robot %s completed task %s at %s", robotID, taskID, timestamp)
    err := stub.PutState(robotID, []byte(taskData))
    if err != nil {
        return shim.Error(err.Error())
    }
    return shim.Success([]byte("Task logged successfully"))
}

func main() {
    err := shim.Start(new(RobotChaincode))
    if err != nil {
        fmt.Printf("Error starting chaincode: %s", err)
    }
}