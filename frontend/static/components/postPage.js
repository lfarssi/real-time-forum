import { errorPage } from "./errorPage.js";
import { navigateTo } from "../js/app.js";
export async function PostsPage() {
    const response = await fetch("/api/getPosts");
    const data= await response.json()
      if (!data.data ) {        
        return errorPage("No Post Available",404 )
       
    }
    console.log(data.data);
    
    let posts=data.data.map(post=>{
    return /*html*/`
        <div class="post">
            <div>${post.username}</div>
            <div>${post.title}</div>
            <div>${post.content}</div>
            <div>
              <span></span>
              <button class="likePost"  data-id="${post.id}">Like</button>
            </div>
            <div>
              <span></span>
            <button class="disLikePost"  data-id="${post.id}">DisLike</button>
            </div>
            <div>
              <span></span>
            <button class="displayComment"  data-id="${post.id}">Comment</button>
            </div>
        </div>
    `
    })
    
  return /*html*/`
        <div class="posts">
            ${posts.join('')}
        </div>

    `
}
export function ReactPost() {
  document.querySelectorAll(".likePost, .disLikePost").forEach(button => {
    button.addEventListener("click", async () => {
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
            button.textContent = "Liked!";
            // button.disabled = true; 
          } else if (status === "dislike") {
            button.textContent = "Disliked!";
            // button.disabled = true; 
          }
        } else {
          console.log(response);
          
          console.error(`Failed to ${status} post`);
        }
      } catch (err) {
        console.error(`Error ${status}ing post:`, err);
      }
    });
  });
}

export function ShowComment() {
  document.querySelectorAll(".displayComment").forEach(button => {
    button.addEventListener("click", CommentSection);
  });
}

export async function LikedPostsPage() {
  const response = await fetch("/api/getLikedPosts");
  const data= await response.json()
    if (!data.data ) {        
      return errorPage("No Post Available",404 )
     
  }
  console.log(data.data);
  
  let posts=data.data.map(post=>{
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
  const data= await response.json()
    if (!data.data ) {        
      return errorPage("No Post Available",404 )
     
  }
  console.log(data.data);
  
  let posts=data.data.map(post=>{
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
  const data= await response.json()
    if (!data.data ) {        
      return errorPage("No Post Available",404 )
     
  }
  console.log(data.data);
  
  let posts=data.data.map(post=>{
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
  const response = await fetch("/getCategory");
  const categories = await response.json();
  console.log(categories);
  
  const categoriesInputs = categories.data.map(category => `
    <label>
      <input type="checkbox" name="categories" value="${category.id}" />
      ${category.name}
    </label>
  `).join("");

  return /*html*/`
    <form id="postForm">
      <h2>Post</h2>
      
      <input type="text" name="title" placeholder="title" />
      <span class="errPost" id="title"></span>
      
      <input type="text" name="content" placeholder="content" />
      <span class="errPost" id="content"></span>
      
      <div class="category-section">
        <h3>Categories</h3>
        ${categoriesInputs}
      </div>
      
      <button type="submit">Register</button>
    </form>
  `;
}




export function AddPosts() {
    const form = document.querySelector("#postForm");
    const spans = document.querySelectorAll("#errPost");
    form.addEventListener("submit",async e=>{
        e.preventDefault()
        const formDataRaw = new FormData(form);
        const formData = Object.fromEntries(formDataRaw.entries());
    
        // ✅ Fix: ensure "categories" is always an array
        formData.categories = formDataRaw.getAll("categories");        try{
            const response = await fetch("/api/addPost", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(formData)
              });
              if (!response.ok) {
                      for (let span of spans) {
                        if (data.hasOwnProperty(span.id))
                          showRegisterInputError(data[span.id], span)
                      }
                    } else {
                      navigateTo("/");
                    }
        }catch(err){
          console.log(err);
          
            document.body.innerHTML = errorPage("Something went wrong!", 500)

        }
       
    })

    
}