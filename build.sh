#!/bin/bash

APP_NAME="mailsender"
BUILD_DIR="./bin"
SRC_DIR="."

PLATFORMS=(
    "linux/amd64"
    "windows/amd64"
)

mkdir -p $BUILD_DIR
#cd $SRC_DIR

for platform in "${PLATFORMS[@]}"; do
    OS=$(echo $platform | cut -d'/' -f1)
    ARCH=$(echo $platform | cut -d'/' -f2)

    export GOOS=$OS
    export GOARCH=$ARCH

    OUTPUT_NAME="$BUILD_DIR/${APP_NAME}-${OS}-${ARCH}"
    if [ "$OS" = "windows" ]; then
        OUTPUT_NAME="$OUTPUT_NAME.exe"
    fi

    echo "build $OS/$ARCH..."

    # 执行构建
    go build -ldflags="-w -s" -o $OUTPUT_NAME $SRC_DIR

    if [ $? -eq 0 ]; then
        echo "✓ $OS/$ARCH build success"
    else
        echo "✗ $OS/$ARCH build failed"
    fi
done

echo -e "\nbuild finished! saved $BUILD_DIR directory!"