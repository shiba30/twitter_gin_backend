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
    loadTweets(currentPage);
  });

  document.getElementById("prevPage").addEventListener("click", function () {
    if (currentPage > 1) {
      currentPage -= 1;
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
        updateTweets(data, avatarImagePath);
        document.getElementById("currentPage").innerText = page;
      })
      .catch((error) => {
        console.error("Error fetching tweets:", error);
      });
  }

  function updateTweets(tweets, avatarImage) {
    const tweetContainer = document.querySelector(
      ".col-md-9 > div:nth-child(3)"
    );
    tweetContainer.innerHTML = ""; // 既存のツイートをクリア

    tweets.forEach((tweet) => {
      const tweetElement = document.createElement("div");
      tweetElement.className = "border-bottom d-flex flex-row m-3 tweet";
      tweetElement.setAttribute("data-tweet-id", tweet.tweet_id); // ツイートIDを設定

      const avatarDiv = document.createElement("div");
      avatarDiv.className = "pr-3";

      const userImage =
        tweet.user_image && tweet.user_image.Valid
          ? tweet.user_image.String
          : avatarImage;
      const avatarImg = document.createElement("img");
      avatarImg.src = userImage;
      avatarImg.width = 50;
      avatarImg.height = 50;

      avatarDiv.appendChild(avatarImg);
      tweetElement.appendChild(avatarDiv);

      const contentDiv = document.createElement("div");
      const userNameDiv = document.createElement("div");
      userNameDiv.innerText = tweet.user_name;
      userNameDiv.style.marginBottom = "15px";

      const tweetContentDiv = document.createElement("div");
      tweetContentDiv.innerText = tweet.tweet_content;
      tweetContentDiv.style.marginBottom = "15px";

      contentDiv.appendChild(userNameDiv);
      contentDiv.appendChild(tweetContentDiv);

      if (tweet.image_path && tweet.image_path.Valid) {
        const imageContainer = document.createElement("div");
        imageContainer.className = "justify-content-center";

        const tweetImage = document.createElement("img");
        tweetImage.src = "/" + tweet.image_path.String;
        tweetImage.style.width = "100%";
        tweetImage.style.height = "100%";
        tweetImage.style.objectFit = "cover";
        tweetImage.style.border = "0.2px solid #ddd";

        imageContainer.appendChild(tweetImage);
        contentDiv.appendChild(imageContainer);
      }

      const todoDiv = document.createElement("div");
      todoDiv.className = "mb-3";
      // TODO: Reply, Retweet, Likes area の実装
      contentDiv.appendChild(todoDiv);

      tweetElement.appendChild(contentDiv);

      // ツイート要素にクリックイベントリスナーを追加
      tweetElement.addEventListener("click", function () {
        const tweetId = this.getAttribute("data-tweet-id");
        window.location.href = `/tweets/${tweetId}`; // ツイート詳細ページにリダイレクト
      });
      tweetContainer.appendChild(tweetElement);
    });
  }

  // 初期ロード時に1ページ目のツイートをロード
  loadTweets(1);
});
