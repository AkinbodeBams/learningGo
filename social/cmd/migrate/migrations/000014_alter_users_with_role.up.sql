-- Step 1: Add the column with a default
ALTER TABLE IF EXISTS users 
ADD COLUMN role_id BIGINT REFERENCES roles(id) DEFAULT 1;

-- Step 2: Update existing users to 'user' role
UPDATE users
SET role_id = (
    SELECT id FROM roles WHERE name = 'user'
);

-- Step 3: Remove default for new rows
ALTER TABLE users
ALTER COLUMN role_id DROP DEFAULT;

-- Step 4: Make column NOT NULL
ALTER TABLE users
ALTER COLUMN role_id SET NOT NULL;
