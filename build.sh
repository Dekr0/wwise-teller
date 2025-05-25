GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc CGO_LDFLAGS='--static -static-libgcc -static-libstdc++' go build -ldflags .
