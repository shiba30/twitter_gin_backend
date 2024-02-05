// ツイートにメニュー追加
export function tweetMenu(
  tweetElement,
  tweetId,
  isFollowing,
  tweetUserId,
  currentUserId
) {
  const menuWrapper = document.createElement("div");
  menuWrapper.className = "tweet-menu-wrapper";

  const menuButton = document.createElement("button");
  menuButton.className = "btn btn-sm tweet-menu ";
  menuButton.innerHTML = "⋮"; // 三点リーダーのアイコン
  menuButton.style.color = "white";

  const menuDropdown = document.createElement("div");
  menuDropdown.className = "tweet-menu-dropdown";
  menuDropdown.style.display = "none";

  // メニューウィンドウ開く
  menuButton.addEventListener("click", function (e) {
    // 他のメニューを閉じる
    document
      .querySelectorAll(".tweet-menu-dropdown")
      .forEach(function (dropdown) {
        dropdown.style.display = "none";
      });

    // このメニューを切り替える
    menuDropdown.style.display =
      menuDropdown.style.display === "block" ? "none" : "block";

    e.stopPropagation(); // メニュー開閉イベントの伝播を停止
  });

  // ページのどこかをクリックしたときに全てのメニューを閉じる
  document.addEventListener("click", function () {
    document
      .querySelectorAll(".tweet-menu-dropdown")
      .forEach(function (dropdown) {
        dropdown.style.display = "none";
      });
  });

  // フォローリンクの作成
  const followLink = document.createElement("a");
  followLink.href = "#";
  if (tweetUserId === currentUserId) {
    followLink.innerText = "Delete";
  } else {
    followLink.innerText = isFollowing ? "UnFollow" : "Follow";
  }
  followLink.addEventListener("click", function (e) {
    e.preventDefault();
    e.stopPropagation();
    followAction(tweetId, isFollowing, tweetUserId, e);
    // フォロー操作後にメニューを閉じる
    menuDropdown.style.display = "none";
  });

  // メッセージリンクの作成
  const messageLink = document.createElement("a");
  messageLink.href = "#";
  messageLink.innerText = "Message";
  messageLink.addEventListener("click", function (e) {
    e.preventDefault();
    e.stopPropagation();
    window.location.href = "/message"; // メッセージページへ遷移
  });

  // メニューボタンとドロップダウンをwrapperに追加
  menuDropdown.appendChild(followLink); // フォローリンク
  menuDropdown.appendChild(messageLink); // メッセージリンク
  menuWrapper.appendChild(menuButton);
  menuWrapper.appendChild(menuDropdown);

  // メニューボタンのwrapperをツイート要素に追加
  tweetElement.appendChild(menuWrapper);
}

// フォロー/フォロー解除処理
function followAction(tweetId, isFollowing, tweetUserId, e) {
  e.stopPropagation();
  const path = isFollowing ? "unfollow" : "follow";
  fetch(`/tweets/${tweetId}/${path}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ userId: tweetUserId }),
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error("フォロー操作に失敗しました。");
      }
      return response.json();
    })
    .then(() => {
      alert(isFollowing ? "フォローを解除しました" : "フォローしました");
      window.location.reload();
    })
    .catch((error) => {
      console.error("エラー:", error);
    });
}
