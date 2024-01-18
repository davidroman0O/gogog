package types

/// Generate file from the original JSON

type UserData struct {
	Country            string         `json:"country"`
	Currencies         []Currency     `json:"currencies"`
	SelectedCurrency   Currency       `json:"selectedCurrency"`
	PreferredLanguage  Language       `json:"preferredLanguage"`
	RatingBrand        string         `json:"ratingBrand"`
	IsLoggedIn         bool           `json:"isLoggedIn"`
	Checksum           UserChecksum   `json:"checksum"`
	Updates            UserUpdates    `json:"updates"`
	UserID             string         `json:"userId"`
	Username           string         `json:"username"`
	GalaxyUserID       string         `json:"galaxyUserId"`
	Email              string         `json:"email"`
	Avatar             string         `json:"avatar"`
	WalletBalance      Wallet         `json:"walletBalance"`
	PurchasedItems     PurchasedItems `json:"purchasedItems"`
	WishlistedItems    int            `json:"wishlistedItems"`
	Friends            []interface{}  `json:"friends"`
	PersonalizedPrices []interface{}  `json:"personalizedProductPrices"`
	PersonalizedSeries []interface{}  `json:"personalizedSeriesPrices"`
}

type Currency struct {
	Code   string `json:"code"`
	Symbol string `json:"symbol"`
}

type Language struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type UserChecksum struct {
	Cart         interface{} `json:"cart"`
	Games        string      `json:"games"`
	Wishlist     interface{} `json:"wishlist"`
	ReviewsVotes interface{} `json:"reviews_votes"`
	GamesRating  interface{} `json:"games_rating"`
}

type UserUpdates struct {
	Messages              int `json:"messages"`
	PendingFriendRequests int `json:"pendingFriendRequests"`
	UnreadChatMessages    int `json:"unreadChatMessages"`
	Products              int `json:"products"`
	Total                 int `json:"total"`
}

type Wallet struct {
	Currency string `json:"currency"`
	Amount   int    `json:"amount"`
}

type PurchasedItems struct {
	Games  int `json:"games"`
	Movies int `json:"movies"`
}
