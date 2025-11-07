// front/js/auth/auth-api.js

(() => {
  // авто-детект адреса бэка как у тебя в init-скриптах
  const BACKEND_ORIGIN =
    (location.port === "8080")
      ? `${location.protocol}//${location.host}`
      : `${location.protocol}//${location.hostname}:8080`;

  async function request(path, options = {}) {
    const resp = await fetch(BACKEND_ORIGIN + path, {
      credentials: "include",
      headers: { "Content-Type": "application/json", ...(options.headers || {}) },
      ...options,
    });
    let data = null;
    try { data = await resp.json(); } catch {} // тело может быть пустым
    if (!resp.ok) {
      const msg = (data && data.error) ? data.error : `HTTP ${resp.status}`;
      const err = new Error(msg);
      err.status = resp.status;
      throw err;
    }
    return data;
  }

  // запрос, который НЕ роняет страницу на 401/404/405 (удобно для /auth/me)
  async function requestOptional(path, options = {}) {
    try {
      return await request(path, options);
    } catch (e) {
      if (e && (e.status === 401 || e.status === 404 || e.status === 405)) {
        return null; // нет сессии / нет эндпоинта — считаем как "не авторизован"
      }
      throw e; // остальные ошибки пробрасываем
    }
  }

  window.AUTH = {
    // базовые методы
    me()        { return request("/api/auth/me"); },
    login(p)    { return request("/api/auth/login",    { method: "POST", body: JSON.stringify(p) }); },
    register(p) { return request("/api/auth/register", { method: "POST", body: JSON.stringify(p) }); },
    logout()    { return request("/api/auth/logout",   { method: "POST" }); },

    // мягкая проверка авторизации (не бросает исключение на 401/404/405)
    async meOptional() { return await requestOptional("/api/auth/me"); },

    // удобные константы для других файлов
    ORIGIN: BACKEND_ORIGIN,
    IMG_BASE: `${BACKEND_ORIGIN}/images_db`,
  };
})();
