import { AddPosts, PostForm, PostsPage } from "./postPage.js"

export function header() {
    return /*html*/`
        <header >
            <nav>
                <a id="logo" href="/">FOR<span class="U">U</span>M</a>
                <ul>
                    <li><a class="home" href="/" data-link=""><i class="fa fa-home" ></i></a></li>
                    <li><a class="createdPosts" href="/createdPosts" data-link=""><i class="fa-solid fa-pen"></i></a></li>
                    <li><a class="likedPosts" href="/likedPosts" data-link="" ><i class="fa-solid fa-thumbs-up"></i></a></li>
                    <li><button class="postsByCategory"><i class="fa-solid fa-tag"></i></button></li>
                </ul>
                <a class="logout" href="/logout"><i class="fa-solid fa-right-from-bracket"></i></a>
            </nav>
        </header>
    `
}



export async function homePage(param) {
    return /*html*/`   
        ${header()}
        <main>
            <section>
                ${await PostForm()}
                ${await PostsPage(param)}
            </section>
        </main>

    `
}