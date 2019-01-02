find-broken-urls
==================
A script to run through course resources to check for broken or redirected urls

### Installation

As the script is written in Go, make sure you have version 1.10.0 or greater installed.

Using [Homebrew](https://brew.sh/) to install go
* Run `brew install go` or `brew upgrade go`
* Set your `GOPATH` environment variable, this specifies the location of your workspace

#### Database

* Run `brew install mongodb`
* Run `brew services restart mongodb`

#### Loading Mongo Data from scripts

For script to work data is needed to be stored in mongo db and so it is 
advised to run the load-data scripts with the latest csv files from [HESA website](https://www.hesa.ac.uk/support/tools-and-downloads/unistats)

See [load data](https://github.com/office-for-students/alpha-scripts/tree/develop/mongo/load-data)

#### Running script

* Run `make debug`

### Contributing

See [CONTRIBUTING](../../CONTRIBUTING.md) for details.

### License

See [LICENSE](../../LICENSE.md) for details.
