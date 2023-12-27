package master

//
// func (m *Master) loadResource() error {
// 	resp, err := m.etcdCli.Get(
// 		context.Background(), "",
// 		clientv3.WithPrefix(),
// 		clientv3.WithSerializable(),
// 	)
// 	if err != nil {
// 		return errors.New("etcd get failed")
// 	}
//
// 	resources := make(map[string]*ResourceSpec)
// 	for _, kv := range resp.Kvs {
// 		r, err := decode(kv.Value)
// 		if err == nil && r != nil {
// 			resources[r.Name] = r
// 		}
// 	}
// 	m.logger.Info("leader init load resource",
// 		zap.Int("lenth", len(m.resources)))
//
// 	m.rlock.Lock()
// 	defer m.rlock.Unlock()
// 	m.resources = resources
//
// 	return nil
// }
