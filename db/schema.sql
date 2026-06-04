-- gochat DB schema (demo 用)
-- 載入: psql "$DB_URL" -f db/schema.sql

CREATE TABLE avatar (
    id       SERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL DEFAULT ''
);

-- 資料表名為保留字, 需加雙引號。
CREATE TABLE "user" (
    id              BIGSERIAL    PRIMARY KEY,
    username        VARCHAR(255) NOT NULL DEFAULT '' ,
    avatar_id       INTEGER      NOT NULL DEFAULT 0,
    is_suspended    BOOLEAN      NOT NULL DEFAULT false,
    code            VARCHAR(32)  NOT NULL,   -- 建立 user 時必帶 (genUserCode), 故不設 default
    create_time     BIGINT       NOT NULL DEFAULT 0,
    update_time     BIGINT       NOT NULL DEFAULT 0,
    last_login_time BIGINT       NOT NULL DEFAULT 0,
    nickname        VARCHAR(20)  NOT NULL DEFAULT '',
    CONSTRAINT uk_user_code UNIQUE (code),
    CONSTRAINT uk_user_username UNIQUE (username)
);

CREATE TABLE chat_group (
    id               SERIAL       PRIMARY KEY,
    is_dismiss       BOOLEAN      NOT NULL DEFAULT false,
    sort             INTEGER      NOT NULL DEFAULT 0,
    visible          BOOLEAN      NOT NULL DEFAULT true,
    code             VARCHAR(20)  NOT NULL,   -- 建立 group 時必帶 (req.Code / genCode), 故不設 default
    user_limit       INTEGER      NOT NULL DEFAULT 0,
    bulletin         VARCHAR(500) NOT NULL DEFAULT '',
    allow_url        BOOLEAN      NOT NULL DEFAULT true,
    pin_chat_id      BIGINT       NOT NULL DEFAULT 0,
    speak_user_level INTEGER      NOT NULL DEFAULT 0,
    join_user_level  INTEGER      NOT NULL DEFAULT 0,
    owner_username   VARCHAR(32)  NOT NULL DEFAULT '',
    owner_user_id    INTEGER      NOT NULL DEFAULT 0,
    title            VARCHAR(32)  NOT NULL DEFAULT '',
    remark           VARCHAR(50)  NOT NULL DEFAULT '',
    create_time      BIGINT       NOT NULL DEFAULT 0,
    update_time      BIGINT       NOT NULL DEFAULT 0,
    cover_filename   VARCHAR(300) NOT NULL DEFAULT '',
    open_join        BOOLEAN      NOT NULL DEFAULT false,
    CONSTRAINT uk_chat_group_code UNIQUE (code)
);

CREATE TABLE chat_record (
    id          BIGSERIAL    PRIMARY KEY,
    group_id    INTEGER      NOT NULL DEFAULT 0,
    user_id     INTEGER      NOT NULL DEFAULT 0,
    type        SMALLINT     NOT NULL DEFAULT 0,
    create_time BIGINT       NOT NULL DEFAULT 0,
    update_time BIGINT       NOT NULL DEFAULT 0,
    content     VARCHAR(500) NOT NULL DEFAULT '',
    deleted     BOOLEAN      NOT NULL DEFAULT false
);

CREATE INDEX idx_chat_record_type_user_group ON chat_record (type, user_id, group_id);
-- 取群組最新發言用 (符合 GetRecentChats 的 WHERE deleted=false ORDER BY id DESC)。
CREATE INDEX idx_chat_record_group_recent ON chat_record (group_id ASC, id DESC) WHERE deleted = false;

CREATE TABLE group_user (
    id                   BIGSERIAL PRIMARY KEY,
    user_id              BIGINT    NOT NULL DEFAULT 0,
    group_id             INTEGER   NOT NULL DEFAULT 0,
    is_ban               BOOLEAN   NOT NULL DEFAULT false,
    join_time            BIGINT    NOT NULL DEFAULT 0,
    role_type            SMALLINT  NOT NULL DEFAULT 0,
    last_visible_chat_id BIGINT    NOT NULL DEFAULT 0,
    last_read_chat_id    BIGINT    NOT NULL DEFAULT 0,
    update_time          BIGINT    NOT NULL DEFAULT 0,
    deleted              BOOLEAN   NOT NULL DEFAULT false,
    CONSTRAINT uk_group_user_user_group UNIQUE (user_id, group_id)
);

CREATE INDEX idx_group_user_group_role ON group_user (group_id, role_type);
CREATE INDEX idx_group_user_lastread ON group_user (last_read_chat_id, user_id, group_id);
