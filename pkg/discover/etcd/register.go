package getcdv3

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
	"websocket-pool/global"
	"websocket-pool/pkg/utils"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type RegEtcd struct {
	cli    *clientv3.Client
	ctx    context.Context
	cancel context.CancelFunc
	key    string
}

var rEtcd *RegEtcd

// "%s:///%s/"
func GetPrefix(schema, serviceName string) string {
	return fmt.Sprintf("%s:///%s/", schema, serviceName)
}

// "%s:///%s"
func GetPrefix4Unique(schema, serviceName string) string {
	return fmt.Sprintf("%s:///%s", schema, serviceName)
}

// "%s:///%s/" ->  "%s:///%s:ip:port"
func RegisterEtcd4Unique(schema, etcdAddr, myHost string, myPort int, serviceName string, ttl int) error {
	serviceName = serviceName + ":" + net.JoinHostPort(myHost, strconv.Itoa(myPort))
	return RegisterEtcd(schema, etcdAddr, myHost, myPort, serviceName, ttl)
}

func GetTarget(schema, myHost string, myPort int, serviceName string) string {
	return GetPrefix4Unique(schema, serviceName) + ":" + net.JoinHostPort(myHost, strconv.Itoa(myPort)) + "/"
}

//etcdAddr separated by commas
func RegisterEtcd(schema, etcdAddr, myHost string, myPort int, serviceName string, ttl int) error {
	operationID := utils.OperationIDGenerator()
	args := schema + " " + etcdAddr + " " + myHost + " " + serviceName + " " + utils.Int32ToString(int32(myPort))
	ttl = ttl * 3
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: strings.Split(etcdAddr, ","), DialTimeout: 5 * time.Second})

	global.GVA_LOG.Info("", zap.String("operationID", operationID), zap.String("args", args), zap.Int("ttl", ttl))
	if err != nil {
		global.GVA_LOG.Error("create etcd clientv3 client failed", zap.String("etcdAddr", etcdAddr), zap.Error(err))
		return fmt.Errorf("create etcd clientv3 client failed, errmsg:%v, etcd addr:%s", err, etcdAddr)
	}

	//lease
	ctx, cancel := context.WithCancel(context.Background())
	resp, err := cli.Grant(ctx, int64(ttl))
	if err != nil {
		global.GVA_LOG.Error("grant failed", zap.String("operationID", operationID), zap.Error(err))
		return fmt.Errorf("grant failed")
	}
	global.GVA_LOG.Info("grant ok", zap.String("operationID", operationID), zap.Any("respID", resp.ID))

	//  schema:///serviceName/ip:port ->ip:port
	serviceValue := net.JoinHostPort(myHost, strconv.Itoa(myPort))
	serviceKey := GetPrefix(schema, serviceName) + serviceValue

	//set key->value
	if _, err := cli.Put(ctx, serviceKey, serviceValue, clientv3.WithLease(resp.ID)); err != nil {
		global.GVA_LOG.Error("put failed", zap.String("key", serviceKey), zap.String("serviceValue", serviceValue), zap.Error(err))
		return fmt.Errorf("put failed, errmsg:%v， key:%s, value:%s", err, serviceKey, serviceValue)
	}

	//keepalive
	kresp, err := cli.KeepAlive(ctx, resp.ID)
	if err != nil {
		global.GVA_LOG.Error("keepalive failed", zap.String("operationID", operationID), zap.Any("leaseId", resp.ID), zap.Error(err))
		return fmt.Errorf("keepalive failed, errmsg:%v, lease id:%d", err, resp.ID)
	}
	// log.Info(operationID, "RegisterEtcd ok ", args)

	go func() {
		for {
			select {
			case pv, ok := <-kresp:
				if ok == true {
					global.GVA_LOG.Debug("KeepAlive kresp ok", zap.String("operationID", operationID), zap.Any("pv", pv), zap.Any("agrs", args))
				} else {
					global.GVA_LOG.Error("KeepAlive kresp failed", zap.String("operationID", operationID), zap.Any("pv", pv), zap.Any("agrs", args))
					// log.Error(operationID, "KeepAlive kresp failed ", pv, args)
					t := time.NewTicker(time.Duration(ttl/2) * time.Second)
					for {
						select {
						case <-t.C:
						}
						ctx, _ := context.WithCancel(context.Background())
						resp, err := cli.Grant(ctx, int64(ttl))
						if err != nil {
							global.GVA_LOG.Error("grant failed", zap.String("operationID", operationID), zap.Any("args", args), zap.Error(err))
							continue
						}

						if _, err := cli.Put(ctx, serviceKey, serviceValue, clientv3.WithLease(resp.ID)); err != nil {
							global.GVA_LOG.Error("etcd Put failed ", zap.String("operationID", operationID), zap.Any("args", args), zap.Any("respID", resp.ID), zap.Error(err))
							continue
						} else {
							global.GVA_LOG.Info("etcd Put ok", zap.String("operationID", operationID), zap.Any("args", args), zap.Any("respID", resp.ID))
						}
					}
				}
			}
		}
	}()

	rEtcd = &RegEtcd{ctx: ctx,
		cli:    cli,
		cancel: cancel,
		key:    serviceKey}

	return nil
}

func UnRegisterEtcd() {
	//delete
	rEtcd.cancel()
	rEtcd.cli.Delete(rEtcd.ctx, rEtcd.key)
}

func registerConf(key, conf string) {
	etcdAddr := strings.Join(global.GVA_CONFIG.Etcd.Addr, ",")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: strings.Split(etcdAddr, ","), DialTimeout: 5 * time.Second})

	if err != nil {
		panic(err.Error())
	}
	//lease
	if _, err := cli.Put(context.Background(), key, conf); err != nil {
		fmt.Println("panic, params: ")
		panic(err.Error())
	}

}

func RegisterConf() {
	// bytes, err := yaml.Marshal(config.Config)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// secretMD5 := utils.Md5(global.GVA_CONFIG.Etcd.Secret)
	// confBytes, err := utils.AesEncrypt(bytes, []byte(secretMD5[0:16]))
	// if err != nil {
	// 	panic(err.Error())
	// }
	// fmt.Println("start register", secretMD5, GetPrefix(config.Config.Etcd.EtcdSchema, config.ConfName))
	// registerConf(GetPrefix(config.Config.Etcd.EtcdSchema, "Websocket-pool"), string(confBytes))
	// fmt.Println("etcd register conf ok")
}
