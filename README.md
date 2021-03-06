### Install

`go get -u -v github.com/gobuffalo/buffalo/buffalo`

### Dependency

`dep ensure`

`yarn install`

### Database

`buffalo db create -a`

### Migration

`buffalo db migrate`

### Seeding

`buffalo task db:seed`

### RSA Key (JWT)

`ssh-keygen -t rsa -b 4096 -f jwtRS256.key`

### Debugging

`buffalo build -t -gcflags="-N" && dlv --listen=:2345 --headless=true --api-version=2 exec bin\coinssh.exe`

Then:

Click debug button on GoLand.

[Reference](https://blog.gobuffalo.io/debugging-a-buffalo-app-in-gogland-b9a00e8076b8)


