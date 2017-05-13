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

```bash
$ # Suppose that $HOME is /home/user, and $GOPATH is /home/user/go.
$ # Note: Please install glide first.
$ go get github.com/Masterminds/glide
$ # Note: If there are any tools what modules that you want to install into clay depend on, please install first like below.
$ # go get github.com/jteeuwen/go-bindata/...
$ mkdir -p $GOPATH/src/github.com/qb0C80aE/
$ cd $GOPATH/src/github.com/qb0C80aE/
$ git clone https://github.com/qb0C80aE/clay.git
$ cd $GOPATH/src/github.com/qb0C80aE/clay
$ # Edit: If you have modules what you want to install into Clay, add lines like below into the import section of main.go.
$ # _ "github.com/qb0C80aE/loam" // Install Loam module by importing
$ # _ "github.com/qb0C80aE/pottery" // Install Pottery module by importing
$ # Note: If you have added modules into main.go, execute glide get to retrieve those modules like below.
$ # glide get github.com/qb0C80aE/loam
$ # glide get github.com/qb0C80aE/pottery
$ glide install
$ go generate -tags=prebuild ./...
$ go build
$ # Note: If you want to build Clay as a statically linked single binary file, add the flag like below.
$ # go build --ldflags '-extldflags "-static"'
$ ./clay &
```

The server runs at http://localhost:8080 by default.

Creating go-sqlite3 build archive makes rebuild time shorter.

```bash
$ go install github.com/qb0C80aE/clay/vendor/github.com/mattn/go-sqlite3
```

You'll see ``$GOPATH/pkg/linux_amd64/github.com/qb0C80aE/clay/vendor/github.com/mattn/go-sqlite3.a``.

## Environmental variables

You can give the environmental variables to Clay.

|Key         |Description                                                                      |Options    |Default  |
|:-----------|:--------------------------------------------------------------------------------|:----------|:--------|
|HOST        |The host to listen.                                                              |-          |localhost|
|PORT        |The port to listen.                                                              |-          |8080     |
|DB_MODE     |The indentifier how the db is managed.                                           |memory/file|memory   |
|DB_FILE_PATH|The path where the db file is located. This value is used if DB_MODE=file is set.|-          |clay.db  |

## Cross-compile

Due to ``mattn/go-sqlite3``, cross-compilers like mingw gcc are required.
For example, you can build Clay for Linux 32bit, Windows 32bit and 64bit on Ubuntu 16.04.2 LTS 64bit.

```bash
$ # Suppose that $HOME is /home/user, GOROOT is /usr/local/go, and $GOPATH is /home/user/go.
$ cd $HOME
$ sudo apt-get update
$ # Install required packages.
$ sudo apt-get install -y git wget tar gcc
$ ## For Linux 32bit.
$ sudo apt-get install -y gcc-multilib
$ ## For Windows 64bit and 32bit.
$ sudo apt-get install -y binutils-mingw-w64 mingw-w64
$ # Install go
$ wget https://storage.googleapis.com/golang/go1.7.5.linux-amd64.tar.gz
$ sudo tar -C /usr/local -xzf  go1.7.5.linux-amd64.tar.gz
$ # Install go go cross-compile environments. It requires go 1.4.
$ wget https://storage.googleapis.com/golang/go1.4.3.linux-amd64.tar.gz
$ mkdir -p $HOME/go1.4
$ tar -C $HOME/go1.4 --strip-components 1 -xzf go1.4.3.linux-amd64.tar.gz
$ cd $GOROOT/src
$ ## For Linux 32bit.
$ GOOS=linux GOARCH=386 ./make.bash
$ ## For Windows 64bit.
$ GOOS=windows GOARCH=amd64 ./make.bash
$ ## For Windows 32bit.
$ GOOS=windows GOARCH=386 ./make.bash
$ # Prepare go directories.
$ mkdir -p $GOPATH/{src, bin}
$ # Install glide, and additional tools like go-bindata if you need.
$ go get github.com/Masterminds/glide
$ ## go get github.com/jteeuwen/go-bindata/...
$ # Build clay
$ mkdir -p $GOPATH/src/github.com/qb0C80aE
$ cd $GOPATH/src/github.com/qb0C80aE
$ git clone https://github.com/qb0C80aE/clay.git
$ cd $GOPATH/src/github.com/qb0C80aE/clay
$ glide install
$ go generate -tags=generate ./...
$ ## For Linux 32bit.
$ CGO_ENABLED=1 GOOS=linux GOARCH=386 go build --ldflags '-extldflags "-static"' -o linux_386/clay
$ ## For Windows 64bit.
$ CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc LD=x86_64-w64-mingw32-ld GOOS=windows GOARCH=amd64 go build --ldflags '-extldflags "-static"' -o windows_amd64/clay.exe
$ ## For Windows 32bit.
$ CGO_ENABLED=1 CC=i686-w64-mingw32-gcc LD=i686-w64-mingw32-ld GOOS=windows GOARCH=386 go build --ldflags '-extldflags "-static"' -o windows_386/clay.exe
```

## Build on Windows

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

You can import and export the models you created through ``designs`` resource.
Clay is designed as a standalone modeling tool, and the created design should be stored as human-readable text files in versioning repositories like git to make it easier to realize Infrastructure-as-Code.

```
$ # Export the design
$ curl -X GET 'localhost:8080/designs/present?pretty' > design.json
$ # Import and overwrite the design
$ curl -X PUT 'localhost:8080/designs/present' -H 'Content-Type: application/json' -d @design.json
```

If you added some models like [Loam](https://github.com/qb0C80aE/loam), you will be able to use those models in Clay.
See [Loam](https://github.com/qb0C80aE/loam).

## Templates

You can register some text templates and generate something using the models in Clay.
Some functions are provided in template processing, see [the example template in Clay](https://github.com/qb0C80aE/clay/blob/develop/examples/sample.template).

```
$ # register template1 and persistent parameters
$ curl -X POST "localhost:8080/templates" -H "Content-Type: multipart/form-data" -F name=sample1 -F template_content=@examples/sample.template
$ curl -X POST "localhost:8080/templates" -H "Content-Type: application/json" -d '{"name": "sample2", "template_content": "sample2"}'
$ curl -X POST "localhost:8080/templates" -H "Content-Type: application/json" -d '{"name": "sample3", "template_content": "sample3"}'
$ curl -X POST "localhost:8080/template_persistent_parameters" -H "Content-Type: application/json" -d '{"template_id": 1, "name": "testParameter11", "value_string": {"String": "TestParameter11", "Valid": true}, "value_int": {"Int64": 1, "Valid": true}}'
$ curl -X POST "localhost:8080/template_persistent_parameters" -H "Content-Type: application/json" -d '{"template_id": 1, "name": "testParameter12", "value_string": {"String": "TestParameter12", "Valid": true}}'
$ curl -X POST "localhost:8080/template_persistent_parameters" -H "Content-Type: application/json" -d '{"template_id": 1, "name": "testParameter1X", "value_int": {"Int64": 100, "Valid": true}}'
$ # register template2 and persistent parameters
$ curl -X POST "localhost:8080/templates" -H "Content-Type: application/json" -d '{"name": "sample2", "template_content": "{{.testParameter1X}}"}'
$ curl -X POST "localhost:8080/template_persistent_parameters" -H "Content-Type: application/json" -d '{"template_id": 2, "name": "testParameter1X", "value_int": {"Int64": 200, "Valid": true}}'
$ # register template3 and persistent parameters
$ curl -X POST "localhost:8080/templates" -H "Content-Type: application/json" -d '{"name": "sample3", "template_content": "{{.testParameter1X}}"}'
$ curl -X POST "localhost:8080/template_persistent_parameters" -H "Content-Type: application/json" -d '{"template_id": 3, "name": "testParameter1X", "value_int": {"Int64": 300, "Valid": true}}'
$ # show generated template
$ curl -X GET "localhost:8080/templates/1"
$ # Geenrate a text from the tempalte
$ curl -X GET "localhost:8080/templates/1/generation"
```

When you generate a text from template, you can give volatile parameters with query parameter to templates, and templates can use those as slice.

```
$ curl -X GET "localhost:8080/templates/1/generation?param1=100&param1=101&param2=200"
$ # In template, those parameters can be accessed as {{.param1}} and {{.param2}}, as slices [100 101] and [200].
```

If you added some models like [Loam](https://github.com/qb0C80aE/loam), you will be able to use those models in templates.
See [the example template in Loam](https://github.com/qb0C80aE/loam/blob/develop/examples/terraform.template).

# API Server

Simple Rest API using [Gin](https://github.com/gin-gonic/gin)(framework) & [GORM](https://github.com/jinzhu/gorm)(orm)

## Endpoint list

### Designs Resource

```
GET    /designs/present
PUT    /designs/present
DELETE /designs/present
```

### TemplatePersistentParameter Resource

```
GET    /template_persistent_parameters
GET    /template_persistent_parameters/:id
POST   /template_persistent_parameters
PUT    /template_persistent_parameters/:id
DELETE /template_persistent_parameters/:id
```

### Template Resource

```
GET    /templates
GET    /templates/:id
POST   /templates
PUT    /templates/:id
DELETE /templates/:id
GET    /templates/:id/generation
```

# Thanks

* The base part of Clay was generated by [apig](https://github.com/wantedly/apig)
* Clay is using [Glide](https://github.com/Masterminds/glide) to manage dependencies of packages
