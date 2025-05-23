import { isLogged, navigateTo } from "./app.js"
import { CommentSection } from "./commentSection.js"
import { popupThrottled as popup } from "./errorPage.js";
import { chatFriend, displayMessage, FriendsPage, sendMessage, sortFriendsList, updateUnreadBadges } from "./friends.js"
import { AddPosts, filterByCategories, PostForm, PostsPage, ReactPost } from "./postPage.js"
import { squareMouseHandler } from "./squares.js";


export let ws;

export function header() {
    return /*html*/`
        <header>
            <nav>
                <a id="logo" href="/" data-link>FOR<span class="U">U</span>M</a>
                <ul>
                    <li><a class="home active" href="/" data-link=""><i class="fa fa-home" ></i></a></li>
                    <li><a class="createdPosts" href="/createdPosts" data-link=""><i class="fa-solid fa-pen"></i></a></li>
                    <li><a class="likedPosts" href="/likedPosts" data-link="" ><i class="fa-solid fa-thumbs-up"></i></a></li>
                    <li><button class="postsByCategory"><i class="fa-solid fa-tag"></i></button></li>
                    <li><button class="chatMessages"><i class="fa-solid fa-message"></i></button></li>
                </ul>
                <a class="logout" href="/logout" data-link><i class="fa-solid fa-right-from-bracket"></i></a>
            </nav>
        </header>
    `
}


export async function homePage(param) {
    let logged = await isLogged()
    if (!logged) {
        return
    }


    let response = await fetch('/api/getCategories')
    let data = await response.json()

    let iconsCategories = ['<i class="fa-solid fa-file-code"></i>', '<i class="fa-solid fa-lightbulb"></i>', ' <i class="fa-solid fa-bitcoin-sign"></i>',
        '<i class="fa-solid fa-child-reaching"></i>', '<i class="fa-solid fa-file-video"></i>', '<i class="fa-solid fa-medal"></i>', ' <i class="fa-solid fa-utensils"></i>'
    ]

    const categoriesInputs = data.data.map((category, index) => /*html*/`
        <input style="display:none;" type="checkbox" name="categories" class="categories" id="filter${category.id}" value="${category.id}" />
    <label for="filter${category.id}">
    ${iconsCategories[index]} <span> ${category.name}</span>
    </label>
    `).join("");

    let isAsideExists = document.querySelector('aside')

    if (!isAsideExists) {
        document.body.innerHTML = /*html*/`   
        ${header()}
        <main class="container">
            <aside class="asideFilter">
                <div class="filter">    
                    <h3>Categories</h3>
                    <form>
                        <ul>
                            ${categoriesInputs}
                        </ul>
                    </form>
                    <button class="closeBtn"><i class="fa-solid fa-xmark"></i></button>
                </div>  
            </aside>
            <section>
                <div class="profile">
                <div>
                    <p><i class="fa-solid fa-user"></i> ${logged.firstName} ${logged.lastName}</p>  
                    </div>
                    <span> ${logged.username}</span>
                </div>
                ${await PostForm()}
                <div class="postContainer">
                    <div class="posts">
                    ${await PostsPage(param)}
                    </div>
                </div>
                
            </section>
            <aside class="asideFriendsProfile">
                <div class="profile">
                    <div>
                    <p><i class="fa-solid fa-user"></i> ${logged.firstName} ${logged.lastName}</p>  
                    </div>
                    <span> ${logged.username}</span>
                </div>
                <div class="friends">
                <ul class="listFriends">
                ${await FriendsPage()}
                </ul>
                <div class="chat">
                    <div class="header">
                        <p><i class="fa-solid fa-user"></i> <span></span></p>
                        <button class="closeChat"><i class="fa-solid fa-xmark"></i></button>
                    </div>
                    <div class="cbody">
                        <div class="messages">
        
                        </div>
                        <form id="chatForm" class="chatForm">
                            <div class="input-container">
                                <input type="text" class="messageInput" name="content" placeholder="Write your message..." required />
                                <button type="submit" class="sendMessage"><i class="fa-solid fa-paper-plane"></i></button>
                            </div>
                            <span class="errChat" id="errChat"></span>
                        </form> 
                    </div>
                </div>
                </div>
            </aside>
        </main>

    `

        document.body.removeEventListener('mousemove', squareMouseHandler);
        AddPosts()
        ReactPost()
        filterByCategories()
        chatFriend()
        sendMessage()
        asideNav()

        ws = new WebSocket(`/ws/messages`);
        ws.onerror = function (event) {
             popup(event.message, "failed");
            
        };

        ws.onclose = function (event) {
            popup("Connection closed unexpectedly.", "warning")
        };

        ws.onmessage = async function (event) {
            const logged = await isLogged()
            if (!logged) {
                ws.close()
                return
            }

            let user = document.querySelector('.chat .header span');
            let openChatUserId = user ? parseInt(user.dataset.id) : null;

            const msg = JSON.parse(event.data);

            if (msg.type === "userStatus") {

                let friendLi = document.querySelector(`.listFriends li[data-id="${msg.userID}"]`);

                if (!friendLi) {
                    if (!msg.userName) return; // Defensive: ignore if info missing

                    friendLi = document.createElement('li');
                    friendLi.setAttribute('data-id', msg.userID);
                    friendLi.innerHTML = /*html*/`
                            <i class="fas fa-user ${msg.isOnline ? 'online' : 'offline'}"></i>
                            <span>${msg.userName}</span>    
                            <div class="loader loader-${msg.userID}"></div>
                    `;
                    document.querySelector('.listFriends').appendChild(friendLi);
                } else {
                    // Already exists: update status icon
                    const icon = friendLi.querySelector("i.fa-user");
                    if (icon) {
                        icon.classList.toggle("online", msg.isOnline);
                        icon.classList.toggle("offline", !msg.isOnline);
                    }
                }


                let isFriendListNotEmpty = document.querySelector('.listFriends li')
                if (isFriendListNotEmpty) {
                    let noFriends = document.querySelector('.listFriends .noPost')
                    if (noFriends) {
                        noFriends.remove()
                    }
                }

                sortFriendsList()
            }

            if (msg.type === "unreadCounts") {
                sortFriendsList()
            }

            if (msg.type === "allMessages") {

                if (msg.data) {
                    msg.data.map(m => displayMessage(m, logged.username, m.username));
                }
                updateUnreadBadges(msg.counts, openChatUserId);

            } else if (msg.type === "newMessage") {
                const senderId = msg.data.senderID;
                const recipientId = msg.data.recipientID;

                const ul = document.querySelector(".listFriends")
                const friend = document.querySelector(`.listFriends li[data-id="${recipientId != logged.id ? recipientId : senderId}"]`);

                ul.prepend(friend)
                if (!friend.classList.contains("has-messages")) {
                    friend.classList.add('has-messages')
                }
                if (openChatUserId !== senderId) {
                    updateUnreadBadges(msg.counts, openChatUserId);
                }
                if (user.dataset.id == recipientId || user.dataset.id == senderId) {
                    let receiverChat = document.querySelector('.chat .header p span')
                    displayMessage(msg.data, logged.username, receiverChat.textContent, msg.isSender, true);
                    if (user.dataset.id == senderId) {
                        ws.send(
                            JSON.stringify({
                                status: openChatUserId == senderId ? "read" : "unread",
                                senderID: parseInt(senderId),
                                recipientID: parseInt(recipientId),
                                type: "updateMessage",
                            })
                        );
                    }
                }

            } else if (msg.type == "errMessage") {
                // console.log("typing...");
                popup(msg.message, "failed")

            } else if (msg.type == "isTyping") {
                let receiverChat = document.querySelector('.chat .header p')
                if (receiverChat.children[1].dataset.id == msg.senderID) {
                    let loaderElement = receiverChat.querySelector('.loader')
                    if (!loaderElement) {
                        receiverChat.innerHTML += /*html*/`
                        <div class="loader"></div>
                    `
                    }
                } else {
                    let sender = document.querySelector(`.listFriends li[data-id="${msg.senderID}"]`)
                    let loaderElement = sender.querySelector('.loader')
                    if (!loaderElement) {
                        sender.innerHTML += /*html*/`
                        <div class="loader"></div>
                        `
                    }
                }

            } else if (msg.type == "pauseTyping") {
                let chattypingElement = document.querySelector('.chat .header p .loader');
                if (chattypingElement) {
                    chattypingElement.remove()
                }

                let listFriendstypingElement = document.querySelector(`.listFriends li[data-id="${msg.senderID}"] .loader`)
                if (listFriendstypingElement) {
                    listFriendstypingElement.remove()
                }

            }

            if (!openChatUserId) {
                updateUnreadBadges(msg.counts, openChatUserId);
            }


        };




    } else {
        let posts = document.querySelector('.posts')
        posts.innerHTML = `${await PostsPage(param)}`
    }

    const urlParams = new URLSearchParams(location.search);
    const myParam = urlParams.getAll('categories');
    let categories = document.querySelectorAll('.categories')

    categories.forEach(category => {
        category.checked = false
        if (myParam.includes(category.value)) {
            category.checked = true
        }
    })

    activePage()

    document.querySelectorAll(".displayComment").forEach(button => {
        button.addEventListener("click", CommentSection);
    });
}

function activePage() {
    let a = document.querySelectorAll('header nav ul a')

    a.forEach(element => {
        element.style.color = "var(--text-light)"
        if (element.pathname === location.pathname) {
            element.style.color = "var(--accent)"
        }
    })
}

function asideNav() {
    let btn = document.querySelector('header .postsByCategory')
    let category = document.querySelector('aside .filter')
    let closeBtn = document.querySelector('.closeBtn')
    let chatMessages = document.querySelector('.chatMessages')
    let friends = document.querySelector('aside .friends')

    btn.addEventListener('click', () => {
        friends.classList.remove('showFriends')
        category.classList.toggle('showFilter')
    })

    closeBtn.addEventListener('click', () => {
        category.classList.remove('showFilter')
    })

    chatMessages.addEventListener('click', () => {
        category.classList.remove('showFilter')
        friends.classList.toggle('showFriends')
    })
}
