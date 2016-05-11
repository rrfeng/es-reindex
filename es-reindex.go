package main

import (
	"flag"
	"gopkg.in/olivere/elastic.v2"
	//esv3 "gopkg.in/olivere/elastic.v3"
	"encoding/json"
	"log"
	"strings"
)

var esAddr *string = flag.String("e", "http://127.0.0.1:9200", "The Elasticsearch API.")
var src *string = flag.String("f", "", "Index which from: `my_index/my_type`")
var dst *string = flag.String("t", "", "Index which to: `dst_index/dst_type`")

func splitBySlash(s string) (a string, b string) {
	res := strings.Split(s, "/")
	a = res[0]
	b = ""
	if len(res) == 2 {
		b = res[1]
	}
	return
}

func main() {
	flag.Parse()
	if *src == "" || *dst == "" {
		log.Fatalln("Please set the source and target index names!")
	} else if *esAddr == "" {
		*esAddr = "http://127.0.0.1:9200"
	}

	srcIndex, srcType := splitBySlash(*src)
	dstIndex, dstType := splitBySlash(*dst)
	if srcType == "" || dstType == "" {
		log.Fatalln("Please set the source and target type!")
	}

	es, err := elastic.NewClient(elastic.SetURL(*esAddr))
	if err != nil {
		log.Fatalln(err.Error())
	}

	exists, err := es.IndexExists(srcIndex).Do()
	if err != nil {
		log.Fatalln(err.Error())
	} else if !exists {
		log.Fatalln("Source index does not exists:", srcIndex)
	}

	exists, err = es.IndexExists(dstIndex).Do()
	if err != nil {
		log.Fatalln(err.Error())
	} else if !exists {
		log.Fatalln("Destnation index does not exists:", dstIndex, "\nYou should create the new index and put the mapping template first!")
	}

	var count int64 = 0

	copyByType := func(hit *elastic.SearchHit, bulkService *elastic.BulkService) error {
		source := make(map[string]interface{})
		if err := json.Unmarshal(*hit.Source, &source); err != nil {
			return err
		}
		if hit.Type == srcType {
			req := elastic.NewBulkIndexRequest().Index(dstIndex).Type(dstType).Id(hit.Id).Doc(source)
			bulkService.Add(req)
			count++
			if count%1000 == 0 {
				log.Println("Progress:", count)
			}
		}
		return nil
	}

	//	showProgress := func(current, total int64) {
	//		if current%1000 == 0 {
	//			log.Println(current, "of", total)
	//		}
	//	}

	task := elastic.NewReindexer(es, srcIndex, copyByType)
	ret, err := task.Do()
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Done!", ret)
}
