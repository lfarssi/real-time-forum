import { PostForm, PostsPage } from "./postPage.js"

export function header() {
    return /*html*/`
        <header>
            <nav>
                <a id="logo" href="/">FOR<span class="U">U</span>M</a>
                <ul>
                    <li><a href="/"><i class="fa fa-home"></i></a></li>
                    <li><a href="/createdPosts"><i class="fa-solid fa-pen"></i></a></li>
                    <li><a href="/likedPosts" ><i class="fa-solid fa-thumbs-up"></i></a></li>
                    <li><button><i class="fa-solid fa-tag"></i></button></li>
                </ul>
                <a class="logout" href="/logout"><i class="fa-solid fa-right-from-bracket"></i></a>
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