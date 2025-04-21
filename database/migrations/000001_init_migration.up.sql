CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE user_role AS ENUM ('employee', 'moderator');

CREATE TYPE city AS ENUM ('Москва', 'Санкт-Петербург', 'Казань');

CREATE TYPE prod_type AS ENUM ('электроника', 'одежда', 'обувь');

CREATE TYPE work_status AS ENUM ('in_progress', 'close');

CREATE TABLE IF NOT EXISTS "user" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "password" VARCHAR(255) NOT NULL,
    "role" user_role NOT NULL
);

CREATE TABLE IF NOT EXISTS pvz (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "registrationDate" TIMESTAMP WITH TIME ZONE DEFAULT NOW() ,
    "city" city NOT NULL
);

CREATE TABLE IF NOT EXISTS reception (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "dateTime" TIMESTAMP WITH TIME ZONE DEFAULT NOW() ,
    "pvzId" uuid NOT NULL,
    "status" work_status NOT NULL,
    CONSTRAINT "fk_pvz" FOREIGN KEY ("pvzId")
    REFERENCES pvz(id)
);

CREATE TABLE IF NOT EXISTS product (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "dateTime" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "type" prod_type NOT NULL,
    "receptionId" uuid NOT NULL,
    CONSTRAINT "fk_reception" FOREIGN KEY ("receptionId")
    REFERENCES reception(id)
);
