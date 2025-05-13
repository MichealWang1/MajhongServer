package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/panjf2000/gnet/v2"
	"kxmj.common/codes"
	"kxmj.common/log"
	"kxmj.common/net"
	"kxmj.common/proto/gateway_pb"
	"kxmj.common/proto/lobby_pb"
	"kxmj.common/redis_cache"
	"kxmj.common/service"
	"kxmj.gateway/config"
	"time"
)

type Outer struct {
	gnet.BuiltinEventEngine                        // gnet接口
	eng                     gnet.Engine            // gnet实例
	gateway                 *Gateway               // 网关管理类
	config                  *config.OuterConfig    // 网关配置
	codec                   *net.OuterCodec        // 外网解包器
	sessions                map[string]net.Session // 客户端链接信息列表
	userSessions            map[uint32]net.Session // 用户链接信息列表
	inputChan               chan net.MsgContext    // 接收消息管道
	outputChan              chan *net.Message      // 发送消息管道
	openSession             chan net.Session       // 客户端连接管道
	closeSession            chan net.Session       // 客户端断开管道
	closeChan               chan struct{}          // 程序结束管道
}

func NewOuter(gateway *Gateway, config *config.OuterConfig) *Outer {
	return &Outer{
		gateway:      gateway,
		config:       config,
		codec:        &net.OuterCodec{},
		sessions:     make(map[string]net.Session, 0),
		userSessions: make(map[uint32]net.Session, 0),
		inputChan:    make(chan net.MsgContext, 10000),
		outputChan:   make(chan *net.Message, 10000),
		openSession:  make(chan net.Session, 1024),
		closeSession: make(chan net.Session, 1024),
		closeChan:    make(chan struct{}, 1),
	}
}

func (o *Outer) ToOuter(msg *net.Message) {
	o.outputChan <- msg
}

func (o *Outer) Start() {
	go func() {
		err := gnet.Run(o, fmt.Sprintf("tcp://:%d", o.config.Port), gnet.WithMulticore(false))
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		for {
			select {
			case <-o.closeChan:
				return
			case msg := <-o.inputChan:
				o.input(msg)
			case msg := <-o.outputChan:
				o.output(msg)
			case session := <-o.openSession:
				o.open(session)
			case session := <-o.closeSession:
				o.close(session)
			}
		}
	}()
}

func (o *Outer) Close() {
	close(o.closeChan)
}

func (o *Outer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	o.eng = eng
	return
}

func (o *Outer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	session := net.NewOuterSession(uuid.New().String(), c, 0)
	c.SetContext(session)
	o.openSession <- session
	return
}

func (o *Outer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	session := c.Context().(*net.OuterSession)
	o.closeSession <- session
	return
}

func (o *Outer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	for {
		msg, err := o.codec.Decode(c)
		if err != nil {
			if err == net.ErrMagicNumberPacket {
				err = c.Close()
				if err != nil {
					log.Sugar().Errorf("Decode err:%v", err)
				}
			}
			break
		}
		session := c.Context().(*net.OuterSession)
		ctx := net.NewOuterContext(session, msg, o.codec)
		o.inputChan <- ctx
	}
	return
}

func (o *Outer) input(ctx net.MsgContext) {
	log.Sugar().Info(fmt.Sprintf("[M:%dT:%dSID:%dU:%dL:%d] <---", ctx.Request().MsgId, ctx.Request().SvrType, ctx.Request().SvrId, ctx.Request().UserId, ctx.Request().Length()))
	switch ctx.Request().SvrType {
	case service.GatewayService: // 网关消息
		if ctx.Request().MsgId == uint16(gateway_pb.MID_LOGIN) {
			// 网关登录
			o.login(ctx)
		} else if ctx.Request().MsgId == uint16(gateway_pb.MID_PING) {
			// ping
			o.pong(ctx)
		}
	default:
		// 丢弃没有 userId 的包
		userId := ctx.Session().UserId()
		if userId == 0 {
			log.Sugar().Error(fmt.Sprintf("err msg:%d", ctx.Request().MsgId))
			return
		}

		session, has := o.userSessions[userId]
		if has == false {
			return
		}

		// 如果用户被踢掉,不回传消息
		if session.SessionId() != ctx.Session().SessionId() {
			return
		}

		ctx.Request().UserId = ctx.Session().UserId()
		o.gateway.ToInner(ctx.Request())
	}
}

func (o *Outer) login(ctx net.MsgContext) {
	payload := &gateway_pb.LoginReq{}
	err := ctx.ShouldBind(payload)
	if err != nil {
		log.Sugar().Error(fmt.Sprintf("Decode err: %v", err))
		return
	}

	ok, _ := redis_cache.GetCache().GetTokenCache().CheckToken(context.Background(), int(payload.UserId), payload.Token)
	if ok == false {
		log.Sugar().Error(fmt.Sprintf("[U:%d] get token: %v\n", payload.UserId, err))
		err = ctx.Send(&net.Message{
			MsgId:   uint16(gateway_pb.MID_ERR),
			SvrType: ctx.Request().SvrType,
			SvrId:   ctx.Request().SvrId,
			UserId:  payload.UserId,
			Data: net.Marshal(&gateway_pb.Err{
				Code:        uint32(codes.AuthorizeFailed),
				Msg:         codes.GetMessage(codes.AuthorizeFailed),
				OriginMsgId: uint32(ctx.Request().MsgId),
			}),
		})

		if err != nil {
			log.Sugar().Infof("Send message err:%v", err)
		}

		// 2秒后关闭连接
		time.AfterFunc(time.Second*2, func() {
			ctx.Session().SetUserId(0)
			err = ctx.Session().Close()
			if err != nil {
				log.Sugar().Errorf("close session err:%v", err)
			}
		})
		return
	}

	oldSession, has := o.userSessions[payload.UserId]
	if has && ctx.Session().SessionId() != oldSession.SessionId() {
		oldSession.SetUserId(0)
		err = oldSession.Close()
		if err != nil {
			log.Sugar().Errorf("close session err:%v", err)
		}
	}

	ctx.Session().SetUserId(payload.UserId)
	o.userSessions[payload.UserId] = ctx.Session()

	// 通知位置服务用户上线
	o.gateway.ToInner(&net.Message{
		MsgId:   uint16(lobby_pb.MID_ON_LINE),
		SvrType: service.LobbyService,
		SvrId:   0,
		UserId:  ctx.Session().UserId(),
		Data: net.Marshal(&lobby_pb.Online{
			UserId:  payload.UserId,
			SvrType: uint32(o.gateway.Inner.self.Type),
			SvrId:   uint32(o.gateway.Inner.self.Id),
		}),
	})

	// 返回登陆成功消息
	o.loginSuccess(ctx, payload)
}

func (o *Outer) loginSuccess(ctx net.MsgContext, payload *gateway_pb.LoginReq) {
	log.Sugar().Info(fmt.Sprintf("[M:%dT:%dSID:%dU:%dL:%d] --->", ctx.Request().MsgId, ctx.Request().SvrType, ctx.Request().SvrId, ctx.Request().UserId, ctx.Request().Length()))
	// 回复登录成功消息
	err := ctx.Send(&net.Message{
		MsgId:   uint16(gateway_pb.MID_LOGIN),
		SvrType: ctx.Request().SvrType,
		SvrId:   ctx.Request().SvrId,
		UserId:  ctx.Session().UserId(),
		Data: net.Marshal(&gateway_pb.LoginResp{
			Code:   codes.Success,
			UserId: payload.UserId,
		}),
	})

	if err != nil {
		log.Sugar().Errorf("Send message err:%v", err)
	}
}

func (o *Outer) pong(ctx net.MsgContext) {
	log.Sugar().Info(fmt.Sprintf("[M:%dT:%dSID:%dU:%dL:%d] --->", ctx.Request().MsgId, ctx.Request().SvrType, ctx.Request().SvrId, ctx.Request().UserId, ctx.Request().Length()))

	payload := &gateway_pb.Ping{}
	err := ctx.ShouldBind(payload)
	if err != nil {
		log.Sugar().Error(fmt.Sprintf("Decode err:%v", err))
		return
	}

	err = ctx.Send(&net.Message{
		MsgId:   uint16(gateway_pb.MID_PING),
		SvrType: ctx.Request().SvrType,
		SvrId:   ctx.Request().SvrId,
		UserId:  ctx.Session().UserId(),
		Data: net.Marshal(&gateway_pb.Pong{
			ClientTime: payload.ClientTime,
		}),
	})

	if err != nil {
		log.Sugar().Errorf("Send message err:%v", err)
	}
}

func (o *Outer) output(msg *net.Message) {
	log.Sugar().Info(fmt.Sprintf("[M:%dT:%dSID:%dU:%d:L:%d] --->", msg.MsgId, msg.SvrType, msg.SvrId, msg.UserId, msg.Length()))
	if msg.UserId <= 0 {
		log.Sugar().Errorf("error output message")
		return
	}

	session, has := o.userSessions[msg.UserId]
	if has == false {
		return
	}

	data, err := o.codec.Encode(msg)
	if err != nil {
		log.Sugar().Errorf("Encode err:%v", err)
	}

	err = session.Send(data)
	if err != nil {
		log.Sugar().Errorf("Send message err:%v", err)
	}
}

func (o *Outer) open(session net.Session) {
	log.Sugar().Infof("client %s is connected", session.RemoteAddr())
	o.sessions[session.SessionId()] = session
}

func (o *Outer) close(session net.Session) {
	log.Sugar().Infof("client %s is disconnected", session.RemoteAddr())
	delete(o.sessions, session.SessionId())
	if session.UserId() > 0 {
		delete(o.userSessions, session.UserId())

		// 通知位置服务用户下线
		o.gateway.ToInner(&net.Message{
			MsgId:   uint16(lobby_pb.MID_OFF_LINE),
			SvrType: service.LobbyService,
			SvrId:   0,
			UserId:  session.UserId(),
			Data: net.Marshal(&lobby_pb.Offline{
				UserId:  session.UserId(),
				SvrType: uint32(o.gateway.Inner.self.Type),
				SvrId:   uint32(o.gateway.Inner.self.Id),
			}),
		})
	}
}
