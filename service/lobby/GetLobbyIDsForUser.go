// service/lobby/lookup.go
package service_lobby

import (
    "context"
    "encoding/json"
    "strings"

    "github.com/nowel-xyz/quiz/database"
    "github.com/nowel-xyz/quiz/database/models"
    "github.com/nowel-xyz/quiz/routers/api/lobby/utils"
)

// GetLobbyIDsForUser scans all lobby:* keys and returns those
// in which the given user appears as a member.
func GetLobbyIDsForUser(ctx context.Context, user models.User) ([]string, error) {
    var (
        cursor uint64
        ids    []string
    )

    // SCAN in batches of 100
    for {
        keys, newCursor, err := database.Redis.Scan(ctx, cursor, "lobby:*", 100).Result()
        if err != nil {
            return nil, err
        }
        cursor = newCursor

        for _, key := range keys {
            raw, err := database.Redis.Get(ctx, key).Result()
            if err != nil {
                return nil, err
            }
            var lob utils.Lobby
            if err := json.Unmarshal([]byte(raw), &lob); err != nil {
                continue
            }

            // check membership
            for _, m := range lob.Members {
                if m.ID == user.ID {
                    // extract the lobbyID from the key
                    // key is "lobby:<lobbyID>"
                    parts := strings.SplitN(key, ":", 2)
                    if len(parts) == 2 {
                        ids = append(ids, parts[1])
                    }
                    break
                }
            }
        }

        if cursor == 0 {
            break
        }
    }

    return ids, nil
}
