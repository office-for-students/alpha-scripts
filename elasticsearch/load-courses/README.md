load-courses
==================

These scripts are to load course resources from mongodb datastore into elasticsearch.

### Installation

As the script is written in Go, make sure you have version 1.10.0 or greater installed.

Using [Homebrew](https://brew.sh/) to install go
* Run `brew install go` or `brew upgrade go`
* Set your `GOPATH` environment variable, this specifies the location of your workspace

#### Database

* Run `brew install mongodb`
* Run `brew services restart mongodb`

* Set environment variable for mongo uri, like so:
`export MONGO_URI=<mongo uri>`

The mongo uri should look something like the following `localhost:27017` (this is the default
when running `make debug`), or `127.0.0.1:27017`. If a username and password are needed follow
this structure `<username>:<password>@<host>:<port>`

#### Elasticsearch

It is expected that any minor version of elastic search version *6* to be installed

* Run `brew install elasticsearch`
* Run `brew services restart elasticsearch`

* Set environment variable for elasticsearch uri, like so:
`export ELASTIC_SEARCH_URI=<elasticsearch uri>`

### How to run scripts

To be able to successfully load data across from mongo db to elasticsearch it is required that you run the following scripts in the [mongo load-data file](https://github.com/office-for-students/alpha-scripts/tree/develop/mongo/load-data); follow the documentation in the README.md.

Once data is in mongo db, then one can do the following:

* Run `make debug` this shall take approximately several seconds to complete

### Contributing

See [CONTRIBUTING](../../CONTRIBUTING.md) for details.

### License

See [LICENSE](../../LICENSE.md) for details.
