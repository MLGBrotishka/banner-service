# Запуск
### 1 вариант Docker
Нужен docker и docker-compose

Запустить:

```make docker``` - Соберет контейнер и запустит его

### 2 вариант Локальный запуск

Заполнить в файле ```.env.example``` настройки к postgresql
```
DB_USER=     #Имя пользователя
DB_PASSWORD= #Пароль пользователя
DB_NAME=     #Имя базы данных
DB_PORT=     #Порт
DB_HOST=     #Имя хоста
```
Переименовать в ```.env```

Запустить:

```make env``` - Экспортирует переменные окружения

``` make run ``` - Соберет проект локально и запустит приложение