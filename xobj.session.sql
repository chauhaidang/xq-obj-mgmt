INSERT INTO xobjects (ref, type, description, created_at)
VALUES (
        'ref:character varying',
        'type:character varying',
        'description:character varying',
        current_timestamp
    );
truncate xobjects restart identity cascade;
SELECT *
FROM xobjects
WHERE ref = ''