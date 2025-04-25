export async function FriendsPage() {
    const response = await fetch("/api/getFriends")
    const data = await response.json()
    let friends = data.data.map(friend => {
        return /*html*/`
        <div>
            <div>${friend.firstName} ${friend.lastName}</div>
        </div>
    `
    })
    return /*html*/`
        ${friends.join('')}
    `
}