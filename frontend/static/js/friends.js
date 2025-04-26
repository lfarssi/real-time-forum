export async function FriendsPage() {
    const response = await fetch("/api/getFriends")
    const data = await response.json()
    let friends = data.data.map(friend => {
        let gender
        if (friend.gender === "male") {
            gender = '<i class="fa-solid fa-user online"></i>'
        } else {
            gender = '<i class="fa-solid fa-user offline"></i>'
        }
        return /*html*/`
            <li>${gender} <span>${friend.firstName} ${friend.lastName}</span></li>
    `
    })
    return /*html*/`
        ${friends.join('')}
    `
}

export function chatFriend() {
    let ul = document.querySelector('.friends ul')
    ul.addEventListener('click', (e) => {
        let friend = e.target.closest("li")
        if (friend) {
            console.log(friend)
        }
    })
}