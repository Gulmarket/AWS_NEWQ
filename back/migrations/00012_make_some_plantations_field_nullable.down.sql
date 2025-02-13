BEGIN;

UPDATE plantations SET description='sdf' WHERE description IS NULL;
UPDATE plantations SET country='sdf' WHERE country IS NULL;
UPDATE plantations SET city='sdf' WHERE city IS NULL;
UPDATE plantations SET logo_url='sdf' WHERE logo_url IS NULL;
UPDATE plantations SET work_schedule='{}' WHERE work_schedule IS NULL;
UPDATE plantations SET name='sdf' WHERE name IS NULL;

-- Step 1: Add the columns as nullable
ALTER TABLE plantations
    ADD COLUMN longitude FLOAT,
    ADD COLUMN latitude FLOAT;

-- Step 2: Populate the columns with default values (if necessary)
UPDATE plantations SET longitude = 0 WHERE longitude IS NULL;
UPDATE plantations SET latitude = 0 WHERE latitude IS NULL;

-- Step 3: Alter the columns to be NOT NULL
ALTER TABLE plantations
    ALTER COLUMN longitude SET NOT NULL,
    ALTER COLUMN latitude SET NOT NULL;

-- Restore the NOT NULL constraints for other columns
ALTER TABLE plantations
    ALTER COLUMN description SET NOT NULL,
    ALTER COLUMN country SET NOT NULL,
    ALTER COLUMN logo_url SET NOT NULL,
    ALTER COLUMN city SET NOT NULL,
    ALTER COLUMN work_schedule SET NOT NULL,
    ALTER COLUMN name SET NOT NULL;

COMMIT;





