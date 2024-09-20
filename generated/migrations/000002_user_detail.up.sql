CREATE TABLE IF NOT EXISTS user_detail (
	detail_id BIGINT PRIMARY KEY,
	nativelingo VARCHAR(3),
    device_os VARCHAR,
    device_brand VARCHAR,
    device_model VARCHAR,
    device_version VARCHAR,
	created_at BIGINT NOT NULL,
	updated_at BIGINT NOT NULL
);