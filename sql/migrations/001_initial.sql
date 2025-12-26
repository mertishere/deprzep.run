CREATE TABLE "project" (
        id text PRIMARY KEY NOT NULL,
        user_id text NOT NULL,
        name text NOT NULL,
        description text,
        docs text,
        tags text, created_at text NOT NULL, updated_at text,
        FOREIGN KEY (user_id) REFERENCES user(id) ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE TABLE "release" (
        id text PRIMARY KEY NOT NULL UNIQUE,
        user_id text NOT NULL,
        url text NOT NULL,
        zig_version text NOT NULL,
        hash text NOT NULL,
        root_file text NOT NULL,
        created_at text NOT NULL,
        updated_at text,
        project_id text NOT NULL,
        release text NOT NULL,
        FOREIGN KEY (user_id) REFERENCES user(id) ON UPDATE CASCADE ON DELETE CASCADE,
        FOREIGN KEY (project_id) REFERENCES project(id) ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE TABLE session (
        id text PRIMARY KEY NOT NULL,
        user_id text NOT NULL,
        secret text NOT NULL, created_at text NOT NULL,
        FOREIGN KEY (user_id) REFERENCES user(id) ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE TABLE "user" (
        id text PRIMARY KEY NOT NULL UNIQUE,
        username text NOT NULL UNIQUE,
        email text NOT NULL UNIQUE,
        password text NOT NULL,
        is_verified boolean NOT NULL,
        created_at text NOT NULL
);
CREATE TABLE verify (
        id text PRIMARY KEY NOT NULL,
        user_id text NOT NULL,
        code text NOT NULL, created_at text NOT NULL,
        FOREIGN KEY (user_id) REFERENCES user(id) ON UPDATE CASCADE ON DELETE CASCADE
);
