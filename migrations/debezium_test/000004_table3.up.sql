create table table3
(
    id uuid not null
        constraint table3_pk
            primary key,
    data text
);

alter table table3
    owner to postgres;

