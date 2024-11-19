cd fabric-samples/chaincode-docker-devmode

#terminal1
docker-compose -f docker-compose-simple.yaml down & docker-compose -f docker-compose-simple.yaml up
#terminal2
docker exec -it chaincode bash
#这里的目录是映射了本地 "../chaincode" 的目录
cd sacc
go build
#运行链码
CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=mycc:0 ./chaincode

#terminal3
docker exec -it cli bash
#安装与实例化链码
peer chaincode install -p chaincodedev/chaincode -n mycc -v 0
peer chaincode instantiate -n mycc -v 0 -c '{"Args":["test","ok"]}' -C myc
#测试链码
peer chaincode query -n mycc -c '{"Args":["queryVideo","test"]}' -C myc

# 注册用户/机构
peer chaincode invoke -n mycc -c '{"Args":["registerUser","alice","alicepk"]}' -C myc
peer chaincode invoke -n mycc -c '{"Args":["registerUser","bob","bobpk"]}' -C myc
peer chaincode invoke -n mycc -c '{"Args":["registerUser","cev","cevpk"]}' -C myc

#查询用户
peer chaincode query -n mycc -c '{"Args":["queryUser","alicepk"]}' -C myc

#作品上链
peer chaincode invoke -n mycc -c '{"Args":["videoUplink","alicepk","bobpk","CID123","1.mp4"]}' -C myc
peer chaincode invoke -n mycc -c '{"Args":["videoUplink","cevpk","bobpk","CID345","2.mp4"]}' -C myc

# 查询上链信息
peer chaincode query -n mycc -c '{"Args":["queryUplink","alicepk"]}' -C myc
peer chaincode query -n mycc -c '{"Args":["queryUplink","bobpk"]}' -C myc
peer chaincode query -n mycc -c '{"Args":["queryUplink","cevpk"]}' -C myc


#查询作品信息
peer chaincode query -n mycc -c '{"Args":["queryVideo","CID123"]}' -C myc
peer chaincode query -n mycc -c '{"Args":["queryVideo","CID345"]}' -C myc




