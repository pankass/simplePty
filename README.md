## golang实现的简易远程终端

### 使用
服务端开启监听
```bash
./server -lport 40000
```

客户端连接
```bash
./client -rhost 127.0.0.1 -rport 40000
```