# Coinssh

## Debugging

Run:

buffalo build -t -gcflags="-N" && dlv --listen=:2345 --headless=true --api-version=2 exec bin\coinssh.exe

Then:

Click Debug on GoLand

Reference:

https://blog.gobuffalo.io/debugging-a-buffalo-app-in-gogland-b9a00e8076b8

## RSA Key (JWT)

Run:

ssh-keygen -t rsa -b 4096 -f rsa/admin/jwtRS256.key
ssh-keygen -t rsa -b 4096 -f rsa/web/jwtRS256.key

## Dependency

Run:

dep ensure
yarn install
