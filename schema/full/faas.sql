create table dockers
(
    id         varchar(64) not null
        primary key,
    host       varchar(200) null,
    user_name  varchar(64) null,
    name_space varchar(64) null,
    secret     text null,
    name       varchar(64) null,
    created_at bigint null,
    updated_at bigint null,
    deleted_at bigint null,
    created_by varchar(64) null,
    updated_by varchar(64) null,
    deleted_by varchar(64) null,
    tenant_id  varchar(64) null
);

create table functions
(
    id           varchar(64) not null
        primary key,
    group_id     varchar(200) null,
    project_id   varchar(200) null,
    version      varchar(200) null,
    describe     text null,
    status       varchar(200) null,
    env          text null,
    created_at   bigint null,
    updated_at   bigint null,
    deleted_at   bigint null,
    created_by   varchar(64) null,
    updated_by   varchar(64) null,
    deleted_by   varchar(64) null,
    tenant_id    varchar(64) null,
    resource_ref varchar(200) null,
    name         varchar(200) null,
    constraint functions_name_uindex
        unique (name)
);

create table gits
(
    id          varchar(64) not null
        primary key,
    host        varchar(200) null,
    token       text null,
    name        varchar(200) null,
    created_at  bigint null,
    updated_at  bigint null,
    deleted_at  bigint null,
    created_by  varchar(64) null,
    updated_by  varchar(64) null,
    deleted_by  varchar(64) null,
    tenant_id   varchar(64) null,
    ssh         text null,
    known_hosts text null
);



