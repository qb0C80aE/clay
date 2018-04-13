#!/bin/bash
curl -X POST localhost:8080/user_defined_model_definitions -H "Content-Type: application/json" -d @node.json
curl -X POST localhost:8080/user_defined_model_definitions -H "Content-Type: application/json" -d @node_group.json
curl -X POST localhost:8080/user_defined_model_definitions -H "Content-Type: application/json" -d @node_group_node_association.json
curl -X POST localhost:8080/user_defined_model_definitions -H "Content-Type: application/json" -d @port.json
curl -X POST localhost:8080/user_defined_model_definitions -H "Content-Type: application/json" -d @port_group.json
curl -X POST localhost:8080/user_defined_model_definitions -H "Content-Type: application/json" -d @port_group_port_association.json
curl -X POST localhost:8080/user_defined_model_definitions -H "Content-Type: application/json" -d @link.json

curl -X POST localhost:8080/user_defined_model_definitions -H "Content-Type: application/json" -d @port_ipv4_address.json
curl -X POST localhost:8080/user_defined_model_definitions -H "Content-Type: application/json" -d @node_route.json

curl -X POST localhost:8080/user_defined_model_definitions -H "Content-Type: application/json" -d @vlan.json
curl -X POST localhost:8080/user_defined_model_definitions -H "Content-Type: application/json" -d @subnet.json
curl -X POST localhost:8080/user_defined_model_definitions -H "Content-Type: application/json" -d @vlan_port_association.json
