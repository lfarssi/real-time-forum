import { isLogged } from "./app.js";
import { errorPage, popup } from "./errorPage.js";
import { ws } from "./homePage.js";

let isScroll = false;
let scrollValue;
let msgID = -1;
let chatMessages;
let timeTyping;

export async function FriendsPage() {
  const response = await fetch("/api/getFriends");
  const data = await response.json();

  if (!data.data) {
    return errorPage("You Don't Have Friends");
  }
  let friends = data.data.map((friend) => {

    let onlineClass = friend.isOnline ? "online" : "offline";
    let msgClass = friend.lastAt.Valid ? "has-messages" : "";
    return /*html*/ `
        <li data-id="${friend.id}" id="friend${friend.id}"  class="${msgClass}">
          <i class="fa-solid fa-user ${onlineClass}"></i> <span>${friend.username}</span>
        </li>
`;

  });
  return /*html*/ `
        ${friends.join("")}
    `;
}
export function chatFriend() {
  let friends = document.querySelector(".friends");
  let closeChat = document.querySelector(".chat .closeChat");
  let chat = document.querySelector(".chat");

  friends.addEventListener("click", (e) => {
    let li = e.target.closest("li");
    if (li) {
      chatMessages = document.querySelector(".chat .messages");
      let input = document.querySelector(".chatForm input");
      chatMessages.innerHTML = "";
      msgID = -1;
      isScroll = false;
      chat.style.display = "flex";
      let span = chat.querySelector(".header span");
      ws.send(
        JSON.stringify({
          recipientID: parseInt(span.dataset.id),
          type: "pauseTyping",
        })
      );
      clearTimeout(timeTyping)
      timeTyping = undefined
      span.textContent = li.children[1].textContent;
      span.dataset.id = li.dataset.id;
      input.removeEventListener('input', debounceTyping)
      Typing()
      GetMessages(span.dataset.id);
      loadMessages();
    }
  });
  closeChat.addEventListener("click", () => {
    chat.style.display = "none";
    let span = chat.querySelector(".header span");
    if (span && span.dataset.id) {
      ws.send(
        JSON.stringify({
          recipientID: parseInt(span.dataset.id),
          type: "pauseTyping",
        })
      );
      clearTimeout(timeTyping)
      timeTyping = undefined
      let typingElement = chat.querySelector('.header p .loader')
      if (typingElement) {
        let sender = document.querySelector(`.listFriends li[data-id="${span.dataset.id}"]`)
        let loaderElement = sender.querySelector('.loader')
        if (!loaderElement) {
          sender.innerHTML += /*html*/`
          <div class="loader"></div>
        `
        }
      }
      span.removeAttribute("data-id");
    }
    let input = document.querySelector(".chatForm input");
    input.removeEventListener('input', debounceTyping)

  });
}

let debounceTyping = leadingDebounceTyping(onTyping, 7000)
export function Typing() {
  let input = document.querySelector(".chatForm input");

  input.addEventListener("input", debounceTyping);
}

async function onTyping() {
  let receiverID = document.querySelector(".header span").dataset.id;
  let logged = await isLogged();
  if (!logged) {
    return;
  }

  ws.send(
    JSON.stringify({
      recipientID: parseInt(receiverID),
      senderID: logged.id,
      type: "Typing",
    })
  );
}

function leadingDebounceTyping(func, timeout) {
  timeTyping;
  return (...args) => {
    if (!timeTyping) {
      func(...args)
    }
    clearTimeout(timeTyping);
    timeTyping = setTimeout(() => {
      let receiverID = document.querySelector(".header span").dataset.id;
      ws.send(
        JSON.stringify({
          recipientID: parseInt(receiverID),
          type: "pauseTyping",
        })
      );
    }, timeout);
  };
}

function stopTyping() {
  let receiverID = document.querySelector('.chat .header span');
  clearTimeout(timeTyping);
  timeTyping = undefined;
  if (span && span.dataset.id) {
    ws.send(
      JSON.stringify({
        recipientID: parseInt(receiverID),
        type: "pauseTyping",
      })
    );

  }
}

export function sendMessage() {
  let form = document.querySelector(".chatForm");
  form.addEventListener("submit", async (e) => {
    e.preventDefault();
    let input = document.querySelector(".chatForm input");
    let receiverID = document.querySelector(".header span").dataset.id;
    let logged = await isLogged();
    if (!logged) {
      return;
    }
    if (input.value.trim() == "") {
      popup("Cannot Send Empty Message", "failed")
      return
    }
    ws.send(
      JSON.stringify({
        content: input.value,
        recipientID: parseInt(receiverID),
        senderID: logged.id,
        type: "addMessage",
      })
    );
    ws.send(
      JSON.stringify({
        recipientID: parseInt(receiverID),
        type: "pauseTyping",
      })
    );
    clearTimeout(timeTyping)
    timeTyping = undefined
    input.value = "";
  });
}

function GetMessages(receiverID) {
  ws.send(
    JSON.stringify({
      recipientID: parseInt(receiverID),
      type: "loadMessage",
      lastID: parseInt(msgID),
    })
  );
}

let debounceScrollEvent = scrollChatDebounce(scrollEventLoadMessages, 500)
function loadMessages() {
  chatMessages = document.querySelector(".chat .messages");

  chatMessages.addEventListener("scroll", debounceScrollEvent);
}

function scrollEventLoadMessages() {
  if (chatMessages.scrollTop === 0 && chatMessages.querySelector('p')) {
    let span = document.querySelector('.chat .header span')
    msgID = chatMessages.querySelector('p').dataset.id
    scrollValue = chatMessages.scrollHeight
    GetMessages(span.dataset.id)
  }
  isScroll = true

  if (Math.ceil(chatMessages.scrollTop + chatMessages.clientHeight) >= chatMessages.scrollHeight - 50) {
    isScroll = false
  }

  // const observer = new IntersectionObserver()

}

function scrollChatDebounce(func, timeout = 300) {
  let timer;
  return (...args) => {
    clearTimeout(timer);
    timer = setTimeout(() => func(...args), timeout);
  };
}


export function updateUnreadBadges(counts, openedUserId = null) {
  const openedUserIdNum = openedUserId !== null ? Number(openedUserId) : null;
  document.querySelectorAll(".listFriends li").forEach((li) => {
    const friendID = Number(li.dataset.id);
    const badge = li.querySelector(".notification");
    if (openedUserIdNum !== null && friendID === openedUserIdNum) {
      if (badge) badge.remove();
      return;
    }

    const count = counts && counts[friendID] ? counts[friendID] : 0;

    if (count > 0) {
      if (!badge) {
        const span = document.createElement("span");
        span.className = "notification";
        span.textContent = count;
        li.append(span);
      } else {
        badge.textContent = count;
      }
    }
  });
}


export function displayMessage(msg, sender, receiver, isSender, isLastMsg = false) {

  const chatMessages = document.querySelector(".chat .messages");

  if (chatMessages) {
    const newdate = new Date(msg.sentAT)

    let html = "";
    if (msg.username === sender || isSender) {
      html = /*html*/ `
                <div class="messagesSender">
                    <div>
                    <span style="color:green;">${sender}</span>
                        <p data-id=${msg.id}>${msg.content
        } <span class="msgTime">${newdate.toLocaleString("en-GB")}</span></p>
                    </div>
                </div>
            `;
    } else {
      html = /*html*/ `
                <div class="messagesReceiver">
                <span style="color:blue;">${receiver}</span>
                    <p data-id=${msg.id}>${msg.content
        } <span class="msgTime">${newdate.toLocaleString("en-GB")}</span></p>
                </div>
            `;
    }


    if (isLastMsg) {
      chatMessages.innerHTML += html
    } else {
      chatMessages.insertAdjacentHTML("afterbegin", html);
    }

    if (!isScroll || isSender) {
      chatMessages.scrollTop = chatMessages.scrollHeight;
    } else if (isScroll && !isLastMsg) {
      chatMessages.scrollTop = chatMessages.scrollHeight - scrollValue

    }

  }
}
export function sortFriendsList() {
  const list = document.querySelector('.listFriends');

  if (!list) return;

  const allFriends = Array.from(list.children);

  // Split into "messaged" and "never messaged"
  const messaged = [];
  const nonMessaged = [];
  allFriends.forEach(li => {
    if (li.classList.contains('has-messages')) {
      messaged.push(li); // Keep their order
    } else {
      nonMessaged.push(li);
    }
  });

  // Sort non-messaged friends by first name (case-insensitive)
  nonMessaged.sort((a, b) => {
    const aName = a.querySelector('span').textContent.trim().toLowerCase();
    const bName = b.querySelector('span').textContent.trim().toLowerCase();
    const aFirstName = aName.split(' ')[0];
    const bFirstName = bName.split(' ')[0];
    return aFirstName.localeCompare(bFirstName);
  });

  // Clear and re-append in order
  list.innerHTML = '';
  messaged.concat(nonMessaged).forEach(li => list.appendChild(li));
}
