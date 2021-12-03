CREATE TABLE goods (
    id         BIGSERIAL PRIMARY KEY,
    goods_id   BIGINT NOT NULL,
    order_id   BIGINT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);