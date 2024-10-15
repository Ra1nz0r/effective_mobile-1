CREATE TABLE IF NOT EXISTS "artist" (
    "id" serial PRIMARY KEY,
    "group" varchar NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS "library" (
    "id" serial PRIMARY KEY,
    "group_id" int NOT NULL,
    "song" varchar NOT NULL DEFAULT '',
    "releaseDate" date NOT NULL DEFAULT 'now()',
    "text" text NOT NULL DEFAULT '',
    "link" varchar NOT NULL DEFAULT '',
    CONSTRAINT unique_group_song UNIQUE (group_id, song),
    FOREIGN KEY ("group_id") REFERENCES "artist" ("id")
);
CREATE INDEX ON "library" ("group_id");
CREATE INDEX ON "library" ("releaseDate");