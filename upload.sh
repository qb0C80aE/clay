#!/bin/bash
curl -X DELETE -H "AuthKey: ${AUTH_KEY}" https://clay-download.herokuapp.com/clay.linux-amd64.tgz
curl -X DELETE -H "AuthKey: ${AUTH_KEY}" https://clay-download.herokuapp.com/clay.linux-386.tgz
curl -X DELETE -H "AuthKey: ${AUTH_KEY}" https://clay-download.herokuapp.com/clay.windows-amd64.zip
curl -X DELETE -H "AuthKey: ${AUTH_KEY}" https://clay-download.herokuapp.com/clay.windows-386.zip
curl -X DELETE -H "AuthKey: ${AUTH_KEY}" https://clay-download.herokuapp.com/clay.darwin-amd64.tgz
curl -X DELETE -H "AuthKey: ${AUTH_KEY}" https://clay-download.herokuapp.com/clay.darwin-386.tgz
curl -X PUT -H "AuthKey: ${AUTH_KEY}" -H "Content-Type: multipart/form-data" https://clay-download.herokuapp.com/clay.linux-amd64.tgz -F content=@clay.linux-amd64.tgz -F description="Clay binary for Linux/x86_64"
curl -X PUT -H "AuthKey: ${AUTH_KEY}" -H "Content-Type: multipart/form-data" https://clay-download.herokuapp.com/clay.linux-386.tgz -F content=@clay.linux-386.tgz -F description="Clay binary for Linux/x86"
curl -X PUT -H "AuthKey: ${AUTH_KEY}" -H "Content-Type: multipart/form-data" https://clay-download.herokuapp.com/clay.windows-amd64.zip -F content=@clay.windows-amd64.zip -F description="Clay binary for Windows/x86_64"
curl -X PUT -H "AuthKey: ${AUTH_KEY}" -H "Content-Type: multipart/form-data" https://clay-download.herokuapp.com/clay.windows-386.zip -F content=@clay.windows-386.zip -F description="Clay binary for Windows/x86"
curl -X PUT -H "AuthKey: ${AUTH_KEY}" -H "Content-Type: multipart/form-data" https://clay-download.herokuapp.com/clay.darwin-amd64.tgz -F content=@clay.darwin-amd64.tgz -F description="Clay binary for Mac OS X/x86_64"
curl -X PUT -H "AuthKey: ${AUTH_KEY}" -H "Content-Type: multipart/form-data" https://clay-download.herokuapp.com/clay.darwin-386.tgz -F content=@clay.darwin-386.tgz -F description="Clay binary for Max OS X/x86"
