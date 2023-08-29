document.addEventListener('DOMContentLoaded', function() {

    const tweetForm = document.getElementById('tweetForm');
    const userId = parseInt(tweetForm.getAttribute('data-user-id'), 10);
    const tweetContent = document.getElementById('tweetContent');
    const uploadImage = document.getElementById('uploadImage');
    const imagePreview = document.getElementById('selectedImagePreview');

    // 画像プレビュー
    uploadImage.addEventListener('change', function() {
        const file = uploadImage.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = function(e) {
                imagePreview.src = e.target.result;
                imagePreview.style.display = 'block';
            };
            reader.onerror = function(error) {
                console.error('Error reading the file:', error);
            };
            reader.readAsDataURL(file);
        } else {
            imagePreview.style.display = 'none';
        }
    });
    
    tweetForm.addEventListener('submit', function(e) {
        e.preventDefault();

        let data = {
            userId: userId,
            content: tweetContent.value
        };

        // 画像をBase64としてエンコード
        const file = uploadImage.files[0];
        if (file) {
            const reader = new FileReader();
            reader.readAsDataURL(file);
            reader.onload = function() {
                data.image = reader.result.split(',')[1];

                sendRequest(data);
            };
            reader.onerror = function(error) {
                console.error('Error reading the file:', error);
            };
        } else {
            sendRequest(data);
        }
    });

    function sendRequest(data) {
        fetch('/api/tweet/post', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        })
        .then(response => {
            if (!response.ok) {
                return response.json().then(err => {
                    throw err;
                });
            }
            return response.json();
        })
        .then(data => {
            if (data.user_id) {
                alert('ツイート完了');
                tweetContent.value = '';
                imagePreview.src = '';
                imagePreview.style.display = 'none';
            } else {
                alert('ツイート失敗');
            }
        })
        .catch(errorData => {
            if (errorData && errorData.error) {
                alert(errorData.error);
            } else {
                console.error('Error:', errorData);
                alert('内部エラー');
            }
        });
    }
});
