package user

import (
	"net/http"

	"github.com/cloud-disk/app/types/result"
	"github.com/gin-gonic/gin"
)

func GetUserInfo(c *gin.Context) {
	userName := c.DefaultQuery("user_name", "")
	if userName == "" {
		c.JSON(http.StatusBadRequest, result.Failure(http.StatusBadRequest, "the request user name is empty"))
		return
	}

	rsp, err := queryUserInfoTable()
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.Failure(http.StatusInternalServerError, "query user info error"))
		return
	}
	c.JSON(http.StatusOK, result.Success(rsp))
}

func queryUserInfoTable() (*UserInfo, error) {
	return nil, nil
}
