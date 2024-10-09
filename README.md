# Mock backend for lk project

## Примеры с curl:

### Создание сессии
```sh
curl --header "Content-Type: application/json" --cookie-jar cookie.txt -b cookie.txt --request POST --data '{"iin": "821019000888", "pin":"82"}' http://localhost:8080/session
```
**Ответ**:
```
{"key":"c9aba8d6-351a-4d85-a8b6-9427ea2f8c8e","iin":"821019000888","individ":null}
```

### Who Am I
```sh
curl -b cookie.txt --request GET http://localhost:8080/wai
```
**Ответ**:
```
{"key":"c9aba8d6-351a-4d85-a8b6-9427ea2f8c8e","iin":"821019000888","individ":null}
```

### Получение данных о пользователе

```
curl -b cookie.txt --request GET http://localhost:8080/i/users/c9aba8d6-351a-4d85-a8b6-9427ea2f8c8e
```

**Ответ**:
```json
{
    "key": "c9aba8d6-351a-4d85-a8b6-9427ea2f8c8e",
    "individKey": "27f74b66-cba7-486d-a263-81b6cb9a3e57",
    "iin": "821019000888",
    "login": "821019000888",
    "name": "Усенбаев Дархан Жаксылыкович"
}
```

### Получение деталей
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

### Получение списка типов рапортов

```sh
curl -b cookie.txt http://localhost:8080/i/reports/types
```

**Ответ**:
```json
[
    {
        "ref": "fcf8e381-ea56-43ea-a83f-c2059a3aa329",
        "parent": "",
        "code": "0001",
        "title": "Об убытии в служебные командировки"
    }
]
```

### Получение данных по одному типу

```sh
curl -b cookie.txt http://localhost:8080/i/reports/types/fcf8e381-ea56-43ea-a83f-c2059a3aa329
```

**Ответ**:

```json
{
    "ref": "fcf8e381-ea56-43ea-a83f-c2059a3aa329",
    "parent": "",
    "code": "0001",
    "title": "Об убытии в служебные командировки"
}
```

### Создание рапорта
```sh
curl -b cookie.txt --data @mockData.json http://localhost:8080/i/reports/0001/save
```

**Содержимое `mockData.json`**:

```json
{
  "head": {
    "type": "fcf8e381-ea56-43ea-a83f-c2059a3aa329",
    "date": "2024.10.06 18:06:00"
  },
  "coordinators": [
    {"coordinator_ref": "8c272f7c-6c2c-4dba-bba5-4062005b2400"},
    {"coordinator_ref": "f31c6a0f-b07c-4632-8949-2f24fde4fc26"}
  ],
  "details": {
    "supervisor": "f31c6a0f-b07c-4632-8949-2f24fde4fc26",
    "basis": "Распоряжение руководства",
    "transport_type": "Железная дорога"
  }
}
```

### Получение рапортов **всех** типов 

```sh
curl -b cookie.txt http://localhost:8080/i/reports/
```

**Ответ**:
```json
[
    {
        "ref": "8d624514-7eda-4b6f-a350-5cce76b58870",
        "type": "fcf8e381-ea56-43ea-a83f-c2059a3aa329",
        "date": "2024.10.06 23:06:00",
        "number": "",
        "reg_number": "",
        "author": "c9aba8d6-351a-4d85-a8b6-9427ea2f8c8e"
    },
    {
        "ref": "2b843c2d-c3e9-4ec7-abff-499fe08aae59",
        "type": "fcf8e381-ea56-43ea-a83f-c2059a3aa329",
        "date": "2024.10.06 23:06:00",
        "number": "",
        "reg_number": "",
        "author": "c9aba8d6-351a-4d85-a8b6-9427ea2f8c8e"
    },
    {
        "ref": "cafb6d75-8a45-4d78-b116-26cfb74f1205",
        "type": "fcf8e381-ea56-43ea-a83f-c2059a3aa329",
        "date": "2024.10.06 23:06:00",
        "number": "",
        "reg_number": "",
        "author": "c9aba8d6-351a-4d85-a8b6-9427ea2f8c8e"
    },
    {
        "ref": "64faa2af-d3ec-4536-91a8-adab0cfa2966",
        "type": "fcf8e381-ea56-43ea-a83f-c2059a3aa329",
        "date": "2024.10.06 23:06:00",
        "number": "",
        "reg_number": "",
        "author": "c9aba8d6-351a-4d85-a8b6-9427ea2f8c8e"
    }
]
```
