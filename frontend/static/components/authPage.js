export function loginPage() {
    return /*html*/`
        <input type="text" placeholder="Enter your email">
        <input type="password" placeholder="Enter your password">
        <button>Login</button>
    `
}

export function registerPage() {
    return /*html*/`
        <form id="registerForm">
            <h2>Register</h2>
          <input type="text" name="username"     placeholder="Username"    />
          <input type="email" name="email"        placeholder="Email"       />
          <input type="text" name="firstName"    placeholder="First Name"  />
          <input type="text" name="lastName"     placeholder="Last Name"   />
          <input type="number" name="age"         placeholder="Age"         />
          <select name="gender"                   >
            <option value="">Select Gender</option>
            <option value="Male">Male</option>
            <option value="Female">Female</option>
          </select>
          <input type="password" name="password" placeholder="Password"    />
          <button type="submit">Register</button>
        </form>

    `
}

export function register() {
    const form = document.querySelector("#registerForm");
    form.addEventListener("submit", async e => {
        e.preventDefault();

        console.log(new FormData(form))

        const data = Object.fromEntries(new FormData(form).entries());

        console.log(data)
        //   try {
        //     const res = await fetch("/api/register", {
        //       method: "POST",
        //       headers: { "Content-Type": "application/json" },
        //       body: JSON.stringify(data)
        //     });
        //     if (!res.ok) throw await res.json();
        //     alert("Registered successfully!");
        //     navigateTo("/login");
        //   } catch (err) {
        //     alert("Error: " + (err.message || JSON.stringify(err)));
        //   }
    });
}