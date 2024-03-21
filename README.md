# Hotel California Backend

Hotel California; online ortamda müşterilerine hizmet veren bir turizm şirketidir. Bu Repository açık kaynaklıdır ve aşağıda sahip olduğu bazı yönergeler yer almaktadır.

indirme:

`git clone https://github.com/erhankrygt/hotel-california-backend.git`

Projeyi local ortamda ayağa kaldırabilmek için bazı environment değerlerine ihtiyacınız bulunmakadır. Aşağıdaki environment değerlerini projenize ekleyiniz.

- HTTP_SERVER_ADDRESS
- JWT_TOKEN_SECRET
- MYSQL_DATABASE
- MYSQL_PASSWORD
- MYSQL_PORT
- MYSQL_URI
- MYSQL_USER_NAME

bu yapıda **HTTP_SERVER_ADDRESS** aşağıdaki şekilde tanımlanmıştır

`localhost:8001`

/docs içerisinde postman collection u yer almaktadır. ancak aşağıda endpoint lere ait curl değerleri paylaşılmaktadır.

>** SignIn endpoint**
Kullanıcı login olur ve token alır

`curl --location 'localhost:8001/v1/account/sign-in' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username":"erhan@gmail.com",
    "password":"123123"
}'`

> **NewReservation**
SignIn den token alan kullanıcı rezervasyon bilgilerini ve token ı göndererek yeni rezervasyon oluşturur

`curl --location 'localhost:8001/v1/reservation/new' \
--header 'Accept-Language: tr' \
--header 'token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIiwiZXhwIjoxNzExMTEwNTMxfQ.z7O21hp91OGmVhpwUffJ2gkUs6niijJChjDFaYDPGwo' \
--header 'Content-Type: application/json' \
--data '{
    "destination": "Istanbul",
    "checkInDate": "2024-03-29",
    "checkOutDate": "2024-03-30",
    "accommodation": "mountain",
    "guestCount": 1
}'`

> **UpdateReservation**
SignIn den token alan kullanıcı, daha önce oluşturduğu rezervasyon pnr numarası ile rezervasyon bilgilerini göndererek güncelleme gerçekleşştirir.

`curl --location 'localhost:8001/v1/reservation/update' \
--header 'Accept-Language: en' \
--header 'token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIiwiZXhwIjoxNzExMTEwNTMxfQ.z7O21hp91OGmVhpwUffJ2gkUs6niijJChjDFaYDPGwo' \
--header 'Content-Type: application/json' \
--data '{
    "pnr":"ObAlIEJP",
    "destination": "zonguldak",
    "checkInDate": "2024-03-26",
    "checkOutDate": "2024-03-30",
    "accommodation": "mountain",
    "guestCount": 1
}'`

> **FindReservations**
SignIn den token alan kullanıcı  token bilgisiyle pnr numaranısını göndererek ilgili rezervasyonun bilgilerini alır

`curl --location --request GET 'localhost:8001/v1/reservation?pnr=pE5TYsDj' \
--header 'Accept-Language: tr' \
--header 'token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIiwiZXhwIjoxNzExMTEwNTMxfQ.z7O21hp91OGmVhpwUffJ2gkUs6niijJChjDFaYDPGwo' \
--header 'Content-Type: application/json' \
--data '{
    "pnr":"2G90DqRo",
    "destination": "sivas",
    "checkInDate": "2024-03-24",
    "checkOutDate": "2024-03-30",
    "accommodation": "mountain",
    "guestCount": 1
}'`

> **FindReservations**
SignIn den token alan kullanıcı  token bilgisiyle kullanıcıya ait tüm rezervasyonlar döner

`curl --location --request GET 'localhost:8001/v1/reservations' \
--header 'Accept-Language: tr' \
--header 'token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIiwiZXhwIjoxNzExMTEwNTMxfQ.z7O21hp91OGmVhpwUffJ2gkUs6niijJChjDFaYDPGwo' \
--header 'Content-Type: application/json' \
--data '{
    "pnr":"2G90DqRo",
    "destination": "sivas",
    "checkInDate": "2024-03-24",
    "checkOutDate": "2024-03-30",
    "accommodation": "mountain",
    "guestCount": 1
}'`


