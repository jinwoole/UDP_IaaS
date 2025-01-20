package vm

import "fmt"

func (d *Data) CreateDomain(xmlConfig string) error {
	domain, err := d.conn.DomainDefineXML(xmlConfig)
	if err != nil {
		return fmt.Errorf("failed to define domain: %v", err)
	}

	err = domain.Create()
	if err != nil {
		return fmt.Errorf("failed to start domain: %v", err)
	}

	return nil
}
