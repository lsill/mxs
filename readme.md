# Zinx最基本的两个模块zoface和znet
- ziface主要是存放一些Zinx框架的全部模块的抽象层接口类，Zinx框架的最基本的是服务类接口iserver，定义在ziface模块中。
- znet模块是zinx框架中网络相关功能的实现，所有网络相关模块都会定义在znet模块中。

1. go1.5暂时不支持quic-go
2. kcp
