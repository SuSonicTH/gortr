#!/bin/sh
BASE_NAME="gortr"

echo installing goz
go install github.com/SuSonicTH/goz@latest

echo clearing ./bin
rm -rf bin &> /dev/null
echo 

case "$(uname -s)" in
    Linux*)     HOST_OS='linux';;
    Darwin*)    HOST_OS='mac';;
    MINGW*)     HOST_OS='windows';;
    MSYS_NT*)   HOST_OS='windows';;
    *)          HOST_OS="UNKNOWN:${unameOut}"
esac

function build() {
    GOARCH=$1
    GOOS=$2

    echo building $GOARCH-$GOOS

    EXE_NAME="$BASE_NAME"
    if [ "$GOOS" == "windows" ]; then
        EXE_NAME="$BASE_NAME.exe"
    fi
    
    GOARCH=$GOARCH GOOS=$GOOS GOZ_SMALL=1 goz build -o bin/$EXE_NAME &> /dev/null

    echo compressing $GOARCH-$GOOS
    ARCHIVE_NAME="$BASE_NAME-$GOARCH-$GOOS"
    cd bin/
    if [ "$GOOS" == "windows" ]; then 
        if [ "$HOST_OS" == "windows" ]; then
            zip a $ARCHIVE_NAME.zip $EXE_NAME > /dev/null
        else
            zip $ARCHIVE_NAME.zip $EXE_NAME > /dev/null
        fi
    else 
        tar -czf $ARCHIVE_NAME.tgz $EXE_NAME
    fi
    rm $EXE_NAME
    cd ..

    echo
}

build "amd64" "windows"
build "amd64" "linux"
build "arm64" "windows"
build "arm64" "linux"

echo 
echo finished building release packages
ls -la bin/
