#! /bin/sh
if [ "$(uname)" == "Darwin" ]; then # Mac OS X
    go build -o convert main.go
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then # Linux
    go build -o convert main.go
elif [ "$(expr substr $(uname -s) 1 10)" == "MINGW32_NT" ];then # Windows
    go build -o convert.exe main.go
fi