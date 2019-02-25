package rescode

/*
10000	操作成功
10001	操作失败
20001	token不存在
20002	token过期
20003	token非法
20004	登录超时
30001	无管理员权限
*/
const (
	Success = "10000"
	Error   = "10001"
	// 验证
	Token_Missing   = "20001"
	Token_Timed_Out = "20002"
	Token_Invalid   = "20003"
	Login_Timed_Out = "20004"
	// 授权
	Access_Denied = "30001"
)
