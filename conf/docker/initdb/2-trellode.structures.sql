USE trellode;

CREATE TABLE logs (
    id CHAR(36) DEFAULT UUID() PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    board_id CHAR(36) NOT NULL,
    action VARCHAR(32) NOT NULL,
    action_target_id CHAR(36) NOT NULL,
    changes TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id CHAR(36) DEFAULT UUID() PRIMARY KEY,
    email VARCHAR(100) NOT NULL UNIQUE,
    firstname VARCHAR(100) NOT NULL,
    lastname VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE backgrounds (
    id CHAR(36) DEFAULT UUID() PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    data MEDIUMTEXT NOT NULL,
    color varchar(7) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Boards table
CREATE TABLE boards (
    id CHAR(36) DEFAULT UUID() PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    background_id CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    archived_at TIMESTAMP NULL,
    opened_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    -- FOREIGN KEY (user_id) REFERENCES users(id),
    -- FOREIGN KEY (background_id) REFERENCES backgrounds(id)
);

-- Lists table
CREATE TABLE lists (
    id CHAR(36) DEFAULT UUID() PRIMARY KEY,
    board_id CHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    position INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    archived_at TIMESTAMP NULL
    -- FOREIGN KEY (board_id) REFERENCES boards(id)
);

-- Cards table
CREATE TABLE cards (
    id CHAR(36) DEFAULT UUID() PRIMARY KEY,
    list_id CHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    position INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    archived_at TIMESTAMP NULL
    -- FOREIGN KEY (list_id) REFERENCES lists(id)
);

-- Comments table
CREATE TABLE comments (
    id CHAR(36) DEFAULT UUID() PRIMARY KEY,
    card_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE checklists (
    id CHAR(36) DEFAULT UUID() PRIMARY KEY,
    card_id CHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    archived_at TIMESTAMP NULL
);

-- Cards table
CREATE TABLE checklistitems (
    id CHAR(36) DEFAULT UUID() PRIMARY KEY,
    checklist_id CHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    position INT NOT NULL,
    checked TINYINT(1) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);