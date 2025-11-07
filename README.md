# Web-lab-3
## База данных

Проект использует **PostgreSQL**. В репозитории есть дамп базы `db/dump.sql`, чтобы можно было быстро развернуть такую же БД локально.

### 1. Требования

- Установленный **PostgreSQL** (желательно 13+).
- Утилиты `psql` и `createdb` доступны в терминале  
  (на Windows — либо добавить `C:\Program Files\PostgreSQL\<версия>\bin` в `PATH`, либо выполнять команды из этой папки).

### 2. Создание базы данных

По умолчанию используется:

- хост: `localhost`
- порт: `5432`
- пользователь: `postgres`
- имя базы данных: `web_database`
- пароль: 'password'

Создайте пустую БД:

```bash
createdb -h localhost -p 5432 -U postgres web_database

psql -h localhost -p 5432 -U postgres -d web_database

```

## Запуск приложения

### Backend

Backend написан на Go.

Запуск из корня проекта:

```bash
cd backend
go run main.go
```

### Уже созданные учетные записи (можно создать новые)
```
Администратор
Логин: SHILOP0P
Пароль: qwertyuiop
Пользователи
Логин: User1
Пароль: qwertyuiop
Логин: User2
Пароль: qwertyuiop
``` 