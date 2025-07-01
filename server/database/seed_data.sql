-- Comprehensive seed data for visit tracker
-- 8 days of schedules (today + 7 days ahead), 5 schedules per day, 7 activities per schedule

-- Clear existing data
DELETE FROM activities;
DELETE FROM visits;
DELETE FROM tasks;
DELETE FROM schedules;

-- Insert schedules for today + 7 days ahead (5 schedules per day)
-- Day 0 (Today)
INSERT INTO schedules (client_name, shift_start, shift_end, latitude, longitude, status) VALUES
('Margaret Thompson', datetime('now', 'start of day', '+8 hours'), datetime('now', 'start of day', '+10 hours'), 40.7128, -74.0060, 'upcoming'),
('Robert Chen', datetime('now', 'start of day', '+9 hours'), datetime('now', 'start of day', '+11 hours'), 40.7589, -73.9851, 'upcoming'),
('Eleanor Rodriguez', datetime('now', 'start of day', '+13 hours'), datetime('now', 'start of day', '+15 hours'), 40.6892, -74.0445, 'upcoming'),
('James Mitchell', datetime('now', 'start of day', '+14 hours'), datetime('now', 'start of day', '+16 hours'), 40.7831, -73.9712, 'upcoming'),
('Dorothy Williams', datetime('now', 'start of day', '+16 hours'), datetime('now', 'start of day', '+18 hours'), 40.7505, -73.9934, 'upcoming');

-- Day 1 (Tomorrow)
INSERT INTO schedules (client_name, shift_start, shift_end, latitude, longitude, status) VALUES
('Benjamin Foster', datetime('now', '+1 day', 'start of day', '+7 hours'), datetime('now', '+1 day', 'start of day', '+9 hours'), 40.7282, -74.0776, 'upcoming'),
('Catherine Davis', datetime('now', '+1 day', 'start of day', '+9 hours'), datetime('now', '+1 day', 'start of day', '+11 hours'), 40.7614, -73.9776, 'upcoming'),
('Harold Johnson', datetime('now', '+1 day', 'start of day', '+12 hours'), datetime('now', '+1 day', 'start of day', '+14 hours'), 40.6782, -73.9442, 'upcoming'),
('Patricia Garcia', datetime('now', '+1 day', 'start of day', '+14 hours'), datetime('now', '+1 day', 'start of day', '+16 hours'), 40.7489, -73.9680, 'upcoming'),
('Frank Anderson', datetime('now', '+1 day', 'start of day', '+17 hours'), datetime('now', '+1 day', 'start of day', '+19 hours'), 40.7505, -73.9934, 'upcoming');

-- Day 2
INSERT INTO schedules (client_name, shift_start, shift_end, latitude, longitude, status) VALUES
('Alice Murphy', datetime('now', '+2 days', 'start of day', '+8 hours'), datetime('now', '+2 days', 'start of day', '+10 hours'), 40.7328, -74.0076, 'upcoming'),
('George Brown', datetime('now', '+2 days', 'start of day', '+10 hours'), datetime('now', '+2 days', 'start of day', '+12 hours'), 40.7549, -73.9840, 'upcoming'),
('Mary O''Connor', datetime('now', '+2 days', 'start of day', '+13 hours'), datetime('now', '+2 days', 'start of day', '+15 hours'), 40.6912, -74.0402, 'upcoming'),
('Thomas Wilson', datetime('now', '+2 days', 'start of day', '+15 hours'), datetime('now', '+2 days', 'start of day', '+17 hours'), 40.7780, -73.9665, 'upcoming'),
('Helen Martinez', datetime('now', '+2 days', 'start of day', '+17 hours'), datetime('now', '+2 days', 'start of day', '+19 hours'), 40.7448, -73.9876, 'upcoming');

-- Day 3
INSERT INTO schedules (client_name, shift_start, shift_end, latitude, longitude, status) VALUES
('Charles Taylor', datetime('now', '+3 days', 'start of day', '+7 hours'), datetime('now', '+3 days', 'start of day', '+9 hours'), 40.7178, -74.0431, 'upcoming'),
('Ruth Jackson', datetime('now', '+3 days', 'start of day', '+9 hours'), datetime('now', '+3 days', 'start of day', '+11 hours'), 40.7690, -73.9653, 'upcoming'),
('William Lee', datetime('now', '+3 days', 'start of day', '+12 hours'), datetime('now', '+3 days', 'start of day', '+14 hours'), 40.6823, -73.9654, 'upcoming'),
('Betty White', datetime('now', '+3 days', 'start of day', '+14 hours'), datetime('now', '+3 days', 'start of day', '+16 hours'), 40.7720, -73.9570, 'upcoming'),
('Arthur Harris', datetime('now', '+3 days', 'start of day', '+16 hours'), datetime('now', '+3 days', 'start of day', '+18 hours'), 40.7377, -74.0059, 'upcoming');

-- Day 4
INSERT INTO schedules (client_name, shift_start, shift_end, latitude, longitude, status) VALUES
('Joan Thompson', datetime('now', '+4 days', 'start of day', '+8 hours'), datetime('now', '+4 days', 'start of day', '+10 hours'), 40.7255, -74.0134, 'upcoming'),
('Edward Clark', datetime('now', '+4 days', 'start of day', '+10 hours'), datetime('now', '+4 days', 'start of day', '+12 hours'), 40.7505, -73.9934, 'upcoming'),
('Evelyn Lewis', datetime('now', '+4 days', 'start of day', '+13 hours'), datetime('now', '+4 days', 'start of day', '+15 hours'), 40.6734, -73.9389, 'upcoming'),
('Joseph Robinson', datetime('now', '+4 days', 'start of day', '+15 hours'), datetime('now', '+4 days', 'start of day', '+17 hours'), 40.7831, -73.9665, 'upcoming'),
('Mildred Walker', datetime('now', '+4 days', 'start of day', '+17 hours'), datetime('now', '+4 days', 'start of day', '+19 hours'), 40.7282, -73.9776, 'upcoming');

-- Day 5
INSERT INTO schedules (client_name, shift_start, shift_end, latitude, longitude, status) VALUES
('Raymond Hall', datetime('now', '+5 days', 'start of day', '+7 hours'), datetime('now', '+5 days', 'start of day', '+9 hours'), 40.7420, -74.0124, 'upcoming'),
('Frances Allen', datetime('now', '+5 days', 'start of day', '+9 hours'), datetime('now', '+5 days', 'start of day', '+11 hours'), 40.7648, -73.9776, 'upcoming'),
('Louis Young', datetime('now', '+5 days', 'start of day', '+12 hours'), datetime('now', '+5 days', 'start of day', '+14 hours'), 40.6901, -73.9567, 'upcoming'),
('Marie King', datetime('now', '+5 days', 'start of day', '+14 hours'), datetime('now', '+5 days', 'start of day', '+16 hours'), 40.7505, -73.9934, 'upcoming'),
('Kenneth Wright', datetime('now', '+5 days', 'start of day', '+16 hours'), datetime('now', '+5 days', 'start of day', '+18 hours'), 40.7178, -74.0431, 'upcoming');

-- Day 6
INSERT INTO schedules (client_name, shift_start, shift_end, latitude, longitude, status) VALUES
('Florence Lopez', datetime('now', '+6 days', 'start of day', '+8 hours'), datetime('now', '+6 days', 'start of day', '+10 hours'), 40.7589, -73.9851, 'upcoming'),
('Albert Hill', datetime('now', '+6 days', 'start of day', '+10 hours'), datetime('now', '+6 days', 'start of day', '+12 hours'), 40.7282, -74.0776, 'upcoming'),
('Gladys Scott', datetime('now', '+6 days', 'start of day', '+13 hours'), datetime('now', '+6 days', 'start of day', '+15 hours'), 40.6823, -73.9654, 'upcoming'),
('Victor Green', datetime('now', '+6 days', 'start of day', '+15 hours'), datetime('now', '+6 days', 'start of day', '+17 hours'), 40.7720, -73.9570, 'upcoming'),
('Lillian Adams', datetime('now', '+6 days', 'start of day', '+17 hours'), datetime('now', '+6 days', 'start of day', '+19 hours'), 40.7377, -74.0059, 'upcoming');

-- Day 7
INSERT INTO schedules (client_name, shift_start, shift_end, latitude, longitude, status) VALUES
('Ralph Baker', datetime('now', '+7 days', 'start of day', '+7 hours'), datetime('now', '+7 days', 'start of day', '+9 hours'), 40.7255, -74.0134, 'upcoming'),
('Rose Gonzalez', datetime('now', '+7 days', 'start of day', '+9 hours'), datetime('now', '+7 days', 'start of day', '+11 hours'), 40.7505, -73.9934, 'upcoming'),
('Carl Nelson', datetime('now', '+7 days', 'start of day', '+12 hours'), datetime('now', '+7 days', 'start of day', '+14 hours'), 40.6734, -73.9389, 'upcoming'),
('Irene Carter', datetime('now', '+7 days', 'start of day', '+14 hours'), datetime('now', '+7 days', 'start of day', '+16 hours'), 40.7831, -73.9665, 'upcoming'),
('Eugene Mitchell', datetime('now', '+7 days', 'start of day', '+16 hours'), datetime('now', '+7 days', 'start of day', '+18 hours'), 40.7282, -73.9776, 'upcoming');

-- Insert visits for each schedule
INSERT INTO visits (schedule_id)
SELECT id FROM schedules ORDER BY id;

-- Insert tasks for each schedule (5 tasks per schedule)
INSERT INTO tasks (schedule_id, description, status) 
SELECT s.id, task_desc, 'pending'
FROM schedules s
CROSS JOIN (
    SELECT 'Assist with morning medication' as task_desc
    UNION SELECT 'Help with personal hygiene'
    UNION SELECT 'Prepare nutritious meal'
    UNION SELECT 'Check vital signs'
    UNION SELECT 'Light housekeeping'
) tasks;

-- Insert activities for each schedule (7 activities per schedule)
INSERT INTO activities (schedule_id, title, description, is_resolved, reason)
SELECT s.id, activity_title, activity_desc, 
       CASE WHEN random() % 4 = 0 THEN 1 ELSE 0 END,
       CASE WHEN random() % 4 = 0 THEN '' ELSE 'Pending completion' END
FROM schedules s
CROSS JOIN (
    SELECT 'Medication Administration' as activity_title, 'Administered prescribed medications according to schedule' as activity_desc
    UNION SELECT 'Physical Therapy Support', 'Assisted client with prescribed physical therapy exercises'
    UNION SELECT 'Meal Preparation', 'Prepared healthy meal according to dietary restrictions'
    UNION SELECT 'Personal Care Assistance', 'Helped with bathing, grooming, and dressing'
    UNION SELECT 'Mobility Support', 'Assisted with walking and movement around the home'
    UNION SELECT 'Health Monitoring', 'Checked blood pressure, pulse, and general wellness'
    UNION SELECT 'Social Engagement', 'Engaged in conversation and recreational activities'
) activities; 