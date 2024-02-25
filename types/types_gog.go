package types

type Configuration[T any] func(params *T) error

func NewDependencyParams[T any](cfgs ...Configuration[T]) (*T, error) {
	data := new(T)
	for i := 0; i < len(cfgs); i++ {
		if err := cfgs[i](data); err != nil {
			return nil, err
		}
	}
	return data, nil
}

type Cookie struct {
	Domain         string  `json:"domain"`
	ExpirationDate float64 `json:"expirationDate"`
	HostOnly       bool    `json:"hostOnly"`
	HTTPOnly       bool    `json:"httpOnly"`
	Name           string  `json:"name"`
	Path           string  `json:"path"`
	SameSite       string  `json:"sameSite"`
	Secure         bool    `json:"secure"`
	Session        bool    `json:"session"`
	StoreID        string  `json:"storeId"`
	Value          string  `json:"value"`
	ID             int     `json:"id"`
}

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

type GameMeta struct {
	Title                  string          `json:"title"`
	BackgroundImage        string          `json:"backgroundImage"`
	CdKey                  string          `json:"cdKey"`
	TextInformation        string          `json:"textInformation"`
	Downloads              [][]interface{} `json:"downloads"`
	GalaxyDownloads        []interface{}   `json:"galaxyDownloads"`
	Extras                 []*Download     `json:"extras"`
	Dlcs                   []interface{}   `json:"dlcs"`
	Tags                   []interface{}   `json:"tags"`
	IsPreOrder             bool            `json:"isPreOrder"`
	ReleaseTimestamp       int             `json:"releaseTimestamp"`
	Messages               []interface{}   `json:"messages"`
	Changelog              string          `json:"changelog"`
	ForumLink              string          `json:"forumLink"`
	IsBaseProductMissing   bool            `json:"isBaseProductMissing"`
	MissingBaseProduct     interface{}     `json:"missingBaseProduct"`
	Features               []interface{}   `json:"features"`
	SimpleGalaxyInstallers []struct {
		Path string `json:"path"`
		Os   string `json:"os"`
	} `json:"simpleGalaxyInstallers"`
	ProductID *int `json:"productID"`
}

type Download struct {
	ManualURL string `json:"manualUrl"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Date      string `json:"date"`
	Size      string `json:"size"`
	Type      string `json:"type"`
}

type Product struct {
	IsGalaxyCompatible bool          `json:"isGalaxyCompatible"`
	Tags               []interface{} `json:"tags"`
	ID                 int           `json:"id"`
	Availability       struct {
		IsAvailable          bool `json:"isAvailable"`
		IsAvailableInAccount bool `json:"isAvailableInAccount"`
	} `json:"availability"`
	Title   string `json:"title"`
	Image   string `json:"image"`
	URL     string `json:"url"`
	WorksOn struct {
		Windows bool `json:"Windows"`
		Mac     bool `json:"Mac"`
		Linux   bool `json:"Linux"`
	} `json:"worksOn"`
	Category     string `json:"category"`
	Rating       int    `json:"rating"`
	IsComingSoon bool   `json:"isComingSoon"`
	IsMovie      bool   `json:"isMovie"`
	IsGame       bool   `json:"isGame"`
	Slug         string `json:"slug"`
	Updates      int    `json:"updates"`
	IsNew        bool   `json:"isNew"`
	DlcCount     int    `json:"dlcCount"`
	ReleaseDate  struct {
		Date         string `json:"date"`
		TimezoneType int    `json:"timezone_type"`
		Timezone     string `json:"timezone"`
	} `json:"releaseDate"`
	IsBaseProductMissing bool          `json:"isBaseProductMissing"`
	IsHidingDisabled     bool          `json:"isHidingDisabled"`
	IsInDevelopment      bool          `json:"isInDevelopment"`
	ExtraInfo            []interface{} `json:"extraInfo"`
	IsHidden             bool          `json:"isHidden"`
}

type Search struct {
	SortBy                     string      `json:"sortBy"`
	Page                       int         `json:"page"`
	TotalProducts              int         `json:"totalProducts"`
	TotalPages                 int         `json:"totalPages"`
	ProductsPerPage            int         `json:"productsPerPage"`
	ContentSystemCompatibility interface{} `json:"contentSystemCompatibility"`
	MoviesCount                int         `json:"moviesCount"`
	Tags                       []struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		ProductCount string `json:"productCount"`
	} `json:"tags"`
	Products                   []Product
	UpdatedProductsCount       int `json:"updatedProductsCount"`
	HiddenUpdatedProductsCount int `json:"hiddenUpdatedProductsCount"`
	AppliedFilters             struct {
		Tags interface{} `json:"tags"`
	} `json:"appliedFilters"`
	HasHiddenProducts bool `json:"hasHiddenProducts"`
}

type PlatformName string

var (
	Windows PlatformName = "1,2,4,8,4096,16384"
	Linux   PlatformName = "1024,2048,8192"
	Mac     PlatformName = "16,32"
)

var AllPlatforms = []PlatformName{Windows, Linux, Mac}

type SearchParams struct {
	Query        *string
	Page         *int
	PlatformName *PlatformName
	Language     *string
}

func NewSearchParams(cfgs ...Configuration[SearchParams]) (*SearchParams, error) {
	return NewDependencyParams[SearchParams](cfgs...)
}

func SearchWithQuery(query string) Configuration[SearchParams] {
	return func(params *SearchParams) error {
		params.Query = &query
		return nil
	}
}

func SearchWithPage(page int) Configuration[SearchParams] {
	return func(params *SearchParams) error {
		params.Page = &page
		return nil
	}
}

func SearchWithPlatformID(platform PlatformName) Configuration[SearchParams] {
	return func(params *SearchParams) error {
		params.PlatformName = &platform
		return nil
	}
}

func SearchWithLanguage(language string) Configuration[SearchParams] {
	return func(params *SearchParams) error {
		params.Language = &language
		return nil
	}
}
