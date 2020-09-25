# URLshortener
HTTP сервис для сокращения URL
На вход принимает JSON с ссылкой, которую необходимо скоратить (POST запрос), возвращает также JSON с сжатой ссылкой.
При переходе по полученной ссылке посетитель перенаправляется по полной ссылке.

# Сокращение ссылок
Для сокращения пользователь должен передать JSON следующего вида:
{
  "address" : "https://github.com/zvdenis/"
}
В ответ пользователь получит такой же JSON, но уже с сокращенной ссылкой.

Для сокращения использован следующий алгоритм:
В базе данных хранится id, short_link, long_link, добавляя новую ссылку long link создается из 
id путем прибавления bias и переводом в систему по основанию 62. bias нужен чтобы ссылка не была слишком короткой

Пример сжатой ссылки:
http://localhost:8000/1t

# Перенаправление
При переходе по сокращенной ссылке, код выделяется с адреса. Затем по коду ищется полная ссылка, по которой и будет направлен пользователь.

# Структура
* LinkController отвечает за логику работы с ссылками (сокращение, добавление в базу и т.д.)
* Handler отвечает за обработку запросов и вызывает соответствующие методы LinkController(сокращение ссылки, получение полной ссылки)
