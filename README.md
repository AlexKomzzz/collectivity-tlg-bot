# Telegram Бот

В качестве хранилища используется <a href="https://github.com/boltdb/bolt">Bolt DB</a>.

Чтобы реализовать авторизацию пользователей, бот отправляет ссылку на регистрацию на основной сайт. После успешной регистрации происходит переадресация на бот. Для того, чтобы бот понимал, кто к нему вернулся, используются файлы cookie


### Стек:
- Go 1.18
- BoltDB
- Docker (для развертывания)



Не кликабельны ссылки в сообщениид