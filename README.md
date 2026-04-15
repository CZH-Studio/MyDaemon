# MyDaemon

这是一个使用Go语言编写的守护进程

This is a daemon program written in Golang.

这个守护进程可以包裹住任意一个进程，当进程结束时，守护进程可以输出最近几条日志，并将这些日志发送到指定的邮箱中。

This daemon can wrap any process, and when the process ends, the daemon can output the last few logs and send them to the specified mailbox.

这个守护进程特别适用于通知一个长时间的训练任务完成，并直接看到训练结果。

This daemon is particularly useful for notifying a long training task of completion and seeing the results directly.

## 使用方法

1. 邮箱

   1. 添加邮箱

      ```sh
      mydaemon email add --from xxx@example.com --pwd SMTP_password --to yyy@example.com
      ```

   2. 列出邮箱

      ```sh
      mydaemon email ls
      ```

   3. 移除邮箱

      ```
      mydaemon email rm
      ```

      然后会进入到cli，根据提示操作即可。

2. 配置

   使用`mydaemon config`可修改配置，在后面添加如下参数以及值可以修改对应的配置

   1. 缓存大小：`--buffer`
   2. 日志文件名：`--log`
   3. 邮件发送人名：`--from`
   4. 邮件主题：`--subject`

3. 启动进程

   ```sh
   mydaemon run <command> [args, ...]
   例如
   mydaemon run python main.py
   ```
## Update Log

### v1.0

1. 时间：2026/3/21
2. 功能：
   1. 支持单邮箱配置
   2. 支持程序主要功能（守护进程，记录日志并发送邮件）

### v1.1

1. 时间：2026/4/15
2. 功能：
   1. 支持多邮箱配置
   2. 支持配置缓存长度、日志名等
