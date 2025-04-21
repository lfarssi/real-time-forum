

export async function PostsPage() {
    const response = await fetch("/api/getPosts");
    const data= await response.json()
    console.log(data);
    console.log(data.data==null);
    console.log(!data.data);
    
    if (!data.data ) {
        return /*html*/`
            <div>There's No Post Right Now</div>
        `;
    }
    
  return /*html*/`
        <div class="posts">

        </div>

    `
}

