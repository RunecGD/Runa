checkAuth();
fetchUsers();

async function fetchUsers() {
    const token = localStorage.getItem('token');
    const response = await fetch(`${BASE_URL}/users`, {
        method: 'GET',
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
        const li = document.createElement('div');
        li.classList.add('user-item'); // Добавляем класс для стилей

        // Создание кнопки для фамилии, имени и отчества
        const nameButton = document.createElement('button');
        nameButton.textContent = `${user.surname} ${user.name} ${user.patronymic}`;
        nameButton.classList.add('name-button'); // Добавляем класс для стилей
        nameButton.onclick = () => {
            location.href = `/frontend/userProfile.html?userId=${user.id}`;
        };
        li.appendChild(nameButton);

        // Создание элемента для специальности
        const detailsText = document.createElement('span');
        detailsText.textContent = user.specialty;
        detailsText.classList.add('specialty-text'); // Класс для стилей
        detailsText.style.display = 'block'; // Переносим специальность на новую строку
        li.appendChild(detailsText);

        userList.appendChild(li);
    });
}

function checkAuth() {
    const token = localStorage.getItem('token');
    if (!token) {
        localStorage.setItem('redirectUrl', window.location.href);
        window.location.href = 'RegOrLog.html';
    }
}