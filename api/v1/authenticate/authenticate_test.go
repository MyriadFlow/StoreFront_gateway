package authenticate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/TheLazarusNetwork/marketplace-engine/api/types"
	"github.com/TheLazarusNetwork/marketplace-engine/api/v1/flowid"
	"github.com/TheLazarusNetwork/marketplace-engine/config"
	"github.com/TheLazarusNetwork/marketplace-engine/config/dbconfig/dbinit"
	"github.com/TheLazarusNetwork/marketplace-engine/util/pkg/logwrapper"
	testingcommmon "github.com/TheLazarusNetwork/marketplace-engine/util/testingcommon"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TODO: Write test to verify expiry
func Test_PostAuthenticate(t *testing.T) {
	config.Init("../../../.env")
	logwrapper.Init("../../../logs")
	dbinit.Init()
	t.Cleanup(testingcommmon.DeleteCreatedEntities())
	gin.SetMode(gin.TestMode)

	url := "/api/v1.0/authenticate"

	t.Run("Should return 200 with correct wallet address", func(t *testing.T) {
		testWallet := testingcommmon.GenerateWallet()
		eula, flowId := callFlowIdApi(testWallet.WalletAddress, t)
		signature := getSignature(eula, flowId, testWallet.PrivateKey)
		body := AuthenticateRequest{Signature: signature, FlowId: flowId}
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		//Request with signature created from correct wallet address
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatal(err)
		}

		c, _ := gin.CreateTestContext(rr)
		c.Request = req
		authenticate(c)
		assert.Equal(t, http.StatusOK, rr.Code, rr.Body.String())
	})
	t.Run("Should return 403 with different wallet address", func(t *testing.T) {
		testWallet := testingcommmon.GenerateWallet()
		eula, flowId := callFlowIdApi(testWallet.WalletAddress, t)
		// Different private key will result in different wallet address
		differentPrivatekey := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		signature := getSignature(eula, flowId, differentPrivatekey)
		body := AuthenticateRequest{Signature: signature, FlowId: flowId}
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		newWalletAddress := testWallet.WalletAddress + "ba"
		callFlowIdApi(newWalletAddress, t)

		rr := httptest.NewRecorder()

		//Request with signature stil created from different walletAddress
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatal(err)
		}
		c, _ := gin.CreateTestContext(rr)
		c.Request = req
		authenticate(c)
		assert.Equal(t, http.StatusForbidden, rr.Code, rr.Body.String())
	})

}

func callFlowIdApi(walletAddress string, t *testing.T) (eula string, flowidString string) {
	// Call flowid api
	u, err := url.Parse("/api/v1.0/flowid")
	q := url.Values{}
	q.Set("walletAddress", walletAddress)
	u.RawQuery = q.Encode()
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", u.String(), nil)
	req.URL.RawQuery = q.Encode()
	if err != nil {
		t.Error(err)
	}
	c, _ := gin.CreateTestContext(rr)
	c.Request = req
	flowid.GetFlowId(c)
	assert.Equal(t, http.StatusOK, rr.Code, "Failed to call flowApi")
	var flowIdPayload flowid.GetFlowIdPayload
	var res types.ApiResponse
	decoder := json.NewDecoder(rr.Result().Body)
	err = decoder.Decode(&res)
	testingcommmon.ExtractPayload(&res, &flowIdPayload)
	if err != nil {
		t.Fatal(err)
	}
	return flowIdPayload.Eula, flowIdPayload.FlowId
}

func getSignature(eula string, flowId string, hexPrivateKey string) string {
	message := eula + flowId
	newMsg := fmt.Sprintf("\x19Ethereum Signed Message:\n%v%v", len(message), message)

	privateKey, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		log.Fatal("HexToECDSA failed ", err)
	}

	// keccak256 hash of the data
	dataBytes := []byte(newMsg)
	hashData := crypto.Keccak256Hash(dataBytes)

	signatureBytes, err := crypto.Sign(hashData.Bytes(), privateKey)
	if err != nil {
		log.Fatal("len", err)
	}

	signature := hexutil.Encode(signatureBytes)

	return signature
}
