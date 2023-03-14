-- -------------------------------------------------------------
-- TablePlus 5.3.4(492)
--
-- https://tableplus.com/
--
-- Database: projects
-- Generation Time: 2023-03-14 02:01:33.5840
-- -------------------------------------------------------------


-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS projects_id_seq;

-- Table Definition
CREATE TABLE "public"."projects" (
                                     "id" int4 NOT NULL DEFAULT nextval('projects_id_seq'::regclass),
                                     "user_id" int4 NOT NULL,
                                     "name" varchar,
                                     "token" varchar,
                                     "access_time" timestamp,
                                     "created_at" timestamp,
                                     "updated_at" timestamp,
                                     PRIMARY KEY ("id")
);

