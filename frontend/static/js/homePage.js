import { isLogged } from "./app.js"
import { filterByCategories, PostForm, PostsPage } from "./postPage.js"

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

    let iconsCategories = ['<i class="fa-solid fa-file-code"></i>','<i class="fa-solid fa-lightbulb"></i>',' <i class="fa-solid fa-bitcoin-sign"></i>',
        '<i class="fa-solid fa-child-reaching"></i>', '<i class="fa-solid fa-video"></i>', '<i class="fa-solid fa-medal"></i>',' <i class="fa-solid fa-utensils"></i>'
    ]

    const categoriesInputs = data.data.map((category, index) => /*html*/`
        <input style="display:none;" type="checkbox" name="categories" class="categories" id="filter${category.id}" value="${category.id}" />
    <label for="filter${category.id}">
    ${iconsCategories[index]} <span> ${category.name}</span>
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

    filterByCategories()
}