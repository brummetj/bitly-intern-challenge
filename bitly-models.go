package bitly

/*
 * Model For Link History
 */

type EncodedUser struct {
  Login string `json:"login"`
  DisplayName string `json:"display_name"`
  FullName string `json:"full_name"`
}
type LinkHistory struct {
  HasLinkDeeplinks bool `json:"has_link_deeplinks"`
  Archived bool `json:"archived"`
  UserTs int64 `json:"user_ts"`
  Title string `json:"title"`
  CreateAt int64 `json:"created_at"`
  Tags []string `json:"tags"`
  ModifiedAt int64 `json:"modified_at"`
  CampaignIds []int64 `json:"campaign_ids"`
  Private bool `json:"private"`
  AggregateLink string `json:"aggregate_link"`
  LongUrl string `json:"long_url"`
  ClientId string `json:"client_id"`
  Link string `json:"link"`
  IsDomainDeeplink bool `json:"is_domain_deeplink"`
  EncodingUser EncodedUser `json:"encoding_user"`
}
type HistoryData struct {
  LinkHistoryList []LinkHistory `json:"link_history"`
  ResultCount int64 `json:"result_count"`
}
type HistoryModel  struct {
  statusCode int64 `json:"status_code"`
  Data HistoryData `json:"data"`
  StatusTxt string `json:"status_txt"`
}

/*
 * Model for Saving links
 */

type LinkSave struct {
  Link string `json:"link"`
  AggregateLink string `json:"aggregate_link"`
  LongUrl string `json:"long_url"`
  NewLink int16 `json:"new_link"`
  UserHash string `json:"user_hash"`
}

type LinkData struct {
  LinkSavedObject LinkSave `json:"link_save"`
}
type LinkModel struct {
  StatusCode int64 `json:"status_code"`
  Data LinkData `json:"data"`
  StatusTxt string `json:"status_txt"`
}

/*
 * Model for User Clicks
 */

type Click struct {
  Clicks int64 `json:"clicks"`
  DayStart int64 `json:"day_start"`
}

type ClickData struct {
  Days int64 `json:"days"`
  TotalClicks int64 `json:"day_start"`
  Clicks []Click `json:"clicks"`
}

type ClickModel struct {
  StatusCode int64 `json:"status_code"`
  Data ClickData `json:"data"`
  StatusTxt string `json:"status_txt"`
}

/*
    Model for a tweet

    NOTE: I kept this model simple just to give back the necessary information to show your new tweet!
 */
type Tweet struct {
  Id uint64
  Tweet string
  User string
}
