package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// 作品版权信息存证：用户信息、机构信息、作品IPFSCID、作品名称、时间戳、txID
// 用户/机构注册：机构的RSA公钥、机构名称

// 查询已上链的作品版权存证信息：输入用户/机构公钥，返回存证列表
// 验证作品上链信息：输入作品CID，查询到用户信息、机构信息、作品IPFSCID、作品名称、时间戳、txID

// 定义合约结构体
type SmartContract struct {
}

// 定义用户结构体
type User struct {
	UserRSAPK string   `json:"userRSAPK"`
	UserName  string   `json:"userName"`
	Videos    []*Video `json:"video"`
}

// 定义作品结构体
type Video struct {
	Owner     string `json:"owner"`
	OrgInfo   string `json:"orgInfo"`
	VideoCID  string `json:"videoCID"`
	VideoName string `json:"videoName"`
	Txid      string `json:"txid"`
	Timestamp string `json:"timestamp"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {

	// 获取调用合约的参数
	args := APIstub.GetStringArgs()
	if len(args) != 2 {
		return shim.Error("Incorrect arguments. Expecting a key and a value")
	}
	// 测试 Putstate 是否正常
	err := APIstub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	}
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// 获取调用合约的函数与参数名称
	function, args := APIstub.GetFunctionAndParameters()
	// 根据调用合约函数名称路由到相应的函数中
	if function == "registerUser" {
		return s.registerUser(APIstub, args)
	} else if function == "queryUser" {
		return s.queryUser(APIstub, args)
	} else if function == "videoUplink" {
		return s.videoUplink(APIstub, args)
	} else if function == "queryVideo" {
		return s.queryVideo(APIstub, args)
	} else if function == "queryUplink" {
		return s.queryUplink(APIstub, args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

// 注册用户
func (s *SmartContract) registerUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	// 验证函数参数个数
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	// 实例化结构体
	var user User
	user.UserName = args[0]
	user.UserRSAPK = args[1]
	// 转化为json格式字符串
	userjson, err := json.Marshal(user)
	if err != nil {
		return shim.Error("json marshal failed")
	}
	// 存储到区块链账本
	err = APIstub.PutState(args[1], userjson)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(userjson)
}

// 查询用户名，传入参数为用户的公钥
func (s *SmartContract) queryUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	// 验证参数个数
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	// 获取用户信息
	userjson, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error("GetState failed")
	}
	var user User
	json.Unmarshal(userjson, &user)
	userName := user.UserName
	return shim.Success([]byte(userName))
}

// 查询用户/机构存证的作品
func (s *SmartContract) queryUplink(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	var user User
	userRecordsAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	json.Unmarshal(userRecordsAsBytes, &user)
	videos := user.Videos
	videosAsBytes, err := json.Marshal(videos)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(videosAsBytes)
}

// 查询作品的信息
func (s *SmartContract) queryVideo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	videoAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(videoAsBytes)
}

// 作品版权存证 传入参数有：作品所有人信息、上传的组织信息、作品IPFS CID、作品名称四个参数
func (s *SmartContract) videoUplink(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//验证函数参数
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	//获取时间戳
	txtime, err := APIstub.GetTxTimestamp()
	if err != nil {
		return shim.Error(err.Error())
	}
	timesecond := txtime.Seconds
	time := time.Unix(timesecond+28800, 0).Format("2006-01-02 15:04:05")

	// 作品信息实例化
	var video = Video{Owner: args[0], OrgInfo: args[1], VideoCID: args[2], VideoName: args[3], Timestamp: time, Txid: APIstub.GetTxID()}
	// fmt.Printf("作品结构体信息：%v\n", video)

	// 将作品信息存储至用户的信息里
	var user User
	userRecordsAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	json.Unmarshal(userRecordsAsBytes, &user)
	user.Videos = append(user.Videos, &video)
	userRecordsAsBytes, err = json.Marshal(user)
	if err != nil {
		return shim.Error(err.Error())
	}
	APIstub.PutState(args[0], userRecordsAsBytes)
	// fmt.Printf("用户信息：%v\n", string(userRecordsAsBytes))

	// 将作品信息存储至机构的信息
	// 将信息查出来
	var org User
	userRecordsAsBytes, err = APIstub.GetState(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	//解析到结构体User
	json.Unmarshal(userRecordsAsBytes, &org)
	//添加新上链的作品信息
	org.Videos = append(org.Videos, &video)
	userRecordsAsBytes, err = json.Marshal(org)
	if err != nil {
		return shim.Error(err.Error())
	}
	APIstub.PutState(args[1], userRecordsAsBytes)
	// fmt.Printf("机构信息：%v\n", string(userRecordsAsBytes))

	// 存储作品信息
	videoAsBytes, err := json.Marshal(video)
	if err != nil {
		return shim.Error(err.Error())
	}

	APIstub.PutState(args[2], videoAsBytes)
	// fmt.Printf("作品信息：%v\n", string(videoAsBytes))

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
