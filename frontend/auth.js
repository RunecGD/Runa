const BASE_URL = 'http://localhost:8000'; // Указываем базовый URL с портом 8000

document.getElementById('registerForm')?.addEventListener('submit', async (e) => {
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

document.getElementById('loginForm')?.addEventListener('submit', async (e) => {
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
        localStorage.setItem('token', data.token);
        console.log("Токен установлен в localStorage:", data.token); // Для отладки
        alert(data.message);

        // После успешного входа перенаправляем на предыдущую страницу
        const redirectUrl = localStorage.getItem('redirectUrl') || 'home.html'; // Указываем страницу по умолчанию
        window.location.href = redirectUrl;
    } else {
        alert(data.error);
    }
});
