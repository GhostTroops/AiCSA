# Ai(ChatGPT-4) Code Security Audit

# test
- 运行前，请先调整./tools/doFernflower.sh文件，确保 java 是11或高版本
- 确定rt.jar的路径，修改./tools/doFernflower.sh文件中的rt.jar路径
```
find /Library/Java/JavaVirtualMachines -name "rt.jar"
```
out
```
/Library/Java/JavaVirtualMachines/jdk1.8.0_181.jdk/Contents/Home/jre/lib/rt.jar
/Library/Java/JavaVirtualMachines/jdk1.8.0_72.jdk/Contents/Home/jre/lib/rt.jar
```
```
./tools/doFernflower.sh $HOME/MyWork/vulScanPro/tools/weblogic/weblogic12.2.1.3/coherence/lib/coherence.jar


```