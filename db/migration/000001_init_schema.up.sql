CREATE TABLE "library" (
    "id" serial PRIMARY KEY,
    "group" varchar NOT NULL,
    "song" varchar NOT NULL,
    "releaseDate" date,
    "text" text,
    "patronymic" varchar
);
CREATE INDEX ON "library" ("group");
CREATE INDEX ON "library" ("releaseDate");