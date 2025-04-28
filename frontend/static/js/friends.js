import { errorPage } from "./errorPage.js"

export async function FriendsPage() {
    const response = await fetch("/api/getFriends")
    const data = await response.json()
    
    if(!data.data){
        return errorPage("You Don't Have Friends", 404)
    }
    let friends = data.data.map(friend => {
        let gender
        if (friend.gender === "male") {
            gender = '<i class="fa-solid fa-user online"></i>'
        } else {
            gender = '<i class="fa-solid fa-user offline"></i>'
        }
        return /*html*/`
            <li data-id="${friend.id}">${gender} <span>${friend.firstName} ${friend.lastName}</span></li>
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
            chat.querySelector('.header span').textContent = li.children[1].textContent
            chat.querySelector('.header span').dataset.id = li.dataset.id
        }
    })

    closeChat.addEventListener('click', () => {
        chat.style.display = 'none';
    })
}

const response= await fetch('https://api.ipify.org?format=json')
const data=await response.json()

const ws = new WebSocket(`ws://${data.data.ip}:8080/ws/messages`);

ws.onmessage = function(event) {
    const msg = JSON.parse(event.data);
    console.log(msg)
};

<<<<<<< HEAD
function sendMessage(content, senderID, receiverID) {
    ws.send(JSON.stringify({
        content,
        senderID,
        receiverID,
        "type":"addMessage"
    }));
}
=======
export function sendMessage() {
    let form = document.querySelector('.chatForm')
    
    form.addEventListener('submit', (e) => {
        e.preventDefault()
        let input = document.querySelector('.chatForm input').value
        let receiverID = document.querySelector('.header span').dataset.id
>>>>>>> 674d7d602abe90e70ab949769f6743033b533c22

        ws.send(JSON.stringify({
            content: input,
            recipientID: parseInt(receiverID),
            type: "addMessage"
        }))
    })
}

function displayMessage(msg) {
    const chatMessages = document.querySelector(".chat .messages");
    if (chatMessages) {
        const messageEl = document.createElement("div");
        messageEl.className = "message";
        messageEl.textContent = ` ${msg.content}`;
        chatMessages.appendChild(messageEl);
    }
}

<<<<<<< HEAD
const sendButton = document.querySelector(".sendMessage");
sendButton.addEventListener("click", () => {
    const messageInput = document.querySelector(".messageInput");
    const receiverID = document.querySelector(".chatForm").dataset.receiverId; 
    const content = messageInput.value.trim();

    if (content) {
        const senderID = 1; 
        sendMessage(content, senderID, receiverID);
        messageInput.value = ""; 
    }
});
=======
>>>>>>> 674d7d602abe90e70ab949769f6743033b533c22
