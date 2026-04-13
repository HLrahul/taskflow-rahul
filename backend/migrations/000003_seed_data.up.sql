-- Seed 1 user (password is "password123" hashed with bcrypt cost 12)
INSERT INTO users (id, name, email, password) VALUES 
('11111111-1111-1111-1111-111111111111', 'Test User', 'test@example.com', '$2a$12$/CqVudPvwsZOybZR2SAzleF598gKsWeadqEkXG5haCbtiuyObsUxy')
ON CONFLICT (id) DO NOTHING;

INSERT INTO projects (id, name, description, owner_id) VALUES 
('22222222-2222-2222-2222-222222222222', 'Ship TaskFlow', 'Complete the engineering take-home assignment', '11111111-1111-1111-1111-111111111111')
ON CONFLICT (id) DO NOTHING;

INSERT INTO tasks (id, title, status, priority, project_id, assignee_id) VALUES 
('33333333-3333-3333-3333-333333333331', 'Setup Database', 'done', 'high', '22222222-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111'),
('33333333-3333-3333-3333-333333333332', 'Implement Auth APIs', 'in_progress', 'medium', '22222222-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111'),
('33333333-3333-3333-3333-333333333333', 'Build Frontend Dashboard', 'todo', 'low', '22222222-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111')
ON CONFLICT (id) DO NOTHING;
