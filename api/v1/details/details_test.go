package details

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/TheLazarusNetwork/marketplace-engine/api/types"
	"github.com/TheLazarusNetwork/marketplace-engine/config"
	"github.com/TheLazarusNetwork/marketplace-engine/config/dbconfig/dbinit"
	"github.com/TheLazarusNetwork/marketplace-engine/models/Org"
	"github.com/TheLazarusNetwork/marketplace-engine/util/pkg/envutil"
	"github.com/TheLazarusNetwork/marketplace-engine/util/pkg/logwrapper"
	"github.com/TheLazarusNetwork/marketplace-engine/util/testingcommon"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_Details(t *testing.T) {
	config.Init("../../../.env")
	logwrapper.Init("../../../logs")
	dbinit.Init()
	t.Cleanup(testingcommon.DeleteCreatedEntities())
	gin.SetMode(gin.TestMode)

	t.Run("Should be able to get org details", func(t *testing.T) {
		_ = "/api/v1.0/details"
		rr := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rr)
		getDetails(c)
		var response types.ApiResponse
		body := rr.Body
		json.NewDecoder(body).Decode(&response)
		var org Org.Org
		testingcommon.ExtractPayload(&response, &org)
		assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
		assert.Equal(t, envutil.MustGetEnv("ORG_NAME"), org.Name)
		assert.Equal(t, envutil.MustGetEnv("HOME_TITLE"), org.HomeTitle)
		assert.Equal(t, envutil.MustGetEnv("HOME_DESCRIPTION"), org.HomeDescription)
		assert.Equal(t, envutil.MustGetEnv("GRAPH_URL"), org.GraphUrl)
		assert.Equal(t, envutil.MustGetEnv("CREATIFY_CONTRACT_ADDRESS"), org.CreatifyAddress)
		assert.Equal(t, envutil.MustGetEnv("MARKETPLACE_CONTRACT_ADDRESS"), org.MarketPlaceAddress)
		assert.Equal(t, envutil.MustGetEnv("FOOTER"), org.Footer)

		assert.Equal(t, envutil.MustGetEnv("TOP_HIGHLIGHTS"), strings.Join(org.TopHighlights, ","))
		assert.Equal(t, envutil.MustGetEnv("TRENDINGS"), strings.Join(org.Trendings, ","))
		assert.Equal(t, envutil.MustGetEnv("TOP_BIDS"), strings.Join(org.TopBids, ","))
		logrus.Debug(org)
	})

	t.Run("Should be able to update org details", func(t *testing.T) {
		url := "/api/v1.0/details"
		rr := httptest.NewRecorder()
		requestBody := Org.Org{
			HomeTitle:     "Max",
			TopHighlights: []string{"43"},
			Trendings:     []string{"47"},
			TopBids:       []string{"42"},
		}
		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		c, _ := gin.CreateTestContext(rr)
		c.Request = req
		postDetails(c)
		assert.Equal(t, http.StatusOK, rr.Result().StatusCode)

		t.Cleanup(func() {
			err = Org.UpdateOrg(
				Org.Org{
					Name:               envutil.MustGetEnv("ORG_NAME"),
					HomeTitle:          envutil.MustGetEnv("HOME_TITLE"),
					HomeDescription:    envutil.MustGetEnv("HOME_DESCRIPTION"),
					GraphUrl:           envutil.MustGetEnv("GRAPH_URL"),
					CreatifyAddress:    envutil.MustGetEnv("CREATIFY_CONTRACT_ADDRESS"),
					MarketPlaceAddress: envutil.MustGetEnv("MARKETPLACE_CONTRACT_ADDRESS"),
					Footer:             envutil.MustGetEnv("FOOTER"),
					TopHighlights:      strings.Split(envutil.MustGetEnv("TOP_HIGHLIGHTS"), ","),
					Trendings:          strings.Split(envutil.MustGetEnv("TRENDINGS"), ","),
					TopBids:            strings.Split(envutil.MustGetEnv("TOP_BIDS"), ","),
				},
			)
		})
	})

}
