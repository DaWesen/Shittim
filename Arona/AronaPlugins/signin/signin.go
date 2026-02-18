package signin

import (
	"fmt"
	"math/rand"
	"time"

	"Shittim/pkg/database"
	"Shittim/pkg/models"
)

// 日常事件类型
type DailyEvent struct {
	Type        string // 事件类型
	Description string // 事件描述
}

// 日常事件列表
var dailyEvents = []DailyEvent{
	{"阿拜多斯日常", "今天遇到了阿拜多斯的小鸟游星野，她向你汇报了最近的学习情况，并邀请你参加她们的活动。"},
	{"歌赫娜邂逅", "在歌赫娜学院遇到了空崎日奈，她正在练习射击，看到你后热情地邀请你指导她的技巧。"},
	{"甜点部时光", "崔尼蒂的甜点部成员柚鸟夏为你准备了美味的甜点，邀请你品尝并分享了她的烘焙心得。"},
	{"茶话会邀请", "圣园未花邀请你参加她们的茶话会，在轻松愉快的氛围中，你们交流了许多有趣的话题。"},
	{"课后辅导", "为学习有困难的学生提供了课后辅导，看到她们理解知识点时的笑容，你感到很欣慰。"},
	{"校园漫步", "在校园里悠闲漫步时，遇到了正在练习的小鸟游星野，她向你展示了最近的训练成果。"},
	{"学生委托", "收到了空崎日奈的委托，帮助她解决了一些训练中的问题，获得了她的信任和感谢。"},
	{"甜点制作", "柚鸟夏邀请你一起制作甜点，虽然过程有些手忙脚乱，但最终做出的甜点非常美味。"},
	{"茶话会准备", "协助圣园未花准备茶话会，她对你的建议非常重视，让你感受到了作为老师的价值。"},
	{"学园祭准备", "和小鸟游星野、空崎日奈、柚鸟夏、圣园未花一起准备学园祭活动，虽然忙碌但充满了乐趣。"},
	{"意外惊喜", "学生们为你准备了小惊喜，小鸟游星野、空崎日奈、柚鸟夏和圣园未花都参与其中，让你深受感动。"},
	{"教师日常", "今天收到了学生们的问候和关心，让你感受到了幸福。"},
}

// 获取随机日常事件
func getRandomDailyEvent() DailyEvent {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return dailyEvents[r.Intn(len(dailyEvents))]
}

// 检查用户今日是否已签到
func CheckSignin(qq int64) bool {
	var signin models.Signin
	today := time.Now().Format("2006-01-02")

	result := database.GetDB().Where("qq = ? AND date = ?", qq, today).First(&signin)
	return result.RowsAffected > 0
}

// 执行签到操作
func DoSignin(qq int64, nickname string) (int, int, DailyEvent, error) {
	//检查是否已签到
	signed := CheckSignin(qq)
	if signed {
		return 0, 0, DailyEvent{}, fmt.Errorf("今天已经签到过了")
	}

	//开始事务
	tx := database.GetDB().Begin()

	//查找或创建用户
	var user models.User
	if err := tx.Where("qq = ?", qq).FirstOrCreate(&user, models.User{
		QQ:       qq,
		Nickname: nickname,
	}).Error; err != nil {
		tx.Rollback()
		return 0, 0, DailyEvent{}, err
	}

	//计算连续签到天数
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	var lastSignin models.Signin
	streak := 1

	if tx.Where("qq = ? AND date = ?", qq, yesterday).First(&lastSignin).Error == nil {
		streak = lastSignin.Streak + 1
	}

	//计算随机奖励
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalReward := r.Intn(16)

	//获取随机日常事件
	event := getRandomDailyEvent()

	//创建签到记录
	signin := models.Signin{
		UserID: user.ID,
		QQ:     qq,
		Date:   time.Now().Format("2006-01-02"),
		Reward: totalReward,
		Streak: streak,
	}
	if err := tx.Create(&signin).Error; err != nil {
		tx.Rollback()
		return 0, 0, DailyEvent{}, err
	}

	//更新用户经验
	user.Exp += totalReward
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return 0, 0, DailyEvent{}, err
	}

	//提交事务
	if err := tx.Commit().Error; err != nil {
		return 0, 0, DailyEvent{}, err
	}

	return totalReward, streak, event, nil
}

// 获取经验排行榜
func GetExpRank(limit int) ([]models.User, error) {
	var users []models.User
	result := database.GetDB().Order("exp desc").Limit(limit).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// 获取用户在排行榜中的排名
func GetUserRank(qq int64) (int, error) {
	//查找用户
	var user models.User
	if err := database.GetDB().Where("qq = ?", qq).First(&user).Error; err != nil {
		return 0, err
	}

	//计算排名
	var count int64
	if err := database.GetDB().Model(&models.User{}).Where("exp > ?", user.Exp).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count) + 1, nil
}
