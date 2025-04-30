import { errorPage } from "./errorPage.js"
import { ws } from "./homePage.js"

let messagesPage = 1

export async function FriendsPage() {
    const response = await fetch("/api/getFriends")
    const data = await response.json()
    
    if(!data.data){
        return errorPage("You Don't Have Friends", 404)
    }
    let friends = data.data.map(friend => {
        let onlineClass = friend.isOnline ? 'online' : 'offline';
        let status =`<i class="fa-solid fa-user ${onlineClass}"></i>`
        return /*html*/`
            <li data-id="${friend.id}">${status} <span>${friend.firstName} ${friend.lastName}</span></li>
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
            messagesPage = 1
            document.querySelector(".chat .messages").innerHTML = ""
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
        page: messagesPage
    }))
}

function loadMessages() {
    const chatMessages = document.querySelector(".chat .messages");

    chatMessages.addEventListener('scroll', () => {
        if (chatMessages.scrollTop === 0) {
            messagesPage++
            console.log(chatMessages.scrollTop)
            let span = document.querySelector('.chat .header span')
            GetMessages(span.dataset.id)
        }
    })
}

export function displayMessage(msg, sender, isSender, isLastMsg = false) {
    const chatMessages = document.querySelector(".chat .messages");

    if (chatMessages) {
        let html = "";
        if (msg.username === sender || isSender) {
            html = /*html*/`
                <div class="messagesSender">
                    <div>
                        <p>${msg.content}  <span class="msgTime">${msg.sentAT.slice(0,5)}</span></p>
                    </div>
                </div>
            `;
        } else {
            html = /*html*/`
                <div class="messagesReceiver">
                    <p>${msg.content} <span class="msgTime">${msg.sentAT.slice(0,5)}</span></p>
                </div>
            `;
        }

        if (isLastMsg) {
            console.log('first')
            chatMessages.innerHTML += html
        } else {
            console.log('second')
            chatMessages.insertAdjacentHTML("afterbegin", html);
        }
        // scroll to top if you want to auto-scroll to latest
        // chatMessages.scrollTop = 0;

        chatMessages.scrollTop = chatMessages.scrollHeight;
    }
}

