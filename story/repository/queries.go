package repository

const FetchSavedPosts = "with post_interests as (select jsonb_agg(jsonb_build_object('id', interests.id, 'name', interests.name)) as interests, " +
	"post_x_interests.post_id " +
	"from posts " +
	"inner join post_x_interests on posts.id = post_x_interests.post_id " +
	"inner join interests on post_x_interests.interest_id = interests.id " +
	"inner join saved_posts sp on posts.id = sp.post_id and sp.user_id = $1 " +
	"group by post_x_interests.post_id) " +
	"select posts.id, " +
	"ap.title, " +
	"ap.tagline, " +
	"count(distinct l.liked_by)                                                                as likes_count, " +
	"count(distinct c.id)                                                                      as comments_count, " +
	"post_interests.interests, " +
	"u.id                                                                                      as author_id, " +
	"u.username                                                                                as author_name, " +
	"ap.preview_image                                                                          as preview_image, " +
	"posts.created_at                                                                          as published_at, " +
	"case when $2 in (l.liked_by) then true else false end as is_viewer_liked, " +
	"case when $3 = u.id then true else false end          as is_viewer_is_author " +
	"from posts " +
	"inner join saved_posts on posts.id = saved_posts.post_id and saved_posts.user_id = $4 " +
	"inner join post_interests on posts.id = post_interests.post_id " +
	"inner join post_x_interests on posts.id = post_x_interests.post_id " +
	"inner join interests on post_x_interests.interest_id = interests.id " +
	"inner join users u on u.id = posts.author_id " +
	"inner join abstract_post ap on posts.id = ap.post_id " +
	"left join likes l on l.post_id = posts.id " +
	"left join comments c on c.post_id = posts.id " +
	"where posts.deleted_at is null group by posts.id, u.id, ap.preview_image, l.liked_by, post_interests.interests, ap.title, ap.tagline limit $5 offset $6"
