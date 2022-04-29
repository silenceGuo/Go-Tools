package estools

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestClinet(t *testing.T) {
	//var article Article
	host := []string{
		"http://192.168.254.11:9200/",
	}
	////i := []string{"test-index"}
	var es EsClient
	es.host = host
	es.GetClinet()
	es.QueryIndex("test-index3")
	//es.CreateIndex("test-index2")
	//es.DelIndex([]string{"test-index3"})
	//es.CreateIndex("test-index3")
}
func TestEsClient_BulkParallel(t *testing.T) {
	var articles []interface{}
	//var article Article
	host := []string{
		"http://192.168.254.11:9200/",
	}
	var es EsClient
	es.host = host
	es.GetClinet()
	for i := 1; i <= 100000; i++ {
		articles = append(articles, &Article{
			ID:        i,
			Title:     strings.Join([]string{"Title2-", strconv.Itoa(i)}, " "),
			Body:      "Lorem ipsum dolor sit amet...",
			Published: time.Now().Round(time.Second).UTC().AddDate(0, 0, i),
			Author:    "test",
		})
	}

	es.BulkParallel("test-index3", 8, articles)
}
func TestEsClient_Bulk(t *testing.T) {
	var articles []interface{}
	//var article Article
	host := []string{
		"http://192.168.254.11:9200/",
	}
	var es EsClient
	es.host = host
	es.GetClinet()
	for i := 1; i <= 100000; i++ {
		articles = append(articles, &Article{
			ID:        i,
			Title:     strings.Join([]string{"Title2-", strconv.Itoa(i)}, " "),
			Body:      "Lorem ipsum dolor sit amet...",
			Published: time.Now().Round(time.Second).UTC().AddDate(0, 0, i),
			Author:    "test",
		})
	}

	starttime := time.Now().Second()
	es.Bulk("test-index2", articles)
	fmt.Println("speed", time.Now().Second()-starttime)
}
