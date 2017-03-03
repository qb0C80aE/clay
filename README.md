# Summary

Clay is an abstract system model store to automate various kind of operations.
It provides some APIs to access the system model store.

# How to build and run

```
$ cd $GOPATH/src/github.com/qb0C80aE/clay
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

## Create node types and node kinds that express the node is virtual or physical.

To create new resources, send `POST` request to specific resource urls of clay. These rules are the same about almost all resources.
In this case, the target resources are `node_types` and `node_pvs`.

```
$ curl -X POST 'localhost:8080/v1/node_types' -H 'Content-Type: application/json' -d '{"id": 1, "name": "L2Switch"}'
$ curl -X POST 'localhost:8080/v1/node_types' -H 'Content-Type: application/json' -d '{"id": 2, "name": "L3Switch"}'
$ curl -X POST 'localhost:8080/v1/node_types' -H 'Content-Type: application/json' -d '{"id": 3, "name": "Firewall"}'
$ curl -X POST 'localhost:8080/v1/node_types' -H 'Content-Type: application/json' -d '{"id": 4, "name": "Router"}'
$ curl -X POST 'localhost:8080/v1/node_types' -H 'Content-Type: application/json' -d '{"id": 5, "name": "LoadBalancer"}'
$ curl -X POST 'localhost:8080/v1/node_types' -H 'Content-Type: application/json' -d '{"id": 6, "name": "Server"}'
$ curl -X POST 'localhost:8080/v1/node_types' -H 'Content-Type: application/json' -d '{"id": 7, "name": "Network"}'
$ curl -X POST 'localhost:8080/v1/node_pvs' -H 'Content-Type: application/json' -d '{"id": 1, "name": "physical"}'
$ curl -X POST 'localhost:8080/v1/node_pvs' -H 'Content-Type: application/json' -d '{"id": 2, "name": "virtual"}'
```

## Create a node as a L2Switch

Send a `POST` request to the `nodes` resource.
Remember that if you omit the `id` parameter, clay generates sequencial ids automatically.

```
$ curl -X POST 'localhost:8080/v1/nodes' -H 'Content-Type: application/json' -d '{"id": 1, "name": "l2sw1", "node_pv_id": 1, "node_type_id": 1}'
```

To inquire the node what you created, send a `GET` request to `nodes`.

```
$ # To get all nodes
$ curl -X GET 'localhost:8080/v1/nodes'
$ # To get a specific node
$ curl -X GET 'localhost:8080/v1/nodes/1'
```

By appending `pretty` parameter, clay outputs the formatted result.

```
$ curl -X GET 'localhost:8080/v1/nodes?pretty'
```

If you want to get a `node` with its child elements like `node_types` and `node_pvs`, append `preloads` parameter.

```
$ curl -X GET 'localhost:8080/v1/nodes?pretty&preloads=NodeType,NodePv'
```

Additionally, there are `sort` parameter to sort the result, and `field` parameter to get specific fields of the result.

```
$ curl -X GET 'localhost:8080/v1/nodes?pretty&preloads=NodeType,NodePv&sort=node_type_id&fields=name,node_type.name,node_pv.name'
```

## Create a port on a node

Send a `POST` request to the `ports` resource. The port resources require their parent node's id as `node_id` parameter.

```
curl -X  POST 'localhost:8080/v1/ports' -H 'Content-Type: application/json' -d '{"id": 1, "node_id": 1, "name": "port0"}'
```

If you want to create a node with its ports at once, you can that like below.

```
$ curl -X POST 'localhost:8080/v1/nodes' -H 'Content-Type: application/json' -d '{"id": 1, "name": "l2sw1", "node_pv_id": 1, "node_type_id": 1, "ports": [{"id": 1, "name": "port0"}]}'
```

In this case, you don't need to specify the `node_id` it the child port parameters because clay automatically assign that to them.


## Update a node

Send a `PUT` request to the `nodes` resource.

```
$ curl -X PUT 'localhost:8080/v1/nodes/1' -H 'Content-Type: application/json' -d '{"name": "changedname", "node_pv_id": 1, "node_type_id": 2}'
```

Remember that if you omit parameters what their type is primitive like string or integer, for example `name` or `node_pv_id`, these fields are set to `nil`, a empty string, or zero.
For example, the case below, the `name` of this node will be set to `nil`.
So even if want to update a specific field, you need to set the other primitive fields to keep their values.

```
$ curl -X PUT 'localhost:8080/v1/nodes/1' -H 'Content-Type: application/json' -d '{"node_pv_id": 1, "node_type_id": 2}'
```

When you omit non-primitive parameters like a `node_type`, a `node_pv` or `ports` in a node, you don't need to set those explicitly, but if you want to update include those, this rule is enabled.
For example, in the case below, the port which have the `id` 1, will have empty name and empty primitive field values after sending the request.

```
$ curl -X PUT 'localhost:8080/v1/nodes/1' -H 'Content-Type: application/json' -d '{"name": "changedname", "node_pv_id": 1, "node_type_id": 2, "ports": [{"id": 1}]}'
```

## Delete a node

Send a `DELETE` request to the `nodes` resource.

```
$ curl -X DELETE 'localhost:8080/v1/nodes/1'
```

## Special relationships

* The relation between nodes and ports are cascading, so if you delete a node, ports under that node will be deleted automatically.
* When a port is created or updated as that it has a valid `destionation_port` parameter, the `destination_port` of the related port will be updated automatically. This occurs when ports are deleted as well.

## Example requests

By executing requests below, you can create the following models.

```
Firewall
  - L2Switch1
    - Server1
    - Server2
  - L2Switch2
    - Server3
    - Server4
```

```
$ # Register `node_types` and `node_pvs`.
$ curl -X POST 'localhost:8080/v1/node_types' -H 'Content-Type: application/json' -d '{"id": 1, "name": "L2Switch"}'
$ curl -X POST 'localhost:8080/v1/node_types' -H 'Content-Type: application/json' -d '{"id": 2, "name": "L3Switch"}'
$ curl -X POST 'localhost:8080/v1/node_types' -H 'Content-Type: application/json' -d '{"id": 3, "name": "Firewall"}'
$ curl -X POST 'localhost:8080/v1/node_types' -H 'Content-Type: application/json' -d '{"id": 4, "name": "Router"}'
$ curl -X POST 'localhost:8080/v1/node_types' -H 'Content-Type: application/json' -d '{"id": 5, "name": "LoadBalancer"}'
$ curl -X POST 'localhost:8080/v1/node_types' -H 'Content-Type: application/json' -d '{"id": 6, "name": "Server"}'
$ curl -X POST 'localhost:8080/v1/node_types' -H 'Content-Type: application/json' -d '{"id": 7, "name": "Network"}'
$ curl -X POST 'localhost:8080/v1/node_pvs' -H 'Content-Type: application/json' -d '{"id": 1, "name": "physical"}'
$ curl -X POST 'localhost:8080/v1/node_pvs' -H 'Content-Type: application/json' -d '{"id": 2, "name": "virtual"}'
$ # Create a firewall as firewall1.
$ curl -X POST 'localhost:8080/v1/nodes' -H 'Content-Type: application/json' -d '{"id": 1, "name": "firewall1", "node_pv_id": 1, "node_type_id": 3}'
$ # Create two ports on firewall1, as layer-3 ports, which means these have its `ipv4_address` or something.
$ curl -X POST 'localhost:8080/v1/ports' -H 'Content-Type: application/json' -d '{"id": 1, "node_id": 1, "name": "eth0", "mac_address": {"String": "00:00:00:00:00:01", "Valid": true}, "ipv4_address": {"String": "10.0.0.1", "Valid": true}, "ipv4_prefix": {"Int64": 24, "Valid": true}}'
$ curl -X POST 'localhost:8080/v1/ports' -H 'Content-Type: application/json' -d '{"id": 2, "node_id": 1, "name": "eth1", "mac_address": {"String": "00:00:00:00:00:02", "Valid": true}, "ipv4_address": {"String": "10.0.1.1", "Valid": true}, "ipv4_prefix": {"Int64": 24, "Valid": true}}'
$ # Create a L2Switch as l2sw2
$ curl -X POST 'localhost:8080/v1/nodes' -H 'Content-Type: application/json' -d '{"id": 2, "name": "l2sw1", "node_pv_id": 1, "node_type_id": 1}'
$ # Create three ports on l2sw1, and connect the port named port0 to the eth0 on fierwall1 by specfying the `destination_port_id` parameter.
$ curl -X POST 'localhost:8080/v1/ports' -H 'Content-Type: application/json' -d '{"id": 3, "node_id": 2, "name": "port0", "destination_port_id": {"Int64": 1, "Valid": true}}'
$ curl -X POST 'localhost:8080/v1/ports' -H 'Content-Type: application/json' -d '{"id": 4, "node_id": 2, "name": "port1"}'
$ curl -X POST 'localhost:8080/v1/ports' -H 'Content-Type: application/json' -d '{"id": 5, "node_id": 2, "name": "port2"}'
$ # Create a L2Switch as l2sw2
$ curl -X POST 'localhost:8080/v1/nodes' -H 'Content-Type: application/json' -d '{"id": 3, "name": "l2sw2", "node_pv_id": 1, "node_type_id": 1}'
$ # Create three ports on l2sw2, and connect the port named port0 to the eth1 on fierwall1 by specfying the `destination_port_id` parameter.
$ curl -X POST 'localhost:8080/v1/ports' -H 'Content-Type: application/json' -d '{"id": 6, "node_id": 3, "name": "port0", "destination_port_id": {"Int64": 2, "Valid": true}}'
$ curl -X POST 'localhost:8080/v1/ports' -H 'Content-Type: application/json' -d '{"id": 7, "node_id": 3, "name": "port1"}'
$ curl -X POST 'localhost:8080/v1/ports' -H 'Content-Type: application/json' -d '{"id": 8, "node_id": 3, "name": "port2"}'
$ # Create a server as server1
$ curl -X POST 'localhost:8080/v1/nodes' -H 'Content-Type: application/json' -d '{"id": 4, "name": "server1", "node_pv_id": 1, "node_type_id": 6}'
$ # Create a port named eth0 on server1 as layer-3 port, and connect that to the port1 on l2sw1.
$ curl -X POST 'localhost:8080/v1/ports' -H 'Content-Type: application/json' -d '{"id": 9, "node_id": 4, "name": "eth0", "destination_port_id": {"Int64": 4, "Valid": true}, "mac_address": {"String": "00:00:00:00:01:01", "Valid": true}, "ipv4_address": {"String": "10.0.0.100", "Valid": true}, "ipv4_prefix": {"Int64": 24, "Valid": true}}'
$ # Create a server as server2
$ curl -X POST 'localhost:8080/v1/nodes' -H 'Content-Type: application/json' -d '{"id": 5, "name": "server2", "node_pv_id": 1, "node_type_id": 6}'
$ # Create a port named eth0 on server2 as layer-3 port, and connect that to the port2 on l2sw1.
$ curl -X POST 'localhost:8080/v1/ports' -H 'Content-Type: application/json' -d '{"id": 10, "node_id": 5, "name": "eth0", "destination_port_id": {"Int64": 5, "Valid": true}, "mac_address": {"String": "00:00:00:00:01:02", "Valid": true}, "ipv4_address": {"String": "10.0.0.101", "Valid": true}, "ipv4_prefix": {"Int64": 24, "Valid": true}}'
$ # Create a server as server3
$ curl -X POST 'localhost:8080/v1/nodes' -H 'Content-Type: application/json' -d '{"id": 6, "name": "server3", "node_pv_id": 1, "node_type_id": 6}'
$ # Create a port named eth0 on server3 as layer-3 port, and connect that to the port1 on l2sw2.
$ curl -X POST 'localhost:8080/v1/ports' -H 'Content-Type: application/json' -d '{"id": 11, "node_id": 6, "name": "eth0", "destination_port_id": {"Int64": 7, "Valid": true}, "mac_address": {"String": "00:00:00:00:02:01", "Valid": true}, "ipv4_address": {"String": "10.0.1.100", "Valid": true}, "ipv4_prefix": {"Int64": 24, "Valid": true}}'
$ # Create a server as server4
$ curl -X POST 'localhost:8080/v1/nodes' -H 'Content-Type: application/json' -d '{"id": 7, "name": "server4", "node_pv_id": 1, "node_type_id": 6}'
$ # Create a port named eth0 on server4 as layer-3 port, and connect that to the port2 on l2sw2.
$ curl -X POST 'localhost:8080/v1/ports' -H 'Content-Type: application/json' -d '{"id": 12, "node_id": 7, "name": "eth0", "destination_port_id": {"Int64": 8, "Valid": true}, "mac_address": {"String": "00:00:00:00:02:02", "Valid": true}, "ipv4_address": {"String": "10.0.1.101", "Valid": true}, "ipv4_prefix": {"Int64": 24, "Valid": true}}'
```

## Segments

Now you can get `segments` resources what are created automatically based on the topology.

```
$ curl -X GET 'localhost:8080/v1/segments?pretty'
```

## Import and export the design

You can import and export the models you created through `design` resource.
Clay is designed as a standalone modeling tool, and the created design should be stored as human-readable text files in versioning repositories like git to make it easier to realize infrastructure-as-code.

```
$ # Import and overwrite the design
$ curl -X PUT 'localhost:8080/v1/designs/present' -H 'Content-Type: application/json' -d @examples/design.json
$ # Export the design
$ curl -X GET 'localhost:8080/v1/designs/present?pretty' > design.json
```

## Templates

You can register some text templates and generate something using the models in clay.

```
$ # register template and external parameters
$ curl -X POST "localhost:8080/v1/templates" -H "Content-Type: multipart/form-data" -F name=terraform -F template_content=@examples/terraform.template
$ curl -X POST "localhost:8080/v1/template_external_parameters" -H "Content-Type: application/json" -d '{"template_id": 1, "name": "dpid", "value": "dp-pica8"}'
$ # show generated template
$ curl -X GET "localhost:8080/v1/templates/1"
$ # Geenrate a text from the tempalte
$ curl -X PATCH "localhost:8080/v1/templates/1"
```

# API Server

Simple Rest API using gin(framework) & gorm(orm)

## Endpoint list

### Nodes Resource

```
GET    /<version>/nodes
GET    /<version>/nodes/:id
POST   /<version>/nodes
PUT    /<version>/nodes/:id
DELETE /<version>/nodes/:id
```

### NodeGroups Resource

```
GET    /<version>/node_groups
GET    /<version>/node_groups/:id
POST   /<version>/node_groups
PUT    /<version>/node_groups/:id
DELETE /<version>/node_groups/:id
```

### NodePvs Resource

```
GET    /<version>/node_pvs
GET    /<version>/node_pvs/:id
POST   /<version>/node_pvs
PUT    /<version>/node_pvs/:id
DELETE /<version>/node_pvs/:id
```

### NodeTypes Resource

```
GET    /<version>/node_types
GET    /<version>/node_types/:id
POST   /<version>/node_types
PUT    /<version>/node_types/:id
DELETE /<version>/node_types/:id
```

### Ports Resource

```
GET    /<version>/ports
GET    /<version>/ports/:id
POST   /<version>/ports
PUT    /<version>/ports/:id
DELETE /<version>/ports/:id
```

### Designs Resource

```
GET    /<version>/designs/present
PUT    /<version>/designs/present
DELETE /<version>/designs/present
```

### Segment Resource

```
GET    /<version>/segments
```

### Template Resource

```
GET    /<version>/templates
GET    /<version>/templates/:id
POST   /<version>/templates
PUT    /<version>/templates/:id
DELETE /<version>/templates/:id
PATCH /<version>/templates/:id
```

# Thanks

* Clay was partially generated by https://github.com/wantedly/apig
