package message

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"example.com/golang_twitter/api/interfaces"
	"example.com/golang_twitter/config"
	"example.com/golang_twitter/constants"
	sqlc "example.com/golang_twitter/db/sqlc"
	"example.com/golang_twitter/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Message struct {
	SenderID   int64  `json:"sender_id"`
	ReceiverID int64  `json:"receiver_id"`
	Content    string `json:"content"`
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func GetMessagePage(cfg config.Config, redisConn *interfaces.RedisConn, queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo, err := utils.CurrentUser(c, redisConn, queries)
		if err != nil {
			log.Printf("Failed to retrieve current user: %v", err)
			c.Redirect(303, "/login")
			return
		}

		avatarImage := constants.DefaultAvatarImage
		if userInfo.AvatarImage.Valid {
			avatarImage = userInfo.AvatarImage.String
		}

		// フォローデータ取得
		follows, err := queries.GetFollows(c, userInfo.ID)
		if err != nil {
			log.Printf("Failed to retrieve messages from DB: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "メッセージの取得に失敗しました"})
			return
		}

		c.HTML(200, "message.html", gin.H{
			"userId":      userInfo.ID,
			"displayName": userInfo.DisplayName,
			"avatarImage": avatarImage,
			"follows":     follows,
		})
	}
}

func PostMessage(cfg config.Config, redisConn *interfaces.RedisConn, queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		form := Message{}
		if err := c.ShouldBind(&form); err != nil {
			log.Printf("failed to bind data: %v", err)
			c.JSON(400, gin.H{"error": "処理に失敗しました"})
			return
		}

		// 過去のメッセージデータを取得
		messages, err := queries.GetMessages(c, sqlc.GetMessagesParams{
			SenderID:   form.SenderID,
			ReceiverID: form.ReceiverID,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				// 初回の場合、メッセージがまだ存在しない場合でもエラーを返さない
			} else {
				log.Printf("Failed to retrieve messages from DB: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "メッセージの取得に失敗しました"})
				return
			}
		}

		c.JSON(200, gin.H{"messages": messages})
	}
}

func UpgradeToWebSocket(cfg config.Config, redisConn *interfaces.RedisConn, queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// upgraderを呼び出すことで通常のhttp通信からwebsocketへupgrade
		// コネクションを作成する
		conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocketへのアップグレードに失敗しました: %+v", err)
			return
		}

		// コネクションをclientsマップへ追加
		clients[conn] = true
		log.Println("新しいWebSocket接続が追加されました")

		// WebSocket接続を保持し、メッセージの送受信を処理
		go func(conn *websocket.Conn) {
			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					log.Printf("ReadMessage Error. ERROR: %v", err)
					break
				}

				var message Message
				if err := json.Unmarshal(msg, &message); err != nil {
					log.Printf("Failed to unmarshal message: %v", err)
					continue
				}

				// デシリアライズされたメッセージをbroadcastチャネルに送信
				broadcast <- message
			}

			// エラーが発生した場合、クライアントの接続を閉じる
			conn.Close()
			delete(clients, conn)
		}(conn)
	}
}

func HandleMessages(queries *sqlc.Queries) {
	for {
		msg := <-broadcast
		messages, err := queries.InsertMessage(context.Background(), sqlc.InsertMessageParams{
			SenderID:   msg.SenderID,
			ReceiverID: msg.ReceiverID,
			Content:    msg.Content,
		})
		if err != nil {
			log.Printf("メッセージの保存に失敗しました: %v", err)
			continue
		}

		for client := range clients {
			if err := client.WriteJSON(messages); err != nil {
				log.Printf("メッセージの送信に失敗しました: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
