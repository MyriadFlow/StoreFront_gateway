package flowid

import (
	"net/http"

	"github.com/TheLazarusNetwork/marketplace-engine/config/dbconfig"
	"github.com/TheLazarusNetwork/marketplace-engine/models"
	"github.com/TheLazarusNetwork/marketplace-engine/util/pkg/flowid"
	"github.com/TheLazarusNetwork/marketplace-engine/util/pkg/httphelper"
	"github.com/TheLazarusNetwork/marketplace-engine/util/pkg/logwrapper"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// ApplyRoutes applies router to gin Router
func ApplyRoutes(r *gin.RouterGroup) {
	g := r.Group("/flowid")
	{
		g.GET("", getFlowId)
	}
}

func getFlowId(c *gin.Context) {
	db := dbconfig.GetDb()
	var user models.User
	walletAddress := c.Query("walletAddress")
	if walletAddress == "" {
		httphelper.ErrResponse(c, http.StatusBadRequest, "Wallet address (walletAddress) is required")
		return
	}
	dbRes := db.Model(&models.User{}).Where("wallet_address = ?", walletAddress).First(&user)
	// If there is an error and that error is not of "record not found"
	if dbRes.Error != nil && dbRes.Error != gorm.ErrRecordNotFound {
		log.Error(dbRes.Error)
		httphelper.ErrResponse(c, http.StatusInternalServerError, "Unexpected error occured")
		return
	}
	// If wallet address exist
	if dbRes.Error != gorm.ErrRecordNotFound {
		flowId, err := flowid.GenerateFlowId(walletAddress, true, models.AUTH, "")
		if err != nil {
			log.Error(err)
			httphelper.ErrResponse(c, http.StatusInternalServerError, "Unexpected error occured")

			return
		}
		payload := GetFlowIdPayload{
			FlowId: flowId,
		}
		httphelper.SuccessResponse(c, "Flowid successfully generated", payload)
	} else {
		//If wallet address doesn't exist
		flowId, err := flowid.GenerateFlowId(walletAddress, false, models.AUTH, "")
		if err != nil {
			log.Error(err)
			httphelper.ErrResponse(c, http.StatusInternalServerError, "Unexpected error occured")

			return
		}
		if err != nil {
			logwrapper.Error(err)
			httphelper.ErrResponse(c, 500, "Unexpected error occured")
			return
		}
		userAuthEULA := "TODO AUTH EULA"
		payload := GetFlowIdPayload{
			FlowId: flowId,
			Eula:   userAuthEULA,
		}
		httphelper.SuccessResponse(c, "Flowid successfully generated", payload)
	}
}
