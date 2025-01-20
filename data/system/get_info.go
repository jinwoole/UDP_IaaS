package system

func (d *Data) GetSystemInfo() (*SystemInfo, error) {
	nodeInfo, err := d.conn.GetNodeInfo()
	if err != nil {
		return nil, err
	}

	hostname, err := d.conn.GetHostname()
	if err != nil {
		return nil, err
	}

	version, err := d.conn.GetVersion()
	if err != nil {
		return nil, err
	}

	virtType, err := d.conn.GetType()
	if err != nil {
		return nil, err
	}

	freeMem, err := d.conn.GetFreeMemory()
	if err != nil {
		return nil, err
	}

	maxVcpus, err := d.conn.GetMaxVcpus("")
	if err != nil {
		return nil, err
	}

	isConnected, _ := d.conn.IsAlive()

	return &SystemInfo{
		CPUCount:    uint(nodeInfo.Cores),
		MemorySize:  nodeInfo.Memory,
		Hostname:    hostname,
		NodeInfo:    nodeInfo,
		CPUModel:    nodeInfo.Model,
		CPUArch:     string(nodeInfo.Model),
		LibVersion:  uint64(version),
		VirtType:    virtType,
		FreeMemory:  freeMem,
		MaxVCPUs:    uint(maxVcpus),
		IsConnected: isConnected,
	}, nil
}
