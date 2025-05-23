package route

import (
	"Runa/api/config"
	"Runa/api/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func GetUsers(c *gin.Context) {
	var users []model.User
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not provided"})
		return
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Проверка токена
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("your_secret_key"), nil // Замените на свой секретный ключ
	})
	if err != nil {
		log.Println("Invalid token:", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Получаем userID из утверждений токена
	userID := uint(claims["id"].(float64))

	// Извлечение пользователей из базы данных
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusOK, []string{}) // Возврат пустого списка
		return
	}

	// Преобразование пользователей в ответ
	var userResponses []struct {
		ID         uint   `json:"id"`
		Name       string `json:"name"`
		Surname    string `json:"surname"`
		Patronymic string `json:"patronymic"`
		StudentID  string `json:"student_id"`
		Faculty    string `json:"faculty"`
		Specialty  string `json:"specialty"`
		GroupName  string `json:"group_name"`
		Course     int    `json:"course"`
	}

	for _, user := range users {
		if user.ID != userID {
			userResponses = append(userResponses, struct {
				ID         uint   `json:"id"`
				Name       string `json:"name"`
				Surname    string `json:"surname"`
				Patronymic string `json:"patronymic"`
				StudentID  string `json:"student_id"`
				Faculty    string `json:"faculty"`
				Specialty  string `json:"specialty"`
				GroupName  string `json:"group_name"`
				Course     int    `json:"course"`
			}{
				ID:         user.ID,
				Name:       user.Name,
				Surname:    user.Surname,
				Patronymic: user.Patronymic,
				StudentID:  user.StudentID,
				Faculty:    user.Faculty,
				Specialty:  user.Specialty,
				GroupName:  user.GroupName,
				Course:     user.Course,
			})
		}
	}

	c.JSON(http.StatusOK, userResponses) // Возвращаем пользователей
}

var clients = make(map[*Client]bool) // Подключенные клиенты
var mu sync.Mutex                    // Мьютекс для синхронизации

type Client struct {
	Conn   *websocket.Conn
	UserID uint
}

func HandleWebSocket(c *gin.Context) {
	log.Println("WebSocket connection attempt")

	// Извлекаем токен из параметров URL
	token := c.Query("token")
	if token == "" {
		log.Println("Token not provided")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Проверка токена
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("your_secret_key"), nil // Замените на свой секретный ключ
	})
	if err != nil {
		log.Println("Invalid token:", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Получаем userID из утверждений токена
	userID := uint(claims["id"].(float64))

	// Устанавливаем заголовки для обновления соединения
	header := http.Header{}
	conn, err := websocket.Upgrade(c.Writer, c.Request, header, 1024, 1024)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	log.Println("WebSocket connection established")

	client := &Client{Conn: conn, UserID: userID}

	mu.Lock()
	clients[client] = true
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(clients, client)
		mu.Unlock()
		conn.Close()
	}()

	for {
		var msg model.Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("Error reading JSON:", err)
			break
		}

		// Устанавливаем ID отправителя
		msg.SenderID = userID
		msg.Timestamp = time.Now()
		// Сохраняем сообщение в базе данных
		if err := config.DB.Create(&msg).Error; err != nil {
			log.Println("Error saving message to the database:", err)
			continue
		}

		// Отправляем сообщение конкретному получателю
		mu.Lock()
		for c := range clients {
			if c.UserID == msg.ReceiverID {
				if err := c.Conn.WriteJSON(msg); err != nil {
					log.Println("Error sending message to client:", err)
					c.Conn.Close()
					delete(clients, c)
				}
			}
		}
		mu.Unlock()
	}
}
func GetProfile(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		log.Println("Token not provided")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Проверка токена
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("your_secret_key"), nil // Замените на свой секретный ключ
	})
	if err != nil {
		log.Println("Invalid token:", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Получаем userID из утверждений токена
	userID := uint(claims["id"].(float64))

	// Получение данных пользователя из базы данных
	var user model.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"name":       user.Name,
		"surname":    user.Surname,
		"patronymic": user.Patronymic,
		"student_id": user.StudentID,
		"faculty":    user.Faculty,
		"specialty":  user.Specialty,
		"group_name": user.GroupName,
		"course":     user.Course,
		"photo":      user.Photo,
	})
}

type UpdateProfileRequest struct {
	Surname    string `json:"surname" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Patronymic string `json:"patronymic"`
	StudentID  string `json:"student_id" binding:"required"`
	Course     int    `json:"course" binding:"required"`
	GroupName  string `json:"group_name" binding:"required"`
	Specialty  string `json:"specialty" binding:"required"`
	Faculty    string `json:"faculty" binding:"required"`
}

func UpdateUserProfile(c *gin.Context) {
	var req UpdateProfileRequest

	// Извлечение токена из заголовка
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not provided"})
		return
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Проверка токена
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("your_secret_key"), nil // Замените на свой секретный ключ
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Получаем userID из утверждений токена
	userID := uint(claims["id"].(float64))

	// Привязка данных из запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Получение текущего пользователя из базы данных
	var user model.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Обновляем только переданные поля
	user.Surname = req.Surname
	user.Name = req.Name
	user.Patronymic = req.Patronymic
	user.StudentID = req.StudentID
	user.Course = req.Course
	user.GroupName = req.GroupName
	user.Specialty = req.Specialty
	user.Faculty = req.Faculty

	// Сохранение обновленных данных
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func GetChats(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not provided"})
		return
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("your_secret_key"), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID := uint(claims["id"].(float64))

	var chats []model.Message
	err = config.DB.Where("sender_id = ? OR receiver_id = ?", userID, userID).Order("timestamp DESC").Find(&chats).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	userIDs := make(map[uint]struct{})
	for _, msg := range chats {
		// Добавляем ID отправителя и получателя
		if msg.SenderID != userID {
			userIDs[msg.SenderID] = struct{}{}
		}
		if msg.ReceiverID != userID {
			userIDs[msg.ReceiverID] = struct{}{}
		}
	}

	ids := make([]uint, 0, len(userIDs))
	for id := range userIDs {
		ids = append(ids, id)
	}

	c.JSON(http.StatusOK, gin.H{"user_ids": ids})
}
func GetMessages(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not provided"})
		return
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("your_secret_key"), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID := uint(claims["id"].(float64))
	targetUserID := c.Param("userID") // Получаем ID пользователя из параметров запроса
	targetUserIDUint, err := strconv.ParseUint(targetUserID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var messages []model.Message
	err = config.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", userID, targetUserIDUint, targetUserIDUint, userID).
		Order("timestamp ASC").Find(&messages).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	if len(messages) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Чат пуст"})
		return
	}

	// Возвращаем сообщения, отсортированные по времени
	c.JSON(http.StatusOK, messages)
}
func GetUserInfo(c *gin.Context) {
	targetUserID := c.Param("userID") // Получаем ID пользователя из параметров запроса
	targetUserIDUint, err := strconv.ParseUint(targetUserID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	var user model.User
	err = config.DB.Where("id = ?", targetUserIDUint).Find(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}
	c.JSON(http.StatusOK, user)
}
