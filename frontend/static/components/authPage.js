import { navigateTo } from "../js/app.js";
import { errorPage } from "./errorPage.js";

export function loginPage() {
  return /*html*/`
          <form id="loginForm">
            <h2>Login</h2>
          <input type="text" placeholder="Enter your email">
          <input type="password" placeholder="Enter your password">
          <button type="submit">Login</button>
          </form>
      `
}

export function registerPage() {
  return /*html*/`
          <form id="registerForm">
              <h2>Register</h2>
            <input type="text" name="username"     placeholder="Username"    />
            <span class="errRgister" id="username"></span>
            <input type="email" name="email"        placeholder="Email"       />
            <span class="errRgister" id="email"></span>
            <input type="text" name="firstName"    placeholder="First Name"  />
            <span class="errRgister" id="firstName"></span>

            <input type="text" name="lastName"     placeholder="Last Name"   />
            <span class="errRgister" id="lastName"></span>

            <input type="number" name="age"         placeholder="Age"         />
            <span class="errRgister" id="age"></span>

            <select name="gender"                   >
              <option value="">Select Gender</option>
              <option value="male">Male</option>
              <option value="female">Female</option>
            </select>
            <span class="errRgister" id="gender"></span>

            <input type="password" name="password" placeholder="Password"    />
            <span class="errRgister" id="password"></span>

            <button type="submit">Register</button>
          </form>
      `
}

export function register() {
  const form = document.querySelector("#registerForm");
  const spans = document.querySelectorAll(".errRgister");

  form.addEventListener("submit", async e => {
    e.preventDefault();

    spans.forEach(span => {
      span.innerHTML = "";
      span.style.display = "none";
    });

    const formData = Object.fromEntries(new FormData(form).entries());
    formData.age = parseInt(formData.age)

    try {
      const response = await fetch("/api/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(formData)
      });

      const data = await response.json()

      if (!response.ok) {
        for (let span of spans) {
          if (data.hasOwnProperty(span.id))
            showRegisterInputError(data[span.id], span)
        }
      } else {
        navigateTo("/");
      }
    } catch (err) {
      console.error(err)
      document.body.innerHTML = errorPage("Something went wrong!", 500)
    }
  });
}

export function login() {
  const form = document.querySelector("#loginForm");

  form.addEventListener('submit', e => {
    e.preventDefault()

    console.log('first')
  })
}

function showRegisterInputError(msg, span) {
  span.style.display = "block"
  span.style.color = "red"
  span.innerHTML = msg
}