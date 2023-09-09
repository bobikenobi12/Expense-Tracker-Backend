ALTER TABLE "workspace_invitations"
ADD COLUMN "added_by" bigInt NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "workspace_invitations"
ADD COLUMN "expires" timestamp with time zone NOT NULL;