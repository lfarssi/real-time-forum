import { PostForm, PostsPage } from "./postPage.js"

export function header() {
    return /*html*/`
        <header>
            <nav>
                <a href="/">FORUM</a>
                <ul>
                    <a href="/"><i class="fa fa-home"></i></a>
                    <a href="/createdPosts"><i class="fa-solid fa-pen"></i></a>
                    <a href="/likedPosts" ><i class="fa-solid fa-thumbs-up"></i></a>
                    <button><i class="fa-solid fa-tag"></i></button>
                </ul>
                <a href="/logout">logout</a>
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