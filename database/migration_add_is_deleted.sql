-- Migration script to add is_deleted column to pekerjaan_alumni table
-- Run this script to add the is_deleted column for soft delete functionality

-- Add is_delete column to pekerjaan_alumni table
ALTER TABLE pekerjaan_alumni 
ADD COLUMN is_delete TIMESTAMP NULL;

-- Add comment to explain the column purpose
COMMENT ON COLUMN pekerjaan_alumni.is_delete IS 'Timestamp when the record was soft deleted. NULL means not deleted.';

-- Optional: Create an index for better performance on queries filtering by is_delete
CREATE INDEX idx_pekerjaan_alumni_is_delete ON pekerjaan_alumni(is_delete);

-- Verify the column was added
SELECT column_name, data_type, is_nullable 
FROM information_schema.columns 
WHERE table_name = 'pekerjaan_alumni' 
AND column_name = 'is_delete';
