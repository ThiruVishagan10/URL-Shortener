SELECT
    conname,
    pg_get_constraintdef (oid)
from
    pg_consstraint
WHERE
    conrelid = 'urls'::regclass;