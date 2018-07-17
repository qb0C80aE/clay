# Clay

[![Build Status](https://travis-ci.org/qb0C80aE/clay.svg?branch=develop)](https://travis-ci.org/qb0C80aE/clay)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Clay is an abstract system model store to automate various kind of operations.  
It provides some APIs to access the system model store.

<br>
<br>

![Logo](https://github.com/qb0C80aE/clay/raw/develop/images/logo.jpg)
<br>
[![License: CC BY-SA 4.0](https://img.shields.io/badge/License-CC%20BY--SA%204.0-lightgrey.svg)](https://creativecommons.org/licenses/by-sa/4.0/)
<br>
logo: By Derzsi Elekes Andor (Own work)
<br>
via Wikimedia Commons

## Concept ans usecases

![Concept](https://github.com/qb0C80aE/clay/raw/develop/images/concept.png)
![Usecase](https://github.com/qb0C80aE/clay/raw/develop/images/usecase.png)

## Download binaries

* [Linux/64bit](http://download.clay.dynu.net/clay.linux-amd64.tgz)
* [Linux/32bit](http://download.clay.dynu.net/clay.linux-386.tgz)
* [Windows/64bit](http://download.clay.dynu.net/clay.windows-amd64.zip)
* [Windows/32bit](http://download.clay.dynu.net/clay.windows-386.zip)
* [MacOS/64bit](http://download.clay.dynu.net/clay.darwin-amd64.tgz)
* [MacOS/32bit](http://download.clay.dynu.net/clay.darwin-386.tgz)

## Demo

### Live (GUI)

You can experience a live demo [here](https://clay-demo.herokuapp.com/ui).

### Screenshots

There are two examples. The first one is api only example, the second one is api and GUI example.  
Clay expresses this GUI by considering its templates or scripts, objects in its model store as html, js, css, or raw image objects.

If you don't need, you don't have to implement GUI like this, it's just an example.
You can use Clay in various ways you want, such as infrastructure design, graph output for [nwdiag](http://blockdiag.com/en/nwdiag/) using templates, etc.
Of course, you can implement external GUI tool using Clay as a backend, but keep it in your mind that Clay is just a client tool.

#### Sample UI - network design (L1)
![Sample UI - network design (L1)](https://github.com/qb0C80aE/clay/raw/develop/images/sample1.png)

#### Sample UI - network design (VLAN)
![Sample UI - network design (VLAN)](https://github.com/qb0C80aE/clay/raw/develop/images/sample2.png)

#### Sample UI - network design (L3)
![Sample UI - network design (L3)](https://github.com/qb0C80aE/clay/raw/develop/images/sample3.png)

This sample GUI is using [inet-henge](https://github.com/codeout/inet-henge) to draw diagrams.

#### Using template to generate nwdiag json
![nwdiag template](https://github.com/qb0C80aE/clay/raw/develop/images/nwdiag1.png)
![nwdiag output](https://github.com/qb0C80aE/clay/raw/develop/images/nwdiag2.png)

# How to use

Just download Clay binary and run.

```
$ ./clay
```

The server runs at http://localhost:8080 by default.

If you give ``-h`` to Clay, you can see the help.
The default running mode is server, and Clay has also client functions to store or retrieve data.

If it's the first time to use Clay, run examples are recommended to know what you can on Clay.
This operation boots Clay with an example.

```
$ CLAY_CONFIG_FILE_PATH=examples/api_only/clay_config.json ./clay
```

If you want to try sample GUI, do like below.

```
$ CLAY_CONFIG_FILE_PATH=examples/api_and_gui/clay_config.json ./clay
```

Then access ``http://localhost:8080/ui`` in your browser.

## Environmental variables

You can give some environmental variables to Clay in order to control its behavior.

|Key                  |Description                                                                      |Options          |Default           |
|:--------------------|:--------------------------------------------------------------------------------|:----------------|:-----------------|
|CLAY_CONFIG_FILE_PATH|The config file path.                                                            |-                |./clay_config.json|
|CLAY_HOST            |The host to listen.                                                              |-                |localhost         |
|CLAY_PORT            |The port to listen.                                                              |-                |8080              |
|CLAY_DB_MODE         |The indentifier how its db is managed.                                           |memory/file      |file              |
|CLAY_DB_FILE_PATH    |The path where its db file is located. This value is used if DB_MODE=file is set.|-                |clay.db           |
|CLAY_ASSET_MODE      |The indentifier how its assets is managed.                                       |external/internal|external          |

## Basic REST API access

After booting Clay, you can know what endpoints exist by accessing root path.

```
$ curl http://localhost:8080/
[
    "designs_url [DELETE] http://localhost:8080/designs/present",
    "designs_url [GET] http://localhost:8080/designs/present",
    "designs_url [PUT] http://localhost:8080/designs/present",
    "ephemeral_binary_object_raws_url [GET] http://localhost:8080/ephemeral_binary_objects/:name/raw",
    "ephemeral_binary_objects_url [DELETE] http://localhost:8080/ephemeral_binary_objects/:name",
    "ephemeral_binary_objects_url [GET] http://localhost:8080/ephemeral_binary_objects",
    "ephemeral_binary_objects_url [GET] http://localhost:8080/ephemeral_binary_objects/:name",
    "ephemeral_binary_objects_url [POST] http://localhost:8080/ephemeral_binary_objects",
    "ephemeral_binary_objects_url [PUT] http://localhost:8080/ephemeral_binary_objects/:name",
    "ephemeral_script_executions_url [DELETE] http://localhost:8080/ephemeral_scripts/:name/execution",
    "ephemeral_script_executions_url [POST] http://localhost:8080/ephemeral_scripts/:name/execution",
    "ephemeral_script_raws_url [GET] http://localhost:8080/ephemeral_scripts/:name/raw",
    "ephemeral_scripts_url [DELETE] http://localhost:8080/ephemeral_scripts/:name",
    "ephemeral_scripts_url [GET] http://localhost:8080/ephemeral_scripts",
    "ephemeral_scripts_url [GET] http://localhost:8080/ephemeral_scripts/:name",
    "ephemeral_scripts_url [POST] http://localhost:8080/ephemeral_scripts",
    "ephemeral_scripts_url [PUT] http://localhost:8080/ephemeral_scripts/:name",
    "ephemeral_template_generations_url [GET] http://localhost:8080/ephemeral_templates/:name/generation",
    "ephemeral_template_raws_url [GET] http://localhost:8080/ephemeral_templates/:name/raw",
    "ephemeral_templates_url [DELETE] http://localhost:8080/ephemeral_templates/:name",
    "ephemeral_templates_url [GET] http://localhost:8080/ephemeral_templates",
    "ephemeral_templates_url [GET] http://localhost:8080/ephemeral_templates/:name",
    "ephemeral_templates_url [POST] http://localhost:8080/ephemeral_templates",
    "ephemeral_templates_url [PUT] http://localhost:8080/ephemeral_templates/:name",
    "template_arguments_url [DELETE] http://localhost:8080/template_arguments/:key_parameter(default=id)",
    "template_arguments_url [GET] http://localhost:8080/template_arguments",
    "template_arguments_url [GET] http://localhost:8080/template_arguments/:key_parameter(default=id)",
    "template_arguments_url [POST] http://localhost:8080/template_arguments",
    "template_arguments_url [PUT] http://localhost:8080/template_arguments/:key_parameter(default=id)",
    "template_generations_url [GET] http://localhost:8080/templates/:key_parameter(default=id)/generation",
    "template_raws_url [GET] http://localhost:8080/templates/:key_parameter(default=id)/raw",
    "templates_url [DELETE] http://localhost:8080/templates/:key_parameter(default=id)",
    "templates_url [GET] http://localhost:8080/templates",
    "templates_url [GET] http://localhost:8080/templates/:key_parameter(default=id)",
    "templates_url [POST] http://localhost:8080/templates",
    "templates_url [PUT] http://localhost:8080/templates/:key_parameter(default=id)",
    "url_alias_definitions_url [GET] http://localhost:8080/url_alias_definitions",
    "url_alias_definitions_url [GET] http://localhost:8080/url_alias_definitions/:name",
    "url_alias_definitions_url [POST] http://localhost:8080/url_alias_definitions",
    "user_defined_model_definitions_url [GET] http://localhost:8080/user_defined_model_definitions",
    "user_defined_model_definitions_url [GET] http://localhost:8080/user_defined_model_definitions/:type_name",
    "user_defined_model_definitions_url [POST] http://localhost:8080/user_defined_model_definitions"
]
```

You can access Clay with ``GET``, ``POST``, ``PUT``, and ``DELETE`` methods.

Here are examples with ``curl``.

```
$ curl -X POST -H "Content-Type: application/json" "localhost:8080/nodes" -d '{"name": "Node1", "node_type_id": 1}'
$ curl -X POST -H "Content-Type: application/json" "localhost:8080/nodes" -d '{"name": "Node2", "node_type_id": 1, "ports": [{"name": "port1", "number": 1}]}'
$ curl -X POST -H "Content-Type: application/json" "localhost:8080/nodes" -d '{"id": 100, "name": "Node100", "node_type_id": 2}'
$ curl -X POST -H "Content-Type: application/x-yaml" "localhost:8080/nodes" --data-binary @- <<EOF
id: 102
name: Node102
description: ""
node_type_id: 3
EOF
$ curl -X POST -H "Content-Type: multipart/form-data" "localhost:8080/templates" -F name="template1" -F template_content=@examples/api_only/sample.template
$ curl -X PUT -H "Content-Type: application/json" "localhost:8080/nodes/1" -d '{"name": "NodeX", "node_type_id": 2}'
$ curl -X GET "localhost:8080/nodes"
$ curl -X GET "localhost:8080/nodes/1"
$ curl -X DELETE "localhost:8080/nodes/1"
```

When you ``POST`` or ``PUT`` data, you have to specify ``Content-Type`` header to ``application/json``, ``application/x-yaml``, or ``multipart/form-data``.
If you want ``GET`` data from Clay in YAML, specify ``Accept`` header to ``applicaton/x-yaml``.

```
$ curl -H "Accept: application/x-yaml" "localhost:8080/nodes"
```

### Queries for single (like ``/templates/1``) and multiple (like ``/templates``) resource urls

* fields
* preloads
* pretty

Here are examples with ``curl``.

```
$ curl "localhost:8080/nodes?preloads=ports.link.destination_port.node,routes.output_port" # Select nodes and their ports with their link.destination_port.node, and routes with their output_port.
$ curl "localhost:8080/nodes?preloads=ports&fields=id,ports.id,ports.name" # Select nodes and their ports, and extract id of nodes, and id, names of ports.
$ curl "localhost:8080/nodes?pretty" # Select nodes and format results.
```

### Queries for single resource url

* key_parameter

Here are examples with ``curl``.

```
$ curl "localhost:8080/nodes/1" # Select a node which has the id 1. It's the default behavior.
$ curl "localhost:8080/nodes/Node1?key_parameter=name" # Select a node which has the name Node1. It identifies a single resource by using name instead of id.
```

You can know resources which can be identified by ``key_parameter``, by accessing root path.

```
curl "localhost:8080/"
...
nodes_url [GET] http://localhost:8080/nodes/:key_parameter(default=id)
...
```

### Queries for multiple resource url

* q[field_name]
* sort
* first
* page
* limit

Here are examples with ``curl``.

```
$ curl "localhost:8080/templates?q\[name\]=Template1,Template2&q\[description\]=" # Select templates which has the name Template1 or Template2, and description is empty, as multiple results
$ curl "localhost:8080/templates?q\[name\]=%25Template%25" # Select templates which has the name like ...Template... as multiple results
$ curl "localhost:8080/templates?q\[name\]=%21Template1&sort=+id&first" # Select templates which name is not Template1, and sort by id ascend, then pick up the first one as a single result
$ curl "localhost:8080/templates?q\[name\]=null" # Select templates which name is null as multiple results
$ curl "localhost:8080/nodes?limit=2&page=2" # Select nodes limiting the count of result records up to 2 and specifying the result set page. This operation retrieves the record 3-4.
```

## User defined models

The one of main features of Clay is defining models through REST API at runtime or boot time without recompiling Clay binary. The model definition is described in JSON or YAML format, and it's handled as Golang struct in Clay.

Clay's Rest API is using [Gin](https://github.com/gin-gonic/gin)(framework) & [GORM](https://github.com/jinzhu/gorm)(orm). Which means, by defining Golang struct field tags, various features like json/yaml marshal, GORM instruction, or field validation, will be available.

```
{
  "type_name": <string>,                   # Type Name. You can refer this name in field definitions to use user defined types.
  "resource_name": <string>,               # Resouce Name. It's used as REST resource name, and DB table name.
  "to_be_migrated": <bool>,                # If it's true, a table it's name is resource_name will be created in DB, and it will be appear in design.
  "is_controller_enabled": <bool>,         # If it's true, a controller will be enabled, and REST endpoint will appear.
  "is_design_access_disabled": <bool>,     # If it's true, this resource will be ignored at design import/export.
  "is_many_to_many_association": <bool>,   # It it's true, this resouce will be handled as many-to-many relationship object like
  "fields": [                              # Struct fields.
    {
      "name": <string>,                    # Field Name.
      "tag": <string>,                     # Golang struct tag like json, yaml, sql, gorm.
      "type_name": <string>,               # Type Name. You can use int, uint, float, string, bool, and also refer a type_name of model definitionto use user defined types.
      "is_slice": <bool>                   # If it's true, this field will be handled as an slice.
    },
    ...
  ],
  "sql_before_migration": [                # A list of string which will be executed as SQL before model migration.
    <string>, ...                          # You can use ';' in order to execute multiple SQL at once.
  ],
  "sql_after_migration": [                 # A list of string which will be executed as SQL after model migration.
    <string>, ...
  ],
  "sql_where_for_design_extraction": [     # A list of string which will be executed as SQL where sentence at design extraction.
    <string>, ...
  ],
  "sql_where_for_design_deletion": [       # A list of string which will be executed as SQL where sentence at design deletion.
    <string>, ...
  ]
}
```

Here is an example.

```
{
  "type_name": "Node",                   # Type Name. You can refer this name in field definitions to use user defined types.
  "resource_name": "nodes",              # Resouce Name. It's used as REST resource name, and DB table name.
  "to_be_migrated": true,                # If it's true, a table it's name is resource_name will be created in DB, and it will be appear in design.
  "is_controller_enabled": true,         # If it's true, a controller will be enabled, and REST endpoint will appear.
  "is_design_access_disabled": false,    # If it's true, this resource will be ignored at design import/export.
  "is_many_to_many_association": false,  # It it's true, this resouce will be handles as many-to-many relationship object like NodeGroupNodeAssociation(see examples).
  "fields": [                            # Field definitions of Node struct.
    {
      "name": "ID",
      "tag": "json:\"id\" yaml:\"id\" gorm:\"primary_key;auto_increment\"",
      "type_name": "int"
    },
    {
      "name": "Name",
      "tag": "json:\"name\" yaml:\"name\" gorm:\"not null;unique\"",
      "type_name": "string"
    },
    {
      "name": "Description",
      "tag": "json:\"description\" yaml:\"description\"",
      "type_name": "string"
    },
    {
      "name": "Ports",
      "tag": "json:\"ports\" yaml:\"ports\" gorm:\"ForeignKey:NodeID\" validate:\"omitempty,dive\"",
      "type_name": "Port",
      "is_slice": true
    },
    {
      "name": "NodeGroupNodeAssociations",
      "tag": "json:\"node_group_node_associations\" yaml:\"node_group_node_associations\" gorm:\"ForeignKey:NodeID\" validate:\"omitempty,dive\"",
      "type_name": "NodeGroupNodeAssociation",
      "is_slice": true
    },
    {
      "name": "Routes",
      "tag": "json:\"routes\" yaml:\"routes\" gorm:\"ForeignKey:NodeID\" validate:\"omitempty,dive\"",
      "type_name": "NodeRoute",
      "is_slice": true
    }
  ]
}
```

In ``sql_*`` fields, You can use any SQL which is allowed in SQLite3, which means, you can define triggers in order to cleanup unnecessary data, or auto generate related data.

If ``is_controller_enabled`` is true, those endpoints will appear in the endpoint list after registration.

```
$ curl -X POST -H "Content-Type: application/json" http://localhost:8080/user_defined_model_definitions -d @examples/api_only/node.json
$ curl http://localhost:8080/
[
    ...
    "nodes_url [DELETE] http://localhost:8080/nodes/:key_parameter(default=id)",
    "nodes_url [GET] http://localhost:8080/nodes",
    "nodes_url [GET] http://localhost:8080/nodes/:key_parameter(default=id)",
    "nodes_url [POST] http://localhost:8080/nodes",
    "nodes_url [PUT] http://localhost:8080/nodes/:key_parameter(default=id)",
    ...
]
```

## Import and export the design

You can import and export the models you created through ``designs`` resource.
Clay is designed as a standalone modeling tool, and the created design can be stored as human-readable text files in version management repositories like git to make it easier to realize Infrastructure-as-Code.

Of course, the data related to user defined models will be also available in the design resource, if their ``is_design_access_disabled`` field is false.

```
$ # Export the design
$ curl -X GET 'localhost:8080/designs/present?pretty' > design.json
$ # Import and overwrite the design
$ curl -X PUT 'localhost:8080/designs/present' -H 'Content-Type: application/json' -d @design.json
```

If you have booted Clay with examples, you can see the file below that includes user defined model types after export.

```
$ cat design.json
{
    "clay_version": "...",
    "content": {
        "links": [...],                             # User-Defined in examples
        "node_group_node_associations": [...],      # User-Defined in examples
        "node_groups": [...],                       # User-Defined in examples
        "node_routes": [...],                       # User-Defined in examples
        "nodes": [...],                             # User-Defined in examples
        "port_group_port_associations": [...],      # User-Defined in examples
        "port_groups": [...],                       # User-Defined in examples
        "port_ipv4_addresses": [...],               # User-Defined in examples
        "ports": [...],                             # User-Defined in examples
        "subnets": [...],                           # User-Defined in examples
        "template_arguments": [...],
        "templates": [...],
        "vlan_port_associations": [...],            # User-Defined in examples
        "vlans": [...]                              # User-Defined in examples
    }
}
```

Note that SQL triggers will be disabled during design access.
One more thing, import and export will be processed following the sequence that models were defined.

## Templates

You can register some text templates and generate something using the models in Clay.
Some functions are provided in template processing, see an example template in Clay like ``examples/api_only/sample.template``.

```
$ # register template and template arguments
$ curl -X POST -H "Content-Type: multipart/form-data" "localhost:8080/templates" -F id=1 -F name="sample" -F template_content=@examples/api_only/sample.template
$ curl -X POST -H "Content-Type: application/json" "localhost:8080/template_arguments" -d '{"id": 1, "template_id": 1, "name": "testParameter11", "type": "int", "default_value": "1"}'
$ curl -X POST -H "Content-Type: application/json" "localhost:8080/template_arguments" -d '{"id": 2, "template_id": 1, "name": "testParameter12", "type": "int", "default_value": "2"}'
$ # show registered template
$ curl -X GET "localhost:8080/templates/1"
$ # show registered template content
$ curl -X GET "localhost:8080/templates/1/raw"
$ # generate a text from the tempalte
$ curl -X GET "localhost:8080/templates/1/generation"
```

When you generate a text from template, you can override the default value of template arguments by giving parameter queries.

```
$ curl -X GET "localhost:8080/templates/1/generation?p\[testParameter12\]=9999"
```

If you want to get raw data or generated text data in a specific ``Content-Type``, specify ``Accept`` and ``Accept-Charset`` headers.

```
$ curl -X GET -H "Accept: application/x-yaml" -H "Accept-Charset: Shift_JIS" "localhost:8080/templates/1/raw"
$ curl -X GET -H "Accept: application/x-yaml" -H "Accept-Charset: Shift_JIS" "localhost:8080/templates/1/generation"
```

## Ephemeral templates, scripts and objects

Ephemeral something is volatile objects unrelated to ``designs`` resource. You can register those items into Clay, but after rebooting, those will disappear.
Those are used in various purposes like using as image, css, js or html files, or those templates.

EphemeralScript is using [otto](https://github.com/robertkrimen/otto) which processes JavaScript.

See ``examples/api_and_gui/ephemeral_*`` directories to know in detail.

## URL aliases

Clay provides url aliases like or request forwarding.
You can register those information to Clay through REST API.

You can use this feature for various kinds of purpose, such as creating url shortcut, pretending static files.

For example, write ``alias.json`` in a specific format like below.

```
{
  "name": <string>,                 # Name
  "from": <string>,                 # Redirect path from
  "to": <string>,                   # Redirect path to
  "query": <string>,                # Queries added when access is redirected
  "methods": [                      # HTTP methods
    {
      "method": <string>,           # HTTP methods (GET|POST|PUT|DELETE)
      "target_url_type": <string>,  # Target resource url type (multi|single)
      "accept": <string>,           # Accept header value
      "accept_charset": <string>    # Accept-Charset header value
    }
  ]
}
```

Then, ``POST`` it into ``localhost:8080/url_alias_definitions`` to register.

```
$ curl -X POST -H "Accept: application/json" "localhost:8080/url_alias_definitions" -d @alias.json
```

When you use URL aliases, ``Content-Type`` in the response will be determined by ``accept`` and ``accept_charset`` fields.
Note that files will be loaded from the top, which means, you have to regard these loading sequences if used defined models have dependencies like foreign keys, views, or triggers.
See ``clay_config,json`` in ``examples`` to know in detail.

## Configuration file

You have seen it so far that you can register various kinds of items into Clay. In addition, by writing the configuration file and specify it at boot time, Clay can register information in that automatically.

```
{
  "general": {
    "user_defined_models_directory": <string>,
    "ephemeral_templates_directory": <string>,
    "ephemeral_binary_objects_directory": <string>,
    "ephemeral_scripts_directory": <string>
  },
  "user_defined_models": [
    {
      "file_name": <string>
    },
    ...
  ],
  "ephemeral_templates": [
    {
      "name": <string>,
      "file_name": <string>
    },
    ...
  ],
  "ephemeral_binary_objects": [
    {
      "name": <string>,
      "file_name": <string>
    },
    ...
  ],
  "ephemeral_scripts": [
    {
      "name": <string>,
      "file_name": <string>
    }
    ...
  ],
  "url_aliases": [
    {
      "name": <string>,                 # Name
      "from": <string>,                 # Redirect path from
      "to": <string>,                   # Redirect path to
      "query": <string>,                # Queries added when access is redirected
      "methods": [                      # HTTP methods
        {
          "method": <string>,           # HTTP methods (GET|POST|PUT|DELETE)
          "target_url_type": <string>,  # Target resource url type (multi|single)
          "accept": <string>,           # Accept header value
          "accept_charset": <string>    # Accept-Charset header value
        }
      ]
    },
    ...
  ]
}
```

If you want to know concrete ways to write, see ``clay_config.json`` in ``examples`` directory.

# How to build

```bash
$ # Note: Before you build Clay, you need to install a C compiler lilke gcc in order to build go-sqlite3
$ # Note: Suppose that $HOME is /home/user, and $GOPATH is /home/user/go.
$ # Note: Please install dep, go-assets-builder first.
$ go get -u github.com/golang/dep/cmd/dep
$ go get -u github.com/jessevdk/go-assets-builder
$ dep ensure
$ go generate -tags=prebuild ./...
$ # Note: If you want to build Clay as a statically linked single binary file, add the flag like below.
$ # go build --ldflags '-extldflags "-static"'
$ go build
$ ./clay &
```

## Cross-compile

For example, you can build Clay for Linux 32bit, Windows 32bit and 64bit on Ubuntu.

```bash
$ # Suppose that Go is installed, $HOME is /home/user, GOROOT is /usr/local/go, and $GOPATH is /home/user/go.
$ cd $HOME
$ sudo apt-get update
$ # Install required packages.
$ sudo apt-get install -y git wget tar gcc
$ ## For Linux 32bit.
$ sudo apt-get install -y gcc-multilib
$ ## For Windows 64bit and 32bit.
$ sudo apt-get install -y binutils-mingw-w64 mingw-w64
$ # Install go cross-compile environments. It requires go 1.4.
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
$ # Install dep, go-assets-builder.
$ go get -u github.com/golang/dep/cmd/dep
$ go get -u github.com/jessevdk/go-assets-builder
$ mkdir -p $GOPATH/src/github.com/qb0C80aE
$ cd $GOPATH/src/github.com/qb0C80aE
$ git clone https://github.com/qb0C80aE/clay.git
$ cd $GOPATH/src/github.com/qb0C80aE/clay
$ dep ensure
$ go generate -tags=prebuild ./...
$ ## For Linux 32bit.
$ CGO_ENABLED=1 GOOS=linux GOARCH=386 go build --ldflags '-extldflags "-static"' -o linux_386/clay
$ ## For Windows 64bit.
$ CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc LD=x86_64-w64-mingw32-ld GOOS=windows GOARCH=amd64 go build --ldflags '-extldflags "-static"' -o windows_amd64/clay.exe
$ ## For Windows 32bit.
$ CGO_ENABLED=1 CC=i686-w64-mingw32-gcc LD=i686-w64-mingw32-ld GOOS=windows GOARCH=386 go build --ldflags '-extldflags "-static"' -o windows_386/clay.exe
```

## Build on Windows

For example, build Clay using MinGW.

  1. Install msys2 https://msys2.github.io/
  2. Run msys2 shell. i.e. ``C:\mingw64\msys2.exe``

```bash
$ pacman -S mingw-w64-x86_64-gcc
$ cd $GOPATH/src/github.com/qb0C80aE/clay
$ go build
```

Powershell

```powershell
PS> C:\msys64\usr\bin\pacman -S mingw-w64-x86_64-gcc
PS> cd $env:GOPATH/src/github.com/qb0C80aE/clay
PS> powershell { $env:PATH+=";C:\msys64\mingw64\bin"; go build }
```

# Etc.

* Clay is using [dep](https://github.com/golang/dep) to manage dependencies of its packages
* The base part of Clay has been generated by [apig](https://github.com/wantedly/apig)
* Assets in Clay will be archived by [go-assets](github.com/jessevdk/go-assets-builder)
