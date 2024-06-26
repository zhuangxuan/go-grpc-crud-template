package main

import (
	"fmt"
	"go_crud/logger"
	"go_crud/server"
	"go_crud/server/crud_rpc"
	crudService "go_crud/server/crud_rpc/service"
	"go_crud/server/files"
	"go_crud/server/midware"
	"go_crud/server/user"
	userService "go_crud/server/user/user_dao/service"
	"go_crud/server/utils"
	"log"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	//配置相关
	viper.AddConfigPath("./conf/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("配置文件错误 %s", err.Error()))
	}

	userService.Init()
	crudService.Init()

	// 服务相关
	r := server.CreateServer()
	r.Use(cors.Default()) //解决跨域

	serverLogger, _ := logger.InitLogger(zap.DebugLevel)
	defer serverLogger.Sync()
	r.Use(logger.GinLogger(serverLogger), logger.GinRecovery(serverLogger, true))

	utils.PingGET(r)

	Router := r.Group("api/refresh", midware.CheckLogin("refresh"))
	utils.RefreshGET(Router)

	userRouter := r.Group("api/user")
	userRouter.Use(gin.Logger(), gin.Recovery())
	user.LoginPost(userRouter)
	user.SignUpPost(userRouter)
	user.LogoutGet(userRouter)
	user.ChangePwdPost(userRouter)
	user.GetPubKey(userRouter)

	crudRpcRouter := r.Group("/api/crud")
	//crudRpcRouter := r.Group("/api/crud", midware.CheckLogin("crud"))
	crudRpcRouter.Use(gin.Logger(), gin.Recovery())
	crud_rpc.AddPOST(crudRpcRouter)
	crud_rpc.QueryGET(crudRpcRouter)
	crud_rpc.QueryPageGET(crudRpcRouter)
	crud_rpc.DeletePOST(crudRpcRouter)
	crud_rpc.UpdatePOST(crudRpcRouter)

	filesRouter := r.Group("/api/files")
	filesRouter.Use(gin.Logger(), gin.Recovery(), midware.CheckLogin("files"))
	files.FileUploadPOST(filesRouter, nil)
	files.BigFileUploadPOST(filesRouter, nil)
	files.FileListGet(filesRouter, nil)
	files.FileDownload(filesRouter, nil)
	files.FileDelete(filesRouter, nil)

	//r.Run("0.0.0.0:8088") // 监听并在 0.0.0.0:8088 上启动服务
	// http://127.0.0.1:8088/ping
	//fmt.Println(r)

	err = r.Run(viper.GetString("server.addr") + ":" + viper.GetString("server.port"))
	if err != nil {
		log.Println(err.Error())
	}

}
