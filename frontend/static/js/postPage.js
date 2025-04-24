import { errorPage } from "./errorPage.js";
import { isLogged, navigateTo } from "./app.js";
import { CommentSection } from "./commentSection.js";
import { showInputError } from "./authPage.js";
export async function PostsPage(params) {
  let response;
  if (params == "") {
    response = await fetch("/api/getPosts");
  } else if (params == "likedPosts") {
    response = await fetch("/api/getLikedPosts");

  }
  else if (params == "createdPosts") {
    response = await fetch("/api/getCreatedPosts");

  }
  const data = await response.json()
  if (!data.data) {
    return errorPage("No Post Available", 404)

  }

  console.log(data.data)
  let i; 
  

  let posts = data.data.map(post => {
    return /*html*/`
        <div class="post" id="${post.id}" data-id="${post.id}">
            <div><i class="fa-solid fa-user"></i> ${post.username}</div>
            <p class="dateCreation">${post.dateCreation}</p>
            <div class="title">${post.title}</div>
            <div class="content">${post.content}</div>
            <div class="categories">
              #${post.categories[0].split(",").join(' #')}
            </div>
            <div class="button-group">
                <div>
                  <button class="likePost"  data-id="${post.id}">${post.Likes} <i class="fa-regular fa-thumbs-up"></i></button>
                </div>
                <div>
                <button class="disLikePost"  data-id="${post.id}">${post.Dislikes} <i class="fa-regular fa-thumbs-down"></i></button>
                </div>
                <div>
                  <button class="displayComment"><i class="fa-regular fa-comment"></i></button>
                </div>
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
export function ReactPost() {

  document.querySelectorAll(".likePost, .disLikePost").forEach(button => {
    button.addEventListener("click", async () => {
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
          console.log(`Post ${status}d:`, result.message);

          if (status === "like") {
            // button.textContent = "Liked!";
            // button.disabled = true; 
          } else if (status === "dislike") {
            // button.textContent = "Disliked!";
            // button.disabled = true; 
          }
        } else {
          console.error(`Failed to ${status} post`);
        }
      } catch (err) {
        console.error(`Error ${status}ing post:`, err);
      }
    });
  });
  document.querySelectorAll(".displayComment").forEach(button => {
    button.addEventListener("click", CommentSection);
  });
}


export async function LikedPostsPage() {
  const response = await fetch("/api/getLikedPosts");
  const data = await response.json()
  if (!data.data) {

    return errorPage("No Post Available", 404)

  }
  console.log(data.data);

  let posts = data.data.map(post => {
    return /*html*/`
      <div class="post">
          <div>${post.username}</div>
          <div>${post.title}</div>
          <div>${post.content}</div>
      </div>
  `
  })

  return /*html*/`
      <div class="posts">
          ${posts.join('')}
      </div>

  `
}
export async function CreatedPostsPage() {
  const response = await fetch("/api/getCreatedPosts");
  const data = await response.json()
  if (!data.data) {
    return errorPage("No Post Available", 404)

  }
  console.log(data.data);

  let posts = data.data.map(post => {
    return /*html*/`
      <div class="post">
          <div>${post.username}</div>
          <div>${post.title}</div>
          <div>${post.content}</div>
      </div>
  `
  })

  return /*html*/`
      <div class="posts">
          ${posts.join('')}
      </div>

  `
}
export async function PostsByCategoriesPage() {
  const response = await fetch("/api/getPostsByCategory");
  const data = await response.json()
  if (!data.data) {
    return errorPage("No Post Available", 404)

  }
  console.log(data.data);

  let posts = data.data.map(post => {
    return /*html*/`
      <div class="post">
          <div>${post.username}</div>
          <div>${post.title}</div>
          <div>${post.content}</div>
      </div>
  `
  })

  return /*html*/`
      <div class="posts">
          ${posts.join('')}
      </div>

  `
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
      
      <input maxlength="100" type="text" name="title" placeholder="title" />
      <span class="errPost" id="title"></span>
      
      <input maxlength="1000" type="text" name="content" placeholder="content" />
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
    formData.categories = formDataRaw.getAll("categories"); try {
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