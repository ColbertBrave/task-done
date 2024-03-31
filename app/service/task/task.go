package task

import "github.com/gin-gonic/gin"

func GetTaskInfo(c *gin.Context) {
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