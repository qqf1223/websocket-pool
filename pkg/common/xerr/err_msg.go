package xerr

var message map[uint32]string

func init() {
	message = make(map[uint32]string)
	message[OK] = "SUCCESS"
	message[ServerCommonError] = "服务器开小差啦,稍后再来试一试"
	message[RequestParamError] = "参数错误"
	message[TokenExpireError] = "token失效，请重新登陆"
	message[TokenGenerateError] = "生成token失败"
	message[DbError] = "系统繁忙,请稍后再试"
	message[GetTeamIdError] = "获取teamId失败"
	message[GetBaseUserIdError] = "无法获取登录状态"
	message[ImError] = "im服务发生错误"
	message[UploadFileToCosError] = "上传图片发生错误"
	message[JsonMarshalError] = "Json加密错误"
	message[JsonUnMarshalError] = "Json解密错误"

	message[GetAppError] = "获取App失败"
	message[GetOAuthQrError] = "获取QrCode失败"
	message[QrCodeExpiredError] = "该二维码已经过期"
	message[QrCodeScannedError] = "该二维码已经被扫描"
	message[InvalidClientError] = "无效的客户端凭证"
	message[GetTokenError] = "获取token失败"
	message[SetLoginStatusError] = "设置登录状态失败"
	message[JwtInvalidError] = "jwt作废失败"
	message[StateCodeError] = "无效的授权码"

	// 注册模块
	message[VerificationCodeError] = "验证码错误"
	message[UserAlreadyRegistered] = "用户已经注册"
	message[RegisterNameError] = "仅支持汉字、数字和字母"

	// 群组
	message[GroupMemberNumberError] = "最少要选择 2 人"
	message[GroupMemberExceedError] = "最多可选择 999 人"
	message[QueryGroupByIdError] = "查询话题异常"
	message[GroupUserStatusError] = "该用户状态异常"
	message[GroupGoalNoFinishTask] = "话题目标还有未完成任务"
	message[GroupGoalFinishByOther] = "本目标已被其他用户完成"
	message[GroupGoalNotAddTask] = "本目标已完成，无法添加任务"
	message[GroupGoalTaskAuthFail] = "暂无权限修改"
	message[GroupGoalInvalidMattersId] = "无效的事项ID"
	message[GroupGoalCanNotUpdatePrincipal] = "有任务设置了金币奖励，暂不支持变更目标负责人"
	message[GroupTaskCanNotMove] = "本任务设置了金币奖励，暂不支持移动"
	message[GroupCodeNotExist] = "话题ID不存在， 请重新输入"

	// chat
	message[NotInGroupError] = "发送人不在群组之内"
	message[ConversationNotExistError] = "会话不存在"
	message[ConversationExError] = "排序失败"
	message[ConversationPinnedOrderError] = "order错误"

	// webhook
	message[ReceiveUserError] = "接收人参数错误"
	message[WebHookUserIdError] = "预警机器人不存在"
	message[WebHookUrlError] = "非法的请求,请检查hook地址"
	message[WebHookUserParamError] = "查询用户异常,请检查用户Id"

	// 云文档@提醒
	message[GetPermissionError] = "获取云文档权限失败"

	// 用户签到
	message[MemberSignedError] = "用户已经签过到了， 请明天再来"

	// 组织事项
	message[OrganizationMattersDelError] = "该组织事项还有子集， 无法删除"
	message[OrganizationMattersCircularError] = "移动到的父节点不能是原节点的子节点"
	message[TaskMoveToGoalFinishedError] = "本目标已完成，任务不能移入"
	message[GoalSameError] = "移入目标是当前目标， 不可移动"
	message[GoldBalanceNotSufficient] = "目标负责人的金币余额不足"
	message[GoldRelationNoData] = "金币关系表未查询到数据， 请求参数错误"
	message[GoldAccountInfoNotFound] = "未查询到账户信息， 请稍后再试"
	message[GoldExecutorNotNull] = "金币有值时任务执行人必传"
	message[GoldAccountNotNull] = "数据错误， 金币账户不存在"
	message[GoldGoalPrincipalNotNull] = "请先设置目标负责人后才可添加金币"
	message[GoldPrincipalNotNull] = "金币有值是目标负责人必传"
	message[PrincipalGoldBalanceNotSufficient] = "目标创建人的金币余额不足"
	message[GoldBalanceGreaterThanZero] = "更新时金币余额必须大于0"

	// 会议
	message[ForbiddenChangePattern] = "仅主持人身份才能切换会议模式"
	message[MeetingPatternError] = "会议模式错误"
	message[ForbiddenProcessSpeakRequest] = "仅主持人身份才能操作"
	message[SetAskCacheError] = "申请失败，请重试"
	message[ProcessSpeakCacheError] = "处理申请失败，请重试"
	message[ForbiddenOperation] = "禁止操作"
	message[MeetingEnded] = "主持人已经结束会议"
	message[MeetingRoomReserveError] = "该时段无法预约，请重新预约"
	message[MeetingRoomPartOfReserverError] = "部分会议室无法预约"
	message[MeetingRoomNoPrivilegeError] = "会议室非本人创建"
	message[MeetingRoomShareNoSelect] = "请选择分享人或话题"
	message[MeetingRoomEditNotifyAdd] = "向新增参与者发送会议邀请"
	message[MeetingRoomEditNotifyDel] = "向被删除的参与者发送通知"
	message[MeetingRoomEditNotifyAddUpdate] = "向新增参与者发送会议邀请，并向现有参与者发送会议更新通知"
	message[MeetingRoomEditNotifyDelUpdate] = "向被删除的参与者发送通知，并向剩余参与者发送日程更新通知"
	message[MeetingRoomEditNotifyUpdate] = "向现有参与者发送会议更新通知"
	message[MeetingRoomEditNotifyAddDel] = "向新增参与者发送会议邀请，向被删除的参与者发送通知"
	message[MeetingRoomEditNotifyAddDelUpdate] = "向新增参与者发送会议邀请，向被删除的参与者发送通知，并向剩余参与者发送日程更新通知"
	message[MeetingRoomEditNoChange] = "编辑详情无修改"
	message[MeetingRoomReserveTimeoutFailed] = "预约失败，请稍后再试"

	// 工作总结
	message[InvalidWorkSummaryId] = "无效的工作总结"
	message[ServicePlatformUserNotFound] = "服务号没有关联业务用户"
}

func MapErrMsg(errCode uint32) string {
	if msg, ok := message[errCode]; ok {
		return msg
	} else {
		return "服务器开小差啦,稍后再来试一试"
	}
}

func IsCodeErr(errCode uint32) bool {
	if _, ok := message[errCode]; ok {
		return true
	} else {
		return false
	}
}
