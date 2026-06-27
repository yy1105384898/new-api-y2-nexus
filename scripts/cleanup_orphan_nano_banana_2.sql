-- Remove orphan nano-banana-2 model that has no channel binding.
DELETE FROM models
WHERE model_name = 'nano-banana-2'
  AND NOT EXISTS (
    SELECT 1 FROM abilities WHERE model = 'nano-banana-2'
  );
