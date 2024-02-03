import { updateTweets } from "./common.js";

document.addEventListener("DOMContentLoaded", function () {
  loadBookmarkedTweets(); // ページロード時にブックマークをロード
});

// ブックマークされたツイートをロードして表示する関数
function loadBookmarkedTweets() {
  fetch("/bookmarks", {
    headers: {
      Accept: "application/json",
    },
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      return response.json();
    })
    .then((data) => {
      const bookmarkContainer = document.getElementById("bookmarkContainer");
      bookmarkContainer.innerHTML = ""; // コンテナをクリア
      const avatarImagePath = "/static/img/default_avatar.png";
      updateTweets(
        data.tweets,
        data.userId,
        avatarImagePath,
        "#bookmarkContainer"
      ); // ツイートを表示
    })
    .catch((error) => {
      console.error("Error fetching tweets:", error);
    });
}
