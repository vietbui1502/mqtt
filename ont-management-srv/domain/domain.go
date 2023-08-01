package domain

type DomainRepository interface {
	FindDomainCategory(domain string) (string, error)
}
