// Базы адресов (корректно работают и с 5500, и с 8080)
 const BACKEND_ORIGIN = (window.AUTH && AUTH.ORIGIN) ||
   ((location.port === '8080')
     ? `${location.protocol}//${location.host}`
     : `${location.protocol}//${location.hostname}:8080`);
 const IMG_BASE = (window.AUTH && AUTH.IMG_BASE) || `${BACKEND_ORIGIN}/images_db`;
 const API_BASE = `${BACKEND_ORIGIN}/api`;


document.addEventListener('DOMContentLoaded', init);

async function init() {
  const id = new URLSearchParams(location.search).get('id');
  if (!id) return;

  try {
    const resp = await fetch(`${API_BASE}/products/${encodeURIComponent(id)}`);
    if (!resp.ok) {
      renderNotFound();
      return;
    }
    const p = await resp.json();
    console.log(p)
    renderProduct(p);
  } catch (e) {
    console.error(e);
    renderNotFound();
  }
}

function renderProduct(p) {
  // title + <title> вкладки
  text('pTitle', p.name || 'Товар');
  document.title = `${p.name || 'Товар'} - Добрострой`;

  // главное изображение — объект product_images или placeholder
  const imgEl = document.getElementById('pImage');

  const rawImg = (p.product_images && p.product_images.image) ? String(p.product_images.image) : '';
  // срежем возможный префикс и пробелы
  const file = rawImg.replace(/^\/?images_db\//, '').trim();

  imgEl.src = file ? `${IMG_BASE}/${file}` : `${IMG_BASE}/placeholder.png`;
  imgEl.onclick = () => window.open(imgEl.src, '_blank');
  imgEl.onerror = () => { imgEl.src = `${IMG_BASE}/placeholder.png`; };



  // цена в стиле cement.html
  text('priceH3', p.price != null ? `${p.price} Рублей` : '—');

  // короткое описание (как абзац)
  text('pShort', p.short_description || '—');

  // характеристики: у тебя это многострочный текст в product_properties.characteristics
  const ul = document.getElementById('pCharsList');

  const charsText = (p.product_properties && typeof p.product_properties.characteristics === 'string')
    ? p.product_properties.characteristics
    : '';

  const items = charsText
    .replace(/\r\n/g, '\n')
    .split('\n')
    .map(s => s.trim())
    .filter(Boolean);

  ul.innerHTML = items.length
    ? items.map(s => `<li>${escapeHtml(s)}</li>`).join('')
    : '<li>—</li>';



  // подробное описание: делим на параграфы по пустой строке
  const desc = String(p.description || '').replace(/\r\n/g, '\n');
  const parts = desc.split(/\n\s*\n/).map(s => s.trim()).filter(Boolean);
  const full = document.getElementById('pFullDesc');
  full.innerHTML = parts.length
    ? parts.map(s => `<p class="full-description">${escapeHtml(s)}</p>`).join('')
    : '<p class="full-description">—</p>';
}

// на случай ошибок
function renderNotFound() {
  text('pTitle', 'Товар не найден');
  const imgEl = document.getElementById('pImage');
  imgEl.src = `${IMG_BASE}/placeholder.png`;
  text('priceH3', '—');
  text('pShort', '—');
  document.getElementById('pCharsList').innerHTML = '<li>—</li>';
  document.getElementById('pFullDesc').innerHTML = '<p class="full-description">—</p>';
}

// утилиты
function text(id, value) {
  const el = document.getElementById(id);
  if (el) el.textContent = value;
}
function escapeHtml(s='') {
  return String(s)
    .replaceAll('&','&amp;')
    .replaceAll('<','&lt;')
    .replaceAll('>','&gt;')
    .replaceAll('"','&quot;')
    .replaceAll("'",'&#039;');
}
