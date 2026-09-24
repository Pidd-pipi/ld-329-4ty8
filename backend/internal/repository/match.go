package repository

import "cyskillswap/internal/model"

// ListMatches 返回智能匹配结果，匹配卡片是发起预约的入口
func ListMatches() []model.Match {
	return []model.Match{
		{ID: 1, Provider: "林澈", Learner: "孟野", OfferSkill: "毕业照人像摄影", WantedSkill: "民谣吉他陪练", Score: 96, CommonSlots: []string{"周三晚", "周六上午"}, Recommendation: "互补技能明确，双方均接受技能交换。"},
		{ID: 2, Provider: "周芮", Learner: "许安", OfferSkill: "Python 数据分析", WantedSkill: "论文数据清洗", Score: 89, CommonSlots: []string{"周二晚"}, Recommendation: "时间匹配且需求描述命中 pandas/可视化。"},
		{ID: 3, Provider: "孟野", Learner: "林澈", OfferSkill: "民谣吉他陪练", WantedSkill: "宣传照拍摄", Score: 91, CommonSlots: []string{"周六上午"}, Recommendation: "互换回报类型一致，信用分权重较高。"},
	}
}

// FindMatch 按 ID 查询匹配，发起预约时校验
func FindMatch(id int) (model.Match, bool) {
	for _, match := range ListMatches() {
		if match.ID == id {
			return match, true
		}
	}
	return model.Match{}, false
}
