// @Contact:    huaxinrui
// @Time:       2019/9/26 上午10:18

package encrypt

import (
	"code.byted.org/byte_security/platform_api/service/conf"
	"code.byted.org/security/gokms"
	"sync"
)

const OSM = "PSM"
const CMK = "CMK"

type KmsClient struct {
	lock   *sync.Mutex
	client *gokms.KMSClient
	env    map[string]string
}

func (k *KmsClient) Decrypt(secret string) (string, error) {
	return k.client.Decrypt(secret)
}

func (k *KmsClient) Encrypt(plain string) (string, error) {
	return k.client.Encrypt(k.env[CMK], plain)
}

func (b *KmsClient) addEnv(k, v string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.env[k] = v
}

func (b *KmsClient) updateEnv(m map[string]string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	for k, v := range m {
		b.env[k] = v
	}
}

func (b *KmsClient) initEnv(osm, cmk string) {
	env := map[string]string{
		OSM: osm,
		CMK: cmk, //,
	}
	b.updateEnv(env)
}

func InitEncrypt() *KmsClient {
	k := new(KmsClient)
	k.lock = new(sync.Mutex)
	k.env = make(map[string]string)
	k.initEnv(conf.Config.Psm, conf.Config.CMK)
	client, err := gokms.NewKMSClient(k.env[OSM])
	if err != nil {
		panic(err)
	}
	k.client = client
	return k
}
