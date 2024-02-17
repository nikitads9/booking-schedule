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
    start_date timestamp with time zone not null,
    end_date timestamp with time zone not null,
    notify_at timestamp with time zone,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone,
    suite_id bigint not null,
    owner_id bigint not null,
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

create index ix_uuid ON events using btree (id);
create index ix_start ON events using brin (start_date);

create index ix_end ON events using brin (end_date);

create index ix_suite ON events using btree (suite_id);
create index ix_owner ON events using btree (owner_id);
