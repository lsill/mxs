# 后端用golang，前端用js
1. 逻辑服连接，不考虑短链接，因为后端无法向前端推送消息
2. 使用websocket，js目前支持websocket，微信小程序也可以对接（websocket基于消息队列传输，不需要考虑封包，简单比较实现，可能很快出现性能瓶颈）
