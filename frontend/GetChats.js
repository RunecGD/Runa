checkAuth();
fetchChats();

async function fetchChats() {
    const token = localStorage.getItem('token');
    const response = await fetch(`${BASE_URL}/users/chats`, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        }
    });

    if (!response.ok) {
        alert("Ошибка при получении пользователей для чата");
        return;
    }

    const data = await response.json();
    const userIds = data.user_ids;

    if (userIds.length === 0) {
        document.getElementById('userList').innerHTML = '<div>Нет пользователей для отображения.</div>';
        return;
    }

    // Получаем информацию о всех пользователях
    const userPromises = userIds.map(userId => fetchUserInfo(userId));
    const users = await Promise.all(userPromises);

    // Сортируем пользователей по фамилии, имени
    users.sort((a, b) => {
        return a.surname.localeCompare(b.surname) || a.name.localeCompare(b.name);
    });

    // Отображаем пользователей
    users.forEach(user => {
        displayUser(user);
    });
}

async function fetchUserInfo(userId) {
    const token = localStorage.getItem('token');
    const response = await fetch(`${BASE_URL}/users/${userId}`, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        }
    });

    if (!response.ok) {
        console.error("Ошибка при получении информации о пользователе:", userId);
        return null; // Возвращаем null, если произошла ошибка
    }

    return await response.json(); // Возвращаем данные пользователя
}

function displayUser(user) {
    if (!user) return; // Проверяем на случай, если user равен null

    const userList = document.getElementById('userList');
    const li = document.createElement('div');
    li.classList.add('user-item');

    // Создание кнопки для фамилии, имени и отчества
    const nameButton = document.createElement('button');
    nameButton.textContent = `${user.surname} ${user.name} ${user.patronymic}`;
    nameButton.classList.add('name-button');
    nameButton.onclick = () => {
        location.href = `/frontend/userChat.html?userId=${user.id}`;
    };
    li.appendChild(nameButton);

    // Создание элемента для специальности
    const detailsText = document.createElement('span');
    detailsText.textContent = user.specialty;
    detailsText.classList.add('specialty-text');
    detailsText.style.display = 'block';
    li.appendChild(detailsText);

    userList.appendChild(li);
}

function checkAuth() {
    const token = localStorage.getItem('token');
    if (!token) {
        localStorage.setItem('redirectUrl', window.location.href);
        window.location.href = 'RegOrLog.html';
    }
}