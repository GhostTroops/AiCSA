# Ai(ChatGPT-4) Code Security Audit

<img width="1106" alt="image" src="https://user-images.githubusercontent.com/18223385/229397981-b0eab8a6-9635-4520-8e1a-d11e1c3ffcfe.png">


# feature
- 相同 jar、相同 java 文件，chatGPT ( GPT-4 ) 只执行一次，结果保留在索引库中,所以不用担心多次重复执行的问题
- 免费的 chatGPT 限速20次/分钟，付费用户可以通过修改 config/config.json 调整频率
- 文件大于 3500 字节自动拆分发送给 chatGPT,避免过长的文件导致 chatGPT 无法处理

# web UI
```
https://127.0.0.1:8080/indexes/
```

# How Test
- 运行前，请先调整 ./tools/doFernflower.sh 文件，确保 java 是 11 或高版本
- 确定 rt.jar 的路径，修改 ./tools/doFernflower.sh 文件中的 rt.jar 路径

```
find /Library/Java/JavaVirtualMachines -name "rt.jar"
```

out
```
/Library/Java/JavaVirtualMachines/jdk1.8.0_181.jdk/Contents/Home/jre/lib/rt.jar
/Library/Java/JavaVirtualMachines/jdk1.8.0_72.jdk/Contents/Home/jre/lib/rt.jar
```

## 反编译jar to java
- 源码将自动保存在 src 目录中
- 不同的 ja r会根据hash构建一个源码目录，避免多个jar的源码冲突

```
find $HOME/MyWork/vulScanPro/tools/weblogic/weblogic12.2.1.3 -type f -name "*.jar" | xargs -I {} ./tools/doFernflower.sh {}
ls $HOME/MyWork/vulScanPro/tools/weblogic/weblogic12.2.1.3/coherence/lib/*.jar|xargs -I {} ./tools/doFernflower.sh {}
./tools/doFernflower.sh $HOME/MyWork/vulScanPro/tools/weblogic/weblogic12.2.1.3/coherence/lib/coherence.jar
```
