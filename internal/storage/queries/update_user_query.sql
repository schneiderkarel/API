UPDATE
    "user_service"."users"
SET
    "name" = $2, "age" = $3
WHERE
    "user_id" = $1;
