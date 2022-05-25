# jarvis

```
[root@www]# ./bin/jarvis 

我是Jarvis，你的得力助理。

Usage:
  jarvis [command]

Available Commands:
  bt          宝塔相关操作
  database    数据库相关操作
  joke        输出一条笑话
  ping        输出pang

Flags:
  -h, --help   help for jarvis
```

运行在持续集成阶段的命令行小助手，典型的使用场景：

1. 提交代码后将项目部署到「宝塔面板」的网站模块
2. 提交代码后自动触发「宝塔面板」的创建数据库功能
3. 提交代码后自动为项目创建数据库

以及正在开发中的更多功能。
