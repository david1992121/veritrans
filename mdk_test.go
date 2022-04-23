package veritrans

import (
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/joho/godotenv"
	assert "github.com/stretchr/testify/require"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No env file for testing")
	}
}

func TestGetCardToken(t *testing.T) {
	cardService := NewMDKService(MDKConfig{
		ApiURL:   os.Getenv("MDK_API_URL"),
		ApiToken: os.Getenv("MDK_API_TOKEN"),
	})

	nowTime := time.Now()
	expiredAt := nowTime.AddDate(0, 1, 0)

	cardToken, err := cardService.GetCardToken(&ClientCardInfo{
		CardNumber:   "4111111111111111",
		CardExpire:   expiredAt.Format("01/06"),
		SecurityCode: "123",
	})
	re := regexp.MustCompile(`[0-9a-z\-]{36}`)

	assert.Nil(t, err)
	assert.Equal(t, true, re.Match([]byte(cardToken)))
}
