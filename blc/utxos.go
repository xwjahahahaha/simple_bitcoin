package blc

//为了UTXOs的排序实现sort.Interface接口的三个方法
type UTXOs []*TxOutput

