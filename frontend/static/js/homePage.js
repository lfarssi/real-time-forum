import { PostForm, PostsPage } from "./postPage.js"

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

export async function homePage() {
    return `
        ${header()}
        ${await PostsPage()}
    `
}