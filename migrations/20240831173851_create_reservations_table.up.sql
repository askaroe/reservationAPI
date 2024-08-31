CREATE TABLE IF NOT EXISTS reservations (
    id bigserial   PRIMARY KEY,
    created_at     timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at     timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    room_id bigint NOT NULL,
    start_date     timestamp(0) NOT NULL,
    end_date       timestamp(0) NOT NULL
);