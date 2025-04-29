import { navigateTo } from "./app.js";
import { errorPage, popup } from "./errorPage.js";

export function loginPage() {
  document.body.innerHTML = /*html*/`
    <div class="auth">
    <form class="authForm" id="loginForm">
      <h2>Login</h2>
      <div class="errLogin"></div>

      <input required name="email" type="email" placeholder="Enter your email" />
      <span class="errLoginField" id="email"></span>

      <input required name="password" type="password" placeholder="Enter your password" />
      <span class="errLoginField" id="password"></span>

      <button type="submit">Login</button>

      <p class="form-switch">Don't have an account? 
        <a href="/register" data-link>Create New Account</a>
      </p>
    </form>
    </div>
  `
}

export function registerPage() {
  document.body.innerHTML = /*html*/`
    <div class="auth">
    <form class="authForm" method="post" id="registerForm">
      <h2>Register</h2>

      <input required type="text" name="username" placeholder="Username" />
      <span class="errRgister" id="username"></span>

      <input required type="email" name="email" placeholder="Email" />
      <span class="errRgister" id="email"></span>

      <input required type="text" name="firstName" placeholder="First Name" />
      <span class="errRgister" id="firstName"></span>

      <input required type="text" name="lastName" placeholder="Last Name" />
      <span class="errRgister" id="lastName"></span>

      <input required type="number" name="age" placeholder="Age" />
      <span class="errRgister" id="age"></span>

      <select required name="gender">
        <option value="">Select Gender</option>
        <option value="male">Male</option>
        <option value="female">Female</option>
      </select>
      <span class="errRgister" id="gender"></span>

      <input required type="password" name="password" placeholder="Password" />
      <span class="errRgister" id="password"></span>

      <button type="submit">Register</button>

      <p class="form-switch">Already have an account? 
        <a href="/login" data-link>Login</a>
      </p>
    </form>
    </div>
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
            showInputError(data[span.id], span)
            popup("Registration Failed",'failed')

        }
      } else {
        navigateTo("/")
        }
    } catch (err) {
      console.error(err)
      document.body.innerHTML = popup("Something went wrong!", "failed")
    }
  });
}

export function login() {
  const form = document.querySelector("#loginForm");
  const errLogin = document.querySelector(".errLogin");

  errLogin.style.color = 'red'

  form.addEventListener('submit', async e => {
    e.preventDefault()


    const formData = Object.fromEntries(new FormData(form).entries());

    try {
      const response = await fetch("/api/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(formData)
      })

      const data = await response.json()

      if (!response.ok) {
        errLogin.innerHTML = data.message
        popup("login Failed",'failed')

      } else {
        navigateTo("/")
            }
    } catch (err) {
      console.error(err)
      document.body.innerHTML = popup("Something went wrong!", "failed")
    }
  })
}

export async function logout() {
  try {
    await fetch("/api/logout")
    navigateTo("/register")
    } catch (err) {
    console.error(err)
    document.body.innerHTML = popup("Something went wrong!", "failed")
  }
}

export function showInputError(msg, span) {
  span.style.display = "block"
  span.style.color = "red"
  span.innerHTML = msg
}