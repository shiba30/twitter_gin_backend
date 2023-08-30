document.addEventListener("DOMContentLoaded", function () {
    const signupForm = document.getElementById('signupForm');

    signupForm.addEventListener('submit', function (e) {
        e.preventDefault();

        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
        const displayName = document.getElementById('displayName').value;

        const xhr = new XMLHttpRequest();
        xhr.open('POST', '/api/user/signup', true);  // エンドポイント設定
        xhr.setRequestHeader('Content-Type', 'application/json');

        xhr.onload = function () {
            if (this.status === 200) {
                window.location.href = '/api/user/verification'; // 登録完了後、verification.html へ遷移
            } else {
                const jsonResponse = JSON.parse(this.responseText);
                alert(jsonResponse.error);
            }
        };

        xhr.onerror = function () {
            alert('入力された情報では登録することができません');
        };

        const data = JSON.stringify({ email: email, password: password, displayName: displayName});

        xhr.send(data);
    });
});
