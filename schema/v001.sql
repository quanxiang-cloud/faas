create table docker
(
    id         varchar(64)  not null
        primary key,
    host       varchar(200) null,
    name       varchar(64)  null,
    secret     text         null,
    created_at bigint       null,
    updated_at bigint       null,
    deleted_at bigint       null,
    created_by varchar(64)  null,
    updated_by varchar(64)  null,
    deleted_by varchar(64)  null,
    tenant_id  varchar(64)  null
);

create table git
(
    id         varchar(64)  not null
        primary key,
    host       varchar(200) null,
    token      text         null,
    created_at bigint       null,
    updated_at bigint       null,
    deleted_at bigint       null,
    created_by varchar(64)  null,
    updated_by varchar(64)  null,
    deleted_by varchar(64)  null,
    tenant_id  varchar(64)  null
);

