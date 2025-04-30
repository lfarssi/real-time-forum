import { errorPage } from "./errorPage.js"
import { ws } from "./homePage.js"

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
            chat.style.display = 'flex';
            let span = chat.querySelector('.header span')
            span.textContent = li.children[1].textContent
            span.dataset.id = li.dataset.id
            GetMessages(span.dataset.id)
        }
    })

    closeChat.addEventListener('click', () => {
        chat.style.display = 'none';
    })
}



export function sendMessage() {
    let form = document.querySelector('.chatForm')
    
    form.addEventListener('submit', (e) => {
        e.preventDefault()
        let input = document.querySelector('.chatForm input')
        let receiverID = document.querySelector('.header span').dataset.id

        ws.send(JSON.stringify({
            content: input.value,
            recipientID: parseInt(receiverID),
            type: "addMessage"
        }))

        input.value = ""
    })
}

function GetMessages(receiverID) {
    ws.send(JSON.stringify({
        recipientID: parseInt(receiverID),
        type: "loadMessage"
    }))
}


export function displayMessage(msg, sender, isSender) {
    const chatMessages = document.querySelector(".chat .messages");
    console.log(msg);
    
    if (chatMessages) {
        if (msg.username === sender || isSender) {
            chatMessages.innerHTML += /*html*/`
            <div class="messagesSender">
                <div>
                <p>${msg.content}</p>
                <span class="msgTime">${msg.sentAT}</span>
                </div>
            </div>
        `
        } else {
            chatMessages.innerHTML += /*html*/`
            <div class="messagesReceiver">
                
                <p>${msg.content}</p>
                <span>${msg.sentAT}</span>
            </div>
        `
        }
    }
}

