package tool

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/bwmarrin/snowflake"
	"io"
	"strings"
)

const SESS_PREFIEX = "HEYAO"

func GetSnowflakeId() string {
	//default node id eq 1,this can modify to different serverId node
	node, _ := snowflake.NewNode(1)
	// Generate a snowflake ID.
	id := node.Generate().String()
	return id
}

func GetRandomToken(length int) string {
	r := make([]byte, length)
	io.ReadFull(rand.Reader, r)
	return base64.URLEncoding.EncodeToString(r)
}

func GetSessionID(token string) string {
	return SESS_PREFIEX + token
}

func GetToken(sessionId string) string {
	return strings.TrimPrefix(sessionId, SESS_PREFIEX)
}
// TODO:Added from the remote
