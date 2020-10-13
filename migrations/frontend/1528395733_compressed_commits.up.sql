BEGIN;

ALTER TABLE lsif_nearest_uploads ALTER COLUMN "commit" SET DATA TYPE bytea USING decode("commit", 'hex');

COMMIT;
