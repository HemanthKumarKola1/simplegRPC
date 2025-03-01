-- name: CreateUser :exec
INSERT INTO public.user (userName, passwordHash, email, userType)
VALUES ($1, $2, $3, $4);

-- name: GetUserByID :one
SELECT * FROM public.user WHERE userID = $1;

-- name: GetUserByUserName :one
SELECT * FROM public.user WHERE userName = $1;

-- name: UpdateUser :exec
UPDATE public.user
SET userName = $1, email = $2
WHERE userID = $3;

-- name: DeleteUser :exec
DELETE FROM public.user WHERE userID = $1;