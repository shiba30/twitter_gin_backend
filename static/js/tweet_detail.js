document.addEventListener("DOMContentLoaded", function () {
  const replyForm = document.getElementById("replyForm");
  const tweetId = parseInt(window.location.pathname.split("/").pop(), 10);
  const userId = parseInt(replyForm.getAttribute("data-user-id"), 10);
  const replyContent = document.getElementById("replyContent");
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

  // 返信（リプライ）を投稿する機能
  replyForm.addEventListener("submit", function (e) {
    e.preventDefault();

    let formData = {
      replyTo: tweetId,
      userId: userId,
      content: replyContent.value,
    };

    // 画像をBase64としてエンコード
    const file = uploadImage.files[0];
    if (file) {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = function () {
        formData.image = reader.result.split(",")[1];

        sendReplyRequest(formData);
      };
      reader.onerror = function (error) {
        console.error("Error reading the file:", error);
      };
    } else {
      sendReplyRequest(formData);
    }
  });

  function sendReplyRequest(formData) {
    fetch(`/reply/${tweetId}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(formData),
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
        if (data.message === "Comment posted successfully") {
          alert("リプライ完了");
          replyContent.value = "";
          imagePreview.src = "";
          imagePreview.style.display = "none";
          loadReplies(); // 返信を再読み込み
        } else {
          alert("リプライ失敗");
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

  // 返信（リプライ）を読み込む機能
  function loadReplies() {
    fetch(`/reply/replies/${tweetId}`)
      .then((response) => response.json())
      .then((data) => {
        const avatarImage = "/static/img/default_avatar.png";
        const repliesListDiv = document.getElementById("repliesList");
        repliesListDiv.innerHTML = "";

        data.replies.forEach((reply) => {
          const replyElement = document.createElement("div");
          replyElement.className = "border-bottom d-flex flex-row m-3 tweet";
          replyElement.setAttribute("data-tweet-id", reply.tweet_id); // ツイートIDを設定

          const avatarDiv = document.createElement("div");
          avatarDiv.className = "pr-3";

          const userImage =
            reply.user_image && reply.user_image.Valid
              ? reply.user_image.String
              : avatarImage;
          const avatarImg = document.createElement("img");
          avatarImg.src = userImage;
          avatarImg.width = 50;
          avatarImg.height = 50;

          avatarDiv.appendChild(avatarImg);
          replyElement.appendChild(avatarDiv);

          const contentDiv = document.createElement("div");
          const userNameDiv = document.createElement("div");
          userNameDiv.innerText = reply.display_name;
          userNameDiv.style.marginBottom = "15px";

          const createdAtText = document.createElement("small");
          createdAtText.innerText = `　- ${new Date(
            reply.tweet_date
          ).toLocaleString()}`;
          createdAtText.className = "text-muted";
          userNameDiv.appendChild(createdAtText);

          const replyContentDiv = document.createElement("div");
          replyContentDiv.innerText = reply.tweet_content;
          replyContentDiv.style.marginBottom = "15px";

          contentDiv.appendChild(userNameDiv);
          contentDiv.appendChild(replyContentDiv);

          if (reply.image_path && reply.image_path.Valid) {
            const imageContainer = document.createElement("div");
            imageContainer.className = "justify-content-center";

            const replyImage = document.createElement("img");
            replyImage.src = "/" + reply.image_path.String;
            replyImage.style.width = "100%";
            replyImage.style.height = "100%";
            replyImage.style.objectFit = "cover";
            replyImage.style.border = "0.2px solid #ddd";

            imageContainer.appendChild(replyImage);
            contentDiv.appendChild(imageContainer);
          }

          const todoDiv = document.createElement("div");
          todoDiv.className = "mb-3";
          contentDiv.appendChild(todoDiv);

          replyElement.appendChild(contentDiv);

          repliesListDiv.appendChild(replyElement);
        });
      })
      .catch((error) => {
        console.error("Error:", error);
      });
  }

  loadReplies(); // ページロード時に返信を読み込む
});
