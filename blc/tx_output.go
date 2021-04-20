package blc

type TxOutput struct {
	// 值
	Value int
	// 解锁规则
	ScriptPubKey string
}


// 验证输出是否可解锁
func (out *TxOutput) CanUnlockedWith(unlockingData string) bool {
	// TODO  复杂的判断（属性密码？）
	if out.ScriptPubKey == unlockingData {
		return true
	}
	return false
}
