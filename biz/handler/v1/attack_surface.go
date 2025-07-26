package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAttackSurface(c *gin.Context) {

	// repository := repositories.NewXidInfoRepository()
	// xidInfo, err := repository.FindOneByXidAndPath(c, xid, xidType)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// query := c.Query("query")
	// if query == "aws" {
	// 	awsInstance := attack_surface.NewAWSAttackSurface(id)
	// 	c.JSON(http.StatusOK, awsInstance)
	// 	return
	// }

	c.JSON(http.StatusOK, nil)
}
