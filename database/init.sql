create table users (
	id bigserial primary key,
	telegram_id bigint not null,
    telegram_nickname text notnull,
    unique(telegram_id),
    unique(telegram_nickname)
);

create table rooms (
    id bigserial primary key,
    capacity int not null,
    name text,
);

create table events (
    id uuid primary key,
    start_date timestamp not null,
    end_date timestamp not null,
    notification_period interval,
    suite_id bigint,
    owner_id bigint,
    constraint fk_rooms
        foreign key(suite_id) 
            references rooms(id) 
            on delete cascade
            on update cascade,
    constraint fk_users
        foreign key(owner_id) 
            references users(id)
            on delete cascade
            on update cascade
);
