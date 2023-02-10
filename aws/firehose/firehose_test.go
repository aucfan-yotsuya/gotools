package firehose

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"log"
	"testing"

	myconfig "rebill/config"

	"github.com/aucfan-yotsuya/gomod/common"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/oklog/ulid"
	"github.com/stretchr/testify/assert"
)

var c *myconfig.Config

func TestMain(m *testing.M) {
	c = new(myconfig.Config)
	if err := myconfig.LoadConfig(common.Pstring("../toml/development.toml"), c); err != nil {
		log.Fatal(err)
	}
	m.Run()
}
func TestPutRecord(t *testing.T) {
	b, err := json.Marshal(map[string]interface{}{
		"key": "バイナリデータ",
	})
	assert.NoError(t, err)
	err = PutRecord(common.Pstring("yotsuya-gmo"), c, &b)
	assert.NoError(t, err)
}
func TestPutRecordWithConfig(t *testing.T) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("test"))
	assert.NoError(t, err)
	cfg.Region = "ap-northeast-1"
	tm := common.NowJST()
	ulidStr := ulid.MustNew(ulid.Now(), rand.Reader).String()
	{
		b, err := json.Marshal(map[string]interface{}{
			"ymd":      tm.Format("20060102"),
			"ulid":     ulidStr,
			"api":      "ExecTran",
			"mode":     "request",
			"shopID":   "xxx",
			"memberID": "xxx",
		})
		assert.NoError(t, err)
		err = PutRecordWithConfig(cfg, c, &b)
		assert.NoError(t, err)
	}
	{
		b, err := json.Marshal(map[string]interface{}{
			"ymd":     tm.Format("20060102"),
			"ulid":    ulidStr,
			"api":     "ExecTran",
			"mode":    "response",
			"errCode": "EX1",
			"ErrInfo": "EX1000302",
		})
		assert.NoError(t, err)
		err = PutRecordWithConfig(cfg, c, &b)
	}
	{
		b, err := json.Marshal(map[string]interface{}{
			"ymd":     tm.Format("20060102"),
			"ulid":    ulidStr,
			"api":     "ExecTran",
			"mode":    "response",
			"errCode": "EX1",
			"ErrInfo": "EX1000302",
		})
		assert.NoError(t, err)
		err = PutRecordWithConfig(cfg, c, &b)
	}
	t.Log(tm.Format("20060102"), "time")
}
