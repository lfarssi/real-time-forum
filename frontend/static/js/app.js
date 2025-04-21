import { login, loginPage, logout, register, registerPage } from "../components/authPage.js"
import { errorPage } from "../components/errorPage.js"
import { AddPosts, CreatedPostsPage, LikedPostsPage, PostForm, PostsByCategoriesPage, PostsPage } from "../components/postPage.js"
export const navigateTo = url => {
    history.pushState(null, null, url)
    router()
}

const router = async () => {
    const response = await fetch("/api/isLogged")
    if (!response.ok && location.pathname !== "/login" && location.pathname !== "/register" && location.pathname !== "/logout") {
        navigateTo("/login")
        return
    } else if (response.ok && (location.pathname === "/login" || location.pathname === "/register" || location.pathname === "/logout")) {
        navigateTo("/")
        return
    }

    const routes = [
        { path: "/", view: PostsPage },
        { path: "/likedPosts", view: LikedPostsPage },
        { path: "/createdPosts", view: CreatedPostsPage },
        { path: "/postsByCategory", view: PostsByCategoriesPage },
        { path: "/createPost", view: PostForm, eventStart: AddPosts },
        { path: "/login", view: loginPage, eventStart: login },
        { path: "/register", view: registerPage, eventStart: register },
        { path: "/logout", eventStart: logout }
    ];

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

    if (match.route.hasOwnProperty("view")) {
        document.body.innerHTML = await match.route.view()
    }

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
