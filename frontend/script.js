const BASE_URL = 'http://localhost:8000'; // Указываем базовый URL с портом 8000
let socket; // Объявляем переменную для WebSocket

document.getElementById('registerForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const username = document.getElementById('registerUsername').value;
    const password = document.getElementById('registerPassword').value;

    const response = await fetch(`${BASE_URL}/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
    });

    const data = await response.json();
    alert(data.message || data.error);
});

document.getElementById('loginForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const username = document.getElementById('loginUsername').value;
    const password = document.getElementById('loginPassword').value;

    const response = await fetch(`${BASE_URL}/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
    });

    const data = await response.json();
    if (data.token) {
        // Сохраняем токен в localStorage
        localStorage.setItem('token', data.token);
        console.log("Токен установлен в localStorage:", data.token); // Для отладки
        alert(data.message);

        // После успешного входа подключаемся к WebSocket
        connectWebSocket(data.token);
    } else {
        alert(data.error);
    }
});

function connectWebSocket(token) {
    socket = new WebSocket('ws://localhost:8000/ws');

    socket.onopen = function () {
        console.log('WebSocket connection established');
        // Отправка токена для авторизации (если это поддерживается на сервере)
        socket.send(JSON.stringify({ token: token }));
    };

    socket.onmessage = function (event) {
        const message = JSON.parse(event.data);
        console.log('Message received:', message);
        // Обновите UI, чтобы отобразить новое сообщение
    };

    socket.onclose = function () {
        console.log('WebSocket connection closed');
    };

    socket.onerror = function (error) {
        console.log('WebSocket error:', error);
    };
}

document.getElementById('fetchUsers').addEventListener('click', async () => {
    // Получаем токен из localStorage
    const token = localStorage.getItem('token'); // Извлекаем токен из localStorage

    if (!token) {
        alert("Токен не найден. Пожалуйста, войдите в систему.");
        return;
    }

    const response = await fetch(`${BASE_URL}/users`, {
        method: 'GET', // Явно указываем метод GET
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        }
    });

    if (!response.ok) {
        alert("Ошибка при получении пользователей");
        return;
    }

    const data = await response.json();
    const userList = document.getElementById('userList');
    userList.innerHTML = ''; // Очистка списка перед добавлением новых пользователей

    if (data.length === 0) {
        alert("Нет пользователей для отображения.");
        return;
    }

    data.forEach(user => {
        const li = document.createElement('li');
        li.textContent = user.username;
        userList.appendChild(li);
    });
});

// Функция для отправки сообщения через WebSocket
function sendMessage(content, receiverId) {
    const message = {
        content: content,
        receiver_id: receiverId // Укажите ID получателя
    };
    socket.send(JSON.stringify(message));
}

// Пример использования отправки сообщения
document.getElementById('sendMessageForm').addEventListener('submit', (e) => {
    e.preventDefault();
    const content = document.getElementById('messageContent').value;
    const receiverId = document.getElementById('receiverId').value; // ID получателя

    // Проверка на наличие содержимого сообщения и ID получателя
    if (!content || !receiverId) {
        alert("Пожалуйста, заполните все поля.");
        return;
    }

    sendMessage(content, receiverId);
});