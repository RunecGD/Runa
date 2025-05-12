fetchUserProfile();

function toggleProfileInfo() {
    const info = document.getElementById("profileInfo");
    info.classList.toggle("show");
}

function populateProfile(data) {
    document.getElementById('profileInfo').innerHTML = `
        <p><strong>ФИО:</strong> ${data.surname} ${data.name} ${data.patronymic}</p>
        <p><strong>Номер студ. билета:</strong> ${data.student_id}</p>
        <p><strong>Факультет:</strong> ${data.faculty}</p>
        <p><strong>Специальность:</strong> ${data.specialty}</p>
        <p><strong>Группа:</strong> ${data.group_name}</p>
        <p><strong>Курс:</strong> ${data.course}</p>
    `;

    // Установка фотографии профиля
    const profilePicElement = document.querySelector('.profile-icon');
    profilePicElement.src = data.photo; // Устанавливаем URL из базы данных
}

async function fetchUserProfile() {
    const token = localStorage.getItem('token');
    const response = await fetch(`${BASE_URL}/users/profile`, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        }
    });

    if (!response.ok) {
        alert("Ошибка при загрузке профиля");
        return;
    }

    const data = await response.json();
    populateProfile(data);
}