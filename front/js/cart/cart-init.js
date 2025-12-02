// Базы адресов такие же, как в product-init.js
const BACKEND_ORIGIN = (window.AUTH && AUTH.ORIGIN) ||
  ((location.port === '8080')
    ? `${location.protocol}//${location.host}`
    : `${location.protocol}//${location.hostname}:8080}`);

const IMG_BASE = (window.AUTH && AUTH.IMG_BASE) || `${BACKEND_ORIGIN}/images_db`;
const API_BASE = `${BACKEND_ORIGIN}/api`;


document.addEventListener('DOMContentLoaded', initCart);

async function initCart() {
  const table = document.getElementById('cartTable');
  const totalEl = document.getElementById('cartTotal');
  const emptyMsg = document.getElementById('cartEmptyMsg');
  const clearBtn = document.getElementById('clearCartBtn');
  const checkoutBtn = document.getElementById('checkoutBtn');

  if (!table || !totalEl || !emptyMsg) return;

  // загрузка корзины
  try {
    const resp = await fetch(`${API_BASE}/cart`, {
      credentials: 'include'
    });

    if (resp.status === 401) {
      emptyMsg.textContent = 'Корзина доступна только авторизованным пользователям. Войдите в аккаунт.';
      emptyMsg.style.display = 'block';
      table.style.display = 'none';
      if (clearBtn) clearBtn.style.display = 'none';
      if (checkoutBtn) checkoutBtn.style.display = 'none';
      return;
    }

    if (!resp.ok) {
      throw new Error('bad status');
    }

    const items = await resp.json();
    renderCart(table, totalEl, emptyMsg, items);
  } catch (e) {
    console.error(e);
    emptyMsg.textContent = 'Ошибка загрузки корзины.';
    emptyMsg.style.display = 'block';
  }

  // обработчик "Очистить корзину"
  if (clearBtn) {
    clearBtn.addEventListener('click', async () => {
      if (!confirm('Очистить корзину?')) return;
      try {
        const resp = await fetch(`${API_BASE}/cart`, {
          method: 'DELETE',
          credentials: 'include'
        });
        if (!resp.ok) throw new Error();
        renderCart(table, totalEl, emptyMsg, []);
      } catch (e) {
        console.error(e);
        alert('Не удалось очистить корзину.');
      }
    });
  }

  // заглушка оформления заказа
  if (checkoutBtn) {
    checkoutBtn.addEventListener('click', () => {
      alert('Оформление заказа: заглушка для лабораторной работы.');
    });
  }
}

function renderCart(table, totalEl, emptyMsg, items) {
  // очищаем строки, кроме заголовка
  while (table.rows.length > 1) {
    table.deleteRow(1);
  }

  if (!Array.isArray(items) || items.length === 0) {
    emptyMsg.textContent = 'Ваша корзина пуста.';
    emptyMsg.style.display = 'block';
    totalEl.textContent = '0';
    return;
  }

  emptyMsg.style.display = 'none';

  items.forEach((item) => {
    const tr = table.insertRow(-1);

    const imgTd    = tr.insertCell(-1);
    const nameTd   = tr.insertCell(-1);
    const priceTd  = tr.insertCell(-1);
    const qtyTd    = tr.insertCell(-1);
    const sumTd    = tr.insertCell(-1);
    const actionTd = tr.insertCell(-1);

    // картинка товара
    const img = document.createElement('img');
    img.className = 'cart-product-image';
    if (item.image) {
    img.src = `${IMG_BASE}/${item.image}`;
    } else {
    img.src = `${IMG_BASE}/placeholder.png`;
    }
    img.alt = item.name || 'Товар';
    imgTd.appendChild(img);



    // ссылка на страницу товара
    // ссылка на страницу товара
    const link = document.createElement('a');
    link.href = `products/product_page.html?id=${item.productId}`;
    link.textContent = item.name;
    nameTd.appendChild(link);


    priceTd.textContent = Number(item.price || 0).toFixed(2);

    // поле количества
    const input = document.createElement('input');
    input.type = 'number';
    input.min = '1';
    input.value = item.quantity;
    input.style.width = '60px';
    qtyTd.appendChild(input);

    // кнопка "Удалить"
    const delBtn = document.createElement('button');
    delBtn.textContent = 'Удалить';
    delBtn.className = 'btn btn--danger';
    actionTd.appendChild(delBtn);

    function updateRowSum() {
      const q = parseInt(input.value, 10) || 0;
      const sum = q * Number(item.price || 0);
      sumTd.textContent = sum.toFixed(2);
      recomputeTotal(table, totalEl);
    }

    updateRowSum();

    // изменение количества
    input.addEventListener('change', async () => {
      const q = parseInt(input.value, 10) || 0;
      try {
        const resp = await fetch(`${API_BASE}/cart/${item.productId}`, {
          method: 'PUT',
          credentials: 'include',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ quantity: q })
        });

        if (!resp.ok) throw new Error();

        if (q <= 0) {
          // сервер удалил товар — перезагрузим корзину
          const reload = await fetch(`${API_BASE}/cart`, {
            credentials: 'include'
          });
          const newItems = await reload.json();
          renderCart(table, totalEl, emptyMsg, newItems);
          return;
        }

        updateRowSum();
      } catch (e) {
        console.error(e);
        alert('Не удалось обновить количество.');
      }
    });

    // удаление товара
    delBtn.addEventListener('click', async () => {
      if (!confirm('Удалить товар из корзины?')) return;
      try {
        const resp = await fetch(`${API_BASE}/cart/${item.productId}`, {
          method: 'DELETE',
          credentials: 'include'
        });
        if (!resp.ok) throw new Error();

        table.deleteRow(tr.rowIndex);
        recomputeTotal(table, totalEl);

        if (table.rows.length === 1) {
          emptyMsg.textContent = 'Ваша корзина пуста.';
          emptyMsg.style.display = 'block';
        }
      } catch (e) {
        console.error(e);
        alert('Не удалось удалить товар.');
      }
    });
  });

  recomputeTotal(table, totalEl);
}

function recomputeTotal(table, totalEl) {
  let total = 0;
  for (let i = 1; i < table.rows.length; i++) {
    const sumCell = table.rows[i].cells[4];
    const val = parseFloat(sumCell.textContent.replace(',', '.')) || 0;
    total += val;
  }
  totalEl.textContent = total.toFixed(2);
}
