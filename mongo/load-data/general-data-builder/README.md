General data builder
==================

This script is to insert course location resources into mongodb datastore using the course location csv gathered from [HESA website](https://www.hesa.ac.uk/support/tools-and-downloads/unistats)

### How to run service
* Extract csv files from `office-for-students/alpha-scripts/mongo/load-data/files.zip` by running the following command:
```
cd <path to zip file>; unzip -r files.zip; cd institution-builder;
``` 
* Run `go build`
* Run `./general-data-builder -mongo-uri=<url>`

The url should look something like the following `localhost:27017` or
`127.0.0.1:27017`. If a username and password are needed follow this structure
`<username>:<password>@<host>:<port>`
