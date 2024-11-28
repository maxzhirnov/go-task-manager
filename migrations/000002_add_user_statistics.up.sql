-- Create view for user statistics
CREATE OR REPLACE VIEW user_statistics AS
SELECT 
    u.id as user_id,
    u.username,
    COUNT(t.id) as total_tasks,
    COUNT(CASE WHEN t.status = 'completed' THEN 1 END) as completed_tasks,
    COUNT(CASE WHEN t.status = 'pending' THEN 1 END) as pending_tasks,
    COUNT(CASE WHEN t.status = 'in_progress' THEN 1 END) as in_progress_tasks,
    COUNT(CASE WHEN t.status = 'deleted' THEN 1 END) as deleted_tasks,
    COUNT(CASE WHEN DATE(t.created_at) = CURRENT_DATE THEN 1 END) as tasks_created_today
FROM 
    users u
LEFT JOIN 
    tasks t ON u.id = t.user_id
GROUP BY 
    u.id, u.username;