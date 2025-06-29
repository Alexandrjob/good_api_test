
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE goods (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    priority INT NOT NULL,
    removed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_projects
        FOREIGN KEY(project_id)
            REFERENCES projects(id)
);

CREATE INDEX idx_goods_project_id ON goods(project_id);
CREATE INDEX idx_goods_name ON goods(name);

INSERT INTO projects (name) VALUES ('Первая запись');
