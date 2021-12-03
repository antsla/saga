Создаем заказ
`curl --request POST \
   --header "Content-Type: application/json" \
   --data '{"user_id":1,"goods_ids":[1,2]}' \
   'http://localhost:8080/v1/orders'`