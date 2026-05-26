-- 000002_add_users_and_org_plan.down.sql
DROP TABLE IF EXISTS users;
-- Note: SQLite does not support ALTER TABLE DROP COLUMN directly in older versions,
-- but since this is development we might just leave it or recreate the table if needed.
-- For simplicity in this environment:
