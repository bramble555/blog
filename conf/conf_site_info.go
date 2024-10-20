package config

type SiteInfo struct {
	CreateAt    string `mapstructure:"created_at" json:"created_at"`
	Beian       string `mapstructure:"bei_an" json:"bei_an"`
	Title       string `mapstructure:"title" json:"title"`
	QQImage     string `mapstructure:"qq_image" json:"qq_image"`
	Version     string `mapstructure:"version" json:"version"`
	Email       string `mapstructure:"email" json:"email"`
	WechatImage string `mapstructure:"wechat_image" json:"wechat_image"`
	Name        string `mapstructure:"name" json:"name"`
	Job         string `mapstructure:"job" json:"job"`
	Addr        string `mapstructure:"addr" json:"addr"`
	Slogan      string `mapstructure:"slogan" json:"slogan"`
	SloganEn    string `mapstructure:"slogan_en" json:"slogan_en"`
	Web         string `mapstructure:"web" json:"web"`
	BilibiliUrl string `mapstructure:"bilibili_url" json:"bilibili_url"`
	GiteeUrl    string `mapstructure:"gitee_url" json:"gitee_url"`
	GithubUrl   string `mapstructure:"github_url" json:"github_url"`
}
