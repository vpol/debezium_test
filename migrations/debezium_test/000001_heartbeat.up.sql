CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table heartbeat
(
    id integer not null
        constraint heartbeat_pk
            primary key,
    ts timestamp with time zone
);

alter table heartbeat
    owner to postgres;

