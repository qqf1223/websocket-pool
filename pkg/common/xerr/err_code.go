package xerr

// OK 成功返回
const OK uint32 = 0

/**(前3位代表业务,后三位代表具体功能, 每个模块分配100个错误码)**/

// 全局错误码

const ServerCommonError uint32 = 100001
const RequestParamError uint32 = 100002
const TokenExpireError uint32 = 100003
const TokenGenerateError uint32 = 100004
const DbError uint32 = 100005
const GetTeamIdError = 100006       // 无法获取team_id
const GetBaseUserIdError = 100007   // 无法获取注册用户user_id
const ImError = 100008              // im端返回错误
const MindDocError = 100009         // minddoc返回错误
const UploadFileToCosError = 100010 // 上传文件到cos发生错误
const JsonMarshalError = 100011     // json加密错误
const JsonUnMarshalError = 100012   // json解密错误

//扫码登录模块 101

const GetAppError = 101001         // 获取App失败
const GetOAuthQrError = 101002     // 获取QrCode失败
const QrCodeExpiredError = 101003  // 该二维码已经过期
const QrCodeScannedError = 101004  // 该二维码已经被扫描
const InvalidClientError = 101005  // 无效的客户端凭证
const GetTokenError = 101006       // 获取token失败
const SetLoginStatusError = 101007 // 设置登录状态失败
const JwtInvalidError = 101008     // jwt作废失败
const StateCodeError = 101009      // 无效的授权码

//注册模块 102

const VerificationCodeError = 102001 // 验证码错误
const UserAlreadyRegistered = 102002 // 用户已经注册
const RegisterNameError = 102003     // 用户名异常

// 群组 103

const GroupMemberNumberError = 103001 // 群成员数不能少于3人
const GroupMemberExceedError = 103002 // 群成员数超过最大限制
const QueryGroupByIdError = 103003    // 查询群组异常
const GroupUserStatusError = 103004   // 群组内用户状态异常
const GroupCodeNotExist = 103005      // 话题ID不存在， 请重新输入

const GroupGoalNoFinishTask = 103005          // 群组目标还有未完成任务
const GroupGoalFinishByOther = 103006         // 本目标已被其他用户完成
const GroupGoalNotAddTask = 103007            // 本目标已完成，无法添加任务
const GroupGoalTaskAuthFail = 103008          // 暂无权限修改
const GroupGoalInvalidMattersId = 103009      // 无效的事项ID
const GroupGoalCanNotUpdatePrincipal = 103010 // 有任务设置了金币奖励，暂不支持变更目标负责人
const GroupTaskCanNotMove = 103011            // 本任务设置了金币奖励，暂不支持移动

// chat 104

const NotInGroupError = 104001              // 非群成员
const ConversationNotExistError = 104002    //会话不存在
const ConversationExError = 104003          // 会话ex错误
const ConversationPinnedOrderError = 104004 // 排序字段错误

// webhook 105

const ReceiveUserError = 105001      // 接收人错误
const WebHookUserIdError = 105002    // 预警机器人不存在
const WebHookUrlError = 105003       // 请求地址不存在
const WebHookUserParamError = 105004 // 请求用户异常

// 云文档@提醒 106

const GetPermissionError = 106001 // 获取云文档权限失败

// 用户签到 107

const MemberSignedError = 107001 // 用户已经签过到了 每天签到一次

// 组织事项 108

const OrganizationMattersDelError = 108001      // 该组织事项还有子集， 无法删除
const OrganizationMattersCircularError = 108002 // 循环引用， 移动的父节点不能是原节点的子节点

// 群组目标、任务 109

const TaskMoveToGoalFinishedError = 109001       // 本目标已完成，任务不能移入
const GoalSameError = 109002                     // 移入目标是当前目标， 不可移动
const GoldBalanceNotSufficient = 109003          // 目标负责人的金币余额不足
const GoldRelationNoData = 109004                // 金币关系表未查询到数据, 请稍后再试
const GoldAccountInfoNotFound = 109005           // 未查询到账户信息， 请稍后再试
const GoldExecutorNotNull = 109006               // 金币传参时任务执行人必传
const GoldAccountNotNull = 109007                // 数据错误， 金币账户不存在
const GoldGoalPrincipalNotNull = 109008          // 请先设置目标负责人后才可添加金币
const GoldPrincipalNotNull = 109009              // 金币传参时目标负责人必传
const PrincipalGoldBalanceNotSufficient = 109010 // 目标创建人的金币余额不足
const GoldBalanceGreaterThanZero = 109011        // 更新时金币余额必须大于0

// 会议

const ForbiddenChangePattern = 110001            // 非主持人禁止操作模式切换
const MeetingPatternError = 110002               // 会议模式错误
const ForbiddenProcessSpeakRequest = 110003      // 非主持人不能同意或拒绝
const SetAskCacheError = 110004                  // 设置申请缓存时发生错误
const ProcessSpeakCacheError = 110005            // 处理用户申请发言时发生错误
const ForbiddenOperation = 110006                // 禁止操作
const MeetingEnded = 110007                      // 会议已结束
const MeetingRoomReserveError = 110008           // 该时段无法预约，请重新预约// 工作总结
const MeetingRoomPartOfReserverError = 110009    // 部分会议室无法预约，确定保存吗？
const MeetingRoomNoPrivilegeError = 110010       // 会议室非本人创建
const MeetingRoomEditNotifyAdd = 110011          //新增用户发送会议邀请
const MeetingRoomEditNotifyDel = 110012          //向被删除的参与者发送通知
const MeetingRoomEditNotifyAddUpdate = 110013    //向新增参与者发送会议邀请，并向现有参与者发送会议更新通知
const MeetingRoomEditNotifyDelUpdate = 110014    //向被删除的参与者发送通知，并向剩余参与者发送日程更新通知
const MeetingRoomEditNotifyUpdate = 110015       //向现有参与者发送会议更新通知
const MeetingRoomEditNotifyAddDel = 110016       //向新增参与者发送会议邀请，向被删除的参与者发送通知
const MeetingRoomEditNotifyAddDelUpdate = 110017 //向新增参与者发送会议邀请，向被删除的参与者发送通知，并向剩余参与者发送日程更新通知
const MeetingRoomShareNoSelect = 110018          //请选择分享人或话题
const MeetingRoomEditNoChange = 110019           // 编辑详情无修改
const MeetingRoomReserveTimeoutFailed = 110020   // 预约失败给友好提示-验证超时

const InvalidWorkSummaryId = 111001 // 无效的工作总结
const ServicePlatformUserNotFound = 999999
