document.addEventListener("DOMContentLoaded", function () {
  document
    .querySelectorAll("a[data-follower-id][data-followee-id]")
    .forEach((element) => {
      element.addEventListener("click", function () {
        const followerId = parseInt(this.getAttribute("data-follower-id"), 10);
        const followeeId = parseInt(this.getAttribute("data-followee-id"), 10);

        // 表示しているメッセージをクリア
        clearMessages();

        const messageInput = document.getElementById("messageInput");
        messageInput.setAttribute("data-sender-id", followerId);
        messageInput.setAttribute("data-receiver-id", followeeId);

        fetch("/message", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            sender_id: followerId,
            receiver_id: followeeId,
          }),
        })
          .then((response) => {
            if (!response.ok) {
              throw new Error("Network response was not ok");
            }
            return response.json();
          })
          .then((data) => {
            data.messages.forEach((msgInfo) => {
              displayMessage(msgInfo, followerId);
            });
          })
          .catch((error) => {
            console.error(
              "There has been a problem with your fetch operation:",
              error
            );
          });
      });
    });
});

// WebSocketサーバーのアドレス
const ws = new WebSocket("ws://" + window.location.host + "/message/ws");

ws.onopen = () => {
  console.log("WebSocket connection established");
};

ws.onerror = (error) => {
  console.error("WebSocket error:", error);
};

ws.onmessage = (event) => {
  try {
    const message = JSON.parse(event.data);
    displayMessage(message, message.sender_id);
  } catch (error) {
    console.error("Error parsing JSON:", error);
  }
};

function clearMessages() {
  const messagesDiv = document.getElementById("messages");
  messagesDiv.innerHTML = ""; // メッセージをクリア
}

function sendMessage() {
  const messageInput = document.getElementById("messageInput");
  const senderId = parseInt(messageInput.getAttribute("data-sender-id"), 10);
  const receiverId = parseInt(
    messageInput.getAttribute("data-receiver-id"),
    10
  );
  const message = {
    sender_id: senderId,
    receiver_id: receiverId,
    content: messageInput.value,
  };
  if (message.content) {
    if (ws.readyState === WebSocket.OPEN) {
      // 送信処理
      ws.send(JSON.stringify(message));
      messageInput.value = ""; // メッセージ送信後に入力フィールドをクリア
    } else {
      console.error("WebSocket is not open. readyState:", ws.readyState);
    }
  }
}

function displayMessage(message, followerId) {
  const messagesDiv = document.getElementById("messages");
  const messageElement = document.createElement("div");
  messageElement.textContent = message.content;

  if (message.sender_id === followerId) {
    messageElement.classList.add("sender_balloon");
  } else {
    messageElement.classList.add("receiver_balloon");
  }

  messagesDiv.appendChild(messageElement);
}
