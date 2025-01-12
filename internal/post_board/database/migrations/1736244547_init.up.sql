BEGIN;

CREATE TABLE "user"
(
    "id" UUID NOT NULL,
    "email" VARCHAR(255) NOT NULL CONSTRAINT "user_email_key" UNIQUE,
    "name" VARCHAR(127) NOT NULL,
    "encrypted_pass" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY ("id")
);

CREATE TABLE "post"
(
    "id" UUID NOT NULL,
    "author_id" UUID NOT NULL REFERENCES "user"("id") ON DELETE CASCADE,
    "content" TEXT NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY ("id")
);

COMMIT;