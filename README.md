# MyDaemon

这是一个使用Go语言编写的守护进程

This is a daemon program written in Golang.

这个守护进程可以包裹住任意一个进程，当进程结束时，守护进程可以输出最近几条日志，并将这些日志发送到指定的邮箱中。

This daemon can wrap any process, and when the process ends, the daemon can output the last few logs and send them to the specified mailbox.

这个守护进程特别适用于通知一个长时间的训练任务完成，并直接看到训练结果。

This daemon is particularly useful for notifying a long training task of completion and seeing the results directly.

## Usage

1. config email
```sh
mydaemon --config --from "srcEmail@example.com" --password "srcEmailSMTPPassword" --to "dstEmail@example.com"
```

2. start process
```sh
mydaemon "python train.py"
```
