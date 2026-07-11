-- Disable stale abilities whose owning channel is not enabled.
-- Safe to rerun. Production uses PostgreSQL; runtime reads also enforce this invariant.

BEGIN;

UPDATE abilities AS a
SET enabled = FALSE
FROM channels AS c
WHERE c.id = a.channel_id
  AND c.status <> 1
  AND a.enabled = TRUE;

COMMIT;

SELECT c.id AS channel_id, c.name, c.status, count(*) AS stale_enabled_abilities
FROM abilities AS a
JOIN channels AS c ON c.id = a.channel_id
WHERE a.enabled = TRUE
  AND c.status <> 1
GROUP BY c.id, c.name, c.status
ORDER BY c.id;
