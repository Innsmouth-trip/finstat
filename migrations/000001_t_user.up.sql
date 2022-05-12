BEGIN;


CREATE TABLE IF NOT EXISTS "t_user"
(
    "user_uid"                 VARCHAR(50)        NOT NULL UNBOUNDED,
    "firsname"                 VARCHAR(50)        NOT NULL DEFAULT '',
    "lastname"                 VARCHAR(50)        NOT NULL DEFAULT '',
    "username"                 VARCHAR(50)        NOT NULL DEFAULT '',
    "balance"                  INTEGER            NOT NULL DEFAULT 0,

    PRIMARY KEY("user_uid")
);

CREATE TABLE IF NOT EXISTS "t_users_transaction"
(
    "user_uid"          VARCHAR(50) NOT NULL REFERENCES "t_user"(user_uid) ON DELETE CASCADE,
    "transaction_id"    VARCHAR(50)          NOT NULL,
    "amount"            INTEGER              NOT NULL,
    "from_user"         VARCHAR(50)          NOT NULL DEFAULT 'GOD',
    "to_user"           VARCHAR(50)          NOT NULL,
    "transaction_date"  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    PRIMARY KEY("transaction_id")

    );

COMMIT;