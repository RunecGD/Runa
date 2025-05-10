let socket;
const userId = new URLSearchParams(window.location.search).get('userId'); // Получаем userId из URL

document.getElementById('receiverId').value = userId;

connectWebSocket();
fetchMessages();

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
        displayMessage(message.content, message.sender_id);
    };

    socket.onclose = function () {
        console.log('WebSocket connection closed');
    };

    socket.onerror = function (error) {
        console.error('WebSocket error:', error);
    };
}

function fetchMessages() {
    const token = localStorage.getItem('token');
    const receiverId = document.getElementById('receiverId').value;

    fetch(`${BASE_URL}/users/chats/messages/${receiverId}`, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        }
    })
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => {
                    throw new Error(`Ошибка: ${response.status} ${text}`);
                });
            }
            return response.json();
        })
        .then(messages => {
            // Проверяем, есть ли сообщения
            if (messages.length === 0) {
                console.log("Нет сообщений."); // Можно убрать или заменить на alert
                return; // Если сообщений нет, ничего не делаем
            }

            messages.forEach(message => {
                displayMessage(message.content, message.sender_id);
            });
        })
        .catch(error => {
            console.error('Ошибка при получении сообщений:', error);
        });
}
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

document.getElementById('sendMessageForm').addEventListener('submit', (e) => {
    e.preventDefault();
    const content = document.getElementById('messageContent').value;
    const receiverId = document.getElementById('receiverId').value; // ID получателя

    if (!content || !receiverId) {
        alert("Пожалуйста, заполните все поля.");
        return;
    }

    sendMessage(content, receiverId);
    displayMessage(content, userId); // Отображение отправленного сообщения
    document.getElementById('messageContent').value = ''; // Очистка поля ввода
});

function displayMessage(content, senderId) {
    const messageList = document.getElementById('messageList');
    const messageDiv = document.createElement('div');
    messageDiv.classList.add('message');
    messageDiv.classList.add(senderId === parseInt(userId) ? 'sent' : 'received'); // Определяем стиль сообщения

    messageDiv.textContent = content;
    messageList.appendChild(messageDiv);
    messageList.scrollTop = messageList.scrollHeight; // Прокрутка вниз
}