// Базы адресов — как в product-init.js
const SLIDER_BACKEND_ORIGIN = (window.AUTH && AUTH.ORIGIN) ||
  ((location.port === '8080')
    ? `${location.protocol}//${location.host}`
    : `${location.protocol}//${location.hostname}:8080}`);

const SLIDER_IMG_BASE = (window.AUTH && AUTH.IMG_BASE) || `${SLIDER_BACKEND_ORIGIN}/images_db`;
const SLIDER_API_BASE = `${SLIDER_BACKEND_ORIGIN}/api`;

let sliderProducts = [];
let sliderIndex = 0;
let sliderTimer = null;   // таймер для авто-слайдшоу

document.addEventListener('DOMContentLoaded', initSlider);

async function initSlider() {
  try {
    const resp = await fetch(`${SLIDER_API_BASE}/products`);
    if (!resp.ok) {
      console.error('Не удалось получить продукты для слайдера');
      return;
    }

    const products = await resp.json();

    // Берём только товары с картинками, ограничим, например, первыми 10
    sliderProducts = products
      .filter(p => p.product_images && p.product_images.image)
      .slice(0, 10);

    if (!sliderProducts.length) {
      console.warn('Нет товаров с картинками для слайдера');
      return;
    }

    renderSlider();
    startSliderAuto();   // запустить авто-перелистывание
  } catch (e) {
    console.error('Ошибка загрузки продуктов для слайдера', e);
  }
}

function startSliderAuto() {
  // чтобы не было нескольких таймеров
  if (sliderTimer) {
    clearInterval(sliderTimer);
  }
  // каждые 5 секунд листаем вперёд
  sliderTimer = setInterval(() => {
    sliderMove(1);
  }, 5000);
}

function renderSlider() {
  if (!sliderProducts.length) return;

  const p = sliderProducts[sliderIndex];

  const imgEl   = document.getElementById('sliderImage');
  const titleEl = document.getElementById('sliderTitle');
  const priceEl = document.getElementById('sliderPrice');

  if (!imgEl || !titleEl || !priceEl) return;

  const rawImg = (p.product_images && p.product_images.image)
    ? String(p.product_images.image)
    : '';
  const file = rawImg.replace(/^\/?images_db\//, '').trim();

  imgEl.src = file ? `${SLIDER_IMG_BASE}/${file}` : `${SLIDER_IMG_BASE}/placeholder.png`;
  imgEl.alt = p.name || 'Товар';

  // Клик по картинке — переход на страницу товара
  imgEl.onclick = function () {
    goToProduct(p.id);
  };

  titleEl.textContent = p.name || 'Товар';
  priceEl.textContent = (p.price != null) ? `${p.price} Рублей` : '';

  // 👉 перезапускаем CSS-анимацию (сдвиг)
  imgEl.classList.remove('slider-animate');
  // «хак» чтобы браузер пересчитал стиль и заново проиграл анимацию
  void imgEl.offsetWidth;
  imgEl.classList.add('slider-animate');
}

// Листаем вперёд/назад по кнопкам и при авто-слайдшоу
function sliderMove(delta) {
  if (!sliderProducts.length) return;

  const n = sliderProducts.length;
  sliderIndex = (sliderIndex + delta + n) % n;

  renderSlider();
  startSliderAuto(); // при ручном клике перезапускаем таймер, чтобы не дёргалось
}

function goToProduct(id) {
  const url = `products/product_page.html?id=${encodeURIComponent(id)}`;
  window.location.href = url;
}
