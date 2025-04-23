import { errorPage } from "./errorPage.js";
export async function CommentSection(event) {
  const postElement = event.target.closest(".post");
  const postId = parseInt(postElement.dataset.id);

  const commentsContainer = postElement.querySelector(".comments");

  // Toggle visibility
  const isHidden = commentsContainer.style.display === "none";
  commentsContainer.style.display = isHidden ? "block" : "none";

  try {
    const response = await fetch(`/api/getComments?postID=${postId}`);
    const data = await response.json();

    const commentsHtml = data.data && data.data.length > 0
      ? data.data.map(comment => `
        <div class="comment">
          <div><strong>${comment.username}</strong></div>
          <div>${comment.content}</div>
        </div>
      `).join("")
      : "<h4>No Comments available</h4>";

    // Populate the existing comments container
    commentsContainer.innerHTML = `
      <h4>Comments</h4>
      ${commentsHtml}
      ${CommentForm(postId)}
    `;

    // Add event listener to the comment form
    await AddComments(postId);
  } catch (err) {
    console.error("Error fetching comments:", err);
    errorPage("Something went wrong!", 500);
  }
}
                                                            
export  function CommentForm(postId) {
  return `
    <form id="commentForm-${postId}" class="commentForm">
      <input type="text" name="content" placeholder="Write your comment..." required />
      <span class="errComment" id="errComment-${postId}"></span>
      <button type="submit">Add Comment</button>
    </form>
  `;
}




export async function AddComments(postId) {
  
  const form = document.querySelector(`#commentForm-${postId}`); // Target the specific form for this post
  const errorSpan = document.querySelector(`#errComment-${postId}`); // Target the error span for this form

  form.addEventListener("submit", async (e) => {
    e.preventDefault();

    const formData = Object.fromEntries(new FormData(form).entries());
    formData.postID  = postId; // Include the post ID in the payload

    try {
      const response = await fetch("/api/addComment", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(formData),
      });

      if (response.ok) {
        const result = await response.json();
        console.log("Comment added:", result);

        // Append the new comment dynamically
        const commentsContainer = form.closest(".comments");
        commentsContainer.innerHTML += `
          <div class="comment">
            <div><strong>You</strong></div>
            <div>${formData.content}</div>
          </div>
        `;

        // form.reset(); // Clear the form
      } else {
        const error = await response.json();
        errorSpan.textContent = error.message || "Failed to add comment.";
      }
    } catch (err) {
      console.error("Error adding comment:", err);
      errorSpan.textContent = "Something went wrong!";
    }
  });
}
