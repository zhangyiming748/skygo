package custom_util

import "strings"

func LikeEscape(likeStr string) string {
	s := strings.ReplaceAll(likeStr, "_", "\\_")
	return s
}
