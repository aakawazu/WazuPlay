create table users (
    username                    text,
    pass                        text,
    access_token                text,
    access_token_expiration     timestamp,
    refresh_token               text,
    refresh_token_expiration    timestamp,
    profile_picture_url         text
)

create table albums (
    id                  text,
    album_owner         text,
    artist              text,
    title               text,
    album_picture_url   text
)

create table images (
    id          text,
    image_owner text
)

create table audiofiles (
    id          text,
    audio_owner text,
    album       text,
    title       text,
    artist      text,
    instrument  boolean
)