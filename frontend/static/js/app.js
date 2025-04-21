import { loginPage, register, registerPage } from "../components/authPage.js"
import { errorPage } from "../components/errorPage.js"

export const navigateTo = url => {
    history.pushState(null, null, url)
    router()
}

const router = async () => {
    if (location.pathname !== "/login" && location.pathname !== "/register") {
        let response = await fetch("/api/isLogged")
        let data = await response.json()
        console.log(data)
        if (!response.ok) {
            navigateTo("/login")
            return 
        }
    }

    const routes = [
        {path : "/", view: () => "home page"},
        {path : "/login", view: () => loginPage()},
        {path : "/register", view: () => registerPage()}
    ]

    const potentialMatches = routes.map(route => {
        return {
            route: route,
            isMatch: location.pathname === route.path
        }
    })
    
    let match = potentialMatches.find(potentialMatche => potentialMatche.isMatch)
    if (!match) {
        document.body.innerHTML = errorPage("Page Not Found", 404)
        return
    } 

    document.body.innerHTML = match.route.view()

    register()
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
