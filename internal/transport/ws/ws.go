package ws

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
	"kirkagram-notification/internal/lib/customResponse"
	"kirkagram-notification/internal/models"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
)

const (
	Like = "like"
	Post = "post"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}

type iLike interface {
	GetUsernameByUserID(userID int) (string, error)
	GetUserIDToSendByPostID(userID int) (int, error)
}

type iPost interface {
	GetUsernameByID(userID int) (string, error)
	GetUsersIDToSendByID(userID int) ([]int, error)
}

type iFollow interface {
	GetUsernameByID(followerID int) (string, error)
}

type WebSocketManager struct {
	clients     map[string]*websocket.Conn
	clientsLock sync.RWMutex
	upgrader    websocket.Upgrader
	like        iLike
	post        iPost
	follow      iFollow
	log         *slog.Logger
}

func NewWebSocketManager(log *slog.Logger, like iLike, post iPost, follow iFollow) *WebSocketManager {
	return &WebSocketManager{
		clients:  make(map[string]*websocket.Conn),
		upgrader: upgrader,
		like:     like,
		post:     post,
		follow:   follow,
		log:      log,
	}
}

func (wm *WebSocketManager) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	const op = "transport.ws.HandleWebSocket"

	log := wm.log.With(slog.String("op", op))
	log.Info("starting handle webSocket")

	userID := chi.URLParam(r, "userID") // Получаем userID из URL параметров
	if userID == "" {
		log.Error("user id is required")

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, customResponse.NewError("user id is required"))

		return
	}

	conn, err := wm.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("failed to upgrade ws", slog.String("error", err.Error()))

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, customResponse.NewError(err.Error()))

		return
	}

	wm.AddClient(userID, conn)

	defer func() {
		wm.RemoveClient(userID)
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (wm *WebSocketManager) AddClient(userID string, conn *websocket.Conn) {
	wm.clientsLock.Lock()
	defer wm.clientsLock.Unlock()
	wm.clients[userID] = conn
}

func (wm *WebSocketManager) RemoveClient(userID string) {
	wm.clientsLock.Lock()
	defer wm.clientsLock.Unlock()
	delete(wm.clients, userID)
}

func (wm *WebSocketManager) SendMessageLike(msg []byte) error {
	var like models.Like

	err := json.Unmarshal(msg, &like)
	if err != nil {
		wm.log.Error("failed to unmarshal post", slog.String("error", err.Error()))

		return err
	}

	userIDToSend, err := wm.like.GetUserIDToSendByPostID(like.PostID)

	wm.clientsLock.RLock()
	conn, exists := wm.clients[strconv.Itoa(userIDToSend)]
	wm.clientsLock.RUnlock()

	if !exists {
		wm.log.Debug("no connection for this user", slog.String("userID", string(userIDToSend)))

		return nil
	}

	username, err := wm.like.GetUsernameByUserID(like.UserID)
	if err != nil {
		wm.log.Error("failed to get username", slog.String("error", err.Error()))

		return err
	}

	message := username + " лайкнул вашу запись"

	return conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (wm *WebSocketManager) SendMessageFollow(msg []byte) error {
	var follow models.Follow

	err := json.Unmarshal(msg, &follow)
	if err != nil {
		wm.log.Error("failed to unmarshal follow", slog.String("error", err.Error()))

		return err
	}

	wm.clientsLock.RLock()
	conn, exists := wm.clients[strconv.Itoa(follow.FollowingID)]
	wm.clientsLock.RUnlock()

	if !exists {
		wm.log.Debug("no connection for this user", slog.Int("followingID", follow.FollowingID))

		return nil
	}

	username, err := wm.follow.GetUsernameByID(follow.FollowerID)

	if err != nil {
		wm.log.Error("failed to get username", slog.String("error", err.Error()))

		return err
	}

	message := username + " подписался на вас"

	return conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (wm *WebSocketManager) SendMessageUnFollow(msg []byte) error {
	var follow models.Follow

	err := json.Unmarshal(msg, &follow)
	if err != nil {
		wm.log.Error("failed to unmarshal follow", slog.String("error", err.Error()))

		return err
	}

	wm.clientsLock.RLock()
	conn, exists := wm.clients[strconv.Itoa(follow.FollowingID)]
	wm.clientsLock.RUnlock()

	if !exists {
		wm.log.Debug("no connection for this user", slog.Int("followingID", follow.FollowingID))

		return nil
	}

	username, err := wm.follow.GetUsernameByID(follow.FollowerID)

	if err != nil {
		wm.log.Error("failed to get username", slog.String("error", err.Error()))

		return err
	}

	message := username + " отписался от вас"

	return conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (wm *WebSocketManager) SendMessagePost(msg []byte) error {
	var post models.Like

	err := json.Unmarshal(msg, &post)
	if err != nil {
		wm.log.Error("failed to unmarshal post")

		return err
	}

	userIDToSendList, err := wm.post.GetUsersIDToSendByID(post.UserID)

	var connList []*websocket.Conn

	for _, userID := range userIDToSendList {
		wm.clientsLock.RLock()
		conn, exists := wm.clients[strconv.Itoa(userID)]
		wm.clientsLock.RUnlock()

		if !exists {
			wm.log.Debug("no connection for this user", slog.Int("userID", userID))
		} else {
			connList = append(connList, conn)
		}
	}

	username, err := wm.like.GetUsernameByUserID(post.UserID)
	if err != nil {
		wm.log.Error("failed to get username", slog.String("error", err.Error()))

		return err
	}

	message := username + " выпустил новую запись"

	var errorsList []error

	for idx, conn := range connList {
		err = conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			wm.log.Error("failed to send message", slog.String("error", err.Error()), slog.Int("idx", idx))
			errorsList = append(errorsList, fmt.Errorf("%s: %w", idx, err))
		}
	}

	if len(errorsList) > 0 {
		// Возвращаем агрегированную ошибку или логируем их
		return fmt.Errorf("encountered %d errors while sending messages: %v", len(errorsList), errorsList)
	}

	return nil
}

func (wm *WebSocketManager) BroadcastMessage(message []byte) {
	wm.clientsLock.RLock()
	defer wm.clientsLock.RUnlock()

	for _, conn := range wm.clients {
		conn.WriteMessage(websocket.TextMessage, message)
	}
}
