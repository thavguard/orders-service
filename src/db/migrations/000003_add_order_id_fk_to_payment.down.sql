alter table payment
drop constraint fk_payment_order;

alter table payment
drop column order_id;