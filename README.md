# mercadona-cli

> <https://tienda.mercadona.es>

```text
                               _                       ___   __   _____
  /\/\   ___ _ __ ___ __ _  __| | ___  _ __   __ _    / __\ / /   \_   \
 /    \ / _ \ '__/ __/ _` |/ _` |/ _ \| '_ \ / _` |  / /   / /     / /\/
/ /\/\ \  __/ | | (_| (_| | (_| | (_) | | | | (_| | / /___/ /___/\/ /_
\/    \/\___|_|  \___\__,_|\__,_|\___/|_| |_|\__,_| \____/\____/\____/

```

⚠️ EXPERIMENTAL CLI!! THERE COULD BE BUGS OR UNHANDLED USE CASES!! USE AT YOUR OWN RISK ⚠️

## Features

- [x] Authentication
- [x] Show current account information
- [ ] Make an order
- [x] Get URL to modify latest active order

## Usage

TODO

## Mercadona API reverse engineering

### Set warehouse by postal code

```shell
$ echo '{"new_postal_code": "<postal_code>"}' | http PUT "https://tienda.mercadona.es/api/postal-codes/actions/change-pc/"

HTTP/1.1 200 OK
Allow: PUT, POST
Alt-Svc: clear
Cache-Control: no-cache
Cache-Control: no-cache
Content-Language: es
Content-Length: 27
Content-Type: application/json
Date: Wed, 28 Apr 2021 13:37:02 GMT
Expires: Wed, 28 Apr 2021 13:37:01 GMT
Server: nginx
Strict-Transport-Security: max-age=86400
Vary: Cookie, Origin
Via: 1.1 google
X-Frame-Options: SAMEORIGIN
X-Request-ID: ce8caac5e57e7bfd05410e699a9f9cce
X-SRE-header: location_only_api_nocache
x-customer-pc: <pc>
x-customer-wh: mad1

{
    "warehouse_changed": true
}
```

### Authenticate

```shell
 $ echo '{"username":"<email>","password": "<password>"}' | http POST "https://tienda.mercadona.es/api/auth/tokens/"

HTTP/1.1 200 OK
Allow: POST, OPTIONS
Alt-Svc: clear
Cache-Control: no-cache
Cache-Control: no-cache
Content-Encoding: gzip
Content-Language: es
Content-Type: application/json
Date: Wed, 28 Apr 2021 13:44:42 GMT
Expires: Wed, 28 Apr 2021 13:44:41 GMT
Server: nginx
Strict-Transport-Security: max-age=86400
Transfer-Encoding: chunked
Vary: Accept-Encoding
Vary: Origin
Via: 1.1 google
X-Frame-Options: SAMEORIGIN
X-Request-ID: dc1643fe8d0a13ec12a996853754c229
X-SRE-header: location_only_api_nocache
x-customer-pc: <pc>
x-customer-wh: mad1

{
    "access_token": "<access_token>",
    "customer_id": "<customer_id>"
}
```

### Get customer information

```shell
$ http GET "https://tienda.mercadona.es/api/customers/<customer_id>/" "Authorization:Bearer <auth_token>"

HTTP/1.1 200 OK
Allow: GET, PUT, PATCH, DELETE, HEAD, OPTIONS
Alt-Svc: clear
Cache-Control: no-cache
Cache-Control: no-cache
Content-Language: es
Content-Length: 242
Content-Type: application/json
Date: Wed, 28 Apr 2021 14:12:47 GMT
Expires: Wed, 28 Apr 2021 14:12:46 GMT
Server: nginx
Strict-Transport-Security: max-age=86400
Vary: Origin
Via: 1.1 google
X-Frame-Options: SAMEORIGIN
X-Request-ID: 41afc2e37f33b20810a1abe2d0ffee85
X-SRE-header: location_only_api_nocache
x-customer-pc: <pc>
x-customer-wh: mad1

{
    "cart_id": "<cart_id>",
    "current_postal_code": "<postal_code>",
    "email": "<email>",
    "id": <id>,
    "last_name": "<last_name>",
    "name": "<name>",
    "send_offers": false,
    "uuid": "<uuid>"
}
```

### Get customer cart

```shell
$ http GET "https://tienda.mercadona.es/api/customers/<customer_id>/cart/" "Authorization:Bearer <auth_token>"

HTTP/1.1 200 OK
Allow: GET, PUT, HEAD, OPTIONS
Alt-Svc: clear
Cache-Control: no-cache
Cache-Control: no-cache
Content-Language: es
Content-Length: 138
Content-Type: application/json
Date: Wed, 28 Apr 2021 14:23:15 GMT
Expires: Wed, 28 Apr 2021 14:23:14 GMT
Server: nginx
Strict-Transport-Security: max-age=86400
Vary: Origin
Via: 1.1 google
X-Frame-Options: SAMEORIGIN
X-Request-ID: 6271d00c1fdf045a2f919118714eac6a
X-SRE-header: location_only_api_nocache
x-customer-pc: <pc>
x-customer-wh: mad1

{
    "id": "<cart_id>",
    "lines": [],
    "open_order_id": <open_order_id>,
    "products_count": 0,
    "summary": {
        "total": "0.00"
    },
    "version": 3
}
```

## List all orders

```shell
$ http GET "https://tienda.mercadona.es/api/customers/<customer_id>/orders/?page=<page_num>" "Authorization:Bearer <auth_token>"

HTTP/1.1 200 OK
Allow: GET, HEAD, OPTIONS
Alt-Svc: clear
Cache-Control: no-cache
Cache-Control: no-cache
Content-Encoding: gzip
Content-Language: es
Content-Type: application/json
Date: Wed, 28 Apr 2021 14:42:54 GMT
Expires: Wed, 28 Apr 2021 14:42:53 GMT
Server: nginx
Strict-Transport-Security: max-age=86400
Transfer-Encoding: chunked
Vary: Accept-Encoding
Vary: Origin
Via: 1.1 google
X-Frame-Options: SAMEORIGIN
X-Request-ID: 71e7981780c0b9d8018365676be60750
X-SRE-header: location_only_api_nocache
x-customer-pc: <pc>
x-customer-wh: mad1
{
    "next_page": null,
    "results": [
        {
            "address": {
                "address": "<address>",
                "address_detail": "<address_detail>",
                "comments": "<comments>",
                "entered_manually": false,
                "id": <address_id>,
                "latitude": "<latitude>",
                "longitude": "<longitude>",
                "permanent_address": true,
                "postal_code": "<postal_code>",
                "town": "Madrid"
            },
            "changes_until": "2021-04-29T17:59:59Z",
            "click_and_collect": false,
            "customer_phone": "<customer_phone>",
            "end_date": "2021-04-30T16:00:00Z",
            "final_price": false,
            "id": 8312430,
            "last_edit_message": "Pedido editado hace 16 horas.",
            "order_id": 8312430,
            "payment_method": {
                "credit_card_number": "<last_4_credit_card_digits>",
                "credit_card_type": 1,
                "default_card": true,
                "expiration_status": "valid",
                "expires_month": "<expire_month>",
                "expires_year": "<expire_year>",
                "id": <payment_method_id>
            },
            "payment_status": 0,
            "phone_country_code": "34",
            "phone_national_number": "<phone_national_number>",
            "price": "65.94",
            "products_count": 28,
            "service_rating_token": null,
            "slot": {
                "available": true,
                "end": "2021-04-30T16:00:00Z",
                "id": <slot_id>,
                "price": "7.21",
                "start": "2021-04-30T15:00:00Z"
            },
            "slot_size": 1,
            "start_date": "2021-04-30T15:00:00Z",
            "status": 2,
            "status_ui": "confirmed",
            "summary": {
                "products": "65.94",
                "slot": "7.21",
                "tax_base": "67.07",
                "taxes": "6.08",
                "total": "73.15",
                "volume_extra_cost": {
                    "cost_by_extra_liter": "0.1",
                    "threshold": 70,
                    "total": "0.00",
                    "total_extra_liters": 0.0
                }
            },
            "warehouse_code": "mad1"
        },
        ...
}
```

### Get order info

```shell
$ http GET "https://tienda.mercadona.es/api/customers/<customer_id>/orders/<order_id>/" "Authorization:Bearer <auth_token>"

HTTP/1.1 200 OK
Allow: GET, DELETE, HEAD, OPTIONS
Alt-Svc: clear
Cache-Control: no-cache
Cache-Control: no-cache
Content-Encoding: gzip
Content-Language: es
Content-Type: application/json
Date: Wed, 28 Apr 2021 14:31:30 GMT
Expires: Wed, 28 Apr 2021 14:31:29 GMT
Server: nginx
Strict-Transport-Security: max-age=86400
Transfer-Encoding: chunked
Vary: Accept-Encoding
Vary: Origin
Via: 1.1 google
X-Frame-Options: SAMEORIGIN
X-Request-ID: c24c5000d1bc877b7993cc8a05cbd52d
X-SRE-header: location_only_api_nocache
x-customer-pc: <pc>
x-customer-wh: mad1

{
    "address": {
        "address": "<address>",
        "address_detail": "<address_detail>",
        "comments": "<comments>",
        "entered_manually": false,
        "id": <address_id>,
        "latitude": "<latitude>",
        "longitude": "<longitude>",
        "permanent_address": true,
        "postal_code": "<postal_code>",
        "town": "Madrid"
    },
    "changes_until": "2021-04-29T17:59:59Z",
    "click_and_collect": false,
    "customer_phone": "<customer_phone>",
    "end_date": "2021-04-30T16:00:00Z",
    "final_price": false,
    "id": 8312430,
    "last_edit_message": "Pedido editado hace 16 horas.",
    "order_id": 8312430,
    "payment_method": {
        "credit_card_number": "<last_4_credit_card_digits>",
        "credit_card_type": 1,
        "default_card": true,
        "expiration_status": "valid",
        "expires_month": "<expire_month>",
        "expires_year": "<expire_year>",
        "id": <payment_method_id>
    },
    "payment_status": 0,
    "phone_country_code": "34",
    "phone_national_number": "<phone_national_number>",
    "price": "65.94",
    "products_count": 28,
    "service_rating_token": null,
    "slot": {
        "available": true,
        "end": "2021-04-30T16:00:00Z",
        "id": <slot_id>,
        "price": "7.21",
        "start": "2021-04-30T15:00:00Z"
    },
    "slot_size": 1,
    "start_date": "2021-04-30T15:00:00Z",
    "status": 2,
    "status_ui": "confirmed",
    "summary": {
        "products": "65.94",
        "slot": "7.21",
        "tax_base": "67.07",
        "taxes": "6.08",
        "total": "73.15",
        "volume_extra_cost": {
            "cost_by_extra_liter": "0.1",
            "threshold": 70,
            "total": "0.00",
            "total_extra_liters": 0.0
        }
    },
    "warehouse_code": "mad1"
}
```

### Get my favorites products

#### Parameters

- `<regular_type>` = `precision` (what I most buy) or `recall` (I also buy)

```shell
> http GET "https://tienda.mercadona.es/api/customers/<customer_id>/recommendations/myregulars/<regular_type>/" "Authorization:Bearer <auth_token>"

HTTP/1.1 200 OK
Allow: GET, HEAD, OPTIONS
Alt-Svc: clear
Cache-Control: no-cache
Cache-Control: no-cache
Content-Encoding: gzip
Content-Language: es
Content-Type: application/json
Date: Wed, 28 Apr 2021 18:54:41 GMT
Expires: Wed, 28 Apr 2021 18:54:40 GMT
Server: nginx
Strict-Transport-Security: max-age=86400
Transfer-Encoding: chunked
Vary: Accept-Encoding
Vary: Origin
Via: 1.1 google
X-Frame-Options: SAMEORIGIN
X-Request-ID: 119edafcf45f3d941068a2e29fc7c6a0
X-SRE-header: location_only_api_nocache
x-customer-pc: <pc>
x-customer-wh: mad1

{
    "next_page": null,
    "results": [
        {
            "product": {
                "badges": {
                    "is_water": false,
                    "requires_age_check": false
                },
                "categories": [
                    {
                        "id": 5,
                        "level": 0,
                        "name": "Panadería y pastelería",
                        "order": 508
                    }
                ],
                "display_name": "Empanada de verduras",
                "id": "84780",
                "limit": 999,
                "packaging": "Pieza",
                "price_instructions": {
                    "approx_size": false,
                    "bulk_price": "5.90",
                    "bunch_selector": false,
                    "drained_weight": null,
                    "increment_bunch_amount": 1.0,
                    "is_new": false,
                    "is_pack": false,
                    "iva": 10,
                    "min_bunch_amount": 1.0,
                    "pack_size": null,
                    "price_decreased": false,
                    "reference_format": "kg",
                    "reference_price": "5.90",
                    "selling_method": 0,
                    "size_format": "kg",
                    "total_units": null,
                    "unit_name": null,
                    "unit_price": "2.95",
                    "unit_selector": true,
                    "unit_size": 0.5
                },
                "published": true,
                "share_url": "https://tienda.mercadona.es/product/84780/empanada-verduras-apto-veganos-pieza",
                "slug": "empanada-verduras-apto-veganos-pieza",
                "thumbnail": "https://prod-mercadona.imgix.net/images/a128a58d8a06618683ba674671800dee.jpg?fit=crop&h=300&w=300"
            },
            "recommended_quantity": 1,
            "selling_method": 0
        },
        ...
    ]
}
```
