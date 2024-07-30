CREATE TABLE IF NOT EXISTS "xusers" (
  "id" bigserial PRIMARY KEY,
  "user_name" varchar UNIQUE NOT NULL,
  "password" bytea NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "token" varchar,
  "created_at" timestamptz DEFAULT (now())
);
CREATE TABLE IF NOT EXISTS "xobjects" (
  "id" bigserial PRIMARY KEY,
  "ref" varchar NOT NULL,
  "type" varchar NOT NULL,
  "description" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);
CREATE TABLE IF NOT EXISTS "xownership" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "obj_id" bigint NOT NULL
);
CREATE INDEX ON "xusers" ("user_name");
CREATE INDEX ON "xusers" ("first_name");
CREATE INDEX ON "xusers" ("last_name");
CREATE INDEX ON "xobjects" ("type");
CREATE INDEX ON "xownership" ("user_id");
CREATE INDEX ON "xownership" ("obj_id");
CREATE INDEX ON "xownership" ("user_id", "obj_id");
ALTER TABLE "xownership"
ADD FOREIGN KEY ("user_id") REFERENCES "xusers" ("id");
ALTER TABLE "xownership"
ADD FOREIGN KEY ("obj_id") REFERENCES "xobjects" ("id");