# es-reindex
A tool for reindexing elasticsearch data because of mapping change. Only work with Elasticsearch version 1.x, in 2.3.0 there is a reindex api.
## Usage
```
Usage of ./es-reindex:
  -e string
        The Elasticsearch API. (default "http://127.0.0.1:9200")
  -f my_index/my_type
        Index which from: my_index/my_type
  -t dst_index/dst_type
        Index which to: dst_index/dst_type
```
