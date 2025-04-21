import { errorPage } from "./errorPage.js";

export async function PostsPage() {
    const response = await fetch("/api/getPosts");
    const data= await response.json()
      if (!data.data ) {
        document.body.innerHTML= errorPage("No Post Available",404 )
        return 
       
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

export function AddPosts() {
    const form = document.querySelector("#registerForm");

    const response = await fetch("/api/addPost", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(fromData)
      });
}