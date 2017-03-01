
## Creating .rpm package file

```
yum install -y yum-utils rpm-build
cd $GOPATH/src/github.com/qb0C80aE/clay
yum-builddep pkg/clay.spec
rpmbuild -bb pkg/clay.spec
```

Then ``$HOME/rpmbuild/RPMS/x86_64/clay-0.1.0-1.el7.centos.x86_64.rpm`` will come out.