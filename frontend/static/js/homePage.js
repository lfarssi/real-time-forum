import { PostForm, PostsPage } from "./postPage"

export function header() {
    return /*html*/`
        <header>
            <nav>
                <a href="/">FORUM</a>
                <ul>
                    <a href="/"></a>
                    <a href="/logout"></a>
                </ul>
            </nav>
        </header>
    `
}

export function homePage() {
    return `
        ${header}
        ${PostsPage}
    `
}