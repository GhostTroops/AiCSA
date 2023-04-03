[![Tweet](https://img.shields.io/twitter/url/http/Hktalent3135773.svg?style=social)](https://twitter.com/intent/follow?screen_name=Hktalent3135773) [![Follow on Twitter](https://img.shields.io/twitter/follow/Hktalent3135773.svg?style=social&label=Follow)](https://twitter.com/intent/follow?screen_name=Hktalent3135773) [![GitHub Followers](https://img.shields.io/github/followers/hktalent.svg?style=social&label=Follow)](https://github.com/hktalent/)
# Ai(ChatGPT-4) Code Security Audit

<img width="1106" alt="image" src="https://user-images.githubusercontent.com/18223385/229397981-b0eab8a6-9635-4520-8e1a-d11e1c3ffcfe.png">


# feature
- 相同 jar、相同 java 文件，chatGPT ( GPT-4 ) 只执行一次，结果保留在索引库中,所以不用担心多次重复执行的问题
- 免费的 chatGPT 限速20次/分钟，付费用户可以通过修改 config/config.json 调整频率
- 文件大于 3500 字节自动拆分发送给 chatGPT,避免过长的文件导致 chatGPT 无法处理
- 支持 若干个 openai api key，提高并发能力
- 基于大数据索引存储结果
- 提供 HTTP/2.0 HTTP/3.0 web 界面

# web UI
<img width="715" alt="image" src="https://user-images.githubusercontent.com/18223385/229487667-acdfdfdb-6125-4806-9666-09ecd349e82a.png">

```
mkdir -p src config
vi config/config.json
 ./AiCSA  
 
open https://127.0.0.1:8080/indexes/
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

## config/config.json example
LimitPerMinute: 建议 api key 个数 * 3
```
{
  "proxy": "socks5://127.0.0.1:7890",
  "LimitPerMinute": 6,
  "HttpPort": 8080,
  "org": "org-xx",
  "api_key": "sk-xxx,sk-xxx2",
  "Prefix": "用中文问答，分析%s java代码存在哪些安全风险,如何验证、确认他们",
  "CheckRpt": true
}
```

# How build
```
go get -u ./...
go mod vendor
go build -o AiCSA main.go
```

## 反编译jar to java
- 源码将自动保存在 src 目录中
- 不同的 ja r会根据hash构建一个源码目录，避免多个jar的源码冲突

```
find $HOME/MyWork/vulScanPro/tools/weblogic/weblogic12.2.1.3 -type f -name "*.jar" | xargs -I {} ./tools/doFernflower.sh {}
ls $HOME/MyWork/vulScanPro/tools/weblogic/weblogic12.2.1.3/coherence/lib/*.jar|xargs -I {} ./tools/doFernflower.sh {}
./tools/doFernflower.sh $HOME/MyWork/vulScanPro/tools/weblogic/weblogic12.2.1.3/coherence/lib/coherence.jar
```

# Tips
- Mac OS 所有子目录图片转换为mp4
```
brew install ffmpeg
brew update && brew upgrade ffmpeg

find $HOME/Downloads/outImg -name '*.png' | sort | sed 's/.*/"&"/' | tr '\n' ' ' | xargs ffmpeg -r 30 -i - -c:v libx264 -pix_fmt yuv420p output.mp4
```

## 💖Star
[![Stargazers over time](https://starchart.cc/hktalent/AiCSA.svg)](https://starchart.cc/hktalent/AiCSA)

# Donation
| Wechat Pay | AliPay | Paypal | BTC Pay |BCH Pay |
| --- | --- | --- | --- | --- |
|<img src=https://raw.githubusercontent.com/hktalent/myhktools/main/md/wc.png>|<img width=166 src=https://raw.githubusercontent.com/hktalent/myhktools/main/md/zfb.png>|[paypal](https://www.paypal.me/pwned2019) **miracletalent@gmail.com**|<img width=166 src=https://raw.githubusercontent.com/hktalent/myhktools/main/md/BTC.png>|<img width=166 src=https://raw.githubusercontent.com/hktalent/myhktools/main/md/BCH.jpg>|

