let socket; // Объявляем переменную для WebSocket

// Функция проверки авторизации
function checkAuth() {
    const token = localStorage.getItem('token');
    if (!token) {
        // Если токен отсутствует, сохраняем текущий URL и перенаправляем на страницу входа и регистрации
        localStorage.setItem('redirectUrl', window.location.href);
        window.location.href = 'RegOrLog.html';
    }

}

// Проверка авторизации при загрузке страницы
checkAuth();
connectWebSocket();

document.getElementById('fetchUsers')?.addEventListener('click', async () => {
    checkAuth(); // Проверка токена перед получением пользователей

    const token = localStorage.getItem('token');
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
function connectWebSocket() {
    const token = localStorage.getItem('token');
    socket = new WebSocket(`ws://localhost:8000/ws?token=${token}`); // Передаем токен в URL

    socket.onopen = function () {
        console.log('WebSocket connection established');
        socket.send(JSON.stringify({ message: 'Client connected' }));
    };

    socket.onmessage = function (event) {
        const message = JSON.parse(event.data);
        console.log('Message received:', message);
    };

    socket.onclose = function () {
        console.log('WebSocket connection closed');
    };

    socket.onerror = function (error) {
        console.error('WebSocket error:', error);
    };
}
// Функция для отправки сообщения через WebSocket
function sendMessage(content, receiverId) {
    if (!socket || socket.readyState !== WebSocket.OPEN) {
        console.error('WebSocket is not connected.');
        alert("Ошибка: соединение с сервером не установлено.");
        return;
    }

    const message = {
        content: content,
        receiver_id: parseInt(receiverId, 10) // Убедитесь, что это число
    };

    socket.send(JSON.stringify(message));
}

// Пример использования отправки сообщения
document.getElementById('sendMessageForm')?.addEventListener('submit', (e) => {
    e.preventDefault();
    checkAuth(); // Проверка токена перед отправкой сообщения
    const content = document.getElementById('messageContent').value;
    const receiverId = document.getElementById('receiverId').value; // ID получателя

    // Проверка на наличие содержимого сообщения и ID получателя
    if (!content || !receiverId) {
        alert("Пожалуйста, заполните все поля.");
        return;
    }

    sendMessage(content, receiverId);
});