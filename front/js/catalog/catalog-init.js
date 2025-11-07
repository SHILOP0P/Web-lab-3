// js/catalog-init.js
const BACKEND_ORIGIN =
  (location.port === '8080')
    ? `${location.protocol}//${location.host}`
    : `${location.protocol}//${location.hostname}:8080`;

const API_BASE = `${BACKEND_ORIGIN}/api`;
const IMG_BASE = `${BACKEND_ORIGIN}/images_db`;

// соответствие секций каталога видимым названиям категорий в БД (поле alias)
const CATEGORY_MAP = {
  materials: 'Строительные материалы',
  finishing: 'Отделочные материалы',
  tools: 'Инструменты'
};

document.addEventListener('DOMContentLoaded', async () => {
  await Promise.all([
    loadCategory('materials', 'materialsList'),
    loadCategory('finishing', 'finishingList'),
    loadCategory('tools', 'toolsList')
  ]);
});

async function loadCategory(sectionKey, targetId) {
  const alias = CATEGORY_MAP[sectionKey];
  const root = document.getElementById(targetId);
  if (!root) return;

  try {
    const resp = await fetch(`${API_BASE}/catalog/cards?category=${encodeURIComponent(alias)}`);
    if (!resp.ok) throw new Error('Failed to load cards');
    const items = await resp.json();

    if (!Array.isArray(items) || items.length === 0) return;

    root.innerHTML = '';
    items.forEach(item => root.appendChild(renderCard(item)));
  } catch (e) {
    console.error(e);
  }
}

// ---------- helpers ----------

// Берём имя файла картинки из разных форматов ответа и нормализуем путь
function pickImageFile(item) {
  // 1) новый формат: { product_images: { image: "..." } }
  let file =
    item?.product_images?.image ??
    // 2) старый/упрощённый формат: { image: "..." }
    item?.image ??
    '';

  // убираем возможные префиксы и пробелы
  file = String(file).replace(/^\/?images_db\//, '').trim();
  return file;
}

function escapeHtml(s = '') {
  return String(s)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;');
}

// ---------- render ----------

function renderCard(item) {
  // item: { id, name, price, short_description, image? } ИЛИ { product_images: { image } }
  const a = document.createElement('a');
  a.className = 'product';
  a.href = `products/product_page.html?id=${encodeURIComponent(item.id)}`;

  const img = document.createElement('img');
  img.loading = 'lazy';
  const imgFile = pickImageFile(item);
  img.src = imgFile ? `${IMG_BASE}/${imgFile}` : `${IMG_BASE}/placeholder.png`;
  img.alt = item.name || 'product';
  img.onerror = () => { img.src = `${IMG_BASE}/placeholder.png`; };

  const info = document.createElement('div');
  info.className = 'product-info';

  const priceText = (item.price != null && item.price !== '')
    ? `${item.price} ₽`
    : '—';

  info.innerHTML = `
    <h4>${escapeHtml(item.name || '—')}</h4>
    <p>${priceText}</p>
  `;

  const descr = document.createElement('div');
  descr.className = 'product-description';
  descr.textContent = item.short_description || '';

  a.appendChild(img);
  a.appendChild(info);
  a.appendChild(descr);
  return a;
}
