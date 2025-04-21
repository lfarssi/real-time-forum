const navigateTo = url => {
    history.pushState(null, null, url)
    router()
}

const router = async () => {
    const routes = [
        {path : "/", view: () => console.log('hello home page')},
        {path : "/posts", view: () => console.log('hello posts page')},
        {path : "/settings", view: () => console.log('hello settings page')}
    ]

    const potentialMatches = routes.map(route => {
        return {
            route: route,
            isMatch: location.pathname === route.path
        }
    })
    
    let match = potentialMatches.find(potentialMatche => potentialMatche.isMatch)
    if (!match) {
        console.log("page not found")
        return
    } 

    // console.log(match.route.view())
}

addEventListener("DOMContentLoaded", () => {
    document.body.addEventListener("click", e => {
        if (e.target.hasAttribute("data-link")) {
            e.preventDefault()
            navigateTo(e.target.href)
        }
    })
    router()
})
