Institution builder
==================

This script is to insert institution resources into mongodb datastore from several csvs gathered from [HESA website](https://www.hesa.ac.uk/support/tools-and-downloads/unistats)

### How to run service
* Extract csv files from `office-for-students/alpha-scripts/mongo/load-data/files.zip` by running the following command:
```
cd <path to zip file>; unzip -r files.zip; cd institution-builder;
``` 
* Run `go build`
* Run `./institution-builder -mongo-url=<url>`

The url should look something like the following `localhost:27017` or
`127.0.0.1:27017`. If a username and password are needed follow this structure
`<username>:<password>@<host>:<port>`
