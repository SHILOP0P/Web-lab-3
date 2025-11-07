const API = {
  // Функция для получения списка продуктов
  async getProducts() {
    const res = await fetch('http://localhost:8080/api/products', { credentials: 'omit' });
    if (!res.ok) throw new Error('Failed to load products');
    return res.json();
  },

  // Функция для добавления нового продукта
  async addProduct(productData) {
    const res = await fetch('http://localhost:8080/api/products', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(productData), // Отправляем данные товара
    });

    if (!res.ok) throw new Error('Failed to add product');
    return res.json(); // Возвращаем ответ от сервера, в котором может быть ID добавленного товара
  },
  
  // Сюда можно добавить другие методы, например PUT/DELETE
};

const deleteProduct = async (productId) => {
  try {
    const response = await fetch(`http://localhost:8080/api/products/${productId}`, {
      method: 'DELETE',
    });

    if (response.ok) {
      return await response.json();
    } else {
      throw new Error('Не удалось удалить товар');
    }
  } catch (error) {
    console.error('Ошибка при удалении товара:', error);
  }
};
