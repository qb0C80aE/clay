Name: clay
Version: 0.1.0
Release: 1%{?dist}
Summary: An abstract system model store to automate something
Vendor: qb0C80aE
URL: https://github.com/qb0C80aE/clay
Source: https://github.com/qb0C80aE/clay.git
License: MIT

BuildArch: x86_64

BuildRequires: rpmdevtools sqlite-devel
BuildRequires: golang >= 1.6

%{systemd_requires}

%description
An abstract system model store to automate something

%build
# rpmbuild resets $PATH so ensure to have "$GOPATH/bin".
export PATH="$PATH:${GOPATH}/bin"
cd "${GOPATH}/src/github.com/qb0C80aE/clay"
(
  go get
  go build
)

%install
cd "${GOPATH}/src/github.com/qb0C80aE/clay"
mkdir -p "$RPM_BUILD_ROOT"/opt/qb0C80aE/clay/bin
cp clay "$RPM_BUILD_ROOT"/opt/qb0C80aE/clay/bin

%files
%dir /opt/qb0C80aE
%dir /opt/qb0C80aE/clay
%dir /opt/qb0C80aE/clay/bin
/opt/qb0C80aE/clay/bin/clay
