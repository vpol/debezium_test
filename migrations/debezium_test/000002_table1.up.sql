create table table1
(
    id uuid not null
        constraint table1_pk
            primary key,
    data text
);

alter table table1
    owner to postgres;

