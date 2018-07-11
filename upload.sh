#!/bin/bash
curl -X DELETE -H "AuthKey: ${AUTH_KEY}" https://clay-download.herokuapp.com/clay.bastet.linux-amd64.tgz
curl -X DELETE -H "AuthKey: ${AUTH_KEY}" https://clay-download.herokuapp.com/clay.bastet.linux-386.tgz
curl -X DELETE -H "AuthKey: ${AUTH_KEY}" https://clay-download.herokuapp.com/clay.bastet.windows-amd64.zip
curl -X DELETE -H "AuthKey: ${AUTH_KEY}" https://clay-download.herokuapp.com/clay.bastet.windows-386.zip

md5sum clay.bastet.linux-amd64.tgz > clay.bastet.linux-amd64.tgz.md5sum
md5sum clay.bastet.linux-386.tgz > clay.bastet.linux-386.tgz.md5sum
md5sum clay.bastet.windows-amd64.zip > clay.bastet.windows-amd64.zip.md5sum
md5sum clay.bastet.windows-386.zip > clay.bastet.windows-386.zip.md5sum

curl -X PUT -H "AuthKey: ${AUTH_KEY}" -H "Content-Type: multipart/form-data" https://clay-download.herokuapp.com/clay.bastet.linux-amd64.tgz -F content=@clay.bastet.linux-amd64.tgz -F description="Clay binary for bastet:Linux/x86_64"
curl -X PUT -H "AuthKey: ${AUTH_KEY}" -H "Content-Type: multipart/form-data" https://clay-download.herokuapp.com/clay.bastet.linux-386.tgz -F content=@clay.bastet.linux-386.tgz -F description="Clay binary for bastet:Linux/x86"
curl -X PUT -H "AuthKey: ${AUTH_KEY}" -H "Content-Type: multipart/form-data" https://clay-download.herokuapp.com/clay.bastet.windows-amd64.zip -F content=@clay.bastet.windows-amd64.zip -F description="Clay binary for bastet:Windows/x86_64"
curl -X PUT -H "AuthKey: ${AUTH_KEY}" -H "Content-Type: multipart/form-data" https://clay-download.herokuapp.com/clay.bastet.windows-386.zip -F content=@clay.bastet.windows-386.zip -F description="Clay binary for bastet:Windows/x86"

curl -X PUT -H "AuthKey: ${AUTH_KEY}" -H "Content-Type: multipart/form-data" https://clay-download.herokuapp.com/clay.bastet.linux-amd64.tgz.md5sum -F content=@clay.bastet.linux-amd64.tgz.md5sum -F description="Clay binary for bastet:Linux/x86_64 MD5SUM"
curl -X PUT -H "AuthKey: ${AUTH_KEY}" -H "Content-Type: multipart/form-data" https://clay-download.herokuapp.com/clay.bastet.linux-386.tgz.md5sum -F content=@clay.bastet.linux-386.tgz.md5sum -F description="Clay binary for bastet:Linux/x86 MD5SUM"
curl -X PUT -H "AuthKey: ${AUTH_KEY}" -H "Content-Type: multipart/form-data" https://clay-download.herokuapp.com/clay.bastet.windows-amd64.zip.md5sum -F content=@clay.bastet.windows-amd64.zip.md5sum -F description="Clay binary for bastet:Windows/x86_64 MD5SUM"
curl -X PUT -H "AuthKey: ${AUTH_KEY}" -H "Content-Type: multipart/form-data" https://clay-download.herokuapp.com/clay.bastet.windows-386.zip.md5sum -F content=@clay.bastet.windows-386.zip.md5sum -F description="Clay binary for bastet:Windows/x86 MD5SUM"


