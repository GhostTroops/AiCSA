#!/usr/bin/env bash

Java11="/usr/local/Cellar/openjdk@11/11.0.16.1/libexec/openjdk.jdk/Contents/Home/bin/java"
RtJar="/Library/Java/JavaVirtualMachines/openjdk-8.jdk/Contents/Home/jre/lib/rt.jar"

# 当前bashwenj目录
szCurFile=$(basename "$1")
# 运行本bash的当前目录
szCurDir=$(dirname "$(readlink -f "$0")")
# 生成src目录，避免不同jar冲突
srcDir=$(md5sum "$1"|sed 's/ .*//g')
srcDir1="${szCurDir}/../src/${srcDir}"
if [ -d "${srcDir1}" ]; then
    echo "${srcDir1} already exists! Exiting script..."
    exit 1
fi


mkdir -p "${srcDir1}"
${Java11} -jar ${szCurDir}/fernflower.jar -din=1 -hdc=0 -dgs=1 -rsy=1 -lit=1 "$1" -e=${RtJar} "${srcDir1}"

cd ${srcDir1}
unzip -o "${szCurFile}"
rm -rf "${szCurFile}"
