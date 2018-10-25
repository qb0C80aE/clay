#!/bin/bash
curl -X PUT -H "AuthKey: ${AUTH_KEY}" -H "Content-Type: multipart/form-data" https://clay-download.herokuapp.com/clay.bastet-old.linux-amd64.tgz -F content=@clay.bastet.linux-amd64.tgz -F description="Clay binary for bastet:Linux/x86_64"

