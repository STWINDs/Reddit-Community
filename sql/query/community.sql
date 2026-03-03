-- name: GetCommunityList :many
SELECT community_id as id, community_name as name FROM `community`;

-- name: GetCommunityDetailByID :one
SELECT community_id as id, community_name as name, introduction , create_time FROM `community` WHERE community_id = ?;