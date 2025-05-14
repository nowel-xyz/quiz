package utils

func ContainsMember(members []string, userID string) bool {
    for _, member := range members {
        if member == userID {
            return true
        }
    }
    return false
}