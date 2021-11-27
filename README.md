# Тестовое задание на позицию стажера-бекендера

## Микросервис для работы с балансом пользователей.

**Проблема:**

**Запуск**

На сервере PostgreSQL создать бд и таблицы в этой бд, скрипт по созданию нужных таблиц находится в файле `script.sql`

```
go run cmd/app/main.go -user="username" -password="password" -db="database"
```

`username` - логин на сервере PostgreSQL, 

`password` - пароль этого пользователя, 

`database` - имя базы данных

**Метод начисления средств на баланс:**

Принимает `id` пользователя и сколько средств зачислить.

```curl --header "Content-Type: application/json" \
--request POST \
--data '{"id": 1, "balance": 200}' \
http://localhost:9000/balance/add
```
`Ответ:` сообщение об успехе, либо код ошибки

**Метод списания средств с баланса:**

Принимает `id` пользователя и сколько средств списать.

```curl --header "Content-Type: application/json" \
--request POST \
--data '{"id": 1, "balance": 200}' \
http://localhost:9000/balance/reduce
```
`Ответ:` сообщение об успехе, либо код ошибки

**Метод перевода средств от пользователя к пользователю:**

`id` - id пользователя, с которого нужно списать средства
`id_to` - id пользователя, которому должны зачислить средства
`balance` - сумма.

```curl --header "Content-Type: application/json" \
--request POST \
--data '{"id": 1, "balance": 200, "id_to":3}' \
http://localhost:9000/balance/transfer
```
`Ответ:` сообщение об успехе, либо код ошибки

**Метод получения текущего баланса пользователя:**

Принимает `id` пользователя. Баланс всегда в рублях.

```curl --header "Content-Type: application/json" \
--request POST \
--data '{"id": 1}' \
http://localhost:9000/user
```
`Ответ:` возвращает баланс пользователя в рублях, либо код ошибки

**метод получения списка транзакций:**

Принимает `id` пользователя, а также сортировку по сумме и дате.

```curl --header "Content-Type: application/json" \
--request POST \
--data '{"id": 1}' \
http://localhost:9000/info
```

*сортировка по сумме*
```curl --header "Content-Type: application/json" \
--request POST \
--data '{"id": 1, "field":"date"}' \
http://localhost:9000/info
```

*сортировка по дате*

Ответ: список всех транзакций для пользователя, с полями:

`to_id` - кому зачислены, 

`from_id` - от кого происходило списание денег, 

`money` - сумма, 

`created` - дата транзакции, либо код ошибки




curl --header "Content-Type: application/json" --request POST --data "{\"id\": 1, \"balance\": 200}" http://localhost:9000/balance/add

curl --header "Content-Type: application/json" --request POST --data '{"id": 1}' http://localhost:9000/info