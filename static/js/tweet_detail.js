document.addEventListener("DOMContentLoaded", function () {
  const tweetId = window.location.pathname.split("/").pop();
  const userId = replyForm.getAttribute("data-user-id");

  // 返信（リプライ）を投稿する機能
  document.getElementById("postReply").addEventListener("click", function () {
    const replyContent = document.getElementById("replyContent").value;
    const uploadImage = document.getElementById("uploadImage").files[0];

    if (!replyContent) {
      alert("Please enter a reply");
      return;
    }

    const formData = new FormData();
    formData.append("user_id", userId);
    formData.append("content", replyContent);
    formData.append("reply_to", tweetId);

    if (uploadImage) {
      formData.append("image", uploadImage);
    }

    fetch(`/reply/${tweetId}`, {
      method: "POST",
      body: formData, // ここでFormDataを使用
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("Reply post failed");
        }
        return response.json();
      })
      .then(() => {
        alert("リプライ完了");
        document.getElementById("replyContent").value = "";
        if (uploadImage) {
          document.getElementById("selectedImagePreview").src = "";
          document.getElementById("selectedImagePreview").style.display =
            "none";
        }
        loadReplies(); // 返信を再読み込み
      })
      .catch((error) => {
        console.error("Error:", error);
        alert("リプライ失敗");
      });
  });

  // 返信（リプライ）を読み込む機能
  function loadReplies() {
    fetch(`/reply/replies/${tweetId}`)
      .then((response) => response.json())
      .then((data) => {
        const repliesListDiv = document.getElementById("repliesList");
        repliesListDiv.innerHTML = "";
        data.replies.forEach((reply) => {
          const replyElement = document.createElement("div");
          replyElement.className = "border-bottom d-flex flex-row m-3 reply";

          // アバター画像の表示
          const avatarDiv = document.createElement("div");
          avatarDiv.className = "pr-3";

          const userImage =
            reply.user_image.String || "/static/img/default_avatar.png";
          const avatarImg = document.createElement("img");
          avatarImg.src = userImage;
          avatarImg.width = 50;
          avatarImg.height = 50;

          avatarDiv.appendChild(avatarImg);
          replyElement.appendChild(avatarDiv);

          // コンテンツ部分の表示
          const contentDiv = document.createElement("div");

          // ユーザー名とリプライ日時の表示
          const userNameDiv = document.createElement("div");
          userNameDiv.innerText = reply.display_name;
          userNameDiv.style.marginBottom = "5px";
          const createdAtText = document.createElement("small");
          createdAtText.innerText = `　- ${new Date(
            reply.tweet_date
          ).toLocaleString()}`;
          createdAtText.className = "text-muted";
          userNameDiv.appendChild(createdAtText);

          // リプライテキストの表示
          const replyContentDiv = document.createElement("div");
          replyContentDiv.innerText = reply.tweet_content;
          replyContentDiv.style.marginBottom = "15px";

          contentDiv.appendChild(userNameDiv);
          contentDiv.appendChild(replyContentDiv);

          // リプライに画像が含まれている場合の表示
          if (reply.image_path && reply.image_path.Valid) {
            const imageContainer = document.createElement("div");
            imageContainer.className = "justify-content-center";

            const replyImage = document.createElement("img");
            replyImage.src = reply.image_path.String;
            replyImage.alt = "Reply image";
            replyImage.className = "mt-2";
            replyImage.style.width = "100%";
            replyImage.style.height = "100%";
            replyImage.style.objectFit = "cover";

            imageContainer.appendChild(replyImage);
            contentDiv.appendChild(imageContainer);
          }

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
