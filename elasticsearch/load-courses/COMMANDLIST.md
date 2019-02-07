Command List (examples)
==================

1) Check health of elasticsearch cluster: `curl -XGET 'localhost:9200/_cluster/health/courses?pretty'`

2) Get index: `curl -XGET 'localhost:9200/courses?pretty`

3) Get index settings only: `curl -XGET 'localhost:9200/courses/_mapping?pretty'`
`
4) Get index mappings only: `curl -XGET 'localhost:9200/courses/_mapping?pretty'`

5) Delete index: `curl -XDELETE 'localhost:9200/courses'`

6) Get a list of indexes: `curl -XGET 'localhost:9200/_aliases?pretty'`

7) Watch index being built - `watch -n 2 "curl -s localhost:9200/courses/_count?pretty"`

8) Basic querying of index: `curl -XGET 'localhost:9200/courses/_search?q=physics&size=5&pretty'`
