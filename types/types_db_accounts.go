package types

// Finalized `Account` of an user
// One real user might have multiple accounts so it can backup all it's games
type Account struct {
	Email        string `json:"email" storm:"id"` // email of the account, should be unique since used for sign in
	Username     string `json:"username"`         // username of the account
	UserID       string `json:"userId"`           // i don't know what is this user id
	GalaxyUserID string `json:"galaxyUserId"`

	// TODO @droman: that's the url, we should store it somehow
	Avatar   string `json:"avatar"`   // url of the avatar
	Avatar64 string `json:"avatar64"` // base64 of the avatar

	Country string `json:"country"` // country of the account

	// Required to be able to leverage the application
	Cookies []Cookie `json:"cookies"` // cookies of the account

	// The cookies's might not work anymore, therefore we need to force the user to sign in OR transfer cookies again
	Expired bool `json:"expired"`
}
