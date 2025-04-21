import { errorPage } from "./errorPage.js";

export async function CommentSection() {
    const response = await fetch("/api/getComments");
    const data= await response.json()
      if (!data.data ) {        
        return errorPage("No Comment Available",404 )
       
    }
    console.log(data.data);
    
    let comments=data.data.map(comment=>{
    return /*html*/`
        <div class="comment">
            <div>${comment.username}</div>
            <div>${comment.content}</div>
        </div>
    `
    })
    
  return /*html*/`
        <div class="comments">
            ${comments.join('')}
        </div>

    `
}
export async function CommentForm() {
  return /*html*/`
    <form id="commentForm">
      

      <input type="text" name="content" placeholder="content" />
      <span class="errComment" id="content"></span>
      
      
      <button type="submit">Comment</button>
    </form>
  `;
}




export function AddComments() {
    const form = document.querySelector("#commentForm");
    const spans = document.querySelector("#errComment");
    form.addEventListener("submit",async e=>{
        e.preventDefault()
        const formData = Object.fromEntries(new FormData(form).entries());
        try{
            const response = await fetch("/api/addComment", {
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