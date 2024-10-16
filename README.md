# Mock backend for lk project

## Примеры с curl:

### Создание сессии
```sh
curl --header "Content-Type: application/json" --cookie-jar cookie.txt -b cookie.txt --request POST --data '{"iin": "821019000888", "pin":"82"}' http://localhost:8080/session
```
**Ответ**:
```
{"key":"c9aba8d6-351a-4d85-a8b6-9427ea2f8c8e","iin":"821019000888","individ":"27f74b66-cba7-486d-a263-81b6cb9a3e57"}
```

### Who Am I
```sh
curl -b cookie.txt --request GET http://localhost:8080/wai
```
**Ответ**:
```
{"key":"c9aba8d6-351a-4d85-a8b6-9427ea2f8c8e","iin":"821019000888","individ":"27f74b66-cba7-486d-a263-81b6cb9a3e57"}
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

### Получение КАТО

**Получение деталей одной записи**:
```sh
curl -b cookie.txt --request GET "http://localhost:8080/i/cat/cato/9f4436e0-80c0-11ef-8da5-000c29e050f8"
```
**Ответ**:
```json
{
  "data": {
    "ref": "9f4436e0-80c0-11ef-8da5-000c29e050f8",
    "parentRef": "",
    "code": "610000000",
    "k1": "61",
    "k2": "00",
    "k3": "00",
    "k4": "000",
    "k5": "0",
    "description": "Туркестанская область",
    "title": "Туркестанская область"
  },
  "len": 1,
  "rows": 16272,
  "meta": {
    "code": "string",
    "description": "string",
    "k1": "string",
    "k2": "string",
    "k3": "string",
    "k4": "string",
    "k5": "string",
    "parentRef": "catalog_cato",
    "ref": "catalog_cato",
    "title": "string"
  }
}
```

**Получение списка**:
```sh
curl -b cookie.txt --request GET http://localhost:8080/i/cat/cato
```
#### Параметры

 - `limit`: Максимальное количество получаемых записей - по умолчанию: 20
 - `offset`: Количество пропускаемых записей - по умолчанию: 0
 - `search`: Строка поиска
 - `parent`: Идентификатор родителя. Не имеет смысла при установке `search`

**Пример**

```sh
curl -b cookie.txt --request GET "http://localhost:8080/i/cat/cato?parent=9f4436e0-80c0-11ef-8da5-000c29e050f8&limit=3"
```
**Ответ**:
```json
{
  "data": [
    {
      "ref": "9f4436e1-80c0-11ef-8da5-000c29e050f8",
      "parentRef": "9f4436e0-80c0-11ef-8da5-000c29e050f8",
      "code": "611000000",
      "k1": "61",
      "k2": "10",
      "k3": "00",
      "k4": "000",
      "k5": "1",
      "description": "Туркестан Г.А.",
      "title": "Туркестанская область, Туркестан Г.А."
    },
    {
      "ref": "9f4436e3-80c0-11ef-8da5-000c29e050f8",
      "parentRef": "9f4436e0-80c0-11ef-8da5-000c29e050f8",
      "code": "611600000",
      "k1": "61",
      "k2": "16",
      "k3": "00",
      "k4": "000",
      "k5": "3",
      "description": "Арысь Г.А.",
      "title": "Туркестанская область, Арысь Г.А."
    },
    {
      "ref": "ab756119-80c0-11ef-8da5-000c29e050f8",
      "parentRef": "9f4436e0-80c0-11ef-8da5-000c29e050f8",
      "code": "612000000",
      "k1": "61",
      "k2": "20",
      "k3": "00",
      "k4": "000",
      "k5": "3",
      "description": "Кентау Г.А.",
      "title": "Туркестанская область, Кентау Г.А."
    }
  ],
  "len": 3,
  "rows": 16272,
  "meta": {
    "code": "string",
    "description": "string",
    "k1": "string",
    "k2": "string",
    "k3": "string",
    "k4": "string",
    "k5": "string",
    "parentRef": "catalog_cato",
    "ref": "catalog_cato",
    "title": "string"
  }
}
```

### Получеие ВУС

**Получение деталей одной записи**
```sh
curl -b cookie.txt --request GET "http://localhost:8080/i/cat/vus/5e982366-826f-4b16-804d-178dac0b4ff9"
```

**Ответ**:
```json
{
  "data": {
    "ref": "5e982366-826f-4b16-804d-178dac0b4ff9",
    "code": "7654321",
    "title": "Применение подразделений автоматизированных средств управления зенитными ракетными комплексами и зенитной артиллерией войсковой противовоздушной обороны"
  },
  "len": 1,
  "rows": 1,
  "meta": {
    "code": "string",
    "description": "string",
    "k1": "string",
    "k2": "string",
    "k3": "string",
    "k4": "string",
    "k5": "string",
    "parentRef": "catalog_cato",
    "ref": "catalog_cato",
    "title": "string"
  }
}
```


```sh
curl -b cookie.txt --request GET "http://localhost:8080/i/cat/vus"
```
#### Параметры

 - `limit`: Максимальное количество получаемых записей - по умолчанию: 20
 - `offset`: Количество пропускаемых записей - по умолчанию: 0
 - `search`: Строка поиска

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
    "date": "2024.10.06 18:06:00",
    "number": "0001"
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

### Получение данных одного рапорта

```sh
curl -b cookie.txt --request GET http://localhost:8080/i/reports/77ff0fd3-22bc-4cd3-be79-f778f4b1845d
```

**Ответ**:
```json
{
    "data": {
        "head": {
            "ref": "7ee27642-170c-45e8-9fdb-d4855653c63e",
            "type": "fcf8e381-ea56-43ea-a83f-c2059a3aa329",
            "date": "2024.10.06 23:06:00",
            "number": "0001",
            "reg_number": "",
            "author": "c9aba8d6-351a-4d85-a8b6-9427ea2f8c8e"
        },
        "coordinators": [
            {
                "ref": "90228de8-2292-4920-8d4d-7f48db4e4b71",
                "report_ref": "7ee27642-170c-45e8-9fdb-d4855653c63e",
                "coordinator_ref": "8c272f7c-6c2c-4dba-bba5-4062005b2400",
                "who_author_ref": "c9aba8d6-351a-4d85-a8b6-9427ea2f8c8e",
                "when_added_ref": "2024.10.16 17:30:58"
            },
            {
                "ref": "758a58c2-0446-4cbd-8cd6-93f3f0f096e7",
                "report_ref": "7ee27642-170c-45e8-9fdb-d4855653c63e",
                "coordinator_ref": "f31c6a0f-b07c-4632-8949-2f24fde4fc26",
                "who_author_ref": "c9aba8d6-351a-4d85-a8b6-9427ea2f8c8e",
                "when_added_ref": "2024.10.16 17:30:58"
            }
        ],
        "details": {
            "report_ref": "7ee27642-170c-45e8-9fdb-d4855653c63e",
            "supervisor": "f31c6a0f-b07c-4632-8949-2f24fde4fc26",
            "acting_supervisor": "",
            "basis": "Распоряжение руководства",
            "transport_type": "Железная дорога"
        }
    },
    "len": 1,
    "rows": -1,
    "meta": {
        "author": "users",
        "date": "date",
        "number": "string",
        "ref": "report",
        "reg_number": "string",
        "type": "report_type"
    }
}
```

### Получение рапортов **всех** типов 

```sh
curl -b cookie.txt http://localhost:8080/i/reports/
```

**Ответ**:
```json
{
    "data": [
        {
            "ref": "7ee27642-170c-45e8-9fdb-d4855653c63e",
            "type": "fcf8e381-ea56-43ea-a83f-c2059a3aa329",
            "date": "2024.10.06 23:06:00",
            "number": "0001",
            "reg_number": "",
            "author": "c9aba8d6-351a-4d85-a8b6-9427ea2f8c8e"
        }
    ],
    "len": 1,
    "rows": -1,
    "meta": null
}
```
