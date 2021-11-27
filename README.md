# Тестовое задание на позицию стажера-бекендера

## Микросервис для работы с балансом пользователей.

**Проблема:**

**Запуск**



```
    docker-compose up
```

**Метод начисления средств на баланс:**

Принимает `id` пользователя и сколько средств зачислить.

```curl --header "Content-Type: application/json" \
--request POST \
--data '{"id": 1, "balance": 200}' \
http://localhost:8000/balance/add
```
`Ответ:` сообщение об успехе, либо код ошибки

**Метод списания средств с баланса:**

Принимает `id` пользователя и сколько средств списать.

```curl --header "Content-Type: application/json" \
--request POST \
--data '{"id": 1, "balance": 200}' \
http://localhost:8000/balance/reduce
```
`Ответ:` сообщение об успехе, либо код ошибки

**Метод перевода средств от пользователя к пользователю:**

`id` - id пользователя, с которого нужно списать средства
`id_to` - id пользователя, которому должны зачислить средства
`balance` - сумма.

```curl --header "Content-Type: application/json" \
--request POST \
--data '{"id": 1, "balance": 200, "id_to":3}' \
http://localhost:8000/balance/transfer
```
`Ответ:` сообщение об успехе, либо код ошибки

**Метод получения текущего баланса пользователя:**

Принимает `id` пользователя. Баланс всегда в рублях.

```curl --header "Content-Type: application/json" \
--request POST \
--data '{"id": 1}' \
http://localhost:8000/user
```

Добавлен к методу получения баланса доп. параметр. Пример: ?currency=USD

`Ответ:` возвращает баланс пользователя в рублях, либо код ошибки

**метод получения списка транзакций:**

Принимает `id` пользователя, а также сортировку по сумме и дате.

```curl --header "Content-Type: application/json" \
--request POST \
--data '{"id": 1}' \
http://localhost:8000/info
```

*сортировка по сумме*
```curl --header "Content-Type: application/json" \
--request POST \
--data '{"id": 1, "field":"date"}' \
http://localhost:8000/info
```

*сортировка по дате*

Ответ: список всех транзакций для пользователя, с полями:

`to_id` - кому зачислены, 

`from_id` - от кого происходило списание денег, 

`money` - сумма, 

`created` - дата транзакции

