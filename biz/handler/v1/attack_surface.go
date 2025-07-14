package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	as "github.com/xid-protocol/xidp/protocols/attack_surface"
)

func GetAttackSurface(c *gin.Context) {
	id := c.Query("xid")
	xidType := c.Query("xidType")

	// repository := repositories.NewXidInfoRepository()
	// xidInfo, err := repository.FindOneByXidAndPath(c, xid, xidType)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	if xidType == "aws-instanceid" {
		awsInstance := as.NewAWSAttackSurface(id)
		c.JSON(http.StatusOK, awsInstance)
		return
	}

	c.JSON(http.StatusOK, nil)
}
