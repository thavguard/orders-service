alter table delivery
drop constraint fk_delivery_order;

alter table delivery
drop column order_id;