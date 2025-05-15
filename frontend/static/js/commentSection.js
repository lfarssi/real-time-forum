// comments.js
import { isLogged } from "./app.js";
import { popupThrottled as popup } from "./errorPage.js";

/**
 * Returns the markup for the comment form under a post.
 */
export function CommentForm(postId) {
  return `
    <form id="commentForm-${postId}" class="commentForm">
      <input type="text" name="content" placeholder="Write your comment..." required />
      <span class="errComment" id="errComment-${postId}"></span>
      <button type="submit">Add Comment</button>
    </form>
  `;
}

/**
 * Attaches submit handler to the Add Comment form under a post.
 * On success: refreshes that post's comments.
 */
export function AddComments(postId) {
  const form = document.querySelector(`#commentForm-${postId}`);
  const errorSpan = document.querySelector(`#errComment-${postId}`);
  if (!form) return;

  // Replace node to clear old listeners
  const newForm = form.cloneNode(true);
  form.parentNode.replaceChild(newForm, form);

  newForm.addEventListener("submit", async e => {
    e.preventDefault();

    if (!await isLogged()) {
      return popup("You must be logged in.", "warning");
    }

    const payload = Object.fromEntries(new FormData(newForm).entries());
    payload.postID = postId;

    try {
      const res = await fetch("/api/addComment", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      if (!res.ok) {
        const err = await res.json();
        popup(err.message, "failed")
        return;
      }

      popup("Comment added successfully", "success");
      // Refresh comments
      const postEl = newForm.closest(".post");
      await fetchAndRenderComments(postEl);
    }
    catch (err) {
      popup("You Can't Comment Right Now!! ", "failed");
    }
  });
}

/**
 * Fetches comments for a given post element, then renders them plus the form.
 */
export async function fetchAndRenderComments(postEl) {
  const postId = parseInt(postEl.dataset.id, 10);
  const commentsContainer = postEl.querySelector(".comments");
  if (!commentsContainer) return;

  try {
    const res = await fetch(`/api/getComments?postID=${postId}`);
    if (!res.ok) throw new Error("Failed to fetch comments");

    const { data } = await res.json();

    const html = (data && data.length)
      ? data.map(c => {
        
          const upClass = c.IsLiked ? "fa-solid" : "fa-regular";
          const downClass = c.IsDisliked ? "fa-solid" : "fa-regular";
          return `
            <div class="comment" data-id="${c.id}">
              <div><strong>${c.username}</strong></div>
              <div>${c.content}</div>
              <div class="button-group">
                <button class="likeComment" data-id="${c.id}" aria-label="Like comment">
                  <span>${c.Likes}</span>
                  <i class="${upClass} fa-thumbs-up"></i>
                </button>
                <button class="disLikeComment" data-id="${c.id}" aria-label="Dislike comment">
                  <span>${c.Dislikes}</span>
                  <i class="${downClass} fa-thumbs-down"></i>
                </button>
              </div>
            </div>
          `;
        }).join("")
      : "<h4>No comments available</h4>";

    commentsContainer.innerHTML = html + CommentForm(postId);
    AddComments(postId);
  }
  catch (err) {
    popup("Failed to load comments.", "failed");
  }
}

/**
 * Click handler for the "Show Comments" button on each post.
 * Toggles visibility and fetches on open.
 */
export async function CommentSection(e) {
  const postEl = e.target.closest(".post");
  if (!postEl) return;
  const commentsContainer = postEl.querySelector(".comments");
  if (!commentsContainer) return;

  const hidden = commentsContainer.style.display === "none" || !commentsContainer.style.display;
  commentsContainer.style.display = hidden ? "block" : "none";
  if (hidden) {
    await fetchAndRenderComments(postEl);
  }
}

/**
 * Delegated click handler for like/dislike on comments.
 * On success: re-fetches that post’s comments.
 * Delegated on document to catch dynamically added buttons.
 */
export function setupCommentReactions() {
  document.addEventListener("click", async e => {
    const btn = e.target.closest(".likeComment, .disLikeComment");
    if (!btn) return;

    const commentDiv = btn.closest(".comment");
    if (!commentDiv) return;

    if (!await isLogged()) {
      return popup("You need to be logged in to react to comments.", "warning");
    }

    const commentID = parseInt(btn.dataset.id, 10);
    const status = btn.classList.contains("likeComment") ? "like" : "dislike";

    try {
      const res = await fetch("/api/addLike", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ commentID, sender: "comment", status }),
      });

      if (!res.ok) {
        const err = await res.json();
        return popup(err.message || "Failed to react to comment.", "error");
      }

      // On success, re-render this post’s comments
      const postEl = commentDiv.closest(".post");
      await fetchAndRenderComments(postEl);
    }
    catch (err) {
      popup("Something went wrong!", "failed");
    }
  });
}

// === Auto-wire after DOM is ready ===
document.addEventListener("DOMContentLoaded", () => {
  // Attach "Show comments" on each .displayComment button
  document.querySelectorAll(".displayComment").forEach(btn => {
    btn.addEventListener("click", CommentSection);
  });

  // Set up like/dislike delegation on document
  setupCommentReactions();
});
