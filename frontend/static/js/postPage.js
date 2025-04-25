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
  } else if (params == "postsByCategory") {
    response = await fetch('/api/getPostsByCategory' + location.search)
  }

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
                <button class="likePost"  data-id="${post.id}">${post.Likes} ${reactLike}</button>
                <button class="disLikePost"  data-id="${post.id}">${post.Dislikes} ${reactDislike}</button>
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
export function ReactPost() {
  document.body.addEventListener("click", async (e) => {

    let button = e.target.closest("button[class='likePost'], button[class='disLikePost']")
    if (button) {
      if (!await isLogged()) {
        return
      }

      console.log(button)
  
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
  
          if (status === "like") {
            button.innerHTML = /*html*/`
             ${result.data.nbLikes} 
            <i class="fa-solid fa-thumbs-up"></i>`;
            document.querySelector("button[class='disLikePost']").innerHTML = /*html*/`
            ${result.data.nbDislikes} 
           <i class="fa-regular fa-thumbs-down"></i>`;
          } else if (status === "dislike") {
            button.innerHTML = /*html*/`
            ${result.data.nbDislikes} 
           <i class="fa-solid fa-thumbs-down"></i>`;
           document.querySelector("button[class='likePost']").innerHTML =  /*html*/`
           ${result.data.nbLikes} 
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
    console.log(formData)
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
  }, { once: true })
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
