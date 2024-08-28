package main

func findSuggestedFriends(username string, friendships map[string][]string) []string {
	directFriends, ok := friendships[username]
	if !ok {
		return nil
	}

	suggestedFriends := map[string]bool{}

	// Using different initializtion for demonstration purposes
	notSuggestedFriends := make(map[string]bool)

	// Add here all users that should not be included in the suggested friends slice
	notSuggestedFriends[username] = true
	for _, friend := range directFriends {
		notSuggestedFriends[friend] = true
	}

	for _, friend := range directFriends {
		secondDegFriends, ok := friendships[friend]
		if !ok {
			continue
		}

		for _, secondDegFriend := range secondDegFriends {
			// Add to suggested friends if it is not the user himself or any direct friend
			if _, ok := notSuggestedFriends[secondDegFriend]; !ok {
				suggestedFriends[secondDegFriend] = true
			}
		}
	}

	// Edge case: return nil slice (instead of slice with length 0)
	if len(suggestedFriends) == 0 {
		return nil
	}

	// We could have appended to a slice whenever we were inserting
	// to the map, but it is simpler to do it in one step

	// Here we make a slice with length 0 and capacity 10
	// This means we can use append to keep filling the slice without caring about index or capacity
	result := make([]string, 0, len(suggestedFriends))

	for friend := range suggestedFriends {
		result = append(result, friend)
	}

	return result
}
