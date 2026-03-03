-- name: GetUserByUserID :one
SELECT * FROM `user` WHERE user_id = ?;

-- name: GetUserByUsername :one
SELECT * FROM `user` WHERE username = ?;

-- name: UpdateUserNameByID :exec
UPDATE `user` SET username = ? and update_time = CURRENT_TIMESTAMP WHERE id = ?;

-- name: CreateUser :exec
INSERT INTO `user` (user_id, username, password, email, gender, create_time, update_time) VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- name: DeleteUserByID :exec
DELETE FROM `user` WHERE id = ?;
