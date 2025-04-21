import { loginPage, register, registerPage } from "../components/authPage.js"
import { errorPage } from "../components/errorPage.js"

export const navigateTo = url => {
    history.pushState(null, null, url)
    router()
}

const router = async () => {
        const response = await fetch("/api/isLogged")
        if (!response.ok && location.pathname !== "/login" && location.pathname !== "/register") {
            navigateTo("/login")
            return
        } else if (response.ok && (location.pathname === "/login" || location.pathname === "/register")) {
            navigateTo("/")
            return
        }

        const routes = [
            { path: "/", view: () => "home page" },
            { path: "/login", view: () => loginPage() },
            { path: "/register", view: () => registerPage(), eventStart: () => register() }
        ]

        const potentialMatches = routes.map(route => {
            return {
                route: route,
                isMatch: location.pathname === route.path
            }
        })

        let match = potentialMatches.find(p => p.isMatch)
        if (!match) {
            document.body.innerHTML = errorPage("Page not found", 404)
            return
        }

        document.body.innerHTML = match.route.view()

        if (match.route.hasOwnProperty("eventStart")) {
            match.route.eventStart()
        }
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
