package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"takeout/common/constant"
	"takeout/common/logger"
)

// 定义websocket的连接配置
var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有的跨域请求，这里可以做限制
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Server websocketServer
type Server struct {
	conns   map[string]*websocket.Conn // 使用map来管理连接
	handler WSEventHandler             // WebSocket事件处理程序
	mutex   sync.Mutex                 // 互斥锁，用于保护conns的并发访问
}

// WSEventHandler 定义WebSocket事件处理接口
type WSEventHandler interface {
	OnOpen(clientId string, conn *websocket.Conn)
	OnMessage(clientId string, message string)
	OnClose(clientId string)
}

// WSServer 创建websocketServer并初始化
var WSServer = Server{
	conns:   make(map[string]*websocket.Conn),
	handler: &DefaultEventHandler{},
	mutex:   sync.Mutex{},
}

// DefaultEventHandler 基本的WebSocket事件处理程序实现
type DefaultEventHandler struct{}

func (h *DefaultEventHandler) OnOpen(clientId string, conn *websocket.Conn) {
	logger.Info("websocket连接成功", zap.String("sid", clientId))
	// 将连接存储到map中
	WSServer.mutex.Lock()
	defer WSServer.mutex.Unlock()
	WSServer.conns[clientId] = conn
}

func (h *DefaultEventHandler) OnMessage(clientId string, message string) {
	logger.Info("websocket客户端发送消息: %s", zap.String("sid", clientId), zap.String("message", message))
}

func (h *DefaultEventHandler) OnClose(clientId string) {
	WSServer.mutex.Lock()
	defer WSServer.mutex.Unlock()
	if conn, ok := WSServer.conns[clientId]; ok {
		err := conn.Close() // 关闭连接
		if err != nil {
			logger.Error("连接关闭错误", zap.Error(err))
		}
		delete(WSServer.conns, clientId) // 从map中删除连接
		logger.Info("websocket连接已断开", zap.String("sid", clientId))
	}
}

// WSHandler WebSocket处理器
func WSHandler(ctx *gin.Context) {
	// 升级HTTP连接到WebSocket连接
	conn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.Error("WebSocket升级失败", zap.Error(err))
		return
	}
	clientId := ctx.Param("sid")
	if clientId == "" {
		logger.Error(constant.MsgMissingRequest)
		return
	}

	// 调用OnOpen事件
	WSServer.handler.OnOpen(clientId, conn)

	// 处理连接关闭
	defer func() {
		WSServer.handler.OnClose(clientId)
	}()

	// 监听客户端消息
	for {
		_, message, e := conn.ReadMessage()
		if e != nil {
			logger.Info("websocket读取消息失败", zap.String("sid", clientId), zap.Error(err))
			break // 读取失败，跳出循环
		}

		// 调用OnMessage事件
		WSServer.handler.OnMessage(clientId, string(message))
	}
}

// SendToAllClients 发送消息给所有客户端
func SendToAllClients(jsonMsg any) {
	WSServer.mutex.Lock()
	defer WSServer.mutex.Unlock()

	for clientId, conn := range WSServer.conns {
		if err := conn.WriteJSON(jsonMsg); err != nil {
			e := conn.Close()
			if e != nil {
				logger.Error("连接关闭错误", zap.Error(err))
			}
			delete(WSServer.conns, clientId)
			logger.Info("websocket发送消息失败，已移除连接", zap.String("sid", clientId))
		}
	}
}
