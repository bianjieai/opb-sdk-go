package test

import (
	"fmt"
	"github.com/irisnet/irismod-sdk-go/record"
	"github.com/stretchr/testify/require"
	"testing"
)

// 创建存证记录。
func TestRecordCreate(t *testing.T) {
	content := record.Content{
		Digest:     "test_digest",
		DigestAlgo: "test_digest-algo",
		URI:        "https://www.baidu.com",
		Meta:       "test_meta",
	}
	contents := []record.Content{content}

	req := record.CreateRecordRequest{
		Contents: contents,
	}

	recordID, err := txClient.Record.CreateRecord(req, baseTx)
	require.NoError(t, err)
	require.NotEmpty(t, recordID)
	fmt.Println(recordID)
}

// 查询指定的存证记录。
func TestRecordRecord(t *testing.T) {
	request := record.QueryRecordReq{
		RecordID: "a833204e361febf99b9a37d25b7d8f27134b75b1c215f8a3f873ab6780205a32",
		Prove:    true,
		Height:   0,
	}

	result, err := txClient.Record.QueryRecord(request)
	require.NoError(t, err)
	fmt.Println(result.Record.Contents)
}
