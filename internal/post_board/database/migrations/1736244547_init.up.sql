BEGIN;

CREATE TABLE "user"
(
    "id" UUID NOT NULL PRIMARY KEY,
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "name" VARCHAR(127) NOT NULL,
    "encrypted_pass" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE "post"
(
    "id" UUID NOT NULL PRIMARY KEY,
    "author_id" UUID NOT NULL REFERENCES "user"("id"),
    "content" TEXT NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

COMMIT;