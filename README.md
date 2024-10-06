# Mock backend for lk project

Пример файла хранения списка пользователей
```json
[
    {
        "name": "Саркелова Арайлым",
        "iin": "1122",
        "pin": "2211",
        "birth_date": "2001-11-22"
    },
    {
        "name": "Арсенов Арман",
        "iin": "22",
        "pin": "11",
        "birth_date": "1993-11-21"
    }
]
```
Пример с curl:

**Создание сессии**:
```sh
curl --header "Content-Type: application/json" --cookie-jar cookie.txt -b cookie.txt --request POST --data '{"iin": "821019000888", "pin":"82"}' http://localhost:8080/session
```
Ответ:
```
{"key":"c9aba8d6-351a-4d85-a8b6-9427ea2f8c8e","iin":"821019000888","individ":null}
```

**Who Am I**:
```sh
curl -b cookie.txt --request GET http://localhost:8080/wai
```
Ответ:
```
{"key":"c9aba8d6-351a-4d85-a8b6-9427ea2f8c8e","iin":"821019000888","individ":null}
```

**Получение деталей**:
```sh
curl -b cookie.txt --request GET http://localhost:8080/i/ind/821019000888 | python3 -c 'import sys, json; print(json.dumps(json.load(sys.stdin), ensure_ascii=False, indent=4))'
```
**Ответ**:
```json
{
    "key": "27f74b66-cba7-486d-a263-81b6cb9a3e57",
    "code": "000000015",
    "first_name": "Дархан",
    "last_name": "Усенбаев",
    "patronymic": "Жаксылыкович",
    "nationality": "Казах",
    "iin": "821019000888",
    "image": "821019000888",
    "birth_date": "1981.11.19 00:00:00",
    "birth_place": "с. Баканас Балхашского района Алма-Атинской области",
    "personal_number": "А-000001",
    "last_rank": {
        "date": "2023.10.23 00:00:00",
        "individual": {
            "key": "27f74b66-cba7-486d-a263-81b6cb9a3e57",
            "code": "000000015",
            "first_name": "Дархан",
            "last_name": "Усенбаев",
            "patronymic": "Жаксылыкович",
            "nationality": "Казах",
            "iin": "821019000888",
            "image": "821019000888",
            "birth_date": "1981.11.19 00:00:00",
            "birth_place": "с. Баканас Балхашского района Алма-Атинской области",
            "personal_number": "А-000001"
        },
        "rank": {
            "key": "86bf503e-9327-46d4-8d6c-35dd19b88cfa",
            "name": "Полковник"
        }
    },
    "rank_history": [
        {
            "date": "2023.10.23 00:00:00",
            "individual": {
                "key": "27f74b66-cba7-486d-a263-81b6cb9a3e57",
                "code": "000000015",
                "first_name": "Дархан",
                "last_name": "Усенбаев",
                "patronymic": "Жаксылыкович",
                "nationality": "Казах",
                "iin": "821019000888",
                "image": "821019000888",
                "birth_date": "1981.11.19 00:00:00",
                "birth_place": "с. Баканас Балхашского района Алма-Атинской области",
                "personal_number": "А-000001"
            },
            "rank": {
                "key": "86bf503e-9327-46d4-8d6c-35dd19b88cfa",
                "name": "Полковник"
            }
        }
    ]
}

```