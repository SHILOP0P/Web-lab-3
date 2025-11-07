// front/js/reviews/reviews-init.js

// Базовый адрес API (учитываем запуск фронта не всегда на :8080)
const BACKEND_ORIGIN =
  (window.AUTH && AUTH.ORIGIN) ||
  ((location.port === '8080')
    ? `${location.protocol}//${location.host}`
    : `${location.protocol}//${location.hostname}:8080`);

const API_BASE = `${BACKEND_ORIGIN}/api`;

let ME = null;

document.addEventListener('DOMContentLoaded', async () => {
  await refreshAll();

  const sendBtn = document.getElementById('sendBtn');
  if (sendBtn) sendBtn.addEventListener('click', submitReview);
});

// Если авторизация поменялась (например, без перезагрузки) — перестроим UI
window.addEventListener('auth:changed', async (e) => {
  ME = e?.detail?.user ?? null;
  toggleForm();
  await loadReviews();
});

async function refreshAll() {
  await loadMe();
  toggleForm();
  await loadReviews();
}

function toggleForm() {
  const form = document.getElementById('reviewFormBox');
  const hint = document.getElementById('loginHint');
  if (!form || !hint) return;

  if (ME) {
    form.style.display = '';
    hint.style.display = 'none';
  } else {
    form.style.display = 'none';
    hint.style.display = '';
  }
}

async function loadMe() {
  try {
    const r = await fetch(`${API_BASE}/auth/me?t=${Date.now()}`, {
      credentials: 'include',
      cache: 'no-store'
    });
    if (r.ok) {
      const data = await r.json();
      ME = data && data.user ? data.user : null; // ожидается { user: { id, username, ... } }
    } else {
      ME = null;
    }
  } catch {
    ME = null;
  }
}

async function loadReviews() {
  const root = document.getElementById('reviewsList');
  if (!root) return;
  root.innerHTML = '<p>Загрузка…</p>';

  try {
    const r = await fetch(`${API_BASE}/reviews?t=${Date.now()}`, {
      credentials: 'include',
      cache: 'no-store'
    });
    if (!r.ok) throw new Error('load failed');

    const data = await r.json();
    const list = Array.isArray(data.reviews) ? data.reviews : [];

    if (!list.length) {
      root.innerHTML = '<p>Пока нет ни одного отзыва.</p>';
      return;
    }

    root.innerHTML = '';
    for (const rev of list) {
      root.appendChild(renderReview(rev));
    }
  } catch (e) {
    console.error(e);
    root.innerHTML = '<p>Не удалось загрузить отзывы.</p>';
  }
}

async function submitReview() {
  const msg = document.getElementById('formMsg');
  const textarea = document.getElementById('reviewText');
  if (!msg || !textarea) return;

  msg.textContent = '';
  const text = (textarea.value || '').trim();
  if (!text) { msg.textContent = 'Введите текст отзыва.'; return; }

  try {
    const r = await fetch(`${API_BASE}/reviews?t=${Date.now()}`, {
      method: 'POST',
      credentials: 'include',
      cache: 'no-store',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ content: text }) // сервер ждёт "content"
    });

    if (!r.ok) {
      msg.textContent = (r.status === 401)
        ? 'Только авторизованные пользователи могут оставлять отзывы.'
        : 'Не удалось сохранить отзыв.';
      return;
    }

    textarea.value = '';
    await loadReviews(); // сразу перерисуем список
  } catch (e) {
    console.error(e);
    msg.textContent = 'Ошибка сети.';
  }
}

function canDelete(rev) {
  if (!ME) return false;
  const isAdmin = String(ME.username).toUpperCase() === 'SHILOP0P';
  const ownerId = Number(rev.user_id ?? rev.userId); // нормализуем имя поля
  const isOwner = Number(ME.id) === ownerId;
  return isAdmin || isOwner;
}

function renderReview(rev) {
  // Ожидаемые поля: id, user_id/userId, username, content/text, created_at/createdAt/created
  const wrap = document.createElement('div');
  wrap.className = 'product-card';

  const head = document.createElement('div');
  head.style.display = 'flex';
  head.style.justifyContent = 'space-between';
  head.style.gap = '8px';

  const left = document.createElement('div');
  left.innerHTML = `<b>${escapeHtml(rev.username || '—')}</b>`;

  const right = document.createElement('div');
  const createdRaw = rev.created_at ?? rev.createdAt ?? rev.created ?? null;
  const ds = formatDate(createdRaw);
  right.innerHTML = ds ? `<span style="opacity:.7">${escapeHtml(ds)}</span>` : '';

  head.appendChild(left);
  head.appendChild(right);

  const body = document.createElement('p');
  const content = rev.content ?? rev.text ?? '';
  body.textContent = String(content);

  const actions = document.createElement('div');
  actions.style.textAlign = 'right';

  if (canDelete(rev)) {
    const btn = document.createElement('button');
    btn.className = 'btn-primary';
    btn.style.background = '#d9534f';
    btn.style.color = '#fff';
    btn.textContent = 'Удалить';
    btn.onclick = () => delReview(rev.id);
    actions.appendChild(btn);
  }

  wrap.appendChild(head);
  wrap.appendChild(body);
  wrap.appendChild(actions);
  return wrap;
}

async function delReview(id) {
  if (!confirm('Удалить этот отзыв?')) return;

  try {
    const r = await fetch(`${API_BASE}/reviews/${encodeURIComponent(id)}?t=${Date.now()}`, {
      method: 'DELETE',
      credentials: 'include',
      cache: 'no-store'
    });

    if (r.ok) {
      await loadReviews();
    } else {
      alert(r.status === 403 ? 'Нет прав на удаление.' : 'Не удалось удалить.');
    }
  } catch (e) {
    console.error(e);
    alert('Ошибка сети.');
  }
}

// Форматирование даты из ISO, timestamp или строк вида "YYYY-MM-DD HH:mm:ss"
function formatDate(raw) {
  if (raw == null) return '';

  if (typeof raw === 'number') {
    const ms = raw > 1e12 ? raw : raw * 1000;
    const d = new Date(ms);
    return isNaN(d) ? '' : d.toLocaleString();
  }

  let s = String(raw).trim();
  if (!s) return '';

  let d = new Date(s);
  if (!isNaN(d)) return d.toLocaleString();

  if (/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}/.test(s)) {
    s = s.replace(' ', 'T');
    if (!/[zZ]|[+\-]\d{2}:\d{2}$/.test(s)) s += 'Z';
    d = new Date(s);
    if (!isNaN(d)) return d.toLocaleString();
  }

  return '';
}

function escapeHtml(str = '') {
  return String(str)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;');
}
