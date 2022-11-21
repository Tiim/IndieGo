package indieauth

type ScopeCheck func(scope string) bool

type TokenVerifier func(token string, minimalScopes []string) (ScopeCheck, error)
