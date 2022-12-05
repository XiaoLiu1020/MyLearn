- [信号](#信号)
- [信号本质](#信号本质)
  - [(1) `SIGHUP`](#1-sighup)
  - [(2) `SIGINT`](#2-sigint)
  - [(9) `SIGKILL`](#9-sigkill)
  - [(20）`SIGTSTP`](#20sigtstp)
- [`Linux`信号量`signal`](#linux信号量signal)


# 信号
`kill -l `可以查看系统所支持的所有信号量

```bash
[root@node1 ~]# kill -l
 1) SIGHUP       2) SIGINT       3) SIGQUIT      4) SIGILL       5) SIGTRAP
 6) SIGABRT      7) SIGBUS       8) SIGFPE       9) SIGKILL     10) SIGUSR1
11) SIGSEGV     12) SIGUSR2     13) SIGPIPE     14) SIGALRM     15) SIGTERM
16) SIGSTKFLT   17) SIGCHLD     18) SIGCONT     19) SIGSTOP     20) SIGTSTP
21) SIGTTIN     22) SIGTTOU     23) SIGURG      24) SIGXCPU     25) SIGXFSZ
26) SIGVTALRM   27) SIGPROF     28) SIGWINCH    29) SIGIO       30) SIGPWR
31) SIGSYS      34) SIGRTMIN    35) SIGRTMIN+1  36) SIGRTMIN+2  37) SIGRTMIN+3
38) SIGRTMIN+4  39) SIGRTMIN+5  40) SIGRTMIN+6  41) SIGRTMIN+7  42) SIGRTMIN+8
43) SIGRTMIN+9  44) SIGRTMIN+10 45) SIGRTMIN+11 46) SIGRTMIN+12 47) SIGRTMIN+13
48) SIGRTMIN+14 49) SIGRTMIN+15 50) SIGRTMAX-14 51) SIGRTMAX-13 52) SIGRTMAX-12
53) SIGRTMAX-11 54) SIGRTMAX-10 55) SIGRTMAX-9  56) SIGRTMAX-8  57) SIGRTMAX-7
58) SIGRTMAX-6  59) SIGRTMAX-5  60) SIGRTMAX-4  61) SIGRTMAX-3  62) SIGRTMAX-2
63) SIGRTMAX-1  64) SIGRTMAX
```

# 信号本质
* --软件层次上对 中断机制 的一种模拟
* 信号是异步的, 是进程间通信机制中唯一的异步通信机制
* 编号为 1 ~ 31 为传统`unix`支持信号, 不可靠信号(非实时信号),剩下为可靠
* 不可靠和可靠的区别在于前者不支持排队, 信号可能丢失
* 编号为`0`的信号量相当于一个`PING`信号，用于检查进程是否存在

## (1) `SIGHUP`
* 在用户终端连接(正常或非正常)结束时发出
* 通常是在终端的控制进程结束时, 通知同一`session`内的各个作业, 这时它们与控制终端不再关联。
> 登录Linux时，系统会分配给登录用户一个终端(Session)。在这个终端运行的所有程序，包括前台进程组和后台进程组，一般都属于这个 Session

* 当用户退出`Linux`登录时，前台进程组和后台有对终端输出的进程将会收到`SIGHUP`信号。这个信号的默认操作为终止进程，因此前台进程组和后台有终端输出的进程就会中止。
> 不过可以捕获这个信号，比如wget能捕获SIGHUP信号，并忽略它，这样就算退出了Linux登录，wget也 能继续下载。此外，对于与终端脱离关系的守护进程，这个信号用于通知它重新读取配置文件。

## (2) `SIGINT`
* 程序终止`(interrupt)`信号, 在用户键入`INTR`字符(通常是`Ctrl-C`)时发出，用于通知前台进程组终止进程

## (9) `SIGKILL`
* 用来立即结束程序的运行. **本信号不能被阻塞、处理和忽略**。如果管理员发现某个进程终止不了，可尝试发送这个信号。

## (20）`SIGTSTP`
停止进程的运行, 但该信号可以被处理和忽略. 用户键入`SUSP`字符时(通常是`Ctrl-Z`)发出这个信号

# `Linux`信号量`signal`
[详解Linux信号量signal--http://www.blogdaren.com/post-1298.html](http://www.blogdaren.com/post-1298.html)