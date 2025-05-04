import { isLogged } from "./app.js";
import { errorPage } from "./errorPage.js"
import {  ws } from "./homePage.js"

let isScroll = false
let scrollValue;
let msgID = -1
let chatMessages;

export async function FriendsPage() {
    const response = await fetch("/api/getFriends")
    const data = await response.json()

    if (!data.data) {
        return errorPage("You Don't Have Friends")
    }
    let friends = data.data.map(friend => {
        let onlineClass = friend.isOnline ? 'online' : 'offline';
        let status = `<i class="fa-solid fa-user ${onlineClass}"></i>`
        return /*html*/`
            <li data-id="${friend.id}" id="friend${friend.id}">${status} <span>${friend.firstName} ${friend.lastName}</span></li>
    `
    })
    return /*html*/`
        ${friends.join('')}
    `
}

export function chatFriend() {
    let friends = document.querySelector('.friends')
    let closeChat = document.querySelector('.chat .closeChat')
    let chat = document.querySelector('.chat')

    friends.addEventListener('click', (e) => {
        let li = e.target.closest("li")
        if (li) {
            chatMessages = document.querySelector(".chat .messages")
            chatMessages.innerHTML = ""
            // chatMessages.removeEventListener('scroll', scrollEventLoadMessages)
            msgID = -1
            isScroll = false
            chat.style.display = 'flex';
            let span = chat.querySelector('.header span')
            span.textContent = li.children[1].textContent
            span.dataset.id = li.dataset.id
            GetMessages(span.dataset.id)
            loadMessages()

        }
    })

    closeChat.addEventListener('click', () => {
        chat.style.display = 'none';
        // isScroll = false
        // msgID = -1
    })
}

export function sendMessage() {
    let form = document.querySelector('.chatForm')

    form.addEventListener('submit', async (e) => {
        e.preventDefault()
        let input = document.querySelector('.chatForm input')
        let receiverID = document.querySelector('.header span').dataset.id
        let logged=await isLogged()
        if (!logged){
            ws.close()
            return
        }
        ws.send(JSON.stringify({
            content: input.value,
            recipientID: parseInt(receiverID),
            type: "addMessage"
        }))
        const ul = document.querySelector(".listFriends")
        ul.innerHTML = `
        ${await FriendsPage()}
    `
        input.value = ""
    })
}

function GetMessages(receiverID) {
    ws.send(JSON.stringify({
        recipientID: parseInt(receiverID),
        type: "loadMessage",
        lastID: parseInt(msgID)
    }))
}

function loadMessages() {
    chatMessages = document.querySelector(".chat .messages");

    chatMessages.addEventListener('scroll', scrollEventLoadMessages)
}

function scrollEventLoadMessages() {
    if (chatMessages.scrollTop === 0 && chatMessages.querySelector('p')) {
        let span = document.querySelector('.chat .header span')
        msgID = chatMessages.querySelector('p').dataset.id
        scrollValue = chatMessages.scrollHeight
        GetMessages(span.dataset.id)
        isScroll = true
    }
}
export const notified = {}

export function notify(sender) {
    const friend = document.getElementById(`friend${sender}`)
    const notification = friend.querySelector(".notification")

    if (!notified[sender]) {
        notified[sender] = 0
    }
    notified[sender] += 1
    if (!notification) {
        let span = document.createElement("span")
        span.className = "notification"
        span.textContent = notified[sender]
        friend.append(span)
    } else {
        notification.textContent = notified[sender]
    }
}
// export function clearNotify(sender) {
//     notified[sender] = 0
//     const friend = document.getElementById(`friend${sender}`)
//     const notification = friend.querySelector(".notification")
//     if (notification) notification.remove()
//   }

export function displayMessage(msg, sender, isSender, isLastMsg = false) {
    const chatMessages = document.querySelector(".chat .messages");

    if (chatMessages) {
        let html = "";
        if (msg.username === sender || isSender) {
            html = /*html*/`
                <div class="messagesSender">
                    <div>
                        <p data-id=${msg.id}>${msg.content} <span class="msgTime">${msg.sentAT.slice(0, 5)}</span></p>
                    </div>
                </div>
            `;
        } else {
            html = /*html*/`
                <div class="messagesReceiver">
                    <p data-id=${msg.id}>${msg.content} <span class="msgTime">${msg.sentAT.slice(0, 5)}</span></p>
                </div>
            `;
        }

        if (isLastMsg) {
            chatMessages.innerHTML += html
        } else {
            chatMessages.insertAdjacentHTML("afterbegin", html);
        }
    

        if (!isScroll) {
            chatMessages.scrollTop = chatMessages.scrollHeight;
        } else {
            chatMessages.scrollTop = chatMessages.scrollHeight - scrollValue
        }

    }
}

