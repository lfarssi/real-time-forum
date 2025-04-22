import { login, loginPage, logout, register, registerPage } from "./authPage.js"
import { errorPage } from "./errorPage.js"
import { AddPosts, CreatedPostsPage, LikedPostsPage, PostForm, PostsByCategoriesPage, PostsPage, ReactPost , ShowComment } from "./postPage.js"
export const navigateTo = url => {
    history.pushState(null, null, url)
    router()
}

const router = async () => {
    const response = await fetch("/api/isLogged")
    if (!response.ok && location.pathname !== "/login" && location.pathname !== "/register") {
        navigateTo("/register")
        return
    } else if (response.ok && (location.pathname === "/login" || location.pathname === "/register")) {
        navigateTo("/")
        return
    }

    const routes = [
        { path: "/", view: PostsPage , eventStart: ReactPost },
        { path: "/likedPosts", view: LikedPostsPage,  eventStart: ReactPost},
        { path: "/createdPosts", view: CreatedPostsPage ,  eventStart: ReactPost},
        { path: "/postsByCategory", view: PostsByCategoriesPage ,  eventStart: ReactPost},
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
