USE trellode;

INSERT INTO users (email, password_hash) VALUES 
('user@example.com', 'hashed_password');

INSERT INTO boards (user_id, title, background_image) VALUES 
(1, 'Project Alpha', 'alpha.jpg'),
(1, 'Project Beta', 'beta.jpg'),
(1, 'Project Gamma', 'gamma.jpg');

INSERT INTO lists (board_id, title, position) VALUES 
(1, 'To Do', 1),
(1, 'In Progress', 2),
(1, 'Done', 3),
(2, 'Backlog', 1),
(2, 'Sprint', 2),
(3, 'Ideas', 1);

INSERT INTO cards (list_id, title, description, position) VALUES 
(1, 'Task 1', 'Description for task 1', 1),
(1, 'Task 2', 'Description for task 2', 2),
(2, 'Task 3', 'Description for task 3', 1),
(3, 'Task 4', 'Description for task 4', 1);

INSERT INTO comments (card_id, user_id, content) VALUES 
(1, 1, 'This is a comment on Task 1'),
(2, 1, 'This is a comment on Task 2');
