package oddsUtil

import "math"

// ConvertEuroToMyanmar convert euro odds to myanmar odds
func ConvertEuroToMyanmar(handicap, euroOdds float64) (goals int, odds float64) {
	hk := ConvertToHK(euroOdds)
	malay := ConvertToMalay(hk)

	return ConvertToMyanmar(handicap, malay)
}

// ConvertToHK converts euro odds to HK odds
func ConvertToHK(euroOdds float64) float64 {
	return euroOdds - 1
}

// ConvertToMalay converts hk odds to malay odds
func ConvertToMalay(hkOdds float64) float64 {
	if hkOdds < 1 {
		return hkOdds
	}

	return -1 / hkOdds
}

// ConvertToMyanmar convert malay odds to myanmar odds
func ConvertToMyanmar(handicap, malayOdds float64) (goals int, odds float64) {
	hdp := math.Abs(handicap)
	goals = int(hdp)
	mmOdds := mmBaseOdds(malayOdds)

	fraction := hdp - float64(goals)

	var price float64
	switch fraction {
	case 0:
		if malayOdds > 0 {
			price = mmOdds * -1
		} else {
			price = mmOdds
		}
	case 0.25:
		if malayOdds > 0 {
			price = (malayOdds / -2) - mmOdds
		} else {
			price = (malayOdds / 2) + mmOdds
		}
	case 0.5:
		if malayOdds > 0 {
			goals++
			price = malayOdds - mmOdds
		} else {
			price = malayOdds + mmOdds
		}
	case 0.75:
		goals++
		if malayOdds > 0 {
			price = (malayOdds / 2) - mmOdds
		} else {
			price = (malayOdds / -2) + mmOdds
		}
	}

	return goals, price
}

func mmBaseOdds(malayOdds float64) float64 {
	var mmOdds float64
	if malayOdds > 0 {
		mmOdds = 1 - malayOdds
	} else {
		mmOdds = 1 + malayOdds
	}

	return mmOdds
}
