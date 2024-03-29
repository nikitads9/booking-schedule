-- +goose Up
create table users (
	id bigserial primary key,
    name text not null,
	telegram_id bigint not null,
    telegram_nickname text not null,
    password text not null,
    created_at timestamp not null,
    updated_at timestamp,
    unique(telegram_id),
    unique(telegram_nickname)
);

create table rooms (
    id bigserial primary key,
    capacity int not null,
    name text
);

create table bookings (
    id uuid primary key,
    start_date timestamp not null,
    end_date timestamp not null,
    notify_at interval default '0s',
    created_at timestamp not null,
    updated_at timestamp,
    suite_id bigint not null,
    user_id bigint not null,
    constraint fk_rooms
        foreign key(suite_id) 
            references rooms(id) 
            on delete cascade
            on update cascade,
    constraint fk_users
        foreign key(user_id) 
            references users(id)
            on delete cascade
            on update cascade
);

create index ix_uuid ON bookings using btree (id);
create index ix_start ON bookings using brin (start_date);

create index ix_end ON bookings using brin (end_date);

create index ix_suite ON bookings using btree (suite_id);
create index ix_owner ON bookings using btree (user_id);

create user otelcol with password 'otelcolpassword';
grant SELECT on pg_stat_database to otelcol;

insert into rooms (capacity, name) values(3, 'Winston Churchill');
insert into rooms (capacity, name) values(2, 'Napoleon');
insert into rooms (capacity, name) values(5, 'Putin');

-- +goose Down
drop table bookings;
drop table users;
drop table rooms;
