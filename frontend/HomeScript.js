let socket; // Объявляем переменную для WebSocket

// Функция проверки авторизации
function checkAuth() {
    const token = localStorage.getItem('token');
    if (!token) {
        localStorage.setItem('redirectUrl', window.location.href);
        window.location.href = 'RegOrLog.html';
    }

}

checkAuth();
connectWebSocket();

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
    document.getElementById('receiverId').value = '';
    document.getElementById('messageContent').value = '';
});
