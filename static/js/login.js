document.addEventListener("DOMContentLoaded", function () {
    const loginForm = document.getElementById('loginForm');

    loginForm.addEventListener('submit', function (e) {
        e.preventDefault();

        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;

        const xhr = new XMLHttpRequest();
        xhr.open('POST', '/api/user/login', true);  // エンドポイント設定
        xhr.setRequestHeader('Content-Type', 'application/json');

        xhr.onload = function () {
            if (this.status === 200) {
                window.location.href = '/api/user/home'; // ログイン後、home.html へ遷移
            } else {
                const jsonResponse = JSON.parse(this.responseText);
                alert(jsonResponse.error);
            }
        };

        xhr.onerror = function () {
            alert('入力された情報ではログインすることができません');
        };

        const data = JSON.stringify({ email: email, password: password });

        xhr.send(data);
    });
});
