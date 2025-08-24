alter table delivery
add order_id int;

alter table delivery
add constraint fk_delivery_order foreign key (order_id) references "order" (id) on delete cascade;