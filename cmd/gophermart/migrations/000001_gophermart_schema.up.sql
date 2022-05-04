CREATE TABLE IF NOT EXISTS "users" (
  "id" BIGSERIAL PRIMARY KEY,
  "login" varchar UNIQUE,
  "password" varchar NOT NULL,
  "authtoken" varchar,
  "current" numeric(7,2),
  "withdrawn" numeric(7,2)
);

CREATE TABLE IF NOT EXISTS "orders" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" int NOT NULL,
  "number" varchar NOT NULL UNIQUE,
  "status" varchar,
  "accrual" numeric(7,2),
  "uploaded_at" TIMESTAMP DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "withdrawals" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" int NOT NULL,
  "order" varchar NOT NULL,
  "sum" numeric(7,2),
  "processed_at" TIMESTAMP DEFAULT (now())
);

CREATE INDEX IF NOT EXISTS "order_status" ON "orders" ("user_id", "number", "status");

ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "withdrawals" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
