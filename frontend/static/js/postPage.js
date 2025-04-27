import { errorPage } from "./errorPage.js";
import { isLogged, navigateTo } from "./app.js";
import { CommentSection } from "./commentSection.js";
import { showInputError } from "./authPage.js";

export async function PostsPage(params, offset=0) {
  let response;
  const limit =10;

  let url = `/api/getPosts?offset=${offset}&limit=${limit}`; // Default case
  if (params === "likedPosts") {
      url = "/api/getLikedPosts";
  } else if (params === "createdPosts") {
      url = "/api/getCreatedPosts";
  } else if (params === "postsByCategory") {
      url = '/api/getPostsByCategory' + location.search;
  }

  response = await fetch(url);


  const data = await response.json()
  if (!data.data) {
    return errorPage("No Post Available", 404)
  }

  let posts = data.data.map(post => {
    let reactLike;
    let reactDislike;
    if (post.IsLiked) {
      reactLike = '<i class="fa-solid fa-thumbs-up"></i>'
    } else {
      reactLike = '<i class="fa-regular fa-thumbs-up"></i>'
    }

    if (post.IsDisliked) {
      reactDislike = '<i class="fa-solid fa-thumbs-down"></i>'
    } else {
      reactDislike = '<i class="fa-regular fa-thumbs-down"></i>'
    }

    return /*html*/`
        <div class="post" id="${post.id}" data-id="${post.id}">
            <div><i class="fa-solid fa-user"></i> ${post.username}</div>
            <p class="dateCreation">${post.dateCreation}</p>
            <div class="title">${post.title}</div>
            <div class="content">${post.content}</div>
            <div class="allCategories">
              #${post.categories.join(' #')}
            </div>
            <div class="button-group">
                <button class="likePost" data-id="${post.id}"><span>${post.Likes}</span> ${reactLike}</button>
                <button class="disLikePost" data-id="${post.id}"><span>${post.Dislikes}</span> ${reactDislike}</button>
                <button class="displayComment"><i class="fa-regular fa-comment"></i></button>
            </div>
            
            <div class="comments" style="display:none;">

            </div>
        </div>
    `
  })

  return /*html*/`
      ${posts.join('')}
    `
}

let offset = 0;
const limit = 10;
const params = ""; 
let loading = false;
let allPostsLoaded = false;

async function loadPosts() {
  console.log(offset);
  console.log(limit);
  
  if (loading || allPostsLoaded) return;
  loading = true;

  const postsContainer = document.querySelector('.posts'); // Select here!
  if (!postsContainer) {
      console.error("postsContainer not found");
      return;
  }

  const postsHTML = await PostsPage(params, offset);
  if (!postsHTML || postsHTML.trim() === "") {
      allPostsLoaded = true;
  } else {
      postsContainer.insertAdjacentHTML('beforeend', postsHTML);
      offset += limit;
  }

  loading = false;
}

window.addEventListener('scroll', () => {
  console.log("fii");
  
  if (loading || allPostsLoaded) return;

  const scrollPosition = window.innerHeight + window.scrollY;
  const threshold = document.body.offsetHeight - 500; // 500px before bottom

  if (scrollPosition >= threshold) {
    loadPosts();
  }
});



export function ReactPost() {
  document.querySelector('.posts').addEventListener("click", async (e) => {

    let button = e.target.closest("button[class='likePost'], button[class='disLikePost']")
    if (button) {
      if (!await isLogged()) {
        return
      }

      const postId = parseInt(button.dataset.id);
      const status = button.classList.contains("likePost") ? "like" : "dislike";

      try {
        const response = await fetch(`/api/addLike`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ postID: postId, sender: "post", status })
        });

        if (response.ok) {
          const result = await response.json();

          let currentNbLikes = parseInt(button.children[0].textContent)

          if (status === "like") {
            if (currentNbLikes < result.data.nbLikes) {
              button.innerHTML = /*html*/`
             <span>${result.data.nbLikes} </span>
            <i class="fa-solid fa-thumbs-up"></i>`;
            } else {
              button.innerHTML = /*html*/`
              <span>${result.data.nbLikes} </span>
             <i class="fa-regular fa-thumbs-up"></i>`;
            }

            let disLikeButton = document.querySelector(`button.disLikePost[data-id='${button.dataset.id}']`)
            disLikeButton.innerHTML = /*html*/`
            <span>${result.data.nbDislikes} </span>
           <i class="fa-regular fa-thumbs-down"></i>`;
          } else if (status === "dislike") {
            if (currentNbLikes < result.data.nbDislikes) {
              button.innerHTML = /*html*/`
            <span>${result.data.nbDislikes}</span>
           <i class="fa-solid fa-thumbs-down"></i>`;
            } else {
              button.innerHTML = /*html*/`
                <span>${result.data.nbDislikes}</span>
               <i class="fa-regular fa-thumbs-down"></i>`;
            }

            let likeButton = document.querySelector(`button.likePost[data-id='${button.dataset.id}']`)
            likeButton.innerHTML =  /*html*/`
           <span>${result.data.nbLikes}</span> 
          <i class="fa-regular fa-thumbs-up"></i>`;
          }
        } else {
          console.error(`Failed to ${status} post`);
        }
      } catch (err) {
        console.error(`Error ${status}ing post:`, err);
      }
    }

  });

  document.querySelectorAll(".displayComment").forEach(button => {
    button.addEventListener("click", CommentSection);
  });
}


export async function PostForm() {
  const response = await fetch("/api/getCategories");
  const categories = await response.json();

  const categoriesInputs = categories.data.map(category => /*html*/`
        <input type="checkbox" name="categories" id="${category.id}" value="${category.id}" />
      <label  for="${category.id}">
      ${category.name}
    </label>
  `).join("");

  return /*html*/`
    <form id="postForm">
      <h2>Create post</h2>
      
      <input maxlength="100" type="text" name="title" required placeholder="title" />
      <span class="errPost" id="title"></span>
      
      <input maxlength="1000" type="text" name="content" required placeholder="content" />
      <span class="errPost" id="content"></span>
      
      <h3>Categories</h3>
      <div class="category-section">
        ${categoriesInputs}
      </div>
      
      <button type="submit">Create</button>
    </form>
  `;
}

export function AddPosts() {
  const form = document.querySelector("#postForm");
  const spans = document.querySelectorAll(".errPost");
  const ipt = document.querySelectorAll('#postForm input')
  form.addEventListener("submit", async e => {
    e.preventDefault()

    if (!await isLogged()) {
      return
    }

    spans.forEach(span => {
      span.innerHTML = "";
      span.style.display = "none";
    });

    const formDataRaw = new FormData(form);
    const formData = Object.fromEntries(formDataRaw.entries());

    // ✅ Fix: ensure "categories" is always an array
    formData.categories = formDataRaw.getAll("categories");
    try {
      const response = await fetch("/api/addPost", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(formData)
      });
      const data = await response.json()
      if (!response.ok) {
        for (let span of spans) {
          if (data.hasOwnProperty(span.id))
            showInputError(data[span.id], span)
        }
      } else {
        ipt[0].value = ""
        ipt[1].value = ""
        navigateTo(location.pathname);
      }
    } catch (err) {
      console.log(err);
      document.body.innerHTML = errorPage("Something went wrong!", 500)

    }
  })
}

export function filterByCategories() {
  const filterForm = document.querySelector(".filter form");

  if (!filterForm) return;

  filterForm.addEventListener('change', async () => {
    let checkedInputs = Array.from(filterForm.querySelectorAll('.categories')).filter(cat => cat.checked);

    if (checkedInputs.length !== 0) {
      let checkedInputsValue = checkedInputs.map(cat => "categories=" + cat.value).join('&');
      if (location.pathname + location.search !== '/postsByCategory?' + checkedInputsValue) {
        navigateTo('/postsByCategory?' + checkedInputsValue);
      }
    } else {
      if (location.pathname !== "/") {
        navigateTo("/");
      }
    }
  });
}
