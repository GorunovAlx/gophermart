CREATE TABLE "users" (
  "id" BIGSERIAL PRIMARY KEY,
  "login" varchar UNIQUE,
  "password" varchar NOT NULL,
  "authtoken" varchar,
  "current" numeric(7,2),
  "withdrawn" numeric(7,2)
);

CREATE TABLE "orders" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" int NOT NULL,
  "number" varchar NOT NULL,
  "status" varchar,
  "accrual" numeric(7,2),
  "uploaded_at" datetime DEFAULT (now())
);

CREATE TABLE "withdrawals" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" int NOT NULL,
  "order" varchar NOT NULL,
  "sum" numeric(7,2),
  "processed_at" datetime DEFAULT (now())
);

CREATE INDEX "order_status" ON "orders" ("user_id", "number", "status");

ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "withdrawals" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
