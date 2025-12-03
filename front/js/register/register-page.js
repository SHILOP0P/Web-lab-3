(() => {
  function val(id) { return document.getElementById(id)?.value?.trim() || ""; }

  async function onRegister(e) {
    e.preventDefault();
    const msg = document.getElementById("message");
    const username   = val("username");
    const email      = val("email");
    const password   = document.getElementById("password")?.value || "";
    const confirmPwd = document.getElementById("confirmPassword")?.value || "";
    const firstName  = val("firstName") || null;
    const lastName   = val("lastName")  || null;
    const phone      = val("phone")     || null;
    const birthdate  = document.getElementById("birthdate")?.value || null;
    const region     = document.getElementById("region")?.value || null;
    const gender     = document.querySelector("input[name='gender']:checked")?.value || null;

    if (password !== confirmPwd) {
      if (msg) { msg.textContent = "Пароли не совпадают!"; msg.style.color = "red"; }
      return;
    }
    if (!username || !email || !password) {
      if (msg) { msg.textContent = "Заполните логин, email и пароль"; msg.style.color = "red"; }
      return;
    }

    const agree = document.getElementById("agreeTerms");
    if (!agree || !agree.checked) {
      if (msg) {
        msg.textContent = "Чтобы зарегистрироваться, необходимо согласиться с Пользовательским соглашением и Политикой конфиденциальности.";
        msg.style.color = "red";
      }
      return;
    }


    try {
      await window.AUTH.register({ username, email, password, firstName, lastName, phone, gender, birthdate, region });
      if (msg) { msg.textContent = "Регистрация успешна!"; msg.style.color = "green"; }
      setTimeout(() => (window.location.href = "index.html"), 700);
    } catch (err) {
      if (msg) { msg.textContent = "Ошибка: " + err.message; msg.style.color = "red"; }
    }
  }

  function bindRegisterForm() {
    const form = document.getElementById("registerForm");
    if (form) {
      form.removeEventListener("submit", onRegister);
      form.addEventListener("submit", onRegister);
    }
  }

  window.addEventListener("DOMContentLoaded", bindRegisterForm);
})();
