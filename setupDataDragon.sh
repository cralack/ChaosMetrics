#!/bin/bash

# 获取当前的版本号
CURRENT_VERSION=$(curl -s https://ddragon.leagueoflegends.com/api/versions.json | grep -oP '^\["\K[^"]+')
echo "The current version is: $CURRENT_VERSION"

# 检查current_version_LOL文件是否存在
if [[ ! -f "current_version_LOL" ]]; then
    echo "current_version_LOL file does not exist. Creating and setting up DataDragon assets..."
    echo "$CURRENT_VERSION" > current_version_LOL
    setupRequired=true
else
    FILE_VERSION=$(cat current_version_LOL)
    echo "Version in current_version_LOL file: $FILE_VERSION"

    # 检查文件中的版本号是否与当前版本号相同
    if [[ "$CURRENT_VERSION" != "$FILE_VERSION" ]]; then
        echo "Version mismatch. Updating the file and setting up DataDragon assets..."
        echo "$CURRENT_VERSION" > current_version_LOL
        setupRequired=true
    else
        echo "Version in the file matches the current version. No download required."
        setupRequired=false
    fi
fi

# 如果需要设置
if [ "$setupRequired" = true ] ; then
    # 定义下载URL
    DOWNLOAD_URL="https://ddragon.leagueoflegends.com/cdn/dragontail-${CURRENT_VERSION}.tgz"
    DOWNLOAD_FILE="dragontail-${CURRENT_VERSION}.tgz"

    # 检查是否已有下载文件
    if [ ! -f "$DOWNLOAD_FILE" ]; then
        echo "Downloading DataDragon assets..."
        wget "${DOWNLOAD_URL}" -O "$DOWNLOAD_FILE" -P
    else
        echo "Download file $DOWNLOAD_FILE already exists. Skipping download."
    fi

    # 创建datadragon目录并解压
    mkdir -p ./datadragon && tar -zxf "$DOWNLOAD_FILE" -C ./datadragon

    # 确保目标目录存在
    mkdir -p "web/src/assets/datadragon"

    # 声明目录数组
    declare -a directories=("champion" "item" "passive" "profileicon" "spell")

    # 循环移动目录
    for dir in "${directories[@]}"; do
        mv "datadragon/${CURRENT_VERSION}/img/${dir}/" "web/src/assets/datadragon/"
    done

    # 移动特定目录
    mv "datadragon/img/champion/" "web/src/assets/datadragon/champion_og"
    mv "datadragon/img/perk-images/" "web/src/assets/datadragon/"

    # 清理下载的压缩包和临时目录
    rm -rf "$DOWNLOAD_FILE" datadragon/

    echo "DataDragon assets have been updated."
fi
