load-data
==================

These scripts are to load raw course and institution data into mongodb from several csvs gathered from [HESA website](https://www.hesa.ac.uk/support/tools-and-downloads/unistats) and to create new resource structures depending on the needs of services. 

### Installation

As the script is written in Go, make sure you have version 1.10.0 or greater installed.

Using [Homebrew](https://brew.sh/) to install go
* Run `brew install go` or `brew upgrade go`
* Set your `GOPATH` environment variable, this specifies the location of your workspace

#### Database

* Run `brew install mongodb`
* Run `brew services restart mongodb`

* Set environment variable for mongo url, like so:
`export MONGO_URI=<mongo uri>`

The mongo uri should look something like the following `localhost:27017` (this is the default
when running `make debug`), or `127.0.0.1:27017`. If a username and password are needed follow
this structure `<username>:<password>@<host>:<port>`

### How to run scripts

To get the latest data, download from [HESA website](https://www.hesa.ac.uk/support/tools-and-downloads/unistats) and either add csvs to files directory or replace csvs found in files directory with those downloaded or use the data taken from 28th November 2018.

* Extract csv files from `ofs-alpha-backend/scripts/files.zip` by running the following command:
```
cd <path to zip file>; unzip -r files.zip
``` 

* Run `make debug` this shall take approximately 14 minutes to complete


### Contributing

See [CONTRIBUTING](../../CONTRIBUTING.md) for details.

### License

See [LICENSE](../../LICENSE.md) for details.
