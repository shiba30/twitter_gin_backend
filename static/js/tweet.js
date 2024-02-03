import { updateTweets } from "./common.js";

document.addEventListener("DOMContentLoaded", function () {
  const tweetForm = document.getElementById("tweetForm");
  const userId = parseInt(tweetForm.getAttribute("data-user-id"), 10);
  const tweetContent = document.getElementById("tweetContent");
  const uploadImage = document.getElementById("uploadImage");
  const imagePreview = document.getElementById("selectedImagePreview");

  // 画像プレビュー
  uploadImage.addEventListener("change", function () {
    const file = uploadImage.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = function (e) {
        imagePreview.src = e.target.result;
        imagePreview.style.display = "block";
      };
      reader.onerror = function (error) {
        console.error("Error reading the file:", error);
      };
      reader.readAsDataURL(file);
    } else {
      imagePreview.style.display = "none";
    }
  });

  tweetForm.addEventListener("submit", function (e) {
    e.preventDefault();

    let data = {
      userId: userId,
      content: tweetContent.value,
    };

    // 画像をBase64としてエンコード
    const file = uploadImage.files[0];
    if (file) {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = function () {
        data.image = reader.result.split(",")[1];

        sendRequest(data);
      };
      reader.onerror = function (error) {
        console.error("Error reading the file:", error);
      };
    } else {
      sendRequest(data);
    }
  });

  // ツイート処理
  function sendRequest(data) {
    fetch("/tweets", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    })
      .then((response) => {
        if (!response.ok) {
          return response.json().then((err) => {
            throw err;
          });
        }
        return response.json();
      })
      .then((data) => {
        if (data.user_id) {
          alert("ツイート完了");
          tweetContent.value = "";
          imagePreview.src = "";
          imagePreview.style.display = "none";
          window.location.reload(true);
        } else {
          alert("ツイート失敗");
        }
      })
      .catch((errorData) => {
        if (errorData && errorData.error) {
          alert(errorData.error);
        } else {
          console.error("Error:", errorData);
          alert("内部エラー");
        }
      });
  }

  // ページネーション
  let currentPage = 1;
  // 初回ページロード時に1ページ目のツイートを取得
  loadTweets(currentPage);

  // ボタンのイベントリスナーを設定
  document.getElementById("nextPage").addEventListener("click", function () {
    currentPage += 1;
    document.getElementById("currentPage").innerText = currentPage;
    loadTweets(currentPage);
  });

  document.getElementById("prevPage").addEventListener("click", function () {
    if (currentPage > 1) {
      currentPage -= 1;
      document.getElementById("currentPage").innerText = currentPage;
      loadTweets(currentPage);
    }
  });

  function loadTweets(page) {
    fetch(`/tweets?page=${page}`)
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return response.json();
      })
      .then((data) => {
        const avatarImagePath = "/static/img/default_avatar.png";
        updateTweets(
          data.tweets,
          data.currentUserId,
          avatarImagePath,
          ".col-md-9 > div:nth-child(3)"
        );
        document.getElementById("currentPage").innerText = page;
      })
      .catch((error) => {
        console.error("Error fetching tweets:", error);
      });
  }

  // 初期ロード時に1ページ目のツイートをロード
  loadTweets(1);
});
