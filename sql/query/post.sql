-- name: CreatePost :exec
INSERT INTO `post` (post_id, author_id, community_id, title, content, create_time, update_time) VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- name: GetPostDetailByID :one
SELECT post_id as id, author_id, community_id, title, status, content, create_time FROM `post` WHERE post_id = ?;


-- name: GetPostListWithDetails :many
SELECT 
  p.post_id as id, 
  p.author_id, 
  u.username as author_name,
  p.community_id, 
  c.community_name,
  p.title, 
  p.content, 
  p.status,
  p.create_time
FROM `post` p
LEFT JOIN `user` u ON p.author_id = u.user_id
LEFT JOIN `community` c ON p.community_id = c.community_id
ORDER BY p.create_time DESC
LIMIT ?,?;


#根据给定的id列表查询帖子数据

-- name: GetPostListByIDs :many
SELECT post_id, title, content, author_id, community_id, create_time 
FROM post 
WHERE FIND_IN_SET(post_id, ?) 
ORDER BY FIND_IN_SET(post_id, ?);