package controller

import (
	"fmt"
	"net/http"
	"os"
	"server/api"
	"server/api/rsa"
	"server/model"

	"github.com/gin-gonic/gin"
)

func GenerateKeyPair(c *gin.Context) {
	sk, pk := api.GenerateKeyPair()
	//sk:RSA私钥 pk：RSA公钥
	c.JSON(http.StatusOK, gin.H{"sk": sk, "pk": pk})

}
func RestoreKey(c *gin.Context) {
	//将私钥转化为公钥
	pk, err := rsa.SktoPub(c.PostForm("sk"))
	if err != nil {
		c.String(http.StatusOK, "no sk input")
		return
	}
	c.String(http.StatusOK, pk)

}

func Registeruser(c *gin.Context) {
	//将用户的私钥转化为公钥
	pk, err := rsa.SktoPub(c.PostForm("sk"))
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	userName := c.PostForm("userName")
	// fmt.Printf("用户公钥：%v", pk)
	//将用户信息上传至Fabric
	var args [][]byte
	args = append(args, []byte(userName))
	args = append(args, []byte(pk))
	_, err = api.ChannelExecute("registerUser", args)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("execute chaincode failed err:%v", err))
		return
	}
	c.String(http.StatusOK, "注册成功！")

}

func Searchuser(c *gin.Context) {
	//获取表单中的orgpk参数
	pk := c.PostForm("pk")
	//将用户信息上传至Fabric
	var args [][]byte
	args = append(args, []byte(pk))
	res, err := api.ChannelExecute("queryUser", args)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("execute chaincode failed err:%v", err))
		return
	}
	c.String(http.StatusOK, string(res.Payload))

}

func Upload(c *gin.Context) {
	//先将文件保存到本地
	file, _ := c.FormFile("file")
	err := c.SaveUploadedFile(file, fmt.Sprintf("./files/uploadfiles/%v", file.Filename))
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file upload failed,err:%v", err))
		return
	}
	// 将文件上传到IPFS，返回一个IPFS CID
	cid := api.IpfsAdd(file.Filename)
	// cid := "testCID"

	//将表单中的数据绑定到video结构体中
	var video model.Video
	c.ShouldBind(&video)
	pk, err := rsa.SktoPub(video.Owner)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("转化私钥失败:%v", err))
		return
	}
	video.VideoName = file.Filename

	// 作品版权存证 传入参数有：作品所有人信息、上传的组织信息、作品IPFS CID、作品名称四个参数
	//upload to fabric
	var args [][]byte
	args = append(args, []byte(pk))
	args = append(args, []byte(video.OrgInfo))
	args = append(args, []byte(cid))
	args = append(args, []byte(video.VideoName))
	res, err := api.ChannelExecute("videoUplink", args)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("execute chaincode failed err:%v", err))
		return
	}

	//return txid
	c.JSON(http.StatusOK, gin.H{
		"txid": res.TransactionID,
	})

}

func QueryUplink(c *gin.Context) {
	pk, err := rsa.SktoPub(c.PostForm("sk"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	//查询链码
	var args [][]byte
	args = append(args, []byte(pk))
	res, err := api.ChannelQuery("queryUplink", args)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("execute chaincode failed err:%v", err))
		return
	}
	//返回查询结果
	c.String(http.StatusOK, string(res.Payload))
}

func GetFile(c *gin.Context) {
	cid := c.PostForm("cid")
	filename := c.PostForm("filename")
	err := api.IpfsGet(cid, filename)
	if err != nil {
		c.String(http.StatusBadGateway, fmt.Sprintf("%v", err))
		return
	}
	//返回下载的文件的路径
	c.String(http.StatusOK, fmt.Sprintf("http://127.0.0.1:9090/downloadfile?filepath=%v", filename))
}

func QueryVideo(c *gin.Context) {
	cid := c.PostForm("cid")
	//查询链码
	var args [][]byte
	args = append(args, []byte(cid))
	res, err := api.ChannelQuery("queryVideo", args)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("execute chaincode failed err:%v", err))
		return
	}
	//返回查询结果
	c.String(http.StatusOK, string(res.Payload))
}

func DownloadFile(c *gin.Context) {
	filename := c.Query("filepath")
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	_, err := os.Stat(fmt.Sprintf("./files/downloadfiles/%v", filename))
	if err != nil {
		c.String(http.StatusBadRequest, "文件已删除！请重新获取链接！")
		return
	}
	//return file
	c.File(fmt.Sprintf("./files/downloadfiles/%v", filename))
	//delete file
	// os.Remove(fmt.Sprintf("./files/downloadfiles/%v", enFilename))
}
//这段代码是一个Golang包，使用Gin框架构建的Web服务器的控制器部分，主要与区块链操作相关。下面是每个函数的简要解释：
//enerateKeyPair：生成RSA密钥对并以JSON格式返回私钥（sk）和公钥（pk）。
//RestoreKey：接收一个私钥，转换成对应的公钥，并返回。如果转换失败，返回错误消息。
//Registeruser：用户注册功能，接收用户的私钥和用户名，将其转换为公钥并使用区块链技术注册用户信息。如果成功，返回“注册成功！”，否则返回错误。
//Searchuser：根据提供的公钥查询用户信息，通过区块链执行查询操作，并返回查询结果。
//Upload：处理文件上传，保存到本地后将文件名上传至IPFS（一种分布式文件存储系统），并将相关信息注册到区块链中。返回交易ID。
//QueryUplink：根据用户的私钥查询其上传信息，执行区块链查询，并返回结果。
//GetFile：通过IPFS下载文件，返回下载链接。
//QueryVideo：查询作品信息，通过传入的内容识别符(cid)在区块链上执行查询操作，并返回结果。
//DownloadFile：提供文件下载服务，检查文件是否存在，设置适当的HTTP头，供客户端下载，并可选择在下载后删除文件。
//整体而言，此代码涉及文件处理、区块链交互和基本的网络请求处理。它结合了区块链技术和IPFS，用于存储和验证数据的分布式应用程序。