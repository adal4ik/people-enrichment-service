-- Enable extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Table
CREATE TABLE IF NOT EXISTS persons (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    patronymic TEXT,
    age INT,
    gender TEXT,
    nationality TEXT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_persons_name ON persons (name);
CREATE INDEX IF NOT EXISTS idx_persons_surname ON persons (surname);
CREATE INDEX IF NOT EXISTS idx_persons_patronymic ON persons (patronymic);
CREATE INDEX IF NOT EXISTS idx_persons_nationality ON persons (nationality);
CREATE INDEX IF NOT EXISTS idx_persons_created_at ON persons (created_at);
