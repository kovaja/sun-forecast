CREATE TABLE IF NOT EXISTS remaining_calls (
    date VARCHAR(255) PRIMARY KEY,
    remaining INT
);

CREATE TABLE if NOT EXISTS forecasts (
    id SERIAL PRIMARY KEY,
    period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    value NUMERIC NOT NULL,
    actual NUMERIC
);

CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
    message TEXT NOT NULL
);

-- Migrations
ALTER TABLE forecasts ADD COLUMN actual_count INT DEFAULT 0;
ALTER TABLE forecasts ADD COLUMN last_actual_at TIMESTAMP WITH TIME ZONE DEFAULT NULL;

ALTER TABLE events ADD COLUMN type INT;