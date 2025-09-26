package authorizationscache

import "sync"

type AuthorizationsCache struct {
	cache                map[string]bool     // key would be userName concatenated with doorId
	usersAuthorizedDoors map[string][]string // key would be userName and value would be list of doorIds (optimization for evicting user from cache)
}

var (
	instance *AuthorizationsCache
	once     sync.Once
)

func GetAuthorizationsCacheInstance() *AuthorizationsCache {
	once.Do(func() {
		instance = &AuthorizationsCache{
			cache:                make(map[string]bool),
			usersAuthorizedDoors: make(map[string][]string),
		}
	})
	return instance
}

func (ac *AuthorizationsCache) SetAuthorization(userName string, doorId string) {
	key := userName + doorId
	exists := ac.cache[key]
	if exists {
		return
	}
	ac.usersAuthorizedDoors[userName] = append(ac.usersAuthorizedDoors[userName], doorId)
	ac.cache[key] = true

}

func (ac *AuthorizationsCache) IsAuthorized(userName string, doorId string) bool {
	key := userName + doorId
	return ac.cache[key]
}

func (ac *AuthorizationsCache) EvictFromCache(userName string) string {
	authorizedDoors := ac.usersAuthorizedDoors[userName]
	if authorizedDoors == nil {
		return "No authorized doors found for user"
	}
	for _, doorId := range authorizedDoors {
		key := userName + doorId
		delete(ac.cache, key)
	}
	delete(ac.usersAuthorizedDoors, userName)
	return "User cache evicted successfully"
}
