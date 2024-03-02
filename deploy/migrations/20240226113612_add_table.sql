-- +goose Up
create table users (
	id bigserial primary key,
	telegram_id bigint not null,
    telegram_nickname text not null,
    unique(telegram_id),
    unique(telegram_nickname)
);

create table rooms (
    id bigserial primary key,
    capacity int not null,
    name text
);

create table events (
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

create index ix_uuid ON events using btree (id);
create index ix_start ON events using brin (start_date);

create index ix_end ON events using brin (end_date);

create index ix_suite ON events using btree (suite_id);
create index ix_owner ON events using btree (user_id);

insert into users (telegram_id, telegram_nickname) values(1234, 'nikitads');
insert into users (telegram_id, telegram_nickname) values(4321, 'kitkeni');
insert into users (telegram_id, telegram_nickname) values(1488, 'noteverlife');
insert into rooms (capacity, name) values(3, 'Winston Churchill');
insert into rooms (capacity, name) values(2, 'Napoleon');
insert into rooms (capacity, name) values(5, 'Putin');

-- +goose Down
drop table events;
drop table users;
drop table rooms;
