create table
    delivery (
        id serial primary key,
        name varchar(255) not null,
        phone varchar(15) not null,
        zip varchar(9) not null,
        city varchar(255) not null,
        address varchar(255) not null,
        region varchar(255) not null,
        email varchar(255) not null
    );

create table
    payment (
        id serial primary key,
        transaction varchar(255) not null,
        request_id varchar(255) not null default '',
        currency varchar(3) not null,
        provider varchar(255) not null,
        amount integer not null,
        payment_dt integer not null,
        bank varchar(255) not null,
        delivery_cost integer not null,
        goods_total integer not null,
        custom_fee integer not null default 0
    );

create table
    "order" (
        id serial primary key,
        order_uid varchar(255) not null,
        track_number varchar(255) not null,
        entry varchar(255) not null,
        locale varchar(2) not null,
        internal_signature varchar(255) not null default '',
        customer_id varchar(255) not null,
        delivery_service varchar(255) not null,
        shardkey varchar(255) not null,
        sm_id integer not null,
        date_created timestamp not null,
        oof_shard varchar(255) not null,
        delivery_id integer references delivery (id) on delete cascade,
        payment_id integer references payment (id) on delete cascade
    );

create table
    item (
        id serial primary key,
        chrt_id integer not null,
        track_number varchar(255) not null,
        price integer not null,
        rid varchar(255) not null,
        name varchar(255) not null,
        sale integer not null,
        size varchar(50) not null,
        total_price integer not null,
        nm_id integer not null,
        brand varchar(255) not null,
        status integer not null,
        order_id integer references "order" (id) on delete cascade
    );