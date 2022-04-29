package estools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/silenceGuo/Go-Tools/utils"
	"strings"
	"sync/atomic"
	"time"
)

var host = "http://192.168.254.11:9200/"

type EsClient struct {
	host   []string
	user   string
	passwd string
	client *elasticsearch.Client
	//indexName string
}

func (es *EsClient) GetClinet() {
	// 获取es连接，
	utils.ZapLogger.Info("connect to :", es.host)
	config := elasticsearch.Config{
		Addresses: es.host,
		Username:  es.user,
		Password:  es.passwd,
	}
	esclient, err := elasticsearch.NewClient(config)
	//esclient, err := elasticsearch.NewDefaultClient()
	if err != nil {
		utils.ZapLogger.Error("Error creating the client: %s", err)
	}
	res, err := esclient.Info()
	if err != nil {
		utils.ZapLogger.Fatalf("Error getting response: %s", err)
	}
	//res.StatusCode
	utils.ZapLogger.Info("connect sucessfull:", res.StatusCode)
	es.client = esclient

}
func (es *EsClient) CreateIndex(indexName string) {
	// 创建索引
	if es.client == nil {
		es.GetClinet()
	}
	res, err := es.client.Indices.Create(indexName)
	if err != nil {
		utils.ZapLogger.Fatalf("Cannot create index: %s", err)
	}
	if res.IsError() {
		utils.ZapLogger.Fatalf("Cannot create index response: %s", res)
	}
	utils.ZapLogger.Info("create index response sucessfull:", res)
}

func (es *EsClient) DelIndex(indexNames []string) {
	// 删除 索引
	if es.client == nil {
		es.GetClinet()
	}
	res, err := es.client.Indices.Delete(indexNames)
	if err != nil {
		utils.ZapLogger.Fatalf("Cannot delete index: %s", err)
	}
	if res.IsError() {
		utils.ZapLogger.Fatalf("Cannot delete index response: %s", res)
	}
	utils.ZapLogger.Info("delete index response sucessfull:", res)
}

func (es *EsClient) QueryIndex(indexName string) {
	if es.client == nil {
		es.GetClinet()
	}
	body := `{
               "query": {
                   "match": { "id": 18030 }
               },
              "aggregations": {
              "top_10_states": { "terms": { "field": "state", "size": 10 } }
              }
             }`

	res, err := es.client.Search(
		es.client.Search.WithIndex(indexName),
		es.client.Search.WithBody(strings.NewReader(body)),
		es.client.Search.WithPretty(),
	)
	if err != nil {
		utils.ZapLogger.Fatalf("Error getting response: %s", err)
	}
	utils.ZapLogger.Info(res)
	defer res.Body.Close()

}

type Article struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Published time.Time `json:"published"`
	Author    string    `json:"author"`
}

func (es *EsClient) Bulk(indexName string, data []interface{}) {
	// 批量提交接口bulk， 10w数据8-9s
	// 数据结构体切片
	if es.client == nil {
		es.GetClinet()
	}
	var buf bytes.Buffer
	for i, a := range data {
		//meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%d" } }%s`, a.ID, "\n"))
		meta := []byte(fmt.Sprintf(`{ "index" : {} }%s`, "\n"))
		data, err := json.Marshal(a)
		if err != nil {
			utils.ZapLogger.Fatalf("Cannot encode date %d: %s", i, err)
		}
		// Append newline to the data payload
		data = append(data, "\n"...)
		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)
	}
	res, err := es.client.Bulk(bytes.NewReader(buf.Bytes()), es.client.Bulk.WithIndex(indexName))
	if err != nil {
		utils.ZapLogger.Fatalf("Failure indexing batch : %s", err)
	}
	utils.ZapLogger.Info(res.StatusCode)
	res.Body.Close()
	buf.Reset()
	es.client.Indices.Refresh()

}
func (es *EsClient) BulkParallel(indexName string, NumWorkers int, data []interface{}) {
	// 批量提交接口bulk 多协程，10w数据8-9s
	// 数据结构体切片
	var countSuccessful uint64
	indexer, _ := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         indexName,
		Client:        es.client,
		NumWorkers:    NumWorkers,
		FlushBytes:    int(100000),
		FlushInterval: 30 * time.Second,
	})
	start := time.Now().UTC()
	for i, a := range data {
		datas, err := json.Marshal(a)
		if err != nil {
			utils.ZapLogger.Fatalf("Cannot encode article %d: %s", i, err)
		}
		err = indexer.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action: "index",
				Body:   bytes.NewReader(datas),
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&countSuccessful, 1)
				},
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						utils.ZapLogger.Info("ERROR: %s", err)
					} else {
						utils.ZapLogger.Info("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			})

	}
	if err := indexer.Close(context.Background()); err != nil {
		utils.ZapLogger.Fatalf("Unexpected error: %s", err)
	}
	// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	biStats := indexer.Stats()
	// Report the results: number of indexed docs, number of errors, duration, indexing rate
	//
	utils.ZapLogger.Info(strings.Repeat("▔", 65))

	dur := time.Since(start)

	if biStats.NumFailed > 0 {
		utils.ZapLogger.Fatalf(
			"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			humanize.Comma(int64(biStats.NumFailed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	} else {
		utils.ZapLogger.Infof(
			"Sucessfuly indexed [%s] documents in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	}
}
