package main

import "time"

var likedTiebaUrl = "https://tieba.baidu.com/mo/q/newmoindex"
var tbsUrl = "http://tieba.baidu.com/dc/common/tbs"
var signUrl = "http://c.tieba.baidu.com/c/c/forum/sign"
var commitUrl = "https://sofire.baidu.com/abot/api/v1/tpl/commit"
var addCotentUrl = "https://tieba.baidu.com/f/commit/post/add"

var _BSK = ""
var BizPo = ""
var NickName = ""
var BDUSS = ""
var cookies = "BDUSS=" + BDUSS

var chooseKw string
var chooseTitle string
var chooseFid string
var chooseTid int
var currentTbs string
var CurrentLine int
var SleepTime time.Duration

type TBSResult struct {
	Tbs     string `json:"tbs"`
	IsLogin int    `json:"is_login"`
	Error   error
}

type LIKEDBars struct {
	No              int    `json:"no"`
	Error           string `json:"error"`
	Error2          error
	Data            LikedBarsData `json:"data"`
	UbsSampleIds    string        `json:"ubs_sample_ids"`
	UbsAbtestConfig interface{}   `json:"ubs_abtest_config"`
}

type LikedBarsData struct {
	Uid       int         `json:"uid"`
	Tbs       string      `json:"tbs"`
	ItbTbs    string      `json:"itb_tbs"`
	LikeForum []ForumInfo `json:"like_forum"`
}

type ForumInfo struct {
	ForumName string `json:"forum_name"`
	UserLevel string `json:"user_level"`
	UserExp   string `json:"user_exp"`
	ForumId   int    `json:"forum_id"`
	IsLike    bool   `json:"is_like"`
	FavoType  int    `json:"favo_type"`
	IsSign    int    `json:"is_sign"`
}

//回复之后的响应
type CommentResponse struct {
	No      int             `json:"no"`
	ErrCode int             `json:"err_code"`
	Error   string          `json:"error"`
	Data    AddResponseData `json:"data"`
}

type AddResponseData struct {
	AutoMsg       string               `json:"autoMsg"`
	Fid           int                  `json:"fid"`
	FName         string               `json:"fname"`
	Tid           int                  `json:"tid"`
	IsLogin       int                  `json:"is_login"`
	Content       string               `json:"content"`
	AccessState   string               //`json:"access_statte"`
	VCode         AddResponseDataVCode `json:"vcode"`
	IsPostVisible int                  `json:"is_post_visible"`
}

type AddResponseDataVCode struct {
	NeedVCode       int    `json:"need_vcode"`
	StrReason       string `json:"str_reason"`
	CaptchaVCodeStr string `json:"captcha_vcode_str"`
	CaptchaCodeType int    `json:"captcha_code_type"`
	UserStateVCode  int    `json:"userstatevcode"`
}

type JtItem struct {
	C int    `json:"c"`
	M string `json:"m"`
	S string `json:"s"`
	T string `json:"t"`
}

type JtResponse struct {
	Code    int    `json:"code"`
	Data    JtItem `json:"data"`
	Id      int    `json:"id"`
	Message string `json:"message"`
	Error   error
}
