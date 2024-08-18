package twitterscraper

import (
	"io"
	"net/url"
	"strings"
)

type FollowAction int

const (
	Follow FollowAction = iota + 1
	Unfollow
)

func (s *Scraper) FollowUser(user string, followAction FollowAction) error {
	var baseUrl string
	if followAction == Follow {
		baseUrl = "https://twitter.com/i/api/1.1/friendships/create.json"
	} else {
		baseUrl = "https://twitter.com/i/api/1.1/friendships/destroy.json"
	}

	userID, err := s.GetUserIDByScreenName(user)
	if err != nil {
		return err
	}

	req, err := s.newRequest("POST", baseUrl)
	if err != nil {
		return err
	}

	params := url.Values{
		"include_profile_interstitial_type": {"1"},
		"include_blocking":                  {"1"},
		"include_blocked_by":                {"1"},
		"include_followed_by":               {"1"},
		"include_want_retweets":             {"1"},
		"include_mute_edge":                 {"1"},
		"include_can_dm":                    {"1"},
		"include_can_media_tag":             {"1"},
		"include_ext_is_blue_verified":      {"1"},
		"include_ext_verified_type":         {"1"},
		"include_ext_profile_image_shape":   {"1"},
		"skip_status":                       {"1"},
		"user_id":                           {userID},
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Body = io.NopCloser(strings.NewReader(params.Encode()))

	err = s.RequestAPI(req, nil)
	if err != nil {
		return err
	}

	return nil
}

// FetchFollowing gets following profiles list for a given user, via the Twitter frontend GraphQL API.
func (s *Scraper) FetchFollowing(user string, maxUsersNbr int, cursor string) ([]*Profile, string, error) {
	userID, err := s.GetUserIDByScreenName(user)
	if err != nil {
		return nil, "", err
	}

	return s.FetchFollowingByUserID(userID, maxUsersNbr, cursor)
}

// FetchFollowingByUserID gets following profiles list for a given userID, via the Twitter frontend GraphQL API.
func (s *Scraper) FetchFollowingByUserID(userID string, maxUsersNbr int, cursor string) ([]*Profile, string, error) {
	if maxUsersNbr > 200 {
		maxUsersNbr = 200
	}

	req, err := s.newRequest("GET", "https://twitter.com/i/api/graphql/g5P4cbXR4ta4oCeE7y2vLQ/Following")
	if err != nil {
		return nil, "", err
	}

	variables := map[string]interface{}{
		"userId":                 userID,
		"includePromotedContent": false,
		"count":                  maxUsersNbr,
	}
	features := map[string]interface{}{
		"responsive_web_graphql_exclude_directive_enabled":                        true,
		"verified_phone_label_enabled":                                            false,
		"creator_subscriptions_tweet_preview_api_enabled":                         true,
		"responsive_web_graphql_timeline_navigation_enabled":                      true,
		"responsive_web_graphql_skip_user_profile_image_extensions_enabled":       false,
		"c9s_tweet_anatomy_moderator_badge_enabled":                               true,
		"tweetypie_unmention_optimization_enabled":                                true,
		"responsive_web_edit_tweet_api_enabled":                                   true,
		"graphql_is_translatable_rweb_tweet_is_translatable_enabled":              true,
		"view_counts_everywhere_api_enabled":                                      true,
		"longform_notetweets_consumption_enabled":                                 true,
		"responsive_web_twitter_article_tweet_consumption_enabled":                true,
		"tweet_awards_web_tipping_enabled":                                        false,
		"freedom_of_speech_not_reach_fetch_enabled":                               true,
		"standardized_nudges_misinfo":                                             true,
		"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true,
		"rweb_video_timestamps_enabled":                                           true,
		"longform_notetweets_rich_text_read_enabled":                              true,
		"longform_notetweets_inline_media_enabled":                                true,
		"responsive_web_enhance_cards_enabled":                                    false,
	}

	if cursor != "" {
		variables["cursor"] = cursor
	}

	query := url.Values{}
	query.Set("variables", mapToJSONString(variables))
	query.Set("features", mapToJSONString(features))
	req.URL.RawQuery = query.Encode()

	var timeline timelineV2
	err = s.RequestAPI(req, &timeline)
	if err != nil {
		return nil, "", err
	}

	users, nextCursor := timeline.parseUsers()

	if strings.HasPrefix(nextCursor, "0|") {
		nextCursor = ""
	}

	return users, nextCursor, nil
}

// FetchFollowers gets following profiles list for a given user, via the Twitter frontend GraphQL API.
func (s *Scraper) FetchFollowers(user string, maxUsersNbr int, cursor string) ([]*Profile, string, error) {
	userID, err := s.GetUserIDByScreenName(user)
	if err != nil {
		return nil, "", err
	}

	return s.FetchFollowersByUserID(userID, maxUsersNbr, cursor)
}

// FetchFollowersByUserID gets followers profiles list for a given userID, via the Twitter frontend GraphQL API.
func (s *Scraper) FetchFollowersByUserID(userID string, maxUsersNbr int, cursor string) ([]*Profile, string, error) {
	if maxUsersNbr > 200 {
		maxUsersNbr = 200
	}

	req, err := s.newRequest("GET", "https://twitter.com/i/api/graphql/jwbfbSzn0FRL_AMZGsYDag/Followers")
	if err != nil {
		return nil, "", err
	}

	variables := map[string]interface{}{
		"userId":                 userID,
		"includePromotedContent": false,
		"count":                  maxUsersNbr,
	}
	features := map[string]interface{}{
		"responsive_web_graphql_exclude_directive_enabled":                        true,
		"verified_phone_label_enabled":                                            false,
		"creator_subscriptions_tweet_preview_api_enabled":                         true,
		"responsive_web_graphql_timeline_navigation_enabled":                      true,
		"responsive_web_graphql_skip_user_profile_image_extensions_enabled":       false,
		"c9s_tweet_anatomy_moderator_badge_enabled":                               true,
		"tweetypie_unmention_optimization_enabled":                                true,
		"responsive_web_edit_tweet_api_enabled":                                   true,
		"graphql_is_translatable_rweb_tweet_is_translatable_enabled":              true,
		"view_counts_everywhere_api_enabled":                                      true,
		"longform_notetweets_consumption_enabled":                                 true,
		"responsive_web_twitter_article_tweet_consumption_enabled":                true,
		"tweet_awards_web_tipping_enabled":                                        false,
		"freedom_of_speech_not_reach_fetch_enabled":                               true,
		"standardized_nudges_misinfo":                                             true,
		"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true,
		"rweb_video_timestamps_enabled":                                           true,
		"longform_notetweets_rich_text_read_enabled":                              true,
		"longform_notetweets_inline_media_enabled":                                true,
		"responsive_web_enhance_cards_enabled":                                    false,
	}

	if cursor != "" {
		variables["cursor"] = cursor
	}

	query := url.Values{}
	query.Set("variables", mapToJSONString(variables))
	query.Set("features", mapToJSONString(features))
	req.URL.RawQuery = query.Encode()

	var timeline timelineV2
	err = s.RequestAPI(req, &timeline)
	if err != nil {
		return nil, "", err
	}

	users, nextCursor := timeline.parseUsers()

	if strings.HasPrefix(nextCursor, "0|") {
		nextCursor = ""
	}

	return users, nextCursor, nil
}
