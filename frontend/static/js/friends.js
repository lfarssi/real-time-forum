export async function FriendsPage() {
    const response = await fetch("/api/getFriends")
    const data = await response.json()
    let friends = data.data.map(friend => {
        let gender
        if (friend.gender === "male") {
            gender = '<i class="fa-solid fa-user male"></i>'
        } else {
            gender = '<i class="fa-solid fa-user female"></i>'
        }
        return /*html*/`
            <li>${gender} <span>${friend.firstName} ${friend.lastName}</span></li>
    `
    })
    return /*html*/`
        ${friends.join('')}
    `
}