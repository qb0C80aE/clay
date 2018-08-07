#!/bin/bash

curl -X DELETE -H "AuthKey: ${AUTH_KEY}" https://clay-download.herokuapp.com/clay.baster-20180808.linux-amd64.tgz
curl -X DELETE -H "AuthKey: ${AUTH_KEY}" https://clay-download.herokuapp.com/clay.baster-20180808.linux-386.tgz
curl -X DELETE -H "AuthKey: ${AUTH_KEY}" https://clay-download.herokuapp.com/clay.baster-20180808.windows-amd64.zip
curl -X DELETE -H "AuthKey: ${AUTH_KEY}" https://clay-download.herokuapp.com/clay.baster-20180808.windows-386.zip
curl -X DELETE -H "AuthKey: ${AUTH_KEY}" https://clay-download.herokuapp.com/clay.baster-20180808.darwin-amd64.tgz
curl -X DELETE -H "AuthKey: ${AUTH_KEY}" https://clay-download.herokuapp.com/clay.baster-20180808.darwin-386.tgz
curl -X PUT -H "AuthKey: ${AUTH_KEY}" -H "Content-Type: multipart/form-data" https://clay-download.herokuapp.com/clay.baster-20180808.linux-amd64.tgz -F content=@clay.bastet.linux-amd64.tgz -F description="Clay binary for bastet:Linux/x86_64"
curl -X PUT -H "AuthKey: ${AUTH_KEY}" -H "Content-Type: multipart/form-data" https://clay-download.herokuapp.com/clay.baster-20180808.linux-386.tgz -F content=@clay.bastet.linux-386.tgz -F description="Clay binary for bastet:Linux/x86"
curl -X PUT -H "AuthKey: ${AUTH_KEY}" -H "Content-Type: multipart/form-data" https://clay-download.herokuapp.com/clay.baster-20180808.windows-amd64.zip -F content=@clay.bastet.windows-amd64.zip -F description="Clay binary for bastet:Windows/x86_64"
curl -X PUT -H "AuthKey: ${AUTH_KEY}" -H "Content-Type: multipart/form-data" https://clay-download.herokuapp.com/clay.baster-20180808.windows-386.zip -F content=@clay.bastet.windows-386.zip -F description="Clay binary for bastet:Windows/x86"
curl -X PUT -H "AuthKey: ${AUTH_KEY}" -H "Content-Type: multipart/form-data" https://clay-download.herokuapp.com/clay.baster-20180808.darwin-amd64.tgz -F content=@clay.bastet.darwin-amd64.tgz -F description="Clay binary for bastet:Mac OS X/x86_64"
curl -X PUT -H "AuthKey: ${AUTH_KEY}" -H "Content-Type: multipart/form-data" https://clay-download.herokuapp.com/clay.baster-20180808.darwin-386.tgz -F content=@clay.bastet.darwin-386.tgz -F description="Clay binary for bastet:Max OS X/x86"
