// add_product.js
const BACKEND_ORIGIN =
  (location.port === '8080')
    ? `${location.protocol}//${location.host}`
    : `${location.protocol}//${location.hostname}:8080`;

const API_BASE = `${BACKEND_ORIGIN}/api`;
const IMG_BASE = `${BACKEND_ORIGIN}/images_db`;

document.addEventListener('DOMContentLoaded', () => {
  const submitBtn = document.getElementById('submitProductBtn');
  if (submitBtn) submitBtn.addEventListener('click', submitProduct);
});

function toInt(id) {
  const v = parseInt(document.getElementById(id).value, 10);
  return Number.isFinite(v) ? v : 0;
}
function toFloat(id) {
  const v = parseFloat(document.getElementById(id).value);
  return Number.isFinite(v) ? v : 0;
}
function val(id) {
  return document.getElementById(id).value.trim();
}

async function submitProduct() {
  const result = document.getElementById('resultBox');
  result.style.display = 'none';
  result.textContent = '';

  const formData = new FormData();

  // Данные о продукте (как JSON в поле "product")
  formData.append('product', JSON.stringify({
    manufacturer_id: toInt('manufacturer_id'),
    name:            val('name'),
    alias:           val('alias'),
    price:           toFloat('price'),
    available:       toInt('available'),
    short_description: val('short_description'),
    description:       val('description'),
    meta_title:        val('meta_title'),
    meta_description:  val('meta_description'),
    meta_keywords:     val('meta_keywords'),
  }));

  // Все свойства теперь одной строкой — "characteristics"
  formData.append('characteristics', val('characteristics'));

  // Изображения
  const images = document.getElementById('product_images').files;
  for (let i = 0; i < images.length; i++) {
    formData.append('images[]', images[i]);
  }

  try {
    const resp = await fetch(`${API_BASE}/products`, {
      method: 'POST',
      body: formData, // Для FormData заголовок Content-Type не задаём вручную
    });

    if (!resp.ok) {
      const txt = await resp.text();
      result.style.display = 'block';
      result.textContent = `Ошибка при добавлении товара: ${resp.status} ${txt}`;
      return;
    }

    const data = await resp.json();
    result.style.display = 'block';
    result.textContent = data && data.id
      ? `Товар успешно добавлен. ID: ${data.id}`
      : 'Товар добавлен, но сервер не вернул ID.';
  } catch (e) {
    result.style.display = 'block';
    result.textContent = `Ошибка сети: ${e.message}`;
  }
}
