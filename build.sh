#! /bin/bash

go mod tidy
if [[ "$1" == "linux" ]]; then
# build for linux
    GOOS=linux GOARCH=386 go build .
else

#embed config.yml
go-bindata -o config.go config

# embed icon in the executable
rsrc -ico assets/lhIcon/favicon.ico
# build for windows
    GOOS=windows GOARCH=386 go build .
    cp cliLighthouse.exe /mnt/c/Users/Phillip.Crandall/Projects/cliLighthouse/
    mv cliLighthouse.exe /mnt/c/Users/Phillip.Crandall/Desktop/
    cd /mnt/c/Users/Phillip.Crandall/Desktop/
    wslpath -w "cliLighthouse.exe" | sed -e 's/.*/"&"/' | xargs cmd.exe /C start ""
    cd -
fi

