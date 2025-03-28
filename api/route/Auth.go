package route

import (
	"Runa/api/config"
	"Runa/api/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Register функция для регистрации пользователя
func Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	// Сохранение пользователя в базе данных
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login функция для входа пользователя
func Login(c *gin.Context) {
	var user model.User
	var foundUser model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Where("username = ?", user.Username).First(&foundUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Генерация токена
	token, err := generateToken(foundUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

// Генерация JWT токена
func generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Hour * 72).Unix(), // Токен будет действовать 72 часа
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your_secret_key")) // Замените на свой секретный ключ
}

// Middleware для проверки авторизации
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		tokenString = tokenString[len("Bearer "):]
		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("your_secret_key"), nil // Замените на свой секретный ключ
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", claims["id"])
		c.Next()
	}
}

// Получение списка пользователей
func GetUsers(c *gin.Context) {
	var users []model.User
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
		ID       uint   `json:"id"`
		Username string `json:"username"`
	}

	for _, user := range users {
		userResponses = append(userResponses, struct {
			ID       uint   `json:"id"`
			Username string `json:"username"`
		}{
			ID:       user.ID,
			Username: user.Username,
		})
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

	// Устанавливаем заголовки для обновления соединения
	header := http.Header{}
	conn, err := websocket.Upgrade(c.Writer, c.Request, header, 1024, 1024)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return // Соединение не было установлено
	}
	log.Println("WebSocket connection established")

	// Получаем userID из контекста
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("User ID not found")
		conn.Close() // Закрываем соединение, если userID не найден
		return
	}

	// Преобразуем userID к типу uint
	var uid uint
	switch id := userID.(type) {
	case uint:
		uid = id
	case float64:
		uid = uint(id) // Преобразуем float64 в uint
	default:
		log.Println("User ID is not of type uint or float64")
		conn.Close() // Закрываем соединение, если тип неверный
		return
	}

	client := &Client{Conn: conn, UserID: uid}

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
			break // Прерываем цикл при ошибке чтения
		}

		// Устанавливаем ID отправителя
		msg.SenderID = uid

		// Сохраняем сообщение в базе данных
		if err := config.DB.Create(&msg).Error; err != nil {
			log.Println("Error saving message to the database:", err)
			continue
		}

		// Отправляем сообщение конкретному получателю
		mu.Lock() // Блокируем мьютекс для безопасной работы с clients
		for c := range clients {
			if c.UserID == msg.ReceiverID {
				if err := c.Conn.WriteJSON(msg); err != nil {
					log.Println("Error sending message to client:", err)
					c.Conn.Close()
					delete(clients, c)
				}
			}
		}
		mu.Unlock() // Разблокируем мьютекс
	}
}
