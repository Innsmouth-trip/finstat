# Тестовое задание для Finstat
Реализовать веб-сервер с двумя обработчиками.

## Возможности API:
* Можно пополнить баланс указанного пользователя на указанную сумму.
* Можно переводить, указанную сумму со счета первого пользователя на счет другого. 
### Дополнительные требования:
* В проекте нужна миграция для создания таблицы.
* Тестирование (не реализовано)

### Дополнительные реализованные возможности:
* Добавление пользователя
* Получение пользователя по UserUid
* Получение истории транзакций пользователя
* Документация

### Клонирование и запуск:
* Склонируйте репозиторий,
*  <code>cd finstat </code>,
*  <code>make run </code>.

После запуска команды <code>make run</code>:
* Сгенерируется Docker Image,
* Отработает docker-compose.

**BaseUrl: http://localhost:8080/v1/**

**Документация: http://localhost:1349/**

