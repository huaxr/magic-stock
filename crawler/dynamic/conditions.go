// @Contact:    huaxinrui
// @Time:       2019/7/5 下午6:14

package dynamic

// 量能-近5天量能站上均线
func Condition(recent_count []float64, count_10_0, count_10_1, count_10_2, count_10_3, count_10_4, count_40_0, count_40_1, count_40_2, count_40_3, count_40_4 float64) (bool, bool) {
	tmp := []int{}
	tmp2 := []int{}

	if recent_count[0] > count_10_0 {
		tmp = append(tmp, 1)
	}
	if recent_count[1] > count_10_1 {
		tmp = append(tmp, 1)
	}
	if recent_count[2] > count_10_2 {
		tmp = append(tmp, 1)
	}
	if recent_count[3] > count_10_3 {
		tmp = append(tmp, 1)
	}

	if recent_count[4] > count_10_4 {
		tmp = append(tmp, 1)
	}

	if recent_count[0] > count_40_0 {
		tmp2 = append(tmp2, 1)
	}

	if recent_count[1] > count_40_1 {
		tmp2 = append(tmp2, 1)
	}

	if recent_count[2] > count_40_2 {
		tmp2 = append(tmp2, 1)
	}
	if recent_count[3] > count_40_3 {
		tmp2 = append(tmp2, 1)
	}
	if recent_count[4] > count_40_4 {
		tmp2 = append(tmp2, 1)
	}
	return len(tmp) >= 4, len(tmp2) >= 4
}

func Condition2(recent_count []float64, count_40_0, count_40_1, count_40_2, count_40_3, count_40_4 float64) bool {
	tmp2 := []int{}
	if (recent_count[0]-count_40_0)/count_40_0 > 1.8 {
		tmp2 = append(tmp2, 1)
	}

	if (recent_count[1]-count_40_1)/count_40_1 > 1.8 {
		tmp2 = append(tmp2, 1)
	}

	if (recent_count[2]-count_40_2)/count_40_2 > 1.8 {
		tmp2 = append(tmp2, 1)
	}
	if (recent_count[3]-count_40_3)/count_40_3 > 1.8 {
		tmp2 = append(tmp2, 1)
	}
	if (recent_count[4]-count_40_4)/count_40_4 > 1.8 {
		tmp2 = append(tmp2, 1)
	}
	return len(tmp2) >= 2
}

func ConditionLimitedPrice(array []float64, less float64) bool {
	if len(array) < 1 {
		return false
	}
	var x float64
	for _, i := range array {
		x += i
	}
	return x < less
}

func ConditionTop(array []float64, flag int) bool {
	if len(array) < 1 {
		return false
	}
	for _, i := range array[0:flag] {
		if i > 9.9 {
			return true
		}
	}
	return false
}

func ConditionTopLine(array_shou, array_high, array_kai, percent []float64, recent int) bool {
	if len(array_shou) < 1 {
		return false
	}
	// 上引线是实体柱5倍
	for i := 0; i <= recent-1; i++ {
		//if ((array_high[i] - array_shou[i])  / (array_shou[i] - array_kai[i])) > 5 && percent[i] > 1 {
		//	return true
		//}
		if ((array_high[i]-array_shou[i])/array_shou[i])*100/percent[i] > 2 && percent[i] > 2 {
			return true
		}
	}
	return false
}
