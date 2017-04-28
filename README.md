# Clay

[![Build Status](https://travis-ci.org/qb0C80aE/clay.svg?branch=develop)](https://travis-ci.org/qb0C80aE/clay)

Clay is an abstract system model store to automate various kind of operations.
It provides some APIs to access the system model store.

## Related modules

* [Loam](https://github.com/qb0C80aE/loam)
  * The basic models and functions work on Clay
* [Pottery](https://github.com/qb0C80aE/pottery)
  * A simple GUI module works with Loam on Clay

### Pottery UI

#### UI - network design
![Network design](https://github.com/qb0C80aE/pottery/raw/develop/images/sample1.png)

#### UI - physial diagram from the system model store
![Physical diagram](https://github.com/qb0C80aE/pottery/raw/develop/images/sample2.png)

#### UI - logical diagram from the system model store
![Logical diagram](https://github.com/qb0C80aE/pottery/raw/develop/images/sample3.png)

# How to build and run

```
$ cd $GOPATH/src/github.com/qb0C80aE/clay
$ # note: if you have any modules what you install into Clay, add lines to the import section of submodules/submodules.go like below.
$ # _ "github.com/qb0C80aE/loam" // Install Loam sub module by importing
$ # _ "github.com/qb0C80aE/pottery" // Install Pottery sub module by importing
$ go get github.com/Masterminds/glide
$ glide install
$ go build
$ ./clay
```

The server runs at http://localhost:8080 by default.

## Environmental variables

You can give the environmental variables to Clay.

|Key         |Description                                                                      |Options    |Default  |
|:-----------|:--------------------------------------------------------------------------------|:----------|:--------|
|HOST        |The host to listen.                                                              |-          |localhost|
|PORT        |The port to listen.                                                              |-          |8080     |
|DB_MODE     |The indentifier how the db is managed.                                           |memory/file|memory   |
|DB_FILE_PATH|The path where the db file is located. This value is used if DB_MODE=file is set.|-          |clay.db  |

## Windows build

Due to ``mattn/go-sqlite3``, mingw gcc is required.

  1. Install msys2 https://msys2.github.io/
  2. Run msys2 shell. i.e. ``C:\mingw64\msys2.exe``

```bash
$ pacman -S mingw-w64-x86_64-gcc
$ cd $GOPATH/src/github.com/qb0C80aE/clay
$ go build
$ ./clay
```

Powershell

```powershell
PS> C:\msys64\usr\bin\pacman -S mingw-w64-x86_64-gcc
PS> cd $env:GOPATH/src/github.com/qb0C80aE/clay
PS> powershell { $env:PATH+=";C:\msys64\mingw64\bin"; go build }
PS> .\clay.exe
```

Creating go-sqlite3 build archive makes rebuild time shorter.

```
PS> powershell { $env:PATH+=";C:\msys64\mingw64\bin"; go install github.com/mattn/go-sqlite3 }
```

You'll see ``$GOPATH\pkg\windows_amd64\github.com\mattn\go-sqlite3.a``.

# How to use

## Import and export the design

You can import and export the models you created through `design` resource.
Clay is designed as a standalone modeling tool, and the created design should be stored as human-readable text files in versioning repositories like git to make it easier to realize infrastructure-as-code.

```
$ # Export the design
$ curl -X GET 'localhost:8080/designs/present?pretty' > design.json
$ # Import and overwrite the design
$ curl -X PUT 'localhost:8080/designs/present' -H 'Content-Type: application/json' -d @design.json
```

## Templates

You can register some text templates and generate something using the models in clay.

```
$ # register template1 and external parameters
$ curl -X POST "localhost:8080/templates" -H "Content-Type: multipart/form-data" -F name=sample1 -F template_content=@examples/sample.template
$ curl -X POST "localhost:8080/templates" -H "Content-Type: application/json" -d '{"name": "sample2", "template_content": "sample2"}'
$ curl -X POST "localhost:8080/templates" -H "Content-Type: application/json" -d '{"name": "sample3", "template_content": "sample3"}'
$ curl -X POST "localhost:8080/template_external_parameters" -H "Content-Type: application/json" -d '{"template_id": 1, "name": "testParameter11", "value_string": {"String": "TestParameter11", "Valid": true}, "value_int": {"Int64": 1, "Valid": true}}'
$ curl -X POST "localhost:8080/template_external_parameters" -H "Content-Type: application/json" -d '{"template_id": 1, "name": "testParameter12", "value_string": {"String": "TestParameter12", "Valid": true}}'
$ curl -X POST "localhost:8080/template_external_parameters" -H "Content-Type: application/json" -d '{"template_id": 1, "name": "testParameter1X", "value_int": {"Int64": 100, "Valid": true}}'
$ # register template2 and external parameters
$ curl -X POST "localhost:8080/templates" -H "Content-Type: application/json" -d '{"name": "sample2", "template_content": "{{.testParameter1X}}"}'
$ curl -X POST "localhost:8080/template_external_parameters" -H "Content-Type: application/json" -d '{"template_id": 2, "name": "testParameter1X", "value_int": {"Int64": 200, "Valid": true}}'
$ # register template3 and external parameters
$ curl -X POST "localhost:8080/templates" -H "Content-Type: application/json" -d '{"name": "sample3", "template_content": "{{.testParameter1X}}"}'
$ curl -X POST "localhost:8080/template_external_parameters" -H "Content-Type: application/json" -d '{"template_id": 3, "name": "testParameter1X", "value_int": {"Int64": 300, "Valid": true}}'
$ # show generated template
$ curl -X GET "localhost:8080/templates/1"
$ # Geenrate a text from the tempalte
$ curl -X PATCH "localhost:8080/templates/1"
```

# API Server

Simple Rest API using gin(framework) & gorm(orm)

## Endpoint list

### Designs Resource

```
GET    /designs/present
PUT    /designs/present
DELETE /designs/present
```

### TemplateExternalParameter Resource

```
GET    /template_external_parameters
GET    /template_external_parameters/:id
POST   /template_external_parameters
PUT    /template_external_parameters/:id
DELETE /template_external_parameters/:id
PATCH /template_external_parameters/:id
```

### Template Resource

```
GET    /templates
GET    /templates/:id
POST   /templates
PUT    /templates/:id
DELETE /templates/:id
PATCH /templates/:id
```

# Thanks

* Clay was generated by https://github.com/wantedly/apig
