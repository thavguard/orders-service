alter table payment
add order_id int;

alter table payment
add constraint fk_payment_order foreign key (order_id) references "order" (id) on delete cascade;