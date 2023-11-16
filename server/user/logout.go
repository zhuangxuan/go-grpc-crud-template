package user

import (
	"github.com/gin-gonic/gin"
	"go_crud/server/user/utils"
	"go_crud/server/utils/token"
	"gorm.io/gorm"
)

func LogoutGet(r *gin.RouterGroup, DB *gorm.DB) {
	r.GET("logout/", func(c *gin.Context) {
		tokenData := c.GetHeader("token")
		err := token.CheckHS(tokenData)
		//fmt.Println(tokenData)
		if err != nil {
			c.JSON(200, gin.H{
				"msg":  "无效登录状态",
				"data": "",
				"code": "444",
			})
			c.Abort()
			return
		}
		claims := token.UserClaims{}
		token.Hs.Decode(tokenData, &claims)
		logoutName := claims.Data.(string)
		userDataList := utils.GetUserByName(logoutName, DB)
		if len(userDataList) == 0 {
			c.JSON(200, gin.H{
				"msg":  "用户不存在",
				"data": logoutName,
				"code": "400",
			})
			return
		}
		utils.SetUserStatus(userDataList[0], DB, "out")
		c.JSON(200, gin.H{
			"msg":  "注销登录",
			"data": logoutName,
			"code": "235",
		})
	})
}
