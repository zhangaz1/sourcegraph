BEGIN;

ALTER TABLE lsif_nearest_uploads ALTER COLUMN "commit" SET DATA TYPE text USING encode("commit", 'hex');

COMMIT;
