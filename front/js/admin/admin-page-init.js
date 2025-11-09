document.addEventListener('DOMContentLoaded', init);

const BACKEND_ORIGIN =
  (location.port === '8080')
    ? `${location.protocol}//${location.host}`
    : `${location.protocol}//${location.hostname}:8080`;

const API_BASE = `${BACKEND_ORIGIN}/api`;       // для данных
const IMG_BASE = `${BACKEND_ORIGIN}/images_db`; // для картинок


async function init() {
  try {
    const products = await getProducts();     // массив продуктов
    renderProducts(products);                 // рисуем список
    const root = document.getElementById('productList');
    root.addEventListener('click', onListClick);
  } catch (e) {
    console.error(e);
    const root = document.getElementById('productList');
    if (root) root.innerHTML = `<div class="error">Не удалось загрузить товары</div>`;
  }
}

async function getProducts() {
  const resp = await fetch(`${API_BASE}/products`);
  if (!resp.ok) throw new Error('Failed to load products');
  return resp.json();
}

async function DeleteProduct(productId) {
  const resp = await fetch(`${API_BASE}/products/${productId}`, { method: 'DELETE' });
  if (!resp.ok) throw new Error('Не удалось удалить продукт');
  return resp.json().catch(() => ({}));
}

function renderProducts(products = []) {
  const root = document.getElementById('productList');        
  if (!root) return;

  root.innerHTML = ''; // очистим перед отрисовкой

  products.forEach(p => {
    const card = renderCard(p);
    root.appendChild(card);
  });
}

function renderCard(product) {
  // Логируем весь объект продукта для проверки

  // Изображения
  // В ответе приходит product_images (множественное). Подстрахуемся и под оба варианта.
  const img = product.product_images || product.product_image;

  const imagesHtml = (img && img.image)
    ? `
      <tr>
        <td>
          <img
            src="${IMG_BASE}/${img.image}"
            alt="${escapeHtml((img.title || product.name || '').toString())}"
            onerror="this.src='${IMG_BASE}/placeholder.png'"
            style="width: 100px;"
          />
        </td>
      </tr>`
    : `
      <tr>
        <td>
          <img src="${IMG_BASE}/placeholder.png" alt="нет изображения" style="width: 100px;" />
        </td>
      </tr>`;



  // Характеристики
const characteristics =
  product.product_property?.characteristics || "—";
const formattedCharacteristics = characteristics.replace(/\r?\n/g, '<br />'); // Заменяем переносы строк на <br />

  // Описание
  const description = product.description || "—";  // Извлекаем описание
  const formattedDescription = description.replace(/\r?\n/g, '<br />'); // Заменяем переносы строк на <br />

  // Метаданные
  const metaDescription = product.meta_description || "—";
  const metaKeywords = product.meta_keywords || "—";
  const metaTitle = product.meta_title || "—";
  
  // Краткое описание
  const shortDescription = product.short_description || "—";

  // Сформируем HTML для свойств, используя таблицы
  const propsHtml = `
    <table class="product-table">
      <tr><th>Характеристики:</th><td>${formattedCharacteristics}</td></tr>
      <tr><th>Описание:</th><td>${formattedDescription}</td></tr>
      <tr><th>Краткое описание:</th><td>${escapeHtml(shortDescription)}</td></tr>
      <tr><th>Метаданные (Meta Title):</th><td>${escapeHtml(metaTitle)}</td></tr>
      <tr><th>Метаданные (Meta Description):</th><td>${escapeHtml(metaDescription)}</td></tr>
      <tr><th>Метаданные (Meta Keywords):</th><td>${escapeHtml(metaKeywords)}</td></tr>
    </table>
  `;

  const card = document.createElement('div');
  card.className = 'product-card';
  card.dataset.productId = product.id; // пригодится при удалении

  card.innerHTML = `
    <h3 class="product-title">${escapeHtml(product.name)}</h3>
    <table class="product-table">
      <tr><th>Цена:</th><td>${product.price ?? '—'} руб.</td></tr>
      <tr><th>Категория:</th><td>${escapeHtml(product.alias || '—')}</td></tr>
    </table>

    <table class="product-table">
      <thead>
        <tr><th colspan="2">Изображения</th></tr>
      </thead>
      <tbody>
        ${imagesHtml}
      </tbody>
    </table>

    <table class="product-table">
      <thead>
        <tr><th colspan="2">Свойства</th></tr>
      </thead>
      <tbody>
        ${propsHtml}
      </tbody>
    </table>

    <table class="product-table">
      <tr><th>Действия</th></tr>
      <tr>
        <td>
          <button class="btn btn--danger" data-action="del-product" data-product-id="${product.id}">Удалить</button>
        </td>
      </tr>
    </table>
  `;

  return card;
}

function escapeHtml(s = '') {
  return String(s)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;');
}



// Делегирование кликов по кнопкам
async function onListClick(e) {
  const btn = e.target.closest('button[data-action]');
  if (!btn) return;

  const action = btn.dataset.action;
  const productId = btn.dataset.productId;
  const imgId = btn.dataset.imgId;

  switch (action) {
    case 'add-image':
      console.log('Добавить изображение для продукта', productId);
      break;

    case 'del-image':
      console.log('Удалить изображение', imgId, 'у продукта', productId);
      break;

    case 'add-prop':
      console.log('Добавить/изменить свойства для продукта', productId);
      break;

    case 'edit':
      console.log('Редактировать продукт', productId);
      break;

    case 'del-product':
      console.log('Удалить продукт', productId);
      try {
        await DeleteProduct(productId);
        const card = document.querySelector(`.product-card[data-product-id="${productId}"]`);
        if (card) card.remove();
      } catch (err) {
        console.error(err);
      }
      break;
  }
}