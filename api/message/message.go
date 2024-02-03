package message

import (
	"context"
	"log"
	"net/http"

	"example.com/golang_twitter/api/interfaces"
	"example.com/golang_twitter/config"
	"example.com/golang_twitter/db"
	sqlc "example.com/golang_twitter/db/sqlc" // このパスは適宜修正してください
	"example.com/golang_twitter/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Message struct {
	SenderID   int64  `json:"sender_id"`
	ReceiverID int64  `json:"receiver_id"`
	Type       int    `json:"type"`
	Content    []byte `json:"content"`
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func MessageRoutes(router *gin.RouterGroup, cfg config.Config, redisConn *interfaces.RedisConn) {
	dbConn := db.DbConn()       // データベース接続を取得
	queries := sqlc.New(dbConn) // sqlcのQueriesオブジェクトを初期化

	router.POST("/message", func(c *gin.Context) {
		message(c, cfg, redisConn)
	})

	router.GET("/ws", func(c *gin.Context) {
		upgradeToWebSocket(c, cfg, redisConn)
	})

	go handleMessages(queries) // sqlc.Queriesオブジェクトを渡す
}

func message(c *gin.Context, cfg config.Config, redisConn *interfaces.RedisConn) {
	var msg Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "メッセージの解析に失敗しました"})
		return
	}
	log.Printf("受信したメッセージ: %+v", msg)
	c.JSON(http.StatusOK, gin.H{"status": "メッセージを受け取りました"})
}

func upgradeToWebSocket(c *gin.Context, cfg config.Config, redisConn *interfaces.RedisConn) {
	// セッションIDからユーザーIDを取得
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		log.Printf("Failed to retrieve session ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "セッションIDの取得に失敗しました"})
		return
	}
	userID, err := utils.GetSessionUserId(c, sessionID)
	if err != nil {
		log.Printf("Failed to retrieve userId from session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "セッションからユーザー情報の取得に失敗しました"})
		return
	}

	// WebSocketへのアップグレード
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocketへのアップグレードに失敗しました: %+v", err)
		return
	}

	// 新しいWebSocket接続をクライアントのリストに追加
	clients[conn] = true
	log.Println("新しいWebSocket接続が追加されました")

	// WebSocket接続を保持し、メッセージの送受信を処理
	go func(conn *websocket.Conn) {
		for {
			t, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("メッセージ読み取りエラー: %v", err)
				break
			}
			broadcast <- Message{
				SenderID:   userID, // 送信者のID
				ReceiverID: 0,      // 受信者のID
				Type:       t,
				Content:    msg,
			}
		}

		// エラーが発生した場合、クライアントの接続を閉じる
		conn.Close()
		delete(clients, conn)
	}(conn)
}

func handleMessages(queries *sqlc.Queries) {
	for {
		msg := <-broadcast
		err := queries.CreateMessage(context.Background(), sqlc.CreateMessageParams{
			SenderID:   int32(msg.SenderID),
			ReceiverID: int32(msg.ReceiverID),
			Content:    string(msg.Content),
		})
		if err != nil {
			log.Printf("メッセージの保存に失敗しました: %v", err)
			continue
		}

		for client := range clients {
			if err := client.WriteMessage(msg.Type, msg.Content); err != nil {
				log.Printf("メッセージの送信に失敗しました: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
