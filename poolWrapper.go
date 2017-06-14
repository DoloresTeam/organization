package organization

import (
	"github.com/doloresteam/organization/ldap-pool"
	ldap "gopkg.in/ldap.v2"
)

func (org *Organization) poolConn() (*ldappool.PoolConn, error) {
	return org.pool.Get()
}

func (org *Organization) Add(addRequest *ldap.AddRequest) error {
	c, err := org.poolConn()
	if err != nil {
		return err
	}
	return c.Add(addRequest)
}

func (org *Organization) Del(delRequest *ldap.DelRequest) error {
	c, err := org.poolConn()
	if err != nil {
		return err
	}
	return c.Del(delRequest)
}

func (org *Organization) Modify(modifyRequest *ldap.ModifyRequest) error {
	c, err := org.poolConn()
	if err != nil {
		return err
	}
	return c.Modify(modifyRequest)
}

func (org *Organization) Compare(dn, attribute, value string) (bool, error) {
	c, err := org.poolConn()
	if err != nil {
		return false, err
	}
	return c.Compare(dn, attribute, value)
}

func (org *Organization) PasswordModify(passwordModifyRequest *ldap.PasswordModifyRequest) (*ldap.PasswordModifyResult, error) {
	c, err := org.poolConn()
	if err != nil {
		return nil, err
	}
	return c.PasswordModify(passwordModifyRequest)
}

func (org *Organization) Search(searchRequest *ldap.SearchRequest) (*ldap.SearchResult, error) {
	c, err := org.poolConn()
	if err != nil {
		return nil, err
	}
	return c.Search(searchRequest)
}
func (org *Organization) SearchWithPaging(searchRequest *ldap.SearchRequest, pagingSize uint32) (*ldap.SearchResult, error) {
	c, err := org.poolConn()
	if err != nil {
		return nil, err
	}
	return c.SearchWithPaging(searchRequest, pagingSize)
}
