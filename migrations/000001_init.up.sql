CREATE TABLE "expense_types" (
    "id" integer PRIMARY KEY,
    "name" varchar,
    "created_at" timestamp,
    "updated_at" timestamp
);
CREATE TABLE "expenses" (
    "id" integer PRIMARY KEY,
    "amount" integer,
    "note" varchar,
    "date" timestamp,
    "expense_type_id" integer,
    "workspaceId" integer
);
CREATE TABLE "users" (
    "id" integer PRIMARY KEY,
    "name" varchar,
    "email" varchar,
    "country_code" varchar,
    "profile_pic_id" integer,
    "user_secrets_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp
);
CREATE TABLE "user_secrets" (
    "id" integer PRIMARY KEY,
    "password" integer,
    "updated_at" timestamp
);
CREATE TABLE "roles" ("id" integer PRIMARY KEY, "name" varchar);
CREATE TABLE "permissions" (
    "id" integer PRIMARY KEY,
    "permission" varchar
);
CREATE TABLE "permission_role" (
    "id" integer PRIMARY KEY,
    "userId" integer,
    "permissionId" integer,
    "roleId" integer
);
CREATE TABLE "currency" (
    "id" integer PRIMARY KEY,
    "currency" varchar
);
CREATE TABLE "currency_workspace" (
    "id" integer PRIMARY KEY,
    "workspaceId" integer,
    "currencyId" integer
);
CREATE TABLE "user_currency" (
    "id" integer PRIMARY KEY,
    "currencyId" integer,
    "userId" integer
);
CREATE TABLE "workspaces" (
    "id" integer PRIMARY KEY,
    "name" varchar,
    "created_at" timestamp,
    "updated_at" timestamp,
    "ownerId" integer
);
CREATE TABLE "user_workspace" (
    "id" integer PRIMARY KEY,
    "userId" integer,
    "workspaceId" integer
);
CREATE TABLE "s3objects" (
    "id" integer PRIMARY KEY,
    "e_tag" varchar,
    "key" varchar,
    "version_id" varchar,
    "location" varchar
);
ALTER TABLE "users"
ADD FOREIGN KEY ("profile_pic_id") REFERENCES "s3objects" ("id");
ALTER TABLE "users"
ADD FOREIGN KEY ("user_secrets_id") REFERENCES "user_secrets" ("id");
ALTER TABLE "permission_role"
ADD FOREIGN KEY ("userId") REFERENCES "user_workspace" ("id");
ALTER TABLE "user_currency"
ADD FOREIGN KEY ("userId") REFERENCES "users" ("id");
ALTER TABLE "user_currency"
ADD FOREIGN KEY ("currencyId") REFERENCES "currency" ("id");
ALTER TABLE "currency_workspace"
ADD FOREIGN KEY ("workspaceId") REFERENCES "workspaces" ("id");
ALTER TABLE "currency_workspace"
ADD FOREIGN KEY ("currencyId") REFERENCES "currency" ("id");
ALTER TABLE "permission_role"
ADD FOREIGN KEY ("roleId") REFERENCES "roles" ("id");
ALTER TABLE "permission_role"
ADD FOREIGN KEY ("permissionId") REFERENCES "permissions" ("id");
ALTER TABLE "expenses"
ADD FOREIGN KEY ("workspaceId") REFERENCES "workspaces" ("id");
ALTER TABLE "user_workspace"
ADD FOREIGN KEY ("userId") REFERENCES "users" ("id");
ALTER TABLE "user_workspace"
ADD FOREIGN KEY ("workspaceId") REFERENCES "workspaces" ("id");
ALTER TABLE "workspaces"
ADD FOREIGN KEY ("ownerId") REFERENCES "users" ("id");
ALTER TABLE "expenses"
ADD FOREIGN KEY ("expense_type_id") REFERENCES "expense_types" ("id");