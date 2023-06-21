package getcdv3

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
	"websocket-pool/entity"
	"websocket-pool/global"
	"websocket-pool/pkg/utils"

	"github.com/spf13/cast"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type RegEtcd struct {
	cli    *clientv3.Client
	ctx    context.Context
	cancel context.CancelFunc
	key    string
}

var rEtcd RegEtcd

// "%s:///%s/"
func GetPrefix(schema, serviceName string) string {
	return fmt.Sprintf("%s:///%s/", schema, serviceName)
}

// "%s:///%s"
func GetPrefix4Unique(schema, serviceName string) string {
	return fmt.Sprintf("%s:///%s", schema, serviceName)
}

// "%s:///%s/" ->  "%s:///%s:ip:port"
func RegisterEtcd4Unique(userNode *entity.NodeEntity) error {
	// serviceName = userNode.NodeName + ":" + net.JoinHostPort(userNode.Ip, strconv.Itoa(userNode.Port))
	return RegisterEtcd(userNode)
}

func GetTarget(schema, myHost string, myPort int, serviceName string) string {
	return GetPrefix4Unique(schema, serviceName) + ":" + net.JoinHostPort(myHost, strconv.Itoa(myPort)) + "/"
}

//etcdAddr separated by commas
func RegisterEtcd(userNode *entity.NodeEntity) (err error) {
	operationID := utils.OperationIDGenerator()
	myHost := userNode.Ip
	myPort := userNode.Port
	serviceName := userNode.NodeName
	schema := global.GVA_CONFIG.Etcd.Schema
	etcdAddr := strings.Join(global.GVA_CONFIG.Etcd.Addr, ",")
	args := schema + " " + etcdAddr + " " + myHost + " " + serviceName + " " + utils.Int32ToString(int32(myPort))
	ttl := global.GVA_CONFIG.Etcd.TTL * 3
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: strings.Split(etcdAddr, ","), DialTimeout: 5 * time.Second})

	// global.GVA_LOG.Info("", zap.String("operationID", operationID), zap.String("args", args), zap.Int("ttl", ttl))
	if err != nil {
		return fmt.Errorf("create etcd clientv3 client failed, errmsg:%v, etcd addr:%s", err, etcdAddr)
	}

	//lease
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ttl)*time.Second)
	defer cancel()
	_, err = cli.Status(ctx, strings.Split(etcdAddr, ",")[0])
	if err != nil {
		return fmt.Errorf("etcd server not found.%v", etcdAddr)
	}

	// if clientv3.StatusResponse.
	resp, err := cli.Grant(context.Background(), int64(ttl))
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
				if ok {
					global.GVA_LOG.Debug("KeepAlive kresp ok", zap.String("operationID", operationID), zap.Any("pv", pv), zap.Any("agrs", args))
				} else {
					// global.GVA_LOG.Error("KeepAlive kresp failed", zap.String("operationID", operationID), zap.Any("pv", pv), zap.Any("agrs", args))
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

	rEtcd = RegEtcd{
		ctx:    ctx,
		cli:    cli,
		cancel: cancel,
		key:    serviceKey,
	}

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
		Endpoints:   strings.Split(etcdAddr, ","),
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		panic(err.Error())
	}
	//lease
	if _, err := cli.Put(context.Background(), key, conf); err != nil {
		fmt.Println("panic, params: ")
		panic(err.Error())
	}

}

func (m *RegEtcd) RegisterConf() {
	// bytes, err := yaml.Marshal(global.GVA_CONFIG.Etcd)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// secretMD5 := utils.Md5(global.GVA_CONFIG.Etcd.Secret)
	// confBytes, err := utils.AesEncrypt(bytes, []byte(secretMD5[0:16]))
	// if err != nil {
	// 	panic(err.Error())
	// }
	// fmt.Println("start register", secretMD5, GetPrefix(global.GVA_CONFIG.Etcd.Schema, config.ConfName))
	// registerConf(GetPrefix(global.GVA_CONFIG.Etcd.Schema, "Websocket-pool"), string(confBytes))
	// fmt.Println("etcd register conf ok")
}

func GetServers(ctx context.Context) (serviceEndpoints []*entity.NodeEntity, err error) {
	resp, err := rEtcd.cli.Get(ctx, rEtcd.key, clientv3.WithPrefix())
	if err != nil {
		global.GVA_LOG.Error("GetServers failed", zap.String("key", rEtcd.key), zap.Error(err))
		return nil, fmt.Errorf("GetServers failed, errmsg:%v， key:%s", err, rEtcd.key)
	}
	for _, v := range resp.Kvs {
		value := string(v.Value)
		socketSlice := strings.Split(value, ":")
		serviceEndpoints = append(serviceEndpoints, &entity.NodeEntity{
			NodeName: global.GVA_CONFIG.System.Name,
			Ip:       socketSlice[0],
			Port:     cast.ToInt(socketSlice[1]),
		})
	}
	return
}
