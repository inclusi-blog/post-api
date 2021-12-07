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
	"ap.url, " +
	"u.id                                                                                      as author_id, " +
	"u.username                                                                                as author_name, " +
	"ap.preview_image                                                                          as preview_image, " +
	"posts.created_at                                                                          as published_at, " +
	"case when $2 in (l.liked_by) then true else false end as is_viewer_liked, " +
	"case when $3 = u.id then true else false end          as is_viewer_is_author, " +
	"true as is_bookmarked " +
	"from posts " +
	"inner join saved_posts on posts.id = saved_posts.post_id and saved_posts.user_id = $4 " +
	"inner join post_interests on posts.id = post_interests.post_id " +
	"inner join post_x_interests on posts.id = post_x_interests.post_id " +
	"inner join interests on post_x_interests.interest_id = interests.id " +
	"inner join users u on u.id = posts.author_id " +
	"inner join abstract_post ap on posts.id = ap.post_id " +
	"left join likes l on l.post_id = posts.id " +
	"left join comments c on c.post_id = posts.id " +
	"where posts.deleted_at is null group by posts.id, ap.url, u.id, ap.preview_image, l.liked_by, post_interests.interests, ap.title, ap.tagline limit $5 offset $6"

const FetchViewedPosts = "with post_interests as (select jsonb_agg(jsonb_build_object('id', interests.id, 'name', interests.name)) as interests, " +
	"post_x_interests.post_id " +
	"from posts " +
	"inner join post_x_interests on posts.id = post_x_interests.post_id " +
	"inner join interests on post_x_interests.interest_id = interests.id " +
	"inner join post_views pv on posts.id = pv.post_id and " +
	"pv.user_id = $1 " +
	"group by post_x_interests.post_id) " +
	"select posts.id, " +
	"ap.title, " +
	"ap.tagline, " +
	"count(distinct l.liked_by)                                                                as likes_count, " +
	"count(distinct c.id)                                                                      as comments_count, " +
	"post_interests.interests, " +
	"u.id                                                                                      as author_id, " +
	"ap.url, " +
	"u.username                                                                                as author_name, " +
	"ap.preview_image                                                                          as preview_image, " +
	"posts.created_at                                                                          as published_at, " +
	"case when $2 in (l.liked_by) then true else false end as is_viewer_liked, " +
	"case when $3 = u.id then true else false end          as is_viewer_is_author, " +
	"case when $4 in (sp.user_id) then true else false end as is_bookmarked " +
	"from posts " +
	"inner join post_views p on posts.id = p.post_id and p.user_id = $5 " +
	"left join saved_posts sp on posts.id = sp.post_id and p.user_id = $6 " +
	"inner join post_interests on posts.id = post_interests.post_id " +
	"inner join post_x_interests on posts.id = post_x_interests.post_id " +
	"inner join interests on post_x_interests.interest_id = interests.id " +
	"inner join users u on u.id = posts.author_id " +
	"inner join abstract_post ap on posts.id = ap.post_id " +
	"left join likes l on l.post_id = posts.id " +
	"left join comments c on c.post_id = posts.id " +
	"where posts.deleted_at is null " +
	"group by posts.id, u.id, ap.preview_image,ap.url, l.liked_by, post_interests.interests, ap.title, ap.tagline, sp.user_id " +
	"limit $7 offset $8"

const FetchPostByInterests = "with ins as ( " +
	"select post_id, jsonb_agg(jsonb_build_object('id', i.id, 'name', i.name)) as interests " +
	"from post_x_interests " +
	"inner join interests i on i.id = post_x_interests.interest_id " +
	"inner join posts p on p.id = post_x_interests.post_id " +
	"group by post_id " +
	")" +
	"select ins.post_id                                                                               as id, " +
	"ap.title, " +
	"ap.url, " +
	"ap.tagline, " +
	"posts.author_id, " +
	"count(distinct l.liked_by)                                                                as likes_count, " +
	"count(distinct c.id)                                                                      as comments_count, " +
	"ins.interests, " +
	"username                                                                                  as author_name, " +
	"preview_image, " +
	"posts.created_at                                                                          as published_at, " +
	"case when $1 in (l.liked_by) then true else false end as is_viewer_liked, " +
	"case when $2 = users.id then true else false end      as is_viewer_is_author, " +
	"case when $3 in (sp.user_id) then true else false end as is_bookmarked " +
	"from post_x_interests " +
	"inner join posts on post_x_interests.post_id = posts.id and deleted_at is null " +
	"inner join abstract_post ap on posts.id = ap.post_id " +
	"inner join users on posts.author_id = users.id " +
	"inner join ins on ins.post_id = post_x_interests.post_id " +
	"left join likes l on posts.id = l.post_id " +
	"left join comments c on posts.id = c.post_id " +
	"left join saved_posts sp on posts.id = sp.post_id " +
	"where interest_id = $4 " +
	"group by ap.title, ins.post_id, ap.tagline, ap.url, posts.author_id, ins.interests, username, preview_image, liked_by, users.id, " +
	"sp.user_id, posts.created_at order by posts.created_at limit $5 offset $6"
