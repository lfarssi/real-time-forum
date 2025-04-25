import { isLogged } from "./app.js"
import { PostForm, PostsPage } from "./postPage.js"

export function header() {
    return /*html*/`
        <header>
            <nav>
                <a id="logo" href="/" data-link>FOR<span class="U">U</span>M</a>
                <ul>
                    <li><a class="home" href="/" data-link=""><i class="fa fa-home" ></i></a></li>
                    <li><a class="createdPosts" href="/createdPosts" data-link=""><i class="fa-solid fa-pen"></i></a></li>
                    <li><a class="likedPosts" href="/likedPosts" data-link="" ><i class="fa-solid fa-thumbs-up"></i></a></li>
                    <li><button class="postsByCategory"><i class="fa-solid fa-tag"></i></button></li>
                </ul>
                <a class="logout" href="/logout" data-link><i class="fa-solid fa-right-from-bracket"></i></a>
            </nav>
        </header>
    `
}


export async function homePage(param) {
    let logged = await isLogged()
    if (!logged) {
        return
    }

    let response = await fetch('/api/getCategories')
    let data = await response.json()

    const categoriesInputs = data.data.map(category => /*html*/`
        <input style="display:none;" type="checkbox" name="categories" id="filter${category.id}" value="${category.id}" />
    <label for="filter${category.id}">
    <i class="fa-solid fa-tag"></i> <span> ${category.name}</span>
    </label>
    `).join("");

    let isAsideExists = document.querySelector('aside')

    if (!isAsideExists) {
        document.body.innerHTML = /*html*/`   
        ${header()}
        <main class="container">
            <aside>
                <div class="filter">    
                    <h3>Categories</h3>
                    <form>
                        <ul>
                            ${categoriesInputs}
                        </ul>
                        <button class="btn" type="submit">Filter</button>
                    </form>
                </div>  
            </aside>
            <section>
                ${await PostForm()}
                <div class="posts">
                ${await PostsPage(param)}
                </div>
            </section>
            <aside>
                <div class="profile">
                    <p><i class="fa-solid fa-user"></i> ${logged.username}</p>
                </div>
                <div class="friends">

                </div>
            </aside>
        </main>

    `
    } else {
        let posts = document.querySelector('.posts')
        posts.innerHTML = `${await PostsPage(param)}`
    }
}