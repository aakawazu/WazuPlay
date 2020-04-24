create table users (
    username                    text,
    pass                        text,
    access_token                text,
    access_token_expiration     timestamp,
    refresh_token               text,
    refresh_token_expiration    timestamp,
    profile_picture_url         text
)