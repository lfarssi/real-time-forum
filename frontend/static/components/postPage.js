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
    form.addEventListener("click",async e=>{
        e.preventDefault()
        const formData = Object.fromEntries(new FormData(form).entries());
        try{
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
            document.body.innerHTML = errorPage("Something went wrong!", 500)

        }
       
    })

    
}