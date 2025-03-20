package main

import (
    "fmt"
    "strconv"
    "github.com/hyperledger/fabric-chaincode-go/shim"
    pb "github.com/hyperledger/fabric-protos-go/peer"
)

type BatteryChaincode struct {
}

// Init initializes the chaincode
func (t *BatteryChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
    return shim.Success(nil)
}

// Invoke handles transactions
func (t *BatteryChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
    function, args := stub.GetFunctionAndParameters()

    if function == "reportBattery" {
        return t.reportBattery(stub, args)
    }
    return shim.Error("Invalid function name")
}

// reportBattery logs a robot's battery level in mV
func (t *BatteryChaincode) reportBattery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    // Check argument count
    if len(args) != 3 {
        return shim.Error("Expected 3 args: robotID, batteryLevelMV, timestamp")
    }

    robotID := args[0]
    batteryLevelMVStr := args[1]
    timestamp := args[2]

    // Validate battery level (0 to 8500 mV)
    batteryLevelMV, err := strconv.ParseFloat(batteryLevelMVStr, 64)
    if err != nil || batteryLevelMV < 0 || batteryLevelMV > 8500 {
        return shim.Error("Battery level must be a number between 0 and 8500 mV")
    }

    // Check previous value for plausibility
    prevData, err := stub.GetState(robotID)
    if err == nil && len(prevData) > 0 {
        var prevBatteryMV float64
        fmt.Sscanf(string(prevData), "Robot %s battery: %f mV at", robotID, &prevBatteryMV)
        if batteryLevelMV > prevBatteryMV+500 { // Max jump of 500 mV
            return shim.Error("Battery increase too large")
        }
    }

    // Format the data to store
    batteryData := fmt.Sprintf("Robot %s battery: %.2f mV at %s", robotID, batteryLevelMV, timestamp)

    // Store in the ledger
    err = stub.PutState(robotID, []byte(batteryData))
    if err != nil {
        return shim.Error(err.Error())
    }

    return shim.Success([]byte("Battery level reported successfully"))
}

// Add a query battery function
func (t *BatteryChaincode) queryBattery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) != 1 {
        return shim.Error("Expected 1 arg: robotID")
    }
    robotID := args[0]
    batteryData, err := stub.GetState(robotID)
    if err != nil {
        return shim.Error(err.Error())
    }
    if batteryData == nil {
        return shim.Error("No battery data for this robot")
    }
    return shim.Success(batteryData)
}

// Update Invoke to handle queryBattery
func (t *BatteryChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
    function, args := stub.GetFunctionAndParameters()
    if function == "reportBattery" {
        return t.reportBattery(stub, args)
    } else if function == "queryBattery" {
        return t.queryBattery(stub, args)
    }
    return shim.Error("Invalid function name")
}

func main() {
    err := shim.Start(new(BatteryChaincode))
    if err != nil {
        fmt.Printf("Error starting chaincode: %s", err)
    }
}