create table table2
(
    id uuid not null
        constraint table2_pk
            primary key,
    data text
);

alter table table2
    owner to postgres;

