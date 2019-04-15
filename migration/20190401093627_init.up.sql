create table site
(
    id          smallserial  not null
        constraint site_pk primary key,
    domain_name varchar(128) not null
);

create unique index site_domain_name_uindex
    on site (domain_name);

create unique index site_id_uindex
    on site (id);

create table cookies
(
    student_id char(8)       not null,
    site_id    smallint      not null
        constraint cookies_site_id_fk
            references site
            on update cascade on delete cascade,
    cookie     varchar(2048) not null,
    constraint cookies_pk
        primary key (student_id, site_id)
);



