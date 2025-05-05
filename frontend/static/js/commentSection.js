import { isLogged } from "./app.js";
import {   popupThrottled as popup } from "./errorPage.js";
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
      ? data.data.map(comment => {
        const reactLike = comment.IsLiked 
          ? '<i class="fa-solid fa-thumbs-up"></i>'
          : '<i class="fa-regular fa-thumbs-up"></i>';

        const reactDislike = comment.IsDisliked 
          ? '<i class="fa-solid fa-thumbs-down"></i>'
          : '<i class="fa-regular fa-thumbs-down"></i>';

        return /*html*/ `
          <div class="comment" data-id="${comment.id}">
            <div><strong>${comment.username}</strong></div>
            <div>${comment.content}</div>
            <div class="button-group">
              <button class="likeComment" data-id="${comment.id}"><span>${comment.Likes}</span> ${reactLike}</button>
              <button class="disLikeComment" data-id="${comment.id}"><span>${comment.Dislikes}</span> ${reactDislike}</button>
            </div>
          </div>
        `;
      }).join("")
      : "<h4>No Comments available</h4>";

    commentsContainer.innerHTML = `
      ${commentsHtml}
      ${CommentForm(postId)}
    `;

    await AddComments(postId);
    attachCommentReactionListeners(); // Attach event listeners for reactions
  } catch (err) {
    console.error("Error fetching comments:", err);
    popup("Something went wrong!", "failed");
  }
}

export function attachCommentReactionListeners() {
  document.querySelectorAll(".comments").forEach(commentsContainer => {
    commentsContainer.addEventListener("click", async (e) => {
      const button = e.target.closest("button.likeComment, button.disLikeComment");
      if (!button) return;

      if (!await isLogged()) {
        return popup("You need to be logged in to react to comments.", "warning");
      }

      const commentId = parseInt(button.dataset.id);
      const status = button.classList.contains("likeComment") ? "like" : "dislike";

      try {
        const response = await fetch(`/api/addLike`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ commentID: commentId, sender: "comment", status })
        });

        if (response.ok) {
          const result = await response.json();
          console.log(result.data);
          
          const currentCount = parseInt(button.children[0].textContent);
          const newCount = status === "like" ? result.data.cnbLikes : result.data.cnbDislikes;

          button.innerHTML = /*html*/ `
            <span>${newCount}</span>
            <i class="fa-solid ${status === "like" ? "fa-thumbs-up" : "fa-thumbs-down"}"></i>
          `;

          // Update the opposite reaction button
          const oppositeButton = button.parentElement.querySelector(
            status === "like" ? ".disLikeComment" : ".likeComment"
          );
          const oppositeCount = status === "like" ? result.data.nbDislikes : result.data.nbLikes;

          oppositeButton.innerHTML = /*html*/ `
            <span>${oppositeCount}</span>
            <i class="fa-regular ${status === "like" ? "fa-thumbs-down" : "fa-thumbs-up"}"></i>
          `;
        } else {
          popup("Failed to react to comment.", "error");
        }
      } catch (err) {
        console.error("Error reacting to comment:", err);
        popup("Something went wrong!", "failed");
      }
    });
  });
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
      let logged = await isLogged()
    const formData = Object.fromEntries(new FormData(form).entries());
    formData.postID = postId; // Include the post ID in the payload
    const username = logged.username || "You"; // Replace with actual logic to get the username
    
    try {
      const response = await fetch("/api/addComment", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(formData),
      });

      if (response.ok) {
        const result = await response.json();
        popup("Comment added successfully","success")
        // Create a new comment HTML element
        const newCommentHtml = /*html*/`
          <div class="comment">
            <div><strong>${username}</strong></div>
            <div>${formData.content}</div>
            
          </div>
        `;
          form.insertAdjacentHTML("beforebegin", newCommentHtml);     
        form.reset(); // Clear the form
        errorSpan.textContent = ""; // Clear any previous error message
      } else {
        const error = await response.json();
        errorSpan.textContent = error.message || "Failed to add comment.";
      }
    } catch (err) {
      console.error("Error adding comment:", err);
      document.body.innerHTML = popup("Something went wrong!","failed");
    }
  });
}
