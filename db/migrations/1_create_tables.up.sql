CREATE TABLE Families (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE Groups (
                        id SERIAL PRIMARY KEY,
                        family_id INT REFERENCES Families(id) ON DELETE CASCADE,
                        name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE Images (
                        id SERIAL PRIMARY KEY,
                        group_id INT REFERENCES Groups(id) ON DELETE CASCADE,
                        name VARCHAR(255) NOT NULL,
                        file_path TEXT NOT NULL,
                        usage_count INT DEFAULT 0,
                        meta_tags TEXT[]
);

CREATE UNIQUE INDEX idx_images_name_group ON Images(name, group_id);
