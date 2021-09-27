package wallpaperstruct

type Response struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data Data   `json:"data"`
}
type Category struct {
	Count        int           `json:"count"`
	Ename        string        `json:"ename"`
	Rname        string        `json:"rname"`
	CoverTemp    string        `json:"cover_temp"`
	Name         string        `json:"name"`
	Cover        string        `json:"cover"`
	Rank         int           `json:"rank"`
	Filter       []interface{} `json:"filter"`
	Sn           int           `json:"sn"`
	Icover       string        `json:"icover"`
	Atime        int           `json:"atime"`
	Type         int           `json:"type"`
	ID           string        `json:"id"`
	PicassoCover string        `json:"picasso_cover"`
}
type Wallpaper struct {
	Views   int           `json:"views"`
	Ncos    int           `json:"ncos"`
	Rank    int           `json:"rank"`
	Tag     []string      `json:"tag"`
	Wp      string        `json:"wp"`
	Xr      bool          `json:"xr"`
	Cr      bool          `json:"cr"`
	Favs    int           `json:"favs"`
	Atime   int           `json:"atime"`
	ID      string        `json:"id"`
	Desc    string        `json:"desc"`
	Thumb   string        `json:"thumb"`
	Img     string        `json:"img"`
	Cid     []string      `json:"cid"`
	URL     []interface{} `json:"url"`
	Preview string        `json:"preview"`
	Store   string        `json:"store"`
}
type Res struct {
	Wallpaper []Wallpaper `json:"wallpaper"`
	Category  []Category  `json:"category"`
}
type Data struct {
	Msg  string `json:"msg"`
	Res  Res    `json:"res"`
	Code int    `json:"code"`
}

type Option struct {
	OperateType string `json:"operaType"`
	UserCode    string `json:"userCode"`
	OperateData string `json:"operateData"`
}

type WallStrategy struct {
	TimeGap      int      `json:"timeGap"`
	StrategyCode string   `json:"strategyCode"`
	Categories   []string `json:"categories"`
	UserCode     string   `json:"userCode"`
	TimeUnit     int      `json:"timeUnit"`
}
