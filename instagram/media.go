//
// media.go
// Copyright 2017 Konstantin Dovnar
//
// Licensed under the Apache License, Version 2.0 (the "License");
// http://www.apache.org/licenses/LICENSE-2.0
//

package instagram

import (
	"encoding/json"
)

// TypeImage is a string that define image type for media.
const TypeImage = "image"

// TypeVideo is a string that define video type for media.
const TypeVideo = "video"

// TypeCarousel is a string that define carousel (collection of media) type for media.
const TypeCarousel = "carousel"

const (
	graphVideo   = "GraphVideo"
	graphSidecar = "GraphSidecar"

	video    = "video"
	carousel = "carousel"
)

type mediaJSON struct {
	Graphql struct {
		User struct {
			Biography       string      `json:"biography"`
			BlockedByViewer bool        `json:"blocked_by_viewer"`
			ConnectedFbPage interface{} `json:"connected_fb_page"`
			CountryBlock    bool        `json:"country_block"`
			EdgeFollow      struct {
				Count int `json:"count"`
			} `json:"edge_follow"`
			EdgeFollowedBy struct {
				Count int `json:"count"`
			} `json:"edge_followed_by"`
			EdgeMediaCollections struct {
				Count    int           `json:"count"`
				Edges    []interface{} `json:"edges"`
				PageInfo struct {
					EndCursor   interface{} `json:"end_cursor"`
					HasNextPage bool        `json:"has_next_page"`
				} `json:"page_info"`
			} `json:"edge_media_collections"`
			EdgeOwnerToTimelineMedia struct {
				Count int `json:"count"`
				Edges []struct {
					Node struct {
						Typename         string `json:"__typename"`
						CommentsDisabled bool   `json:"comments_disabled"`
						Dimensions       struct {
							Height int `json:"height"`
							Width  int `json:"width"`
						} `json:"dimensions"`
						DisplayURL  string `json:"display_url"`
						EdgeLikedBy struct {
							Count int `json:"count"`
						} `json:"edge_liked_by"`
						EdgeMediaPreviewLike struct {
							Count int `json:"count"`
						} `json:"edge_media_preview_like"`
						EdgeMediaToCaption struct {
							Edges []struct {
								Node struct {
									Text string `json:"text"`
								} `json:"node"`
							} `json:"edges"`
						} `json:"edge_media_to_caption"`
						EdgeMediaToComment struct {
							Count int `json:"count"`
						} `json:"edge_media_to_comment"`
						GatingInfo   interface{} `json:"gating_info"`
						ID           string      `json:"id"`
						IsVideo      bool        `json:"is_video"`
						MediaPreview string      `json:"media_preview"`
						Owner        struct {
							ID string `json:"id"`
						} `json:"owner"`
						Shortcode          string `json:"shortcode"`
						TakenAtTimestamp   int    `json:"taken_at_timestamp"`
						ThumbnailResources []struct {
							ConfigHeight int    `json:"config_height"`
							ConfigWidth  int    `json:"config_width"`
							Src          string `json:"src"`
						} `json:"thumbnail_resources"`
						ThumbnailSrc string `json:"thumbnail_src"`
					} `json:"node"`
				} `json:"edges"`
				PageInfo struct {
					EndCursor   string `json:"end_cursor"`
					HasNextPage bool   `json:"has_next_page"`
				} `json:"page_info"`
			} `json:"edge_owner_to_timeline_media"`
			EdgeSavedMedia struct {
				Count    int           `json:"count"`
				Edges    []interface{} `json:"edges"`
				PageInfo struct {
					EndCursor   interface{} `json:"end_cursor"`
					HasNextPage bool        `json:"has_next_page"`
				} `json:"page_info"`
			} `json:"edge_saved_media"`
			ExternalURL            string `json:"external_url"`
			ExternalURLLinkshimmed string `json:"external_url_linkshimmed"`
			FollowedByViewer       bool   `json:"followed_by_viewer"`
			FollowsViewer          bool   `json:"follows_viewer"`
			FullName               string `json:"full_name"`
			HasBlockedViewer       bool   `json:"has_blocked_viewer"`
			HasRequestedViewer     bool   `json:"has_requested_viewer"`
			ID                     string `json:"id"`
			IsPrivate              bool   `json:"is_private"`
			IsVerified             bool   `json:"is_verified"`
			ProfilePicURL          string `json:"profile_pic_url"`
			ProfilePicURLHd        string `json:"profile_pic_url_hd"`
			RequestedByViewer      bool   `json:"requested_by_viewer"`
			Username               string `json:"username"`
		} `json:"user"`
	} `json:"graphql"`
}

// A Media describes an Instagram media info.
type Media struct {
	Caption       string
	Code          string
	CommentsCount uint32
	Date          uint64
	ID            string
	AD            bool
	LikesCount    uint32
	Type          string
	MediaURL      string
	Owner         Account
	MediaList     []mediaItem
}

type mediaItem struct {
	Type string
	URL  string
	Code string
}

// Update try to update media data
func (m *Media) Update() error {
	media, err := GetMediaByCode(m.Code)
	if err != nil {
		return err
	}
	*m = media
	return nil
}

func getFromMediaPage(data []byte) (Media, error) {
	var m mediaJSON

	err := json.Unmarshal(data, &m)
	if err != nil {
		return Media{}, err
	}

	media := Media{}
	media.Code = m.Graphql.User.EdgeOwnerToTimelineMedia.Edges[0].Node.Shortcode
	media.ID = m.Graphql.User.ID
	// media.AD = m.Graphql.User.IsAd
	media.Date = uint64(m.Graphql.User.EdgeOwnerToTimelineMedia.Edges[0].Node.TakenAtTimestamp)
	media.CommentsCount = uint32(m.Graphql.User.EdgeOwnerToTimelineMedia.Edges[0].Node.EdgeMediaToComment.Count)
	media.LikesCount = uint32(m.Graphql.User.EdgeOwnerToTimelineMedia.Edges[0].Node.EdgeMediaPreviewLike.Count)

	if len(m.Graphql.User.EdgeOwnerToTimelineMedia.Edges[0].Node.EdgeMediaToCaption.Edges) > 0 {
		media.Caption = m.Graphql.User.EdgeOwnerToTimelineMedia.Edges[0].Node.EdgeMediaToCaption.Edges[0].Node.Text
	}

	// FIXME: This endpoint doesn't allow for any other media types (videos, etc) other than images
	// var mediaType = m.Graphql.User.EdgeOwnerToTimelineMedia.Edges[0].Node.Typename
	media.Type = TypeImage
	media.MediaURL = m.Graphql.User.EdgeOwnerToTimelineMedia.Edges[0].Node.DisplayURL
	var item mediaItem
	item.Code = media.Code
	item.Type = media.Type
	item.URL = media.MediaURL
	media.MediaList = append(media.MediaList, item)

	media.Owner.ID = m.Graphql.User.ID
	media.Owner.ProfilePicURL = m.Graphql.User.ProfilePicURL
	media.Owner.Username = m.Graphql.User.Username
	media.Owner.FullName = m.Graphql.User.FullName
	media.Owner.Private = m.Graphql.User.IsPrivate

	return media, nil
}

func getFromAccountMediaList(data []byte) ([]Media, error) {
	var m mediaJSON

	err := json.Unmarshal(data, &m)
	if err != nil {
		return []Media{}, err
	}

	num := len(m.Graphql.User.EdgeOwnerToTimelineMedia.Edges)
	medias := make([]Media, num)

	for i := 0; i < num; i++ {
		media := Media{}
		media.Code = m.Graphql.User.EdgeOwnerToTimelineMedia.Edges[i].Node.Shortcode
		media.ID = m.Graphql.User.ID
		// media.AD = m.Graphql.User.IsAd
		media.Date = uint64(m.Graphql.User.EdgeOwnerToTimelineMedia.Edges[i].Node.TakenAtTimestamp)
		media.CommentsCount = uint32(m.Graphql.User.EdgeOwnerToTimelineMedia.Edges[i].Node.EdgeMediaToComment.Count)
		media.LikesCount = uint32(m.Graphql.User.EdgeOwnerToTimelineMedia.Edges[i].Node.EdgeMediaPreviewLike.Count)

		if len(m.Graphql.User.EdgeOwnerToTimelineMedia.Edges[i].Node.EdgeMediaToCaption.Edges) > 0 {
			media.Caption = m.Graphql.User.EdgeOwnerToTimelineMedia.Edges[i].Node.EdgeMediaToCaption.Edges[0].Node.Text
		}

		// FIXME: This endpoint doesn't allow for any other media types (videos, etc) other than images
		// var mediaType = m.Graphql.User.EdgeOwnerToTimelineMedia.Edges[i.Node.Typename
		media.Type = TypeImage
		media.MediaURL = m.Graphql.User.EdgeOwnerToTimelineMedia.Edges[i].Node.DisplayURL
		var item mediaItem
		item.Code = media.Code
		item.Type = media.Type
		item.URL = media.MediaURL
		media.MediaList = append(media.MediaList, item)

		media.Owner.ID = m.Graphql.User.ID
		media.Owner.ProfilePicURL = m.Graphql.User.ProfilePicURL
		media.Owner.Username = m.Graphql.User.Username
		media.Owner.FullName = m.Graphql.User.FullName
		media.Owner.Private = m.Graphql.User.IsPrivate

		medias = append(medias, media)
	}

	return medias, nil
}
func getFromSearchMediaList(data []byte) (Media, error) {
	var mediaJSON struct {
		CommentsDisabled bool   `json:"comments_disabled"`
		ID               string `json:"id"`
		Owner            struct {
			ID string `json:"id"`
		} `json:"owner"`
		ThumbnailSrc string  `json:"thumbnail_src"`
		IsVideo      bool    `json:"is_video"`
		Code         string  `json:"code"`
		Date         float64 `json:"date"`
		DisplaySrc   string  `json:"display_src"`
		Caption      string  `json:"caption"`
		Comments     struct {
			Count float64 `json:"count"`
		} `json:"comments"`
		Likes struct {
			Count float64 `json:"count"`
		} `json:"likes"`
	}

	err := json.Unmarshal(data, &mediaJSON)
	if err != nil {
		return Media{}, err
	}

	media := Media{}
	media.ID = mediaJSON.ID
	media.Code = mediaJSON.Code
	media.MediaURL = mediaJSON.DisplaySrc
	media.Caption = mediaJSON.Caption
	media.Date = uint64(mediaJSON.Date)
	media.LikesCount = uint32(mediaJSON.Likes.Count)
	media.CommentsCount = uint32(mediaJSON.Comments.Count)
	media.Owner.ID = mediaJSON.Owner.ID

	if mediaJSON.IsVideo {
		media.Type = TypeVideo
	} else {
		media.Type = TypeImage
	}

	return media, nil
}
