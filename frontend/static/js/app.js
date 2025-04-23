import { login, loginPage, logout, register, registerPage } from "./authPage.js"
import { errorPage } from "./errorPage.js"
import { homePage } from "./homePage.js"
import { AddPosts, PostForm, ReactPost } from "./postPage.js"
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
        { path: "/", view: homePage, eventStart: ReactPost },
        { path: "/likedPosts", view: homePage, eventStart: ReactPost },
        { path: "/createdPosts", view: homePage, eventStart: ReactPost },
        { path: "/postsByCategory", view: homePage, eventStart: ReactPost },
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
        document.body.innerHTML = await match.route.view(match.route.path.slice(1,))
    }

    if (match.route.hasOwnProperty("eventStart")) {
        match.route.eventStart()
        AddPosts()
    }
}

addEventListener("DOMContentLoaded", () => {
    document.body.addEventListener("click", e => {
        const link = e.target.closest("a[data-link]");
        if (link) {
            e.preventDefault();
            navigateTo(link.href);
        }
    });

    router();
});
