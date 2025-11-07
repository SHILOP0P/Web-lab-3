(() => {
  // Сообщаем странице, что авторизация изменилась
  function emitAuthChanged(user) {
    window.dispatchEvent(new CustomEvent('auth:changed', { detail: { user } }));
  }

  // DOM-элементы хедера
  function getHeaderEls() {
    const cell  = document.querySelector(".header-login");
    const table = cell?.querySelector(".login-table"); // текущая таблица с полями
    const uEl   = document.getElementById("loginUsername");
    const pEl   = document.getElementById("loginPassword");

    let state = cell?.querySelector("#authState");
    if (!state && cell) {
      state = document.createElement("div");
      state.id = "authState";
      state.style.display = "none";
      state.style.padding = "6px 0";
      cell.appendChild(state);
    }

    let logoutBtn = cell?.querySelector("#logoutBtn");
    if (!logoutBtn && cell) {
      logoutBtn = document.createElement("button");
      logoutBtn.id = "logoutBtn";
      logoutBtn.className = "login-button";
      logoutBtn.textContent = "Выйти";
      logoutBtn.style.display = "none";
      logoutBtn.addEventListener("click", async (e) => {
        e.preventDefault();
        try { await window.AUTH.logout(); } catch {}
        showLoggedOut();       // моментально переключаем шапку
        emitAuthChanged(null); // говорим странице: пользователь вышел
        // selfReload();       // запасной план — не нужен
      });
      cell.appendChild(logoutBtn);
    }

    return { cell, table, uEl, pEl, state, logoutBtn };
  }

  function showLoggedIn(user) {
    const { table, state, logoutBtn, uEl, pEl } = getHeaderEls();
    if (uEl) uEl.value = "";
    if (pEl) pEl.value = "";
    if (table) table.style.display = "none";
    if (state) {
      state.innerHTML = `Здравствуйте, <b>${user.username || user.email}</b>`;
      state.style.display = "block";
    }
    if (logoutBtn) logoutBtn.style.display = "inline-block";

    const name = (user?.username || user?.email || '').toString();
    if (name.toUpperCase() === 'SHILOP0P') addAdminLink();
    else removeAdminLink();
  }

  function showLoggedOut() {
    const { table, state, logoutBtn } = getHeaderEls();
    if (table) table.style.display = "";
    if (state) { state.style.display = "none"; state.textContent = ""; }
    if (logoutBtn) logoutBtn.style.display = "none";
    removeAdminLink();
  }

  // === Админ-линк в сайдбаре (для SHILOP0P) ===
  function findSidebarContainer() {
    // ищем блок, где лежат пункты "Главная/О нас/Каталог/Контакты"
    const anchor = Array.from(document.querySelectorAll('a')).find(a => {
      const t = (a.textContent || '').trim();
      return t === 'Главная' || t === 'О нас' || t === 'Каталог' || t === 'Контакты';
    });
    if (!anchor) return null;
    // берём ближайший "контейнер" для ссылок
    return anchor.closest('nav, aside, div, ul, ol') || anchor.parentElement;
  }

  function adminHrefBase() {
    // чтобы работало и при /front/... и при корне
    return location.pathname.includes('/front/') ? '/front/' : '/';
  }

  function addAdminLink() {
    const box = findSidebarContainer();
    if (!box) return;
    if (box.querySelector('[data-admin-link]')) return; // уже есть
    const a = document.createElement('a');
    a.setAttribute('data-admin-link', '');
    a.href = adminHrefBase() + 'admin/main_admin.html'; // твоя админка
    a.textContent = 'Админ-панель';
    a.style.display = 'block'; // под твой вертикальный список
    box.appendChild(a);
  }

  function removeAdminLink() {
    document.querySelectorAll('[data-admin-link]').forEach(n => n.remove());
  }

  // Глобальная функция, на которую уже навешан onclick="login()" в твоих html
  async function login(ev) {
    // если вызвали как onclick="login(event)" — гасим сабмит
    if (ev && typeof ev.preventDefault === 'function') ev.preventDefault();
    // если вызвали без event (onclick="login()"), притушим возможный window.event
    if (typeof window !== 'undefined' && window.event && typeof window.event.preventDefault === 'function') {
      window.event.preventDefault();
    }

    // ВОТ ЭТИХ СТРОК НЕ ХВАТАЛО
    const usernameOrEmail = document.getElementById("loginUsername")?.value?.trim();
    const password = document.getElementById("loginPassword")?.value ?? "";

    const data = await window.AUTH.login({ usernameOrEmail, password });
    const user = data?.user || data;
    showLoggedIn(user);     // моментально переключаем шапку и admin-линк
    emitAuthChanged(user);  // говорим странице: пользователь вошёл/сменился
  }

  function selfReload() {
    const a = ensureSelfReloadLink(); // href = текущая страница с ?t=...
    a.click();                        // реальный переход по ссылке
  }

  // Создаёт/обновляет скрытую ссылку на текущую страницу с анти-кэшем
  function ensureSelfReloadLink() {
    let a = document.getElementById('selfReloadLink');
    if (!a) {
      a = document.createElement('a');
      a.id = 'selfReloadLink';
      a.style.display = 'none';
      a.rel = 'nofollow';
      document.body.appendChild(a);
    }
    const url = new URL(window.location.href);
    url.searchParams.set('t', Date.now().toString()); // анти-кэш
    a.href = url.toString();
    return a;
  }

  // Экспорт и подготовка
  window.ensureSelfReloadLink = ensureSelfReloadLink;
  ensureSelfReloadLink();
  window.login = login;

  // Инициализация при загрузке: если уже авторизован — покажем приветствие
  async function initHeaderAuth() {
    try {
      const resp = await fetch('/api/auth/me', { credentials: 'include', cache: 'no-store' });

      if (resp.status === 401) {
        // гость — это норма
        showLoggedOut();
        emitAuthChanged(null);
        return;
      }
      if (!resp.ok) {
        console.warn('auth/me:', resp.status);
        showLoggedOut();
        emitAuthChanged(null);
        return;
      }

      const data = await resp.json();
      const user = data.user || data;
      showLoggedIn(user);
      emitAuthChanged(user);
    } catch (e) {
      console.warn('auth/me failed:', e);
      showLoggedOut();
      emitAuthChanged(null);
    }
  }

  // Автостарт
  if (document.readyState === 'complete') initHeaderAuth();
  else window.addEventListener('load', initHeaderAuth, { once: true });
})();
