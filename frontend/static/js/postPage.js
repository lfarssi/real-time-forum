import {errorPage, popupThrottled as popup } from "./errorPage.js";
import { isLogged, navigateTo } from "./app.js";
import { showInputError } from "./authPage.js";
import {  CommentSection } from "./commentSection.js";

let currentParams = getParamsFromLocation();
let page = 2;
let loading = false;
let allPostsLoaded = false;

let throttle = false;

export async function PostsPage(params, page=1) {
  let url;
  if (params == "") {
    url =`/api/getPosts?page=${page}`
  } else if (params == "likedPosts") {
    url =`/api/getLikedPosts?page=${page}`
  } else if (params == "createdPosts") {
    url =`/api/getCreatedPosts?page=${page}`
  } else if (params == "postsByCategory") {
    url = `/api/getPostsByCategory` + location.search + `&page=${page}`
  }
  const response = await fetch(url)
  const data = await response.json()
  if (!data.data && page==1 ) {
     
     return errorPage("No Post Available")
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
                <div class="comts">
                   
                </div>
            </div>
        </div>
    `
    
  })

  return /*html*/`
      ${posts.join('')}

    ` 
}

function getParamsFromLocation() {
  if (location.pathname.startsWith('/postsByCategory')) {
    return 'postsByCategory';
  }
  if (location.pathname.startsWith('/likedPosts')) {
    return 'likedPosts';
  }
  if (location.pathname.startsWith('/createdPosts')) {
    return 'createdPosts';
  }
  return '';
}

async function loadPosts() {
  if (loading || allPostsLoaded) return;

  loading = true;

  try {
    const postsContainer = document.querySelector('.posts');
    const params = getParamsFromLocation();
    const postsHTML = await PostsPage(params, page);

    if (!postsHTML || postsHTML.trim() === "") {
      allPostsLoaded = true;
    } else {
      postsContainer.insertAdjacentHTML('beforeend', postsHTML);
      page += 1;
    }
  } catch (error) {
      popup("No Post Available")
    
  } finally {
    loading = false;
    document.querySelectorAll(".displayComment").forEach(button => {
      button.addEventListener("click", CommentSection);
    });
  }

  
}

window.addEventListener('scroll',async () => {
  if (throttle || loading || allPostsLoaded ) return;
  throttle = true;
  setTimeout(() => {
    const scrollPosition = window.innerHeight + window.scrollY;
    const bottomOfPage = document.body.offsetHeight  ;

    if (scrollPosition >= bottomOfPage) {
      loadPosts();
    }

    throttle = false;
  }, 500);
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
          return popup(`Failed to ${status} post`)
        }
      } catch (err) {
        return popup(err, "failed")
      }
    }

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
      
      <input maxlength="10000" type="text" name="content" required placeholder="content" />
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
          if(data["category"]){
            popup(data["category"])
          }
          if (data.hasOwnProperty(span.id))
            showInputError(data[span.id], span)
        }
        popup(data.message, "failed")
      } else {
        ipt[0].value = ""
        ipt[1].value = ""
        popup("Post Created successfully", 'success');
        navigateTo(location.pathname)
      }
    } catch (err) {
     popup("Something went wrong!", "failed")

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
        page = 1;
        allPostsLoaded = false;
        loading = false;
        currentParams = 'postsByCategory';
        document.querySelector('.posts').innerHTML = ""; // clear old posts
        navigateTo('/postsByCategory?' + checkedInputsValue);
        loadPosts();
      }
    } else {
      if (location.pathname !== "/") {
        page = 1;
        allPostsLoaded = false;
        loading = false;
        currentParams = '';
        document.querySelector('.posts').innerHTML = "";
        navigateTo("/");
        loadPosts();
      }
    }
  });
}
