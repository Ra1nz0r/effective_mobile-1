CREATE TABLE "library" (
    "id" serial PRIMARY KEY,
    "group" varchar NOT NULL DEFAULT '',
    "song" varchar NOT NULL DEFAULT '',
    "releaseDate" date NOT NULL DEFAULT CURRENT_DATE,
    "text" text NOT NULL DEFAULT '',
    "link" varchar NOT NULL DEFAULT ''
);
CREATE INDEX ON "library" ("group");
CREATE INDEX ON "library" ("releaseDate");