BEGIN;

-- password: "secret123"
INSERT INTO "user" ("id", "email", "name", "encrypted_pass")
  VALUES ('8ff6fe28-14c5-4dc8-a0bf-749fa8a212a0', 'admin@google.com', 'admin', 'fcf730b6d95236ecd3c9fc2d92d7b6b2bb061514961aec041d6c7a7192f592e4');

COMMIT;