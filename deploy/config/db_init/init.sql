CREATE DATABASE IF NOT EXISTS toktik_user;
CREATE DATABASE IF NOT EXISTS toktik_interaction;
CREATE DATABASE IF NOT EXISTS toktik_video;
CREATE DATABASE IF NOT EXISTS toktik_chat;
CREATE DATABASE IF NOT EXISTS toktik_favor;
CREATE DATABASE IF NOT EXISTS toktik_comment;


USE toktik_user;

create table if not exists user
(
    id               bigint unsigned auto_increment
        primary key,
    created_at       datetime(3)  null,
    updated_at       datetime(3)  null,
    deleted_at       datetime(3)  null,
    username         varchar(40)  not null,
    password         longtext     not null,
    avatar           varchar(255) not null,
    background_image varchar(255) not null,
    is_follow        tinyint(1)   not null,
    signature        varchar(255) null,
    constraint username
        unique (username)
);

create table if not exists user_count
(
    id              bigint unsigned auto_increment
        primary key,
    created_at      datetime(3)     null,
    updated_at      datetime(3)     null,
    deleted_at      datetime(3)     null,
    user_id         bigint unsigned null,
    follow_count    bigint          null,
    follower_count  bigint          null,
    total_favorited bigint          null,
    work_count      bigint          null,
    favorite_count  bigint          null,
    constraint fk_user_count_user
        foreign key (user_id) references user (id)
);


USE toktik_interaction;

-- auto-generated definition
create table if not exists relation
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)     null,
    updated_at datetime(3)     null,
    deleted_at datetime(3)     null,
    user_id    bigint unsigned null,
    target_id  bigint unsigned null
);

create index idx_relation
    on relation (user_id, target_id);



USE toktik_video;

create table if not exists video
(
    id             bigint unsigned auto_increment
        primary key,
    created_at     datetime(3)     null,
    updated_at     datetime(3)     null,
    deleted_at     datetime(3)     null,
    user_id        bigint unsigned null,
    title          varchar(255)    not null,
    play_url       varchar(255)    not null,
    cover_url      varchar(255)    not null,
    favorite_count bigint          null,
    comment_count  bigint          null
);

create index user_id
    on video (user_id);


USE toktik_chat;

create table if not exists messages
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)     null,
    updated_at datetime(3)     null,
    deleted_at datetime(3)     null,
    user_id    bigint unsigned not null,
    to_user_id bigint unsigned not null,
    content    longtext        not null
);

create index idx_message
    on messages (user_id, to_user_id);



USE toktik_favor;

create table if not exists favorite
(
    id       bigint unsigned auto_increment
        primary key,
    user_id  bigint unsigned null,
    video_id bigint unsigned null,
    constraint user_video_id
        unique (user_id, video_id)
);


USE toktik_comment;

create table if not exists comment
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)     null,
    updated_at datetime(3)     null,
    deleted_at datetime(3)     null,
    video_id   bigint unsigned not null,
    user_id    bigint unsigned not null,
    content    longtext        not null
);

create index comment_video
    on comment (video_id);
