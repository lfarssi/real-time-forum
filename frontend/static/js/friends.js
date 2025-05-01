import { errorPage } from "./errorPage.js"
import { ws } from "./homePage.js"

let isScroll = false
let scrollValue;
let msgID = -1

export async function FriendsPage() {
    const response = await fetch("/api/getFriends")
    const data = await response.json()

    if (!data.data) {
        return errorPage("You Don't Have Friends", 404)
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
            document.querySelector(".chat .messages").innerHTML = ""
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
        isScroll = false
        msgID = -1
    })
}

export function sendMessage() {
    let form = document.querySelector('.chatForm')

    form.addEventListener('submit', async (e) => {
        e.preventDefault()
        let input = document.querySelector('.chatForm input')
        let receiverID = document.querySelector('.header span').dataset.id

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
    const chatMessages = document.querySelector(".chat .messages");

    chatMessages.addEventListener('scroll', () => {
        if (chatMessages.scrollTop === 0) {
            let span = document.querySelector('.chat .header span')
            msgID = chatMessages.querySelector('p').dataset.id
            scrollValue = chatMessages.scrollHeight
            GetMessages(span.dataset.id)
            isScroll = true
        }
    })
}
export const nbrMsg={nbr:0};
export function notify(id){
    console.log(id);
    nbrMsg.nbr++
    const friend= document.getElementById(`friend${id}`)
    console.log(friend);
    
    let span=document.createElement("span")
    span.style.backgroundColor="green"
    span.style.borderRadius="50%"
    span.style.marginLeft = "10px";
    span.style.padding = "0.2em 0.5em";
    span.style.fontSize = "0.8em";

    span.className="notification"
    span.textContent=nbrMsg.nbr
    friend.append(span)
}

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

