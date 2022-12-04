SELECT
    "user_id",
    "name",
    "age"
FROM
    "user_service"."users"
WHERE
    "user_id" = $1;
